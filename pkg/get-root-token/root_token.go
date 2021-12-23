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

package token

import (
	"errors"

	vaultapi "kubevault.dev/apimachinery/apis/kubevault/v1alpha1"
	"kubevault.dev/cli/pkg/get-root-token/api"
	aws_kms_ssm "kubevault.dev/cli/pkg/get-root-token/aws-kms-ssm"
	azure_key_vault "kubevault.dev/cli/pkg/get-root-token/azure-key-vault"
	google_kms_gcs "kubevault.dev/cli/pkg/get-root-token/google-kms-gcs"
	kubernetes_secret "kubevault.dev/cli/pkg/get-root-token/kubernetes-secret"

	"k8s.io/client-go/kubernetes"
)

func NewTokenInterface(vs *vaultapi.VaultServer, kubeClient kubernetes.Interface) (api.TokenInterface, error) {
	if vs.Spec.Unsealer == nil {
		return nil, errors.New("vaultServer unsealer spec is empty")
	}
	mode := vs.Spec.Unsealer.Mode

	switch true {
	case mode.AwsKmsSsm != nil:
		return aws_kms_ssm.New(vs, kubeClient)
	case mode.GoogleKmsGcs != nil:
		return google_kms_gcs.New(vs, kubeClient)
	case mode.AzureKeyVault != nil:
		return azure_key_vault.New(vs, kubeClient)
	case mode.KubernetesSecret != nil:
		return kubernetes_secret.New(vs, kubeClient)
	}

	return nil, errors.New("unknown/unsupported unsealing mode")
}
