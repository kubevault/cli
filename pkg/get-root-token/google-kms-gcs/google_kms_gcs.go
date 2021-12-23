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

package google_kms_gcs

import (
	"context"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	vaultapi "kubevault.dev/apimachinery/apis/kubevault/v1alpha1"
	"kubevault.dev/cli/pkg/get-root-token/api"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	ServiceAccountJSON    = "sa.json"
	GoogleApplicationCred = "GOOGLE_APPLICATION_CREDENTIALS"
)

type TokenInfo struct {
	storageClient *storage.Client
	kubeClient    kubernetes.Interface
	vs            *vaultapi.VaultServer
	path          string
}

var _ api.TokenInterface = &TokenInfo{}

func New(vs *vaultapi.VaultServer, kubeClient kubernetes.Interface) (*TokenInfo, error) {
	if vs == nil {
		return nil, errors.New("vs spec is empty")
	}

	if vs.Spec.Unsealer.Mode.GoogleKmsGcs == nil {
		return nil, errors.New("GoogleKmsGcs mode is nil")
	}

	if kubeClient == nil {
		return nil, errors.New("kubeClient is nil")
	}

	secret, err := kubeClient.CoreV1().Secrets(vs.Namespace).Get(context.TODO(), vs.Spec.Unsealer.Mode.GoogleKmsGcs.CredentialSecret, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if _, ok := secret.Data[ServiceAccountJSON]; !ok {
		return nil, errors.Errorf("%s not found in secret", ServiceAccountJSON)
	}

	path := filepath.Join("/tmp", fmt.Sprintf("google-sa-cred-%s", randomString(6)))
	if err = os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, err
	}

	saFile := filepath.Join(path, ServiceAccountJSON)
	if err = os.WriteFile(saFile, secret.Data[ServiceAccountJSON], os.ModePerm); err != nil {
		return nil, err
	}

	if err = os.Setenv(GoogleApplicationCred, saFile); err != nil {
		return nil, err
	}

	client, err := storage.NewClient(context.TODO())
	if err != nil {
		return nil, err
	}

	return &TokenInfo{
		storageClient: client,
		kubeClient:    kubeClient,
		vs:            vs,
		path:          path,
	}, nil
}

func (ti *TokenInfo) Token() (string, error) {
	defer func() {
		_ = os.RemoveAll(ti.path)
	}()

	token := ti.TokenName()
	googleKmsGcsSpec := ti.vs.Spec.Unsealer.Mode.GoogleKmsGcs
	rc, err := ti.storageClient.Bucket(googleKmsGcsSpec.Bucket).Object(token).NewReader(context.TODO())
	if err != nil {
		return "", err
	}
	defer rc.Close()

	body, err := ioutil.ReadAll(rc)
	if err != nil {
		return "", err
	}

	name := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		googleKmsGcsSpec.KmsProject, googleKmsGcsSpec.KmsLocation,
		googleKmsGcsSpec.KmsKeyRing, googleKmsGcsSpec.KmsCryptoKey)

	decryptedToken, err := decryptSymmetric(name, body)
	if err != nil {
		return "", err
	}

	return decryptedToken, nil
}

func decryptSymmetric(name string, ciphertext []byte) (string, error) {
	client, err := kms.NewKeyManagementClient(context.TODO())
	if err != nil {
		return "", errors.Errorf("failed to create kms client: %v", err)
	}
	defer client.Close()

	crc32c := func(data []byte) uint32 {
		t := crc32.MakeTable(crc32.Castagnoli)
		return crc32.Checksum(data, t)
	}
	ciphertextCRC32C := crc32c(ciphertext)

	req := &kmspb.DecryptRequest{
		Name:             name,
		Ciphertext:       ciphertext,
		CiphertextCrc32C: wrapperspb.Int64(int64(ciphertextCRC32C)),
	}

	result, err := client.Decrypt(context.TODO(), req)
	if err != nil {
		return "", errors.Errorf("failed to decrypt ciphertext with %s", err.Error())
	}

	if int64(crc32c(result.Plaintext)) != result.PlaintextCrc32C.Value {
		return "", errors.Errorf("decrypt response corrupted in-transit")
	}

	return string(result.Plaintext), nil
}

func (ti *TokenInfo) TokenName() string {
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

func randomString(n int) string {
	rand.Seed(time.Now().Unix())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
