/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha2

import (
	"time"

	core "k8s.io/api/core/v1"
)

// UnsealerSpec contain the configuration for auto vault initialize/unseal
type UnsealerSpec struct {
	// Total count of secret shares that exist
	// +optional
	SecretShares int64 `json:"secretShares,omitempty"`

	// Minimum required secret shares to unseal
	// +optional
	SecretThreshold int64 `json:"secretThreshold,omitempty"`

	// How often to attempt to unseal the vault instance
	// +optional
	RetryPeriodSeconds time.Duration `json:"retryPeriodSeconds,omitempty"` //nolint

	// overwrite existing unseal keys and root tokens, possibly dangerous!
	// +optional
	OverwriteExisting bool `json:"overwriteExisting,omitempty"`

	// should the root token be stored in the key store (default true)
	// +optional
	StoreRootToken bool `json:"storeRootToken,omitempty"`

	// mode contains unseal mechanism
	// +optional
	Mode ModeSpec `json:"mode,omitempty"`
}

// ModeSpec contain unseal mechanism
type ModeSpec struct {
	// +optional
	KubernetesSecret *KubernetesSecretSpec `json:"kubernetesSecret,omitempty"`

	// +optional
	GoogleKmsGcs *GoogleKmsGcsSpec `json:"googleKmsGcs,omitempty"`

	// +optional
	AwsKmsSsm *AwsKmsSsmSpec `json:"awsKmsSsm,omitempty"`

	// +optional
	AzureKeyVault *AzureKeyVault `json:"azureKeyVault,omitempty"`
}

// KubernetesSecretSpec contain the fields that required to unseal using kubernetes secret
type KubernetesSecretSpec struct {
	SecretName string `json:"secretName"`
}

// GoogleKmsGcsSpec contain the fields that required to unseal vault using google kms
type GoogleKmsGcsSpec struct {
	// The name of the Google Cloud KMS crypto key to use
	KmsCryptoKey string `json:"kmsCryptoKey"`

	// The name of the Google Cloud KMS key ring to use
	KmsKeyRing string `json:"kmsKeyRing"`

	// The Google Cloud KMS location to use (eg. 'global', 'europe-west1')
	KmsLocation string `json:"kmsLocation"`

	// The Google Cloud KMS project to use
	KmsProject string `json:"kmsProject"`

	// The name of the Google Cloud Storage bucket to store values in
	Bucket string `json:"bucket"`

	// Secret containing Google application credential
	// secret data:
	//  - sa.json:<value>
	// +optional
	CredentialSecretRef *core.LocalObjectReference `json:"credentialSecretRef,omitempty"`
}

// AwsKmsSsmSpec contain the fields that required to unseal vault using aws kms ssm
type AwsKmsSsmSpec struct {
	// The ID or ARN of the AWS KMS key to encrypt values
	KmsKeyID string `json:"kmsKeyID"`

	// +optional
	// An optional Key prefix for SSM Parameter store
	SsmKeyPrefix string `json:"ssmKeyPrefix,omitempty"`

	Region string `json:"region,omitempty"`

	// Specifies the secret name containing AWS access key and AWS secret key
	// secret data:
	//  - access_key:<value>
	//  - secret_key:<value>
	// +optional
	CredentialSecretRef *core.LocalObjectReference `json:"credentialSecretRef,omitempty"`

	// Used to make AWS KMS requests. This is useful,
	// for example, when connecting to KMS over a VPC Endpoint.
	// If not set, Vault will use the default API endpoint for your region.
	Endpoint string `json:"endpoint,omitempty"`
}

// AzureKeyVault contain the fields that required to unseal vault using azure key vault
type AzureKeyVault struct {
	// Azure key vault url, for example https://myvault.vault.azure.net
	VaultBaseURL string `json:"vaultBaseURL"`

	// The cloud environment identifier
	// default: "AZUREPUBLICCLOUD"
	// +optional
	Cloud string `json:"cloud,omitempty"`

	// The AAD Tenant ID
	TenantID string `json:"tenantID"`

	// Specifies the name of secret containing client cert and client cert password
	// secret data:
	//  - client-cert:<value>
	// 	- client-cert-password: <value>
	// +optional
	TLSSecretRef *core.LocalObjectReference `json:"tlsSecretRef,omitempty"`

	// Specifies the name of secret containing client id and client secret of AAD application
	// secret data:
	//  - client-id:<value>
	//  - client-secret:<value>
	// +optional
	CredentialSecretRef *core.LocalObjectReference `json:"credentialSecretRef,omitempty"`

	// Use managed service identity for the virtual machine
	// +optional
	UseManagedIdentity bool `json:"useManagedIdentity,omitempty"`
}
