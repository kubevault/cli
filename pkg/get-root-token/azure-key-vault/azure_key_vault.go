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

package azure_key_vault

import (
	"context"
	"encoding/base64"
	"os"

	vaultapi "kubevault.dev/apimachinery/apis/kubevault/v1alpha1"
	"kubevault.dev/cli/pkg/get-root-token/api"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	ContentTypePassword = "password"
	ClientID            = "AZURE_CLIENT_ID"
	ClientSecret        = "AZURE_CLIENT_SECRET"
	TenantID            = "AZURE_TENANT_ID"
)

type TokenInfo struct {
	cred         *azidentity.DefaultAzureCredential
	vaultBaseUrl string
}

var _ api.TokenInterface = &TokenInfo{}

func New(vs *vaultapi.VaultServer, kubeClient kubernetes.Interface) (*TokenInfo, error) {
	if vs == nil {
		return nil, errors.New("vs spec is empty")
	}

	if vs.Spec.Unsealer.Mode.AzureKeyVault == nil {
		return nil, errors.New("AzureKeyVault mode is nil")
	}

	if kubeClient == nil {
		return nil, errors.New("kubeClient is nil")
	}

	secret, err := kubeClient.CoreV1().Secrets(vs.Namespace).Get(context.TODO(), vs.Spec.Unsealer.Mode.AzureKeyVault.AADClientSecret, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if _, ok := secret.Data["client-id"]; !ok {
		return nil, errors.Errorf("%s not found in secret", ClientID)
	}
	if _, ok := secret.Data["client-secret"]; !ok {
		return nil, errors.Errorf("%s not found in secret", ClientSecret)
	}

	if err = os.Setenv(ClientID, string(secret.Data["client-id"])); err != nil {
		return nil, err
	}
	if err = os.Setenv(ClientSecret, string(secret.Data["client-secret"])); err != nil {
		return nil, err
	}
	if err = os.Setenv(TenantID, vs.Spec.Unsealer.Mode.AzureKeyVault.TenantID); err != nil {
		return nil, err
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	return &TokenInfo{
		cred:         cred,
		vaultBaseUrl: vs.Spec.Unsealer.Mode.AzureKeyVault.VaultBaseURL,
	}, nil
}

func (ti *TokenInfo) Token() (string, error) {
	client, err := azsecrets.NewClient(ti.vaultBaseUrl, ti.cred, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.GetSecret(context.Background(), ti.TokenName(), nil)
	if err != nil {
		return "", err
	}

	if *resp.ContentType != ContentTypePassword {
		return "", errors.Errorf("content type not matched with %v", *resp.ContentType)
	}

	token, err := base64.StdEncoding.DecodeString(*resp.Value)
	if err != nil {
		return "", err
	}

	return string(token), nil
}

func (ti *TokenInfo) TokenName() string {
	return "vault-root-token"
}
