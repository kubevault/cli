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

package cmds

import (
	"context"
	"fmt"
	"os"

	engineapi "kubevault.dev/apimachinery/apis/engine/v1alpha1"
	vaultapi "kubevault.dev/apimachinery/apis/kubevault/v1alpha2"
	enginecs "kubevault.dev/apimachinery/client/clientset/versioned/typed/engine/v1alpha1"
	engineutil "kubevault.dev/apimachinery/client/clientset/versioned/typed/engine/v1alpha1/util"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	kmapi "kmodules.xyz/client-go/api/v1"
	condutil "kmodules.xyz/client-go/conditions"
)

const (
	VAULT_ADDR = "127.0.0.1:8200"
)

func Fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func modifyStatusCondition(clientGetter genericclioptions.RESTClientGetter, cond kmapi.Condition) error {
	var resourceName string
	switch ResourceName {
	case engineapi.ResourceSecretAccessRequest, engineapi.ResourceSecretAccessRequests:
		resourceName = engineapi.ResourceSecretAccessRequest
	case "":
		resourceName = ""
	default:
		return errors.New("unknown/unsupported resource")
	}

	cfg, err := clientGetter.ToRESTConfig()
	if err != nil {
		return errors.Wrap(err, "failed to read kubeconfig")
	}

	namespace, _, err := clientGetter.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	builder := cmdutil.NewFactory(clientGetter).NewBuilder()

	engineClient, err := enginecs.NewForConfig(cfg)
	if err != nil {
		return err
	}

	r := builder.
		WithScheme(clientsetscheme.Scheme, clientsetscheme.Scheme.PrioritizedVersionsAllGroups()...).
		ContinueOnError().
		NamespaceParam(namespace).DefaultNamespace().
		FilenameParam(false, &FilenameOptions).
		ResourceNames(resourceName, ObjectNames...).
		RequireObject(true).
		Flatten().
		Latest().
		Do()

	err = r.Visit(func(info *resource.Info, err error) error {
		if err != nil {
			return err
		}

		var err2 error
		switch info.Object.(type) {
		case *engineapi.SecretAccessRequest:
			obj := info.Object.(*engineapi.SecretAccessRequest)

			if err = isApplicable(engineClient, obj, cond, obj.Status.Conditions); err != nil {
				return err
			}

			cond.ObservedGeneration = obj.Generation
			err2 = UpdateSecretAccessRequestCondition(engineClient, obj.ObjectMeta, cond)
		default:
			err2 = errors.New("unknown/unsupported type")
		}
		return err2
	})
	return err
}

func UpdateSecretAccessRequestCondition(c enginecs.EngineV1alpha1Interface, req metav1.ObjectMeta, cond kmapi.Condition) error {
	_, err := engineutil.UpdateSecretAccessRequestStatus(
		context.TODO(),
		c,
		req,
		func(in *engineapi.SecretAccessRequestStatus) *engineapi.SecretAccessRequestStatus {
			in.Conditions = condutil.SetCondition(in.Conditions, cond)
			in.ObservedGeneration = req.Generation
			return in
		}, metav1.UpdateOptions{})
	return err
}

func isApplicable(engineClient *enginecs.EngineV1alpha1Client, req *engineapi.SecretAccessRequest, cond kmapi.Condition, conditions []kmapi.Condition) error {
	if cond == secretAccessRevokeCond && !condutil.IsConditionTrue(conditions, engineapi.ConditionRequestExpired) && req.Spec.RoleRef.Kind == engineapi.ResourceKindGCPRole {
		role, err := engineClient.GCPRoles(req.Namespace).Get(context.TODO(), req.Spec.RoleRef.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		if role.Spec.SecretType == engineapi.GCPSecretAccessToken {
			return errors.New("access token is non revocable")
		}
	}

	if cond == secretAccessApprovedCond && condutil.IsConditionTrue(conditions, engineapi.ConditionRequestExpired) {
		return errors.New("failed to approve, request already expired")
	}

	if cond == secretAccessDeniedCond && condutil.IsConditionTrue(conditions, engineapi.ConditionRequestExpired) {
		return errors.New("failed to deny, request already expired")
	}

	if cond == secretAccessDeniedCond && condutil.IsConditionTrue(conditions, condutil.ConditionRequestApproved) {
		return errors.New("failed to deny, request already approved")
	}

	if cond == secretAccessRevokeCond && condutil.IsConditionTrue(conditions, engineapi.ConditionRequestExpired) {
		return errors.New("request already revoked")
	}

	if cond == secretAccessDeniedCond && condutil.IsConditionTrue(conditions, condutil.ConditionRequestDenied) {
		return errors.New("request already denied")
	}

	if cond == secretAccessApprovedCond && condutil.IsConditionTrue(conditions, condutil.ConditionRequestApproved) {
		return errors.New("request already approved")
	}

	return nil
}

func NewVaultClient(vs *vaultapi.VaultServer) (*api.Client, error) {
	cfg := api.DefaultConfig()

	tlsConfig := &api.TLSConfig{
		Insecure: true,
	}

	cfg.Address = fmt.Sprintf("%s://%s", vs.Scheme(), VAULT_ADDR)
	err := cfg.ConfigureTLS(tlsConfig)
	if err != nil {
		return nil, err
	}

	return api.NewClient(cfg)
}
