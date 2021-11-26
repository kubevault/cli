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

package aws_kms_ssm

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	vaultapi "kubevault.dev/apimachinery/apis/kubevault/v1alpha1"
	"kubevault.dev/cli/pkg/get-root-token/api"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	AccessKey = "AWS_ACCESS_KEY_ID"
	SecretKey = "AWS_SECRET_ACCESS_KEY"
)

type TokenInfo struct {
	ssmService *ssm.SSM
	kmsService *kms.KMS
	vs         *vaultapi.VaultServer
}

var _ api.TokenInterface = &TokenInfo{}

func New(vs *vaultapi.VaultServer, kubeClient kubernetes.Interface) (*TokenInfo, error) {
	if vs == nil {
		return nil, errors.New("vs spec is empty")
	}

	if vs.Spec.Unsealer.Mode.AwsKmsSsm == nil {
		return nil, errors.New("AwsKmsSsm mode is nil")
	}

	if kubeClient == nil {
		return nil, errors.New("kubeClient is nil")
	}

	secret, err := kubeClient.CoreV1().Secrets(vs.Namespace).Get(context.TODO(), vs.Spec.Unsealer.Mode.AwsKmsSsm.CredentialSecret, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if _, ok := secret.Data["access_key"]; !ok {
		return nil, errors.Errorf("%s not found in secret", AccessKey)
	}
	if _, ok := secret.Data["secret_key"]; !ok {
		return nil, errors.Errorf("%s not found in secret", SecretKey)
	}

	if err = os.Setenv(AccessKey, string(secret.Data["access_key"])); err != nil {
		return nil, err
	}
	if err = os.Setenv(SecretKey, string(secret.Data["secret_key"])); err != nil {
		return nil, err
	}

	sess, err := session.NewSession(&aws.Config{
		CredentialsChainVerboseErrors: func() *bool {
			f := true
			return &f
		}(),
		Region: aws.String(vs.Spec.Unsealer.Mode.AwsKmsSsm.Region)},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create session")
	}

	return &TokenInfo{
		kmsService: kms.New(sess),
		ssmService: ssm.New(sess),
		vs:         vs,
	}, nil
}
func (ti *TokenInfo) Token() (string, error) {
	token := ti.TokenName()
	req := &ssm.GetParametersInput{
		Names: []*string{
			aws.String(token),
		},
		WithDecryption: aws.Bool(false),
	}
	params, err := ti.ssmService.GetParameters(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to get token from ssm")
	}
	if len(params.Parameters) == 0 {
		return "", errors.New("failed to get token from ssm; empty response")
	}
	// Since len of the params is greater than zero
	sDec, err := base64.StdEncoding.DecodeString(*params.Parameters[0].Value)
	if err != nil {
		return "", errors.Wrap(err, "failed to base64-decode")
	}

	awsKmsSsmSpec := ti.vs.Spec.Unsealer.Mode.AwsKmsSsm
	decryptOutput, err := ti.kmsService.Decrypt(&kms.DecryptInput{
		CiphertextBlob: sDec,
		EncryptionContext: map[string]*string{
			"Tool": aws.String("vault-unsealer"),
		},
		GrantTokens: []*string{},
		KeyId:       aws.String(awsKmsSsmSpec.KmsKeyID),
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to kms decrypt")
	}

	return string(decryptOutput.Plaintext), nil
}

func (ti *TokenInfo) TokenName() string {
	awsKmsSsmSpec := ti.vs.Spec.Unsealer.Mode.AwsKmsSsm
	if awsKmsSsmSpec.SsmKeyPrefix != "" {
		return fmt.Sprintf("%s%s", awsKmsSsmSpec.SsmKeyPrefix, ti.vs.RootTokenID())
	}
	return ti.vs.RootTokenID()
}
