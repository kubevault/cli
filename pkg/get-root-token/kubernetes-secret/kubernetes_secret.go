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

	vaultapi "kubevault.dev/apimachinery/apis/kubevault/v1alpha1"
	"kubevault.dev/cli/pkg/get-root-token/api"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type TokenInfo struct {
	kubeClient kubernetes.Interface
	vs         *vaultapi.VaultServer
}

var _ api.TokenInterface = &TokenInfo{}

func New(vs *vaultapi.VaultServer, kubeClient kubernetes.Interface) (*TokenInfo, error) {
	if vs == nil {
		return nil, errors.New("vs spec is empty")
	}

	if vs.Spec.Unsealer.Mode.KubernetesSecret == nil {
		return nil, errors.New("kubernetes-secret mode is nil")
	}

	if kubeClient == nil {
		return nil, errors.New("kubeClient is nil")
	}

	return &TokenInfo{
		kubeClient: kubeClient,
		vs:         vs,
	}, nil
}

func (ti *TokenInfo) Token() (string, error) {
	secretName := ti.vs.Spec.Unsealer.Mode.KubernetesSecret.SecretName
	secretNamespace := ti.vs.Namespace
	secret, err := ti.kubeClient.CoreV1().Secrets(secretNamespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	token := ti.TokenName()
	if _, ok := secret.Data[token]; !ok {
		return "", errors.Errorf("%s not found in secret %s/%s", token, secretNamespace, secretName)
	}

	return string(secret.Data[token]), nil
}

func (ti *TokenInfo) TokenName() string {
	return ti.vs.RootTokenID()
}
