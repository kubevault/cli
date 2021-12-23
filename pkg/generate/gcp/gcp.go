/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Community License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Community-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gcp

import (
	"context"
	"fmt"
	"strings"

	engineapi "kubevault.dev/apimachinery/apis/engine/v1alpha1"
	vaultapi "kubevault.dev/apimachinery/apis/kubevault/v1alpha1"
	enginecs "kubevault.dev/apimachinery/client/clientset/versioned/typed/engine/v1alpha1"
	vaultcs "kubevault.dev/apimachinery/client/clientset/versioned/typed/kubevault/v1alpha1"
	policycs "kubevault.dev/apimachinery/client/clientset/versioned/typed/policy/v1alpha1"
	"kubevault.dev/cli/pkg/generate/api"

	"github.com/go-errors/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

var available = map[string]bool{
	"token":            true,
	"private_key_data": true,
	"key_type":         true,
	"key_algorithm":    true,
}

type SecretObject struct {
	ObjectName string                 `json:"objectName,omitempty"`
	SecretPath string                 `json:"secretPath,omitempty"`
	SecretKey  string                 `json:"secretKey,omitempty"`
	Method     string                 `json:"method,omitempty"`
	SecretArgs map[string]interface{} `json:"secretArgs,omitempty"`
}

type GCPGenerator struct {
	role         []string
	srb          *engineapi.SecretRoleBinding
	se           *engineapi.SecretEngine
	keys         map[string]string
	engineClient *enginecs.EngineV1alpha1Client
	vaultClient  *vaultcs.KubevaultV1alpha1Client
	policyClient *policycs.PolicyV1alpha1Client
	clusterName  string
}

var _ api.GeneratorInterface = &GCPGenerator{}

func NewGCPGenerator(role []string, srb *engineapi.SecretRoleBinding, keys map[string]string, engineClient *enginecs.EngineV1alpha1Client, vaultClient *vaultcs.KubevaultV1alpha1Client, policyClient *policycs.PolicyV1alpha1Client, kubeClient *kubernetes.Clientset) (*GCPGenerator, error) {
	gcpRole, err := engineClient.GCPRoles(srb.Namespace).Get(context.TODO(), role[1], metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	se, err := engineClient.SecretEngines(srb.Namespace).Get(context.TODO(), gcpRole.Spec.SecretEngineRef.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	sts, err := kubeClient.AppsV1().StatefulSets(se.Spec.VaultRef.Namespace).Get(context.TODO(), se.Spec.VaultRef.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var clName string
	for _, cont := range sts.Spec.Template.Spec.Containers {
		if cont.Name != vaultapi.VaultUnsealerContainerName {
			continue
		}
		for _, arg := range cont.Args {
			if strings.HasPrefix(arg, "--cluster-name=") {
				clName = arg[1+strings.Index(arg, "="):]
			}
		}
	}

	return &GCPGenerator{
		role:         role,
		srb:          srb,
		se:           se,
		keys:         keys,
		engineClient: engineClient,
		vaultClient:  vaultClient,
		policyClient: policyClient,
		clusterName:  clName,
	}, nil
}

func (g *GCPGenerator) Generate() (string, error) {
	gcpRole, err := g.engineClient.GCPRoles(g.srb.Namespace).Get(context.TODO(), g.role[1], metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	for key := range g.keys {
		if _, ok := available[key]; !ok {
			var klist []string
			for k := range available {
				klist = append(klist, k)
			}
			return "", errors.Errorf("key %s not available for roleKind %s\navailable keys are: %s", key, g.role[0], strings.Join(klist, ", "))
		}
	}

	var object []SecretObject
	for key, mapping := range g.keys {
		doc := g.GetSecretObject(key, mapping, gcpRole)
		object = append(object, *doc)
	}

	data, err := yaml.Marshal(object)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (g *GCPGenerator) GetVaultServerURL() (string, error) {
	vs, err := g.vaultClient.VaultServers(g.se.Spec.VaultRef.Namespace).Get(context.TODO(), g.se.Spec.VaultRef.Name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	address := fmt.Sprintf("%s://%s.%s:8200", vs.Scheme(), vs.Name, vs.Namespace)
	return address, nil
}

func (g *GCPGenerator) GetVaultRoleName() (string, error) {
	vpb, err := g.policyClient.VaultPolicyBindings(g.se.Spec.VaultRef.Namespace).Get(context.TODO(), g.srb.VaultPolicyBindingName(), metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return vpb.Spec.VaultRoleName, nil
}

func (g *GCPGenerator) GetSecretObject(key, mapping string, gcpRole *engineapi.GCPRole) *SecretObject {
	var path string
	sePath := fmt.Sprintf("k8s.%s.%s.%s.%s", g.clusterName, g.se.GetSecretEngineType(), g.se.Namespace, g.se.Name)
	roleName := fmt.Sprintf("k8s.%s.%s.%s", g.clusterName, gcpRole.Namespace, gcpRole.Name)

	if gcpRole.Spec.SecretType == engineapi.GCPSecretServiceAccountKey {
		path = fmt.Sprintf("/%s/roleset/%s/key", sePath, roleName)
	} else {
		path = fmt.Sprintf("/%s/roleset/%s/token", sePath, roleName)
	}

	return &SecretObject{
		ObjectName: mapping,
		SecretPath: path,
		SecretKey:  key,
	}
}
