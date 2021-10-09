package aws_kms_ssm

import (
	"encoding/base64"
	"fmt"

	vaultapi "kubevault.dev/apimachinery/apis/kubevault/v1alpha1"
	"kubevault.dev/cli/pkg/get-root-token/api"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pkg/errors"
)

type TokenInfo struct {
	awskmsSsmSpec *vaultapi.AwsKmsSsmSpec
	ssmService    *ssm.SSM
	kmsService    *kms.KMS
}

var _ api.TokenInterface = &TokenInfo{}

func New(spec *vaultapi.AwsKmsSsmSpec) (*TokenInfo, error) {
	if spec == nil {
		return nil, errors.New("awsKmsSsm spec is empty")
	}

	// TODO: Check region for empty

	sess, err := session.NewSession(&aws.Config{
		CredentialsChainVerboseErrors: func() *bool {
			f := true
			return &f
		}(),
		Region: aws.String(spec.Region)},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create session")
	}

	return &TokenInfo{
		awskmsSsmSpec: spec,
		kmsService:    kms.New(sess),
		ssmService:    ssm.New(sess),
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

	decryptOutput, err := ti.kmsService.Decrypt(&kms.DecryptInput{
		CiphertextBlob: sDec,
		EncryptionContext: map[string]*string{
			"Tool": aws.String("vault-unsealer"),
		},
		GrantTokens: []*string{},
		KeyId:       aws.String(ti.awskmsSsmSpec.KmsKeyID),
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to kms decrypt")
	}

	return string(decryptOutput.Plaintext), nil
}

func (ti *TokenInfo) TokenName() string {
	if ti.awskmsSsmSpec.SsmKeyPrefix != "" {
		return fmt.Sprintf("%s%s", ti.awskmsSsmSpec.SsmKeyPrefix, "vault-root-token")
	}
	return "vault-root-token"
}
