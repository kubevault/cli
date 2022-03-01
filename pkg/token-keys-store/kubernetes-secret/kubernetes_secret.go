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

package kubernetes_secret

import (
	"context"
	"fmt"
	"strings"

	vaultapi "kubevault.dev/apimachinery/apis/kubevault/v1alpha1"
	"kubevault.dev/cli/pkg/token-keys-store/api"

	"github.com/pkg/errors"
	errors2 "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type TokenKeyInfo struct {
	kubeClient kubernetes.Interface
	vs         *vaultapi.VaultServer
}

var _ api.TokenKeyInterface = &TokenKeyInfo{}

func New(vs *vaultapi.VaultServer, kubeClient kubernetes.Interface) (*TokenKeyInfo, error) {
	if vs == nil {
		return nil, errors.New("vs spec is empty")
	}

	if vs.Spec.Unsealer.Mode.KubernetesSecret == nil {
		return nil, errors.New("kubernetes-secret mode is nil")
	}

	if kubeClient == nil {
		return nil, errors.New("kubeClient is nil")
	}

	return &TokenKeyInfo{
		kubeClient: kubeClient,
		vs:         vs,
	}, nil
}

func (ti *TokenKeyInfo) Get(key string) (string, error) {
	secretName := ti.vs.Spec.Unsealer.Mode.KubernetesSecret.SecretName
	secretNamespace := ti.vs.Namespace
	secret, err := ti.kubeClient.CoreV1().Secrets(secretNamespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if _, ok := secret.Data[key]; !ok {
		return "", errors.Errorf("%s not found in secret %s/%s", key, secretNamespace, secretName)
	}

	return string(secret.Data[key]), nil
}

func (ti *TokenKeyInfo) Delete(key string) error {
	secretName := ti.vs.Spec.Unsealer.Mode.KubernetesSecret.SecretName
	secretNamespace := ti.vs.Namespace
	secret, err := ti.kubeClient.CoreV1().Secrets(secretNamespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		if errors2.IsNotFound(err) {
			return nil
		}
		return err
	}

	if _, ok := secret.Data[key]; ok {
		secret.Data[key] = nil
		delete(secret.Data, key)
	}

	_, err = ti.kubeClient.CoreV1().Secrets(secretNamespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	return err
}

func (ti *TokenKeyInfo) Set(key, value string) error {
	secretName := ti.vs.Spec.Unsealer.Mode.KubernetesSecret.SecretName
	secretNamespace := ti.vs.Namespace
	secret, err := ti.kubeClient.CoreV1().Secrets(secretNamespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		if errors2.IsNotFound(err) {
			return nil
		}
		return err
	}

	secret.Data[key] = []byte(value)

	_, err = ti.kubeClient.CoreV1().Secrets(secretNamespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	return err
}

func (ti *TokenKeyInfo) NewTokenName() string {
	sts, err := ti.kubeClient.AppsV1().StatefulSets(ti.vs.Namespace).Get(context.TODO(), ti.vs.Name, metav1.GetOptions{})
	if err != nil {
		return ""
	}

	var keyPrefix string
	for _, cont := range sts.Spec.Template.Spec.Containers {
		if cont.Name != vaultapi.VaultUnsealerContainerName {
			continue
		}
		for _, arg := range cont.Args {
			if strings.HasPrefix(arg, "--key-prefix=") {
				keyPrefix = arg[1+strings.Index(arg, "="):]
			}
		}
	}

	return fmt.Sprintf("%s-root-token", keyPrefix)
}

func (ti *TokenKeyInfo) OldTokenName() string {
	return "vault-root-token"
}

func (ti *TokenKeyInfo) NewUnsealKeyName(id int) (string, error) {
	sts, err := ti.kubeClient.AppsV1().StatefulSets(ti.vs.Namespace).Get(context.TODO(), ti.vs.Name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if int64(id) >= ti.vs.Spec.Unsealer.SecretShares {
		return "", errors.Errorf("unseal-key-%d not available, available id range 0 to %d", id, ti.vs.Spec.Unsealer.SecretShares-1)
	}

	var keyPrefix string
	for _, cont := range sts.Spec.Template.Spec.Containers {
		if cont.Name != vaultapi.VaultUnsealerContainerName {
			continue
		}
		for _, arg := range cont.Args {
			if strings.HasPrefix(arg, "--key-prefix=") {
				keyPrefix = arg[1+strings.Index(arg, "="):]
			}
		}
	}

	return fmt.Sprintf("%s-unseal-key-%d", keyPrefix, id), nil
}

func (ti *TokenKeyInfo) OldUnsealKeyName(id int) (string, error) {
	if int64(id) >= ti.vs.Spec.Unsealer.SecretShares {
		return "", errors.Errorf("unseal-key-%d not available, available id range 0 to %d", id, ti.vs.Spec.Unsealer.SecretShares-1)
	}

	return fmt.Sprintf("vault-unseal-key-%d", id), nil
}

func (ti *TokenKeyInfo) Clean() {
}
