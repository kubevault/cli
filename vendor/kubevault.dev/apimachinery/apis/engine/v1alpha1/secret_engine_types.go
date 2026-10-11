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

	// Namespace specifies the OpenBao namespace for this SecretEngine.
	// Only applicable when using OpenBao distribution.
	// If specified, overrides the namespace from VaultServer.
	// Empty string means use VaultServer's namespace (or root if VaultServer has none).
	// Supports hierarchical namespaces (e.g., "tenant-1/project-a").
	// +optional
	Namespace string `json:"namespace,omitempty"`

	SecretEngineConfiguration `json:",inline"`

	DefaultLeaseTTL string `json:"defaultLeaseTTL,omitempty"`
	MaxLeaseTTL     string `json:"maxLeaseTTL,omitempty"`
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
	Druid         *DruidConfiguration         `json:"druid,omitempty"`
	Etcd          *EtcdConfiguration          `json:"etcd,omitempty"`
	HanaDB        *HanaDBConfiguration        `json:"hanadb,omitempty"`
	Ignite        *IgniteConfiguration        `json:"ignite,omitempty"`
	Kafka         *KafkaConfiguration         `json:"kafka,omitempty"`
	Milvus        *MilvusConfiguration        `json:"milvus,omitempty"`
	MSSQLServer   *MSSQLServerConfiguration   `json:"mssqlserver,omitempty"`
	Neo4j         *Neo4jConfiguration         `json:"neo4j,omitempty"`
	Oracle        *OracleConfiguration        `json:"oracle,omitempty"`
	Qdrant        *QdrantConfiguration        `json:"qdrant,omitempty"`
	RabbitMQ      *RabbitMQConfiguration      `json:"rabbitmq,omitempty"`
	Solr          *SolrConfiguration          `json:"solr,omitempty"`
	Weaviate      *WeaviateConfiguration      `json:"weaviate,omitempty"`
}

// DruidConfiguration defines an Apache Druid app configuration. The
// OpenBao `druid-database-plugin` (sigilr/openbao#12) provisions
// credentials via Druid's BasicSecurity coordinator API
// (authenticator user + role bindings), so the connection payload
// uses `url` (Druid coordinator REST endpoint) rather than a
// connection_url. Druid is dynamic: NewUser/UpdateUser/DeleteUser
// are all supported.
// https://druid.apache.org/docs/latest/operations/security-overview/
type DruidConfiguration struct {
	// Specifies the Druid database appbinding reference. The
	// AppBinding URL points at the Druid coordinator REST endpoint
	// (e.g. `http://druid-coordinator.demo.svc:8081`); the secret
	// contributes username/password used to authenticate against
	// Druid's BasicSecurity coordinator API.
	DatabaseRef appcat.AppReference `json:"databaseRef"`

	// Specifies the name of the plugin to use for this connection.
	// Default plugin:
	//  - for druid: druid-database-plugin
	// +optional
	PluginName string `json:"pluginName,omitempty"`

	// List of the roles allowed to use this connection.
	// Defaults to empty (no roles), if contains a "*" any role can use this connection.
	// +optional
	AllowedRoles []string `json:"allowedRoles,omitempty"`

	// Authenticator is the Druid BasicSecurity authenticator that will hold
	// dynamically issued credentials. Defaults to MyBasicMetadataAuthenticator
	// (Druid's Quick Start name).
	// +optional
	Authenticator string `json:"authenticator,omitempty"`

	// Authorizer is the Druid BasicSecurity authorizer that holds the roles
	// referenced by creationStatements. Defaults to MyBasicMetadataAuthorizer.
	// +optional
	Authorizer string `json:"authorizer,omitempty"`
}

// EtcdConfiguration defines an etcd app configuration. The OpenBao
// `etcd-database-plugin` (sigilr/openbao#50) issues dynamic credentials
// against etcd's built-in v3 Auth API: NewUser creates a native etcd
// user and grants it every role named in creationStatements. The
// connection payload uses `endpoints` (comma-separated etcd client
// URLs) rather than a connection_url. Etcd is dynamic:
// NewUser/UpdateUser/DeleteUser are all supported.
// https://etcd.io/docs/latest/op-guide/authentication/
type EtcdConfiguration struct {
	// Specifies the etcd database appbinding reference. The AppBinding's
	// URL is forwarded as the etcd client endpoint (`endpoints=`); the
	// secret contributes the root username/password used to
	// authenticate against etcd's v3 Auth API.
	DatabaseRef appcat.AppReference `json:"databaseRef"`

	// Specifies the name of the plugin to use for this connection.
	// Default plugin:
	//  - for etcd: etcd-database-plugin
	// +optional
	PluginName string `json:"pluginName,omitempty"`

	// List of the roles allowed to use this connection.
	// Defaults to empty (no roles), if contains a "*" any role can use this connection.
	// +optional
	AllowedRoles []string `json:"allowedRoles,omitempty"`
}

// HanaDBConfiguration defines a SAP HANA app configuration. The OpenBao
// `hana-database-plugin` (sigilr/openbao#3) was ported from pre-BUSL
// HashiCorp Vault and is a dynamic SQL-statement based plugin
// (postgres-style): the connection payload is the HANA DSN
// (`connection_url`) plus pool-tuning knobs.
// https://www.vaultproject.io/api/secret/databases/hanadb.html
type HanaDBConfiguration struct {
	// Specifies the HANA database appbinding reference. The AppBinding
	// URL is the HANA DSN (e.g. `hdb://<host>:<port>`); the secret
	// contributes username/password.
	DatabaseRef appcat.AppReference `json:"databaseRef"`

	// Specifies the name of the plugin to use for this connection.
	// Default plugin:
	//  - for hanadb: hana-database-plugin
	// +optional
	PluginName string `json:"pluginName,omitempty"`

	// List of the roles allowed to use this connection.
	// Defaults to empty (no roles), if contains a "*" any role can use this connection.
	// +optional
	AllowedRoles []string `json:"allowedRoles,omitempty"`

	// Specifies the maximum number of open connections to the database.
	// +optional
	MaxOpenConnections int64 `json:"maxOpenConnections,omitempty"`

	// Specifies the maximum number of idle connections to the database.
	// A zero uses the value of max_open_connections and a negative value disables idle connections.
	// If larger than max_open_connections it will be reduced to be equal.
	// +optional
	MaxIdleConnections int64 `json:"maxIdleConnections,omitempty"`

	// Specifies the maximum amount of time a connection may be reused.
	// If <= 0s connections are reused forever.
	// +optional
	MaxConnectionLifetime string `json:"maxConnectionLifetime,omitempty"`
}

// IgniteConfiguration defines an Apache Ignite app configuration. The OpenBao
// `ignite-database-plugin` (sigilr/openbao#14) executes dynamic SQL DDL
// (`CREATE USER` / `ALTER USER` / `DROP USER`) over Ignite's REST API
// (`cmd=qryfldexe`). The connection payload is the Ignite REST endpoint
// (`url`) plus Basic Auth username/password sourced from the AppBinding
// secret. Pool-tuning knobs are intentionally omitted because the REST API
// does not expose them.
// https://ignite.apache.org/docs/latest/security/authentication
type IgniteConfiguration struct {
	// Specifies the Ignite database appbinding reference. The AppBinding URL
	// is the Ignite REST endpoint (e.g. `http://ignite.demo.svc:8080`); the
	// secret contributes username/password.
	DatabaseRef appcat.AppReference `json:"databaseRef"`

	// Specifies the name of the plugin to use for this connection.
	// Default plugin:
	//  - for ignite: ignite-database-plugin
	// +optional
	PluginName string `json:"pluginName,omitempty"`

	// List of the roles allowed to use this connection.
	// Defaults to empty (no roles), if contains a "*" any role can use this connection.
	// +optional
	AllowedRoles []string `json:"allowedRoles,omitempty"`
}

// KafkaConfiguration defines an Apache Kafka app configuration. The
// OpenBao `kafka-database-plugin` (sigilr/openbao#15) writes SCRAM
// user records via the franz-go AdminClient, so the connection
// payload uses `brokers` (a comma-separated broker list) rather than
// `connection_url`. Kafka is dynamic: NewUser/UpdateUser/DeleteUser are
// all supported. ACLs are not yet implemented — the `acls` field on
// the role JSON is reserved and must currently be empty.
// https://kafka.apache.org/documentation/#security_sasl_scram
type KafkaConfiguration struct {
	// Specifies the Kafka cluster appbinding reference. The AppBinding
	// URL is the broker CSV (e.g. `broker1:9092,broker2:9092`) and the
	// referenced Secret contributes username/password for the SASL
	// principal used to manage SCRAM user records.
	DatabaseRef appcat.AppReference `json:"databaseRef"`

	// Specifies the name of the plugin to use for this connection.
	// Default plugin:
	//  - for kafka: kafka-database-plugin
	// +optional
	PluginName string `json:"pluginName,omitempty"`

	// List of the roles allowed to use this connection.
	// Defaults to empty (no roles), if contains a "*" any role can use this connection.
	// +optional
	AllowedRoles []string `json:"allowedRoles,omitempty"`
}

// MSSQLServerConfiguration defines a Microsoft SQL Server app
// configuration. The OpenBao `mssql-database-plugin` (sigilr/openbao#5)
// was ported from pre-BUSL HashiCorp Vault and is a dynamic
// SQL-statement based plugin (postgres-style): the connection payload is
// the SQL Server DSN (`connection_url`) plus pool-tuning knobs and the
// MSSQL-specific `contained_db` toggle.
// https://www.vaultproject.io/api/secret/databases/mssql.html
type MSSQLServerConfiguration struct {
	// Specifies the SQL Server database appbinding reference. The
	// AppBinding URL is the SQL Server DSN
	// (`sqlserver://<user>:<pass>@<host>:1433`); the secret contributes
	// username/password.
	DatabaseRef appcat.AppReference `json:"databaseRef"`

	// Specifies the name of the plugin to use for this connection.
	// Default plugin:
	//  - for mssqlserver: mssql-database-plugin
	// +optional
	PluginName string `json:"pluginName,omitempty"`

	// List of the roles allowed to use this connection.
	// Defaults to empty (no roles), if contains a "*" any role can use this connection.
	// +optional
	AllowedRoles []string `json:"allowedRoles,omitempty"`

	// Specifies the maximum number of open connections to the database.
	// +optional
	MaxOpenConnections int64 `json:"maxOpenConnections,omitempty"`

	// Specifies the maximum number of idle connections to the database.
	// A zero uses the value of max_open_connections and a negative value disables idle connections.
	// If larger than max_open_connections it will be reduced to be equal.
	// +optional
	MaxIdleConnections int64 `json:"maxIdleConnections,omitempty"`

	// Specifies the maximum amount of time a connection may be reused.
	// If <= 0s connections are reused forever.
	// +optional
	MaxConnectionLifetime string `json:"maxConnectionLifetime,omitempty"`

	// ContainedDB enables SQL Server contained-database authentication mode
	// (CREATE USER ... WITH PASSWORD against the user database). Defaults to
	// false (logins on the master DB).
	// +optional
	ContainedDB bool `json:"containedDB,omitempty"`
}

// OracleConfiguration defines an Oracle Database app configuration. The
// OpenBao `oracle-database-plugin` (sigilr/openbao#6) uses the pure-Go
// `sijms/go-ora/v2` driver and is a dynamic SQL-statement based plugin
// (postgres-style): the connection payload is the Oracle DSN
// (`connection_url`) plus pool-tuning knobs.
// https://docs.oracle.com/en/database/
type OracleConfiguration struct {
	// Specifies the Oracle database appbinding reference. The AppBinding
	// URL is the Oracle DSN (e.g. `oracle://<user>:<pass>@<host>:1521/<service>`);
	// the secret contributes username/password.
	DatabaseRef appcat.AppReference `json:"databaseRef"`

	// Specifies the name of the plugin to use for this connection.
	// Default plugin:
	//  - for oracle: oracle-database-plugin
	// +optional
	PluginName string `json:"pluginName,omitempty"`

	// List of the roles allowed to use this connection.
	// Defaults to empty (no roles), if contains a "*" any role can use this connection.
	// +optional
	AllowedRoles []string `json:"allowedRoles,omitempty"`

	// Specifies the maximum number of open connections to the database.
	// +optional
	MaxOpenConnections int64 `json:"maxOpenConnections,omitempty"`

	// Specifies the maximum number of idle connections to the database.
	// A zero uses the value of max_open_connections and a negative value disables idle connections.
	// If larger than max_open_connections it will be reduced to be equal.
	// +optional
	MaxIdleConnections int64 `json:"maxIdleConnections,omitempty"`

	// Specifies the maximum amount of time a connection may be reused.
	// If <= 0s connections are reused forever.
	// +optional
	MaxConnectionLifetime string `json:"maxConnectionLifetime,omitempty"`
}

// QdrantConfiguration defines a Qdrant app configuration.
type QdrantConfiguration struct {
	// Specifies the Qdrant database appbinding reference. The
	// AppBinding's URL is forwarded as the Qdrant HTTP endpoint
	// (`url=`); the secret's `api_key` (or `password`) field is
	// forwarded as the Qdrant API key (`api_key=`).
	DatabaseRef appcat.AppReference `json:"databaseRef"`

	// Specifies the name of the plugin to use for this connection.
	// Default plugin:
	//  - for qdrant: qdrant-database-plugin
	// +optional
	PluginName string `json:"pluginName,omitempty"`

	// List of the roles allowed to use this connection.
	// Defaults to empty (no roles), if contains a "*" any role can use this connection.
	// +optional
	AllowedRoles []string `json:"allowedRoles,omitempty"`
}

// SolrConfiguration defines an Apache Solr app configuration. The
// OpenBao `solr-database-plugin` (sigilr/openbao#11) provisions
// credentials via Solr's Security Plugin API (Basic Auth Plugin user
// + Rule-Based Authorization role bindings), so the connection
// payload uses `url` (Solr HTTP endpoint) rather than a
// connection_url. Solr is dynamic: NewUser/UpdateUser/DeleteUser
// are all supported.
// https://solr.apache.org/guide/solr/latest/deployment-guide/authentication-and-authorization-plugins.html
type SolrConfiguration struct {
	// Specifies the Solr database appbinding reference. The
	// AppBinding URL points at the Solr HTTP endpoint
	// (e.g. `http://solr.demo.svc:8983`); the secret contributes
	// username/password used to authenticate against Solr's
	// Security Plugin API.
	DatabaseRef appcat.AppReference `json:"databaseRef"`

	// Specifies the name of the plugin to use for this connection.
	// Default plugin:
	//  - for solr: solr-database-plugin
	// +optional
	PluginName string `json:"pluginName,omitempty"`

	// List of the roles allowed to use this connection.
	// Defaults to empty (no roles), if contains a "*" any role can use this connection.
	// +optional
	AllowedRoles []string `json:"allowedRoles,omitempty"`
}

// WeaviateConfiguration defines a Weaviate app configuration.
type WeaviateConfiguration struct {
	// Specifies the Weaviate database appbinding reference. The
	// AppBinding's URL is forwarded as the Weaviate HTTP endpoint
	// (`url=`); the secret's `AUTHENTICATION_APIKEY_ALLOWED_KEYS` field
	// is forwarded as the Weaviate API key (`api_key=`).
	DatabaseRef appcat.AppReference `json:"databaseRef"`

	// Specifies the name of the plugin to use for this connection.
	// Default plugin:
	//  - for weaviate: weaviate-database-plugin
	// +optional
	PluginName string `json:"pluginName,omitempty"`

	// List of the roles allowed to use this connection.
	// Defaults to empty (no roles), if contains a "*" any role can use this connection.
	// +optional
	AllowedRoles []string `json:"allowedRoles,omitempty"`
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

// RabbitMQConfiguration defines a RabbitMQ app configuration.
type RabbitMQConfiguration struct {
	// Specifies the primary RabbitMQ database AppBinding reference.
	DatabaseRef appcat.AppReference `json:"databaseRef"`

	// Specifies the name of the plugin to use for this connection.
	// Default plugin:
	//  - for rabbitmq: rabbitmq-database-plugin
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
}

type SecretEnginePhase string

type SecretEngineStatus struct {
	Phase SecretEnginePhase `json:"phase,omitempty"`

	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	Conditions []kmapi.Condition `json:"conditions,omitempty"`

	// Path defines the path used to enable this secret engine
	// Visible to user but immutable
	Path string `json:"path,omitempty"`

	// EffectiveNamespace is the OpenBao namespace this secret engine is actually
	// provisioned in (empty means the root namespace). It is the single source of truth
	// that every dependent {db}Role inherits, and the sticky anchor that prevents an
	// accidental re-mount when the desired namespace changes — migration to a new
	// namespace is admin-authorized. See design/tenant-namespace-design.md §5.2.
	// +optional
	EffectiveNamespace string `json:"effectiveNamespace,omitempty"`
}
