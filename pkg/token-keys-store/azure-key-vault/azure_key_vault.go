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

	vaultapi "kubevault.dev/apimachinery/apis/kubevault/v1alpha2"
	"kubevault.dev/cli/pkg/token-keys-store/api"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/pkg/errors"
	"gomodules.xyz/pointer"
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

	if vs.Spec.Unsealer.Mode.AzureKeyVault.CredentialSecretRef != nil {
		secret, err := kubeClient.CoreV1().Secrets(vs.Namespace).Get(context.TODO(), vs.Spec.Unsealer.Mode.AzureKeyVault.CredentialSecretRef.Name, metav1.GetOptions{})
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
	} else {
		if _, ok := os.LookupEnv(ClientID); !ok {
			_, _ = fmt.Fprintf(os.Stderr, "WARNING!!! missing env variable %s", ClientID)
		}
		if _, ok := os.LookupEnv(ClientSecret); !ok {
			_, _ = fmt.Fprintf(os.Stderr, "WARNING!!! missing env variable %s", ClientSecret)
		}
	}

	if err := os.Setenv(TenantID, vs.Spec.Unsealer.Mode.AzureKeyVault.TenantID); err != nil {
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
	client := azsecrets.NewClient(vaultBaseUrl, ti.cred, nil)

	version, err := ti.getLatestVersion(key)
	if err != nil {
		return "", err
	}

	idx := strings.LastIndex(version, "/")
	if idx == -1 {
		return "", errors.New("version id not found")
	}

	resp, err := client.GetSecret(context.Background(), strings.ReplaceAll(key, ".", "-"), version[idx+1:], nil)
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
	key = strings.ReplaceAll(key, ".", "-")

	vaultBaseUrl := ti.vs.Spec.Unsealer.Mode.AzureKeyVault.VaultBaseURL
	client := azsecrets.NewClient(vaultBaseUrl, ti.cred, nil)

	_, err := client.DeleteSecret(context.TODO(), key, nil)
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
	key = strings.ReplaceAll(key, ".", "-")

	vaultBaseUrl := ti.vs.Spec.Unsealer.Mode.AzureKeyVault.VaultBaseURL
	client := azsecrets.NewClient(vaultBaseUrl, ti.cred, nil)

	_, err := client.SetSecret(context.TODO(), key, azsecrets.SetSecretParameters{
		Value:       pointer.StringP(base64.StdEncoding.EncodeToString([]byte(value))),
		ContentType: pointer.StringP("password"),
	}, nil)
	if err != nil {
		return errors.Wrap(err, "unable to set secrets in key vault")
	}

	return nil
}

func (ti *TokenKeyInfo) getLatestVersion(key string) (string, error) {
	key = strings.ReplaceAll(key, ".", "-")
	vaultBaseUrl := ti.vs.Spec.Unsealer.Mode.AzureKeyVault.VaultBaseURL
	client := azsecrets.NewClient(vaultBaseUrl, ti.cred, nil)

	var version azsecrets.ID
	var dur time.Duration
	pager := client.NewListSecretVersionsPager(key, nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		if err != nil {
			return "", err
		}
		for _, ver := range resp.Value {
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

	return string(version), nil
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
	_ = os.Unsetenv(ClientID)
	_ = os.Unsetenv(ClientSecret)
	_ = os.Unsetenv(TenantID)
}
