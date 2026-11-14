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
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
)

// vault doc: https://www.vaultproject.io/docs/configuration/storage/mysql.html
//
// MySQLSpec defines configuration to set up MySQL Storage as backend storage in vault
type MySQLSpec struct {
	// Specifies the address of the MySQL host.
	// if DatabaseRef is set then Address will be generated from it
	// This must be set if DatabaseRef is empty, validate from ValidatingWebhook
	// host example: <db-name>.<db-ns>.svc:3306
	// +optional
	Address string `json:"address"`

	// Specifies the name of the database. If the database does not exist, Vault will attempt to create it.
	// +optional
	// +kubebuilder:default:="vault"
	Database string `json:"database,omitempty"`

	// Specifies the name of the table. If the table does not exist, Vault will attempt to create it.
	// +optional
	// +kubebuilder:default:="vault"
	Table string `json:"table,omitempty"`

	// Specifies the MySQL username and password to connect to the database
	// secret data:
	//  - username=<value>
	//  - password=<value>
	CredentialSecretRef *core.LocalObjectReference `json:"credentialSecretRef,omitempty"`

	// Specifies the name of the secret containing the CA certificate to connect using TLS.
	// secret data:
	//  - ca.crt=<value>
	// +optional
	TLSSecretRef *core.LocalObjectReference `json:"tlsSecretRef,omitempty"`

	//  Specifies the maximum number of concurrent requests to take place.
	// +optional
	MaxParallel int64 `json:"maxParallel,omitempty"`

	// DatabaseRef contains the info of KubeDB managed Database
	// This will be used to generate the "Address" field
	// +optional
	DatabaseRef *appcat.AppReference `json:"databaseRef,omitempty"`

	PlaintextCredentialTransmission string `json:"plaintextCredentialTransmission,omitempty"`

	// Specifies the maximum number of idle connections to the database.
	// A zero uses value defaults to 2 idle connections and a negative value disables idle connections.
	// If larger than max_parallel it will be reduced to be equal.
	// +optional
	MaxIdleConnection int64 `json:"maxIdleConnection,omitempty"`

	// Specifies the maximum amount of time in seconds that a connection may be reused. If <= 0s connections are reused forever.
	// +optional
	MaxConnectionLifetime int64 `json:"maxConnectionLifetime,omitempty"`

	// High Availability Parameter
	// Specifies if high availability mode is enabled. This is a boolean value, but it is specified as a string like "true" or "false".
	// +optional
	// +kubebuilder:default:="true"
	HAEnabled string `json:"haEnabled,omitempty"`

	// High Availability Parameter
	// Specifies the name of the table to use for storing high availability information.
	// By default, this is the name of the table suffixed with _lock. If the table does not exist, Vault will attempt to create it.
	// +optional
	// +kubebuilder:default:="vault_lock"
	LockTable string `json:"lockTable,omitempty"`
}
