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

// vault doc: https://openbao.org/docs/configuration/storage/cockroachdb/
//
// CockroachDBSpec defines configuration to set up CockroachDB as backend storage in vault.
// CockroachDB speaks the PostgreSQL wire protocol, but is configured as a dedicated backend
// (rather than reusing PostgreSQLSpec) so its schema and locking queries can account for
// CockroachDB-specific behavior.
type CockroachDBSpec struct {
	// Specifies the address of the CockroachDB host.
	// if DatabaseRef is set then Address will be generated from it
	// This must be set if DatabaseRef is empty, validate from ValidatingWebhook
	// +optional
	Address string `json:"address"`

	// Specifies the name of the database in which the tables for Vault data are created.
	// The database must already exist: Vault creates the tables, not the database.
	// +optional
	// +kubebuilder:default:="vault"
	Database string `json:"database,omitempty"`

	// Specifies the secret name containing the username and the password to connect to
	// the CockroachDB cluster. The password can be left out when the user logs in with a
	// client certificate (see tlsSecretRef).
	// Optionally, the secret holds a full connection_url, which then replaces address,
	// database and sslMode.
	// secret data:
	//  - username=<value>
	//  - password=<value>
	//  - connection_url="postgresql://<username>:<password>@<host>:<port>/<db_name>?sslmode=<sslMode>"
	// +optional
	CredentialSecretRef *core.LocalObjectReference `json:"credentialSecretRef,omitempty"`

	// DatabaseRef contains the info of KubeDB managed Database
	// This will be used to generate the "Address" field
	// +optional
	DatabaseRef *appcat.AppReference `json:"databaseRef,omitempty"`

	// SSLMode for both standalone and clusters. [disable;require;verify-ca;verify-full]
	// +optional
	SSLMode PostgresSSLMode `json:"sslMode,omitempty"`

	// Specifies the name of the table in which to write Vault data.
	// +optional
	// +kubebuilder:default:="openbao_kv_store"
	Table string `json:"table,omitempty"`

	// Specifies the maximum number of concurrent requests to take place.
	// +optional
	MaxParallel int64 `json:"maxParallel,omitempty"`

	// Specifies the maximum number of concurrent interactive transactions to take place.
	// +optional
	TransactionMaxParallel int64 `json:"transactionMaxParallel,omitempty"`

	// When set, Vault will not attempt to create the storage table(s); they must already exist.
	// +optional
	SkipCreateTable bool `json:"skipCreateTable,omitempty"`

	// High Availability Parameter
	// Specifies if high availability mode is enabled.
	// +optional
	HAEnabled bool `json:"haEnabled,omitempty"`

	// Specifies the name of the table to use for storing high availability information.
	// +optional
	// +kubebuilder:default:="openbao_ha_locks"
	HATable string `json:"haTable,omitempty"`

	// Specifies the secret name containing the TLS material used to connect to CockroachDB.
	// secret data:
	//  - ca.crt:  CA certificate used to verify the server (required)
	//  - tls.crt: client certificate (optional, for mutual TLS; needs tls.key)
	//  - tls.key: client private key (optional, for mutual TLS; needs tls.crt)
	// If DatabaseRef is set and this is empty, it is generated from the CA bundle and
	// client certificate of the referenced AppBinding.
	// +optional
	TLSSecretRef *core.LocalObjectReference `json:"tlsSecretRef,omitempty"`
}
