package token

import (
	"errors"

	vaultapi "kubevault.dev/apimachinery/apis/kubevault/v1alpha1"
	"kubevault.dev/cli/pkg/get-root-token/api"
	aws_kms_ssm "kubevault.dev/cli/pkg/get-root-token/aws-kms-ssm"
)

func NewTokenInterface(vs *vaultapi.VaultServer) (api.TokenInterface, error) {
	if vs.Spec.Unsealer == nil {
		return nil, errors.New("vaultServer unsealer spec is empty")
	}
	mode := vs.Spec.Unsealer.Mode

	switch true {
	case mode.AwsKmsSsm != nil:
		return aws_kms_ssm.New(mode.AwsKmsSsm)
	}

	return nil, errors.New("unknown/unsupported unsealing mode")
}
