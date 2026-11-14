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

// vault doc: https://www.vaultproject.io/docs/configuration/storage/dynamodb.html
//
// DynamoDBSpec defines configuration to set up DynamoDB Storage as backend storage in vault
type DynamoDBSpec struct {
	// Specifies an alternative, AWS compatible, DynamoDB endpoint.
	// +optional
	Endpoint string `json:"endpoint,omitempty"`

	// Specifies the AWS region
	// +optional
	Region string `json:"region,omitempty"`

	// Specifies whether this backend should be used to run Vault in high availability mode.
	// +optional
	HAEnabled bool `json:"haEnabled,omitempty"`

	// Specifies the maximum number of reads consumed per second on the table
	// +optional
	ReadCapacity int64 `json:"readCapacity,omitempty"`

	// Specifies the maximum number of writes performed per second on the table.
	// +optional
	WriteCapacity int64 `json:"writeCapacity,omitempty"`

	// Specifies the name of the DynamoDB table in which to store Vault data.
	// If the specified table does not yet exist, it will be created during initialization.
	// default: vault-dynamodb-backend
	// +optional
	Table string `json:"table,omitempty"`

	// Specifies the secret name containing AWS session token, AWS access key and AWS secret key
	// secret data:
	//  - access_key=<value>
	//  - secret_key=<value>
	//  - session_token=<value>
	// +optional
	CredentialSecretRef *core.LocalObjectReference `json:"credentialSecretRef,omitempty"`

	// Specifies the maximum number of parallel operations to take place.
	// +optional
	MaxParallel int64 `json:"maxParallel,omitempty"`
}
