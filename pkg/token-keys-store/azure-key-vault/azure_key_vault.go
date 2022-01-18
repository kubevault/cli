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
	"fmt"
	"os"
	"strings"
	"time"

	vaultapi "kubevault.dev/apimachinery/apis/kubevault/v1alpha1"
	"kubevault.dev/cli/pkg/token-keys-store/api"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
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

type TokenKeyInfo struct {
	cred       *azidentity.DefaultAzureCredential
	kubeClient kubernetes.Interface
	vs         *vaultapi.VaultServer
}

var _ api.TokenKeyInterface = &TokenKeyInfo{}

func New(vs *vaultapi.VaultServer, kubeClient kubernetes.Interface) (*TokenKeyInfo, error) {
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

	clientID, ok := secret.Data["client-id"]
	if ok {
		if err = os.Setenv(ClientID, string(clientID)); err != nil {
			return nil, err
		}
	}

	clientSecret, ok := secret.Data["client-secret"]
	if ok {
		if err = os.Setenv(ClientSecret, string(clientSecret)); err != nil {
			return nil, err
		}
	}

	if err = os.Setenv(TenantID, vs.Spec.Unsealer.Mode.AzureKeyVault.TenantID); err != nil {
		return nil, err
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	return &TokenKeyInfo{
		cred:       cred,
		vs:         vs,
		kubeClient: kubeClient,
	}, nil
}

func (ti *TokenKeyInfo) Get(key string) (string, error) {
	vaultBaseUrl := ti.vs.Spec.Unsealer.Mode.AzureKeyVault.VaultBaseURL
	client, err := azsecrets.NewClient(vaultBaseUrl, ti.cred, nil)
	if err != nil {
		return "", err
	}

	version, err := ti.getLatestVersion(key)
	if err != nil {
		return "", err
	}

	idx := strings.LastIndex(version, "/")
	if idx == -1 {
		return "", errors.New("version id not found")
	}

	resp, err := client.GetSecret(context.Background(), strings.Replace(key, ".", "-", -1), &azsecrets.GetSecretOptions{
		Version: version[idx+1:],
	})
	if err != nil {
		return "", err
	}

	if *resp.ContentType != ContentTypePassword {
		return "", errors.Errorf("content type not matched with %v", *resp.ContentType)
	}

	decoded, err := base64.StdEncoding.DecodeString(*resp.Value)
	if err != nil {
		return "", err
	}

	return string(decoded), nil
}

func (ti *TokenKeyInfo) Delete(key string) error {
	key = strings.Replace(key, ".", "-", -1)

	vaultBaseUrl := ti.vs.Spec.Unsealer.Mode.AzureKeyVault.VaultBaseURL
	client, err := azsecrets.NewClient(vaultBaseUrl, ti.cred, nil)
	if err != nil {
		return err
	}

	_, err = client.BeginDeleteSecret(context.TODO(), key, nil)
	if err != nil {
		return err
	}

	for i := 0; i < 15; i++ {
		_, err = client.PurgeDeletedSecret(context.TODO(), key, nil)
		if err == nil {
			return nil
		}
		time.Sleep(2 * time.Second)
	}

	return err
}

func (ti *TokenKeyInfo) Set(key, value string) error {
	key = strings.Replace(key, ".", "-", -1)

	vaultBaseUrl := ti.vs.Spec.Unsealer.Mode.AzureKeyVault.VaultBaseURL
	client, err := azsecrets.NewClient(vaultBaseUrl, ti.cred, nil)
	if err != nil {
		return err
	}

	_, err = client.SetSecret(context.TODO(), key, base64.StdEncoding.EncodeToString([]byte(value)), &azsecrets.SetSecretOptions{
		ContentType: to.StringPtr("password"),
	})
	if err != nil {
		return errors.Wrap(err, "unable to set secrets in key vault")
	}

	return nil
}

func (ti *TokenKeyInfo) getLatestVersion(key string) (string, error) {
	key = strings.Replace(key, ".", "-", -1)
	vaultBaseUrl := ti.vs.Spec.Unsealer.Mode.AzureKeyVault.VaultBaseURL
	client, err := azsecrets.NewClient(vaultBaseUrl, ti.cred, nil)
	if err != nil {
		return "", err
	}

	var version string
	var dur time.Duration
	pager := client.ListSecretVersions(key, nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		for _, ver := range resp.Secrets {
			cur := time.Since(*ver.Attributes.Created)
			if version == "" {
				version = *ver.ID
				dur = cur
			} else if cur < dur {
				dur = cur
				version = *ver.ID
			}
		}
	}

	return version, nil
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

func (ti *TokenKeyInfo) NewUnsealKeyName(id int) string {
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

	return fmt.Sprintf("%s-unseal-key-%d", keyPrefix, id)
}

func (ti *TokenKeyInfo) OldUnsealKeyName(id int) string {
	return fmt.Sprintf("vault-unseal-key-%d", id)
}

func (ti *TokenKeyInfo) Clean() error {
	return nil
}
