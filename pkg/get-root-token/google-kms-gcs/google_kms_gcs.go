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

	vaultapi "kubevault.dev/apimachinery/apis/kubevault/v1alpha1"
	"kubevault.dev/cli/pkg/get-root-token/api"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type TokenInfo struct {
	googleKmsGcsSpec *vaultapi.GoogleKmsGcsSpec
	storageClient    *storage.Client
}

var _ api.TokenInterface = &TokenInfo{}

func New(spec *vaultapi.GoogleKmsGcsSpec) (*TokenInfo, error) {
	if spec == nil {
		return nil, errors.New("googleKmsGcs spec is empty")
	}

	client, err := storage.NewClient(context.TODO())
	if err != nil {
		return nil, err
	}

	return &TokenInfo{
		googleKmsGcsSpec: spec,
		storageClient:    client,
	}, nil
}

func (ti *TokenInfo) Token() (string, error) {
	token := ti.TokenName()
	rc, err := ti.storageClient.Bucket(ti.googleKmsGcsSpec.Bucket).Object(token).NewReader(context.TODO())
	if err != nil {
		return "", err
	}
	defer rc.Close()

	body, err := ioutil.ReadAll(rc)
	if err != nil {
		return "", err
	}

	name := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		ti.googleKmsGcsSpec.KmsProject, ti.googleKmsGcsSpec.KmsLocation,
		ti.googleKmsGcsSpec.KmsKeyRing, ti.googleKmsGcsSpec.KmsCryptoKey)

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
	return "vault-root-token"
}
