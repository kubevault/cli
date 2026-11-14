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

// vault doc: https://openbao.org/docs/configuration/storage/mssql/
//
// MSSQLSpec defines configuration to set up Microsoft SQL Server as backend storage in vault
type MSSQLSpec struct {
	// Specifies the address of the MSSQL host.
	// +optional
	Server string `json:"server,omitempty"`

	// Specifies the port of the MSSQL host. Defaults to the driver's standard port
	// (1433) when unset.
	// +optional
	Port string `json:"port,omitempty"`

	// Specifies the secret name containing the MSSQL username and password to connect
	// with.
	// secret data:
	//  - username=<value>
	//  - password=<value>
	// +optional
	CredentialSecretRef *core.LocalObjectReference `json:"credentialSecretRef,omitempty"`

	// Specifies the name of the database to use. Vault will attempt to create it if it
	// does not already exist.
	// +optional
	// +kubebuilder:default:="openbao"
	Database string `json:"database,omitempty"`

	// Specifies the name of the table in which to write Vault data. Vault will attempt
	// to create it if missing.
	// +optional
	// +kubebuilder:default:="openbao"
	Table string `json:"table,omitempty"`

	// Specifies the schema within the database that the table lives in. Vault will
	// attempt to create it if missing (requires permission to run CREATE SCHEMA).
	// +optional
	// +kubebuilder:default:="dbo"
	Schema string `json:"schema,omitempty"`

	// Specifies the application name to report to the server.
	// +optional
	// +kubebuilder:default:="openbao"
	AppName string `json:"appName,omitempty"`

	// Specifies the connection timeout, in seconds.
	// +optional
	// +kubebuilder:default:="30"
	ConnectionTimeout string `json:"connectionTimeout,omitempty"`

	// Specifies the driver's internal log level bitmask.
	// +optional
	LogLevel string `json:"logLevel,omitempty"`

	// Specifies the maximum number of concurrent requests to MSSQL.
	// +optional
	MaxParallel int64 `json:"maxParallel,omitempty"`

	// DatabaseRef contains the info of a KubeDB managed MSSQLServer.
	// If set, Server will be generated from the referenced AppBinding.
	// +optional
	DatabaseRef *appcat.AppReference `json:"databaseRef,omitempty"`

	// Specifies how the connection is encrypted:
	//  - true:    encrypt the whole connection, and fail if the server doesn't support it
	//  - false:   encrypt only the login packet, unless the server forces encryption
	//  - strict:  use TDS 8.0, where TLS is set up before any TDS traffic (SQL Server 2022+)
	//  - disable: don't encrypt anything
	// Unless it is disable, the server certificate is verified (see TrustServerCertificate).
	// When unset, it is true if a CA is available (from TLSSecretRef or the DatabaseRef
	// AppBinding) or any other TLS option is set, and otherwise left to the driver, which
	// encrypts only the login packet and does not verify the server certificate.
	// +optional
	Encrypt MSSQLEncrypt `json:"encrypt,omitempty"`

	// Specifies the secret name containing the CA certificate used to verify the server.
	// secret data:
	//  - ca.crt=<value>
	// If DatabaseRef is set and this is empty, it is generated from the CA bundle of the
	// referenced AppBinding.
	// +optional
	TLSSecretRef *core.LocalObjectReference `json:"tlsSecretRef,omitempty"`

	// Specifies whether to skip verifying the server certificate. Not recommended in
	// production.
	// +optional
	TrustServerCertificate bool `json:"trustServerCertificate,omitempty"`

	// Specifies the name the server certificate must be valid for. Defaults to Server.
	// +optional
	TLSServerName string `json:"tlsServerName,omitempty"`

	// Specifies the minimum TLS version to use.
	// +optional
	TLSMinVersion MSSQLTLSVersion `json:"tlsMinVersion,omitempty"`
}

// +kubebuilder:validation:Enum="true";"false";strict;disable
type MSSQLEncrypt string

const (
	MSSQLEncryptTrue    MSSQLEncrypt = "true"
	MSSQLEncryptFalse   MSSQLEncrypt = "false"
	MSSQLEncryptStrict  MSSQLEncrypt = "strict"
	MSSQLEncryptDisable MSSQLEncrypt = "disable"
)

// +kubebuilder:validation:Enum=tls10;tls11;tls12;tls13
type MSSQLTLSVersion string
