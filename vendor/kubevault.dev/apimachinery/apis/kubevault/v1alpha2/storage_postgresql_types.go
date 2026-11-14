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

// vault doc: https://www.vaultproject.io/docs/configuration/storage/postgresql.html
//
// PostgreSQLSpec defines configuration to set up PostgreSQL storage as backend storage in vault
type PostgreSQLSpec struct {
	// Specifies the address of the Postgres host.
	// if DatabaseRef is set then Address will be generated from it
	// This must be set if DatabaseRef is empty, validate from ValidatingWebhook
	// host example: <db-name>.<db-ns>.svc:3306
	// +optional
	Address string `json:"address"`

	//  - username=<value>
	//  - password=<value>
	//  - connection_url="postgres://<username>:<password>@<host>:<port>/<db_name>"
	CredentialSecretRef *core.LocalObjectReference `json:"credentialSecretRef,omitempty"`

	// DatabaseRef contains the info of KubeDB managed Database
	// This will be used to generate the "Address" field
	DatabaseRef *appcat.AppReference `json:"databaseRef,omitempty"`

	// SSLMode for both standalone and clusters. [disable;require;verify-ca;verify-full]
	SSLMode PostgresSSLMode `json:"sslMode,omitempty"`

	// Specifies the name of the table in which to write Vault data.
	// This table must already exist (Vault will not attempt to create it).
	// +optional
	// +kubebuilder:default:="vault_kv_store"
	Table string `json:"table,omitempty"`

	//  Specifies the maximum number of concurrent requests to take place.
	// +optional
	MaxParallel int64 `json:"maxParallel,omitempty"`

	// Default not set. Sets the maximum number of connections in the idle connection pool.
	// See golang docs on SetMaxIdleConns(https://pkg.go.dev/database/sql#DB.SetMaxIdleConns) for more information. Requires 1.2 or later.
	// +optional
	MaxIdleConnection int64 `json:"maxIdleConnection,omitempty"`

	// High Availability Parameter
	// Default not enabled, requires 9.5 or later
	// Specifies if high availability mode is enabled. This is a boolean value, but it is specified as a string like "true" or "false".
	// +optional
	// +kubebuilder:default:="false"
	HAEnabled string `json:"haEnabled,omitempty"`

	// Specifies the name of the table to use for storing high availability information. This table must already exist (Vault will not attempt to create it).
	// +optional
	// +kubebuilder:default:="vault_ha_locks"
	HATable string `json:"haTable,omitempty"`

	// Specifies the secret name containing the TLS material used to connect to PostgreSQL.
	// secret data:
	//  - ca.crt:  CA certificate used to verify the server (required)
	//  - tls.crt: client certificate (optional, for mutual TLS; needs tls.key)
	//  - tls.key: client private key (optional, for mutual TLS; needs tls.crt)
	// If DatabaseRef is set and this is empty, it is generated from the CA bundle and
	// client certificate of the referenced AppBinding.
	// +optional
	TLSSecretRef *core.LocalObjectReference `json:"tlsSecretRef,omitempty"`
}
