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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
)

const (
	ResourceKindSecretEngine = "SecretEngine"
	ResourceSecretEngine     = "secretengine"
	ResourceSecretEngines    = "secretengines"
	EngineTypeAWS            = "aws"
	EngineTypeGCP            = "gcp"
	EngineTypeAzure          = "azure"
	EngineTypeDatabase       = "database"
	EngineTypeKV             = "kv"
	EngineTypePKI            = "pki"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=secretengines,singular=secretengine,categories={vault,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type SecretEngine struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SecretEngineSpec   `json:"spec,omitempty"`
	Status            SecretEngineStatus `json:"status,omitempty"`
}

type SecretEngineSpec struct {
	VaultRef kmapi.ObjectReference `json:"vaultRef"`

	SecretEngineConfiguration `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

type SecretEngineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []SecretEngine `json:"items,omitempty"`
}

type SecretEngineConfiguration struct {
	AWS           *AWSConfiguration           `json:"aws,omitempty"`
	Azure         *AzureConfiguration         `json:"azure,omitempty"`
	GCP           *GCPConfiguration           `json:"gcp,omitempty"`
	Postgres      *PostgresConfiguration      `json:"postgres,omitempty"`
	MongoDB       *MongoDBConfiguration       `json:"mongodb,omitempty"`
	Redis         *RedisConfiguration         `json:"redis,omitempty"`
	MySQL         *MySQLConfiguration         `json:"mysql,omitempty"`
	MariaDB       *MariaDBConfiguration       `json:"mariadb,omitempty"`
	KV            *KVConfiguration            `json:"kv,omitempty"`
	Elasticsearch *ElasticsearchConfiguration `json:"elasticsearch,omitempty"`
	PKI           *PKIConfiguration           `json:"pki,omitempty"`
}

// https://developer.hashicorp.com/vault/api-docs/secret/pki#generate-root
// PKIConfiguration contains information about PKI Secret Engine
type PKIConfiguration struct {
	IsRootCA          bool                   `json:"isRootCA"`
	ParentCARef       *kmapi.ObjectReference `json:"parentCARef,omitempty"`
	URLs              *ConfigURL             `json:"urls,omitempty"`
	CommonName        string                 `json:"commonName,omitempty"`
	AltNames          string                 `json:"altNames,omitempty"`
	IssuerName        string                 `json:"issuerName,omitempty"`
	CAType            string                 `json:"type,omitempty"`
	TTL               string                 `json:"ttl,omitempty"`
	MaxPathLength     *int                   `json:"maxPathLength,omitempty"`
	OU                string                 `json:"ou,omitempty"`
	Organization      string                 `json:"organization,omitempty"`
	Country           string                 `json:"country,omitempty"`
	AdditionalPayload map[string]string      `json:"additionalPayload,omitempty"`
}

type ConfigURL struct {
	IssuingCertificates   []string `json:"issuingCertificates,omitempty"`
	CRLDistributionPoints []string `json:"crlDistributionPoints,omitempty"`
	OCSPServers           []string `json:"ocspServers,omitempty"`
	EnableTemplating      bool     `json:"enableTemplating,omitempty"`
}

// https://www.vaultproject.io/api/secret/aws/index.html#configure-root-iam-credentials
// AWSConfiguration contains information to communicate with AWS
type AWSConfiguration struct {
	// Specifies the secret containing AWS access key ID and secret access key
	// secret.Data:
	//	- access_key=<value>
	//  - secret_key=<value>
	CredentialSecret string `json:"credentialSecret"`

	// Specifies the AWS region
	Region string `json:"region"`

	// Specifies a custom HTTP IAM enminidpoint to use
	IAMEndpoint string `json:"iamEndpoint,omitempty"`

	// Specifies a custom HTTP STS endpoint to use
	STSEndpoint string `json:"stsEndpoint,omitempty"`

	// Number of max retries the client should use for recoverable errors.
	// The default (-1) falls back to the AWS SDK's default behavior
	MaxRetries *int64 `json:"maxRetries,omitempty"`

	LeaseConfig *LeaseConfig `json:"leaseConfig,omitempty"`
}

// https://www.vaultproject.io/api/secret/aws/index.html#configure-lease
// LeaseConfig contains lease configuration
type LeaseConfig struct {
	// Specifies the lease value provided as a string duration with time suffix.
	// "h" (hour) is the largest suffix.
	Lease string `json:"lease"`

	// Specifies the maximum lease value provided as a string duration with time suffix.
	// "h" (hour) is the largest suffix
	LeaseMax string `json:"leaseMax"`
}

// https://www.vaultproject.io/api/secret/gcp/index.html#write-config
// GCPConfiguration contains information to communicate with GCP
type GCPConfiguration struct {
	// Specifies the secret containing GCP credentials
	// secret.Data:
	//	- sa.json
	CredentialSecret string `json:"credentialSecret"`

	// Specifies default config TTL for long-lived credentials
	// (i.e. service account keys).
	// +optional
	TTL string `json:"ttl,omitempty"`

	// Specifies the maximum config TTL for long-lived
	// credentials (i.e. service account keys).
	// +optional
	MaxTTL string `json:"maxTTL,omitempty"`
}

// ref:
//	- https://www.vaultproject.io/api/secret/azure/index.html#configure-access

// AzureConfiguration contains information to communicate with Azure
type AzureConfiguration struct {
	// Specifies the secret name containing Azure credentials
	// secret.Data:
	// 	- subscription-id: <value>, The subscription id for the Azure Active Directory.
	//	- tenant-id: <value>, The tenant id for the Azure Active Directory.
	//	- client-id: <value>, The OAuth2 client id to connect to Azure.
	//	- client-secret: <value>, The OAuth2 client secret to connect to Azure.
	CredentialSecret string `json:"credentialSecret"`

	// The Azure environment.
	// If not specified, Vault will use Azure Public Cloud.
	// +optional
	Environment string `json:"environment,omitempty"`
}

// PostgresConfiguration defines a PostgreSQL app configuration.
// https://www.vaultproject.io/api/secret/databases/index.html
// https://www.vaultproject.io/api/secret/databases/postgresql.html#configure-connection
type PostgresConfiguration struct {
	// Specifies the Postgres database appbinding reference
	DatabaseRef appcat.AppReference `json:"databaseRef"`

	// Specifies the name of the plugin to use for this connection.
	// Default plugin:
	//	- for postgres: postgresql-database-plugin
	PluginName string `json:"pluginName,omitempty"`

	// List of the roles allowed to use this connection.
	// Defaults to empty (no roles), if contains a "*" any role can use this connection.
	AllowedRoles []string `json:"allowedRoles,omitempty"`

	// Specifies the maximum number of open connections to the database.
	MaxOpenConnections int64 `json:"maxOpenConnections,omitempty"`

	// Specifies the maximum number of idle connections to the database.
	// A zero uses the value of max_open_connections and a negative value disables idle connections.
	// If larger than max_open_connections it will be reduced to be equal.
	MaxIdleConnections int64 `json:"maxIdleConnections,omitempty"`

	// Specifies the maximum amount of time a connection may be reused.
	// If <= 0s connections are reused forever.
	MaxConnectionLifetime string `json:"maxConnectionLifetime,omitempty"`
}

// MongoDBConfiguration defines a MongoDB app configuration.
// https://www.vaultproject.io/api/secret/databases/index.html
// https://www.vaultproject.io/api/secret/databases/mongodb.html#configure-connection
type MongoDBConfiguration struct {
	// Specifies the database appbinding reference
	DatabaseRef appcat.AppReference `json:"databaseRef"`

	// Specifies the name of the plugin to use for this connection.
	// Default plugin:
	//  - for mongodb: mongodb-database-plugin
	PluginName string `json:"pluginName,omitempty"`

	// List of the roles allowed to use this connection.
	// Defaults to empty (no roles), if contains a "*" any role can use this connection.
	AllowedRoles []string `json:"allowedRoles,omitempty"`

	// Specifies the MongoDB write concern. This is set for the entirety
	// of the session, maintained for the lifecycle of the plugin process.
	WriteConcern string `json:"writeConcern,omitempty"`
}

// RedisConfiguration defines a Redis app configuration.
// https://www.vaultproject.io/api/secret/databases/index.html
// https://developer.hashicorp.com/vault/api-docs/secret/databases/redis#configure-connection
type RedisConfiguration struct {
	// Specifies the database appbinding reference
	DatabaseRef appcat.AppReference `json:"databaseRef"`

	// Specifies the name of the plugin to use for this connection.
	// Default plugin:
	//  - for redis: redis-database-plugin
	PluginName string `json:"pluginName,omitempty"`

	// List of the roles allowed to use this connection.
	// Defaults to empty (no roles), if contains a "*" any role can use this connection.
	AllowedRoles []string `json:"allowedRoles,omitempty"`
}

// MySQLConfiguration defines a MySQL app configuration.
// https://www.vaultproject.io/api/secret/databases/index.html
// https://www.vaultproject.io/api/secret/databases/mysql-maria.html#configure-connection
type MySQLConfiguration struct {
	// DatabaseRef refers to a MySQL/MariaDB database AppBinding in any namespace
	DatabaseRef appcat.AppReference `json:"databaseRef"`

	// Specifies the name of the plugin to use for this connection.
	// Default plugin:
	//  - for mysql: mysql-database-plugin
	PluginName string `json:"pluginName,omitempty"`

	// List of the roles allowed to use this connection.
	// Defaults to empty (no roles), if contains a "*" any role can use this connection.
	AllowedRoles []string `json:"allowedRoles,omitempty"`

	// Specifies the maximum number of open connections to the database.
	MaxOpenConnections int64 `json:"maxOpenConnections,omitempty"`

	// Specifies the maximum number of idle connections to the database.
	// A zero uses the value of max_open_connections and a negative value disables idle connections.
	// If larger than max_open_connections it will be reduced to be equal.
	MaxIdleConnections int64 `json:"maxIdleConnections,omitempty"`

	// Specifies the maximum amount of time a connection may be reused.
	// If <= 0s connections are reused forever.
	MaxConnectionLifetime string `json:"maxConnectionLifetime,omitempty"`
}

// MariaDBConfiguration defines a MariaDB app configuration.
// https://www.vaultproject.io/api/secret/databases/index.html
// https://www.vaultproject.io/api/secret/databases/mysql-maria.html#configure-connection
type MariaDBConfiguration struct {
	// DatabaseRef refers to a MariaDB database AppBinding in any namespace
	DatabaseRef appcat.AppReference `json:"databaseRef"`

	// Specifies the name of the plugin to use for this connection.
	// Default plugin:
	//  - for mysql: mysql-database-plugin
	PluginName string `json:"pluginName,omitempty"`

	// List of the roles allowed to use this connection.
	// Defaults to empty (no roles), if contains a "*" any role can use this connection.
	AllowedRoles []string `json:"allowedRoles,omitempty"`

	// Specifies the maximum number of open connections to the database.
	MaxOpenConnections int64 `json:"maxOpenConnections,omitempty"`

	// Specifies the maximum number of idle connections to the database.
	// A zero uses the value of max_open_connections and a negative value disables idle connections.
	// If larger than max_open_connections it will be reduced to be equal.
	MaxIdleConnections int64 `json:"maxIdleConnections,omitempty"`

	// Specifies the maximum amount of time a connection may be reused.
	// If <= 0s connections are reused forever.
	MaxConnectionLifetime string `json:"maxConnectionLifetime,omitempty"`
}

// KVConfiguration defines a Key-Value engine configuration
// TODO: fill in doc links
type KVConfiguration struct {
	// The version of the KV engine to enable. Defaults to "1", can be either "1" or "2"
	Version int64 `json:"version,omitempty"`

	// The maximum number of versions to keep for any given key. Defaults to 0, which indicates that the Vault default (10) should be
	// used.
	MaxVersions int64 `json:"maxVersions,omitempty"`

	// If true, then all operations on the KV store require the cas (Compare-and-Swap) parameter to be set.
	// https://www.vaultproject.io/api-docs/secret/kv/kv-v2#cas_required
	// https://www.vaultproject.io/docs/secrets/kv/kv-v2#usage
	CasRequired bool `json:"casRequired,omitempty"`

	// If set, keys will be automatically deleted after this length of time. Accepts a Go duration format
	// string.
	// https://golang.org/pkg/time/#ParseDuration
	DeleteVersionsAfter string `json:"deleteVersionsAfter,omitempty"`
}

// ElasticsearchConfiguration defines a Elasticsearch app configuration.
// https://www.vaultproject.io/api-docs/secret/databases/elasticdb
// TODO: Fill in the fields
type ElasticsearchConfiguration struct {
	// Specifies the Elasticsearch database appbinding reference
	DatabaseRef appcat.AppReference `json:"databaseRef"`

	// List of the roles allowed to use this connection.
	// Defaults to empty (no roles), if contains a "*" any role can use this connection.
	AllowedRoles []string `json:"allowedRoles,omitempty"`

	// Specifies the name of the plugin to use for this connection.
	// Default plugin:
	//  - for elasticsearch: elasticsearch-database-plugin
	PluginName string `json:"pluginName,omitempty"`

	// The URL for Elasticsearch's API ("http://localhost:9200").
	// +kubebuilder:validation:Required
	Url string `json:"url,omitempty"`

	// The username to be used in the connection URL ("vault").
	// +kubebuilder:validation:Required
	Username string `json:"username,omitempty"`

	// The password to be used in the connection URL ("pa55w0rd").
	// +kubebuilder:validation:Required
	Password string `json:"password,omitempty"`

	// The path to a PEM-encoded CA cert file to use to verify the Elasticsearch server's identity.
	CACert string `json:"caCert,omitempty"`

	// The path to a directory of PEM-encoded CA cert files to use to verify the Elasticsearch server's identity.
	CAPath string `json:"caPath,omitempty"`

	// The path to the certificate for the Elasticsearch client to present for communication.
	ClientCert string `json:"clientCert,omitempty"`

	// The path to the key for the Elasticsearch client to use for communication.
	ClientKey string `json:"clientKey,omitempty"`

	// This, if set, is used to set the SNI host when connecting via 1TLS.
	TLSServerName string `json:"tlsServerName,omitempty"`

	// Not recommended. Default to false. Can be set to true to disable SSL verification.
	// +kubebuilder:default:=false
	Insecure bool `json:"insecure,omitempty"`
}

type SecretEnginePhase string

type SecretEngineStatus struct {
	Phase SecretEnginePhase `json:"phase,omitempty"`

	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	Conditions []kmapi.Condition `json:"conditions,omitempty"`

	// Path defines the path used to enable this secret engine
	// Visible to user but immutable
	Path string `json:"path,omitempty"`
}
