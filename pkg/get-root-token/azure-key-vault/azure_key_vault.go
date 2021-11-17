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

	vaultapi "kubevault.dev/apimachinery/apis/kubevault/v1alpha1"
	"kubevault.dev/cli/pkg/get-root-token/api"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/pkg/errors"
)

const (
	ContentTypePassword = "password"
)

type TokenInfo struct {
	cred         *azidentity.DefaultAzureCredential
	vaultBaseUrl string
}

var _ api.TokenInterface = &TokenInfo{}

func New(spec *vaultapi.AzureKeyVault) (*TokenInfo, error) {
	if spec == nil {
		return nil, errors.New("azureKeyVault spec is empty")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	return &TokenInfo{
		cred:         cred,
		vaultBaseUrl: spec.VaultBaseURL,
	}, nil
}

// Token will require exportation of AZURE_CLIENT_ID, AZURE_CLIENT_SECRET, AZURE_TENANT_ID
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
