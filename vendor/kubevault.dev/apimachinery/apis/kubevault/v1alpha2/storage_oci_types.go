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
	core "k8s.io/api/core/v1"
)

// vault doc: https://openbao.org/docs/configuration/storage/oci/
//
// OCISpec defines configuration to set up OCI Object Storage as backend storage in vault
type OCISpec struct {
	// Specifies the name of the OCI Object Storage bucket to store data in. It must
	// already exist; Vault does not create it.
	BucketName string `json:"bucketName"`

	// Specifies the Object Storage namespace the bucket belongs to.
	NamespaceName string `json:"namespaceName"`

	// Specifies the OCI region to use. Defaults to the region configured in the
	// resolved OCI configuration provider.
	// +optional
	Region string `json:"region,omitempty"`

	// When true, authenticates using an OCI API-key configuration file (see
	// CredentialSecretRef). When false (the default), authenticates using instance
	// principal credentials, intended for Vault instances running on an OCI compute
	// instance.
	// +optional
	AuthTypeAPIKey bool `json:"authTypeAPIKey,omitempty"`

	// Specifies the secret name containing an OCI API-key configuration file (in the
	// same format as the default "~/.oci/config" file) and the private key it
	// references. Required when AuthTypeAPIKey is true, and only consulted then.
	// Every key of the secret is mounted in /etc/vault/storage/oci/.oci, so key_file
	// in the config must point there by full path, e.g.
	// key_file=/etc/vault/storage/oci/.oci/oci_api_key.pem.
	// secret data:
	//  - config=<value>
	//  - <key file name>=<private key, PEM>
	// +optional
	CredentialSecretRef *core.LocalObjectReference `json:"credentialSecretRef,omitempty"`

	// Enables High Availability support. Requires LockBucketName to also be set.
	// +optional
	HAEnabled bool `json:"haEnabled,omitempty"`

	// Specifies the name of a second bucket used to store HA lock records. Required
	// when HAEnabled is true. Must already exist and should generally be a different
	// bucket than BucketName.
	// +optional
	LockBucketName string `json:"lockBucketName,omitempty"`
}
