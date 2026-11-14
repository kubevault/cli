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

// TODO : set defaults and validation
// BackendStorageSpec defines storage backend configuration of vault
//
// FoundationDB is intentionally not represented here: unlike every other
// backend, it requires CGo and the native libfdb_c client library, so it
// isn't compiled into standard OpenBao/Vault images and can't be assumed
// available in a generic VaultServer deployment.
//
// New backends added here (CockroachDB, Cassandra, Zookeeper, CouchDB,
// MSSQL, Spanner, Aerospike, OCI, AlicloudOSS) are v1alpha2-only for now;
// v1alpha1 already predates several v1alpha2-only conventions used by
// their spec types (SecretRef-typed credentials, DatabaseRef, HAEnabled),
// so mirroring them would need bespoke field-mapping conversion functions
// rather than the mechanical 1:1 copy used elsewhere. v1alpha2 is the
// storage/hub version and the only one the operator's storage package
// consumes; a v1alpha1 backport can follow the existing
// Convert_v1alpha1_<X>Spec_To_v1alpha2_<X>Spec pattern if ever needed.
type BackendStorageSpec struct {
	// ref: https://www.vaultproject.io/docs/configuration/storage/in-memory.html
	// +optional
	Inmem *InmemSpec `json:"inmem,omitempty"`

	// +optional
	Etcd *EtcdSpec `json:"etcd,omitempty"`

	// +optional
	Gcs *GcsSpec `json:"gcs,omitempty"`

	// +optional
	S3 *S3Spec `json:"s3,omitempty"`

	// +optional
	Azure *AzureSpec `json:"azure,omitempty"`

	// +optional
	PostgreSQL *PostgreSQLSpec `json:"postgresql,omitempty"`

	// +optional
	MySQL *MySQLSpec `json:"mysql,omitempty"`

	// +optional
	File *FileSpec `json:"file,omitempty"`

	// +optional
	DynamoDB *DynamoDBSpec `json:"dynamodb,omitempty"`

	// +optional
	Swift *SwiftSpec `json:"swift,omitempty"`

	// +optional
	Consul *ConsulSpec `json:"consul,omitempty"`

	// +optional
	Raft *RaftSpec `json:"raft,omitempty"`

	// +optional
	CockroachDB *CockroachDBSpec `json:"cockroachdb,omitempty"`

	// +optional
	Cassandra *CassandraSpec `json:"cassandra,omitempty"`

	// +optional
	Zookeeper *ZookeeperSpec `json:"zookeeper,omitempty"`

	// +optional
	CouchDB *CouchDBSpec `json:"couchdb,omitempty"`

	// +optional
	MSSQL *MSSQLSpec `json:"mssql,omitempty"`

	// +optional
	Spanner *SpannerSpec `json:"spanner,omitempty"`

	// +optional
	Aerospike *AerospikeSpec `json:"aerospike,omitempty"`

	// +optional
	OCI *OCISpec `json:"oci,omitempty"`

	// +optional
	AlicloudOSS *AlicloudOSSSpec `json:"alicloudoss,omitempty"`
}

// +kubebuilder:validation:Enum=disable;require;verify-ca;verify-full
type PostgresSSLMode string

const (
	// PostgresSSLModeDisable represents `disable` sslMode. It ensures that the server does not use TLS/SSL.
	PostgresSSLModeDisable PostgresSSLMode = "disable"

	// Always SSL (skip verification)
	PostgressSSLModeRequire PostgresSSLMode = "require"

	// Always SSL (verify that the certificate presented by the server was signed by a trusted CA)
	PostgressSSLModeVerifyCA PostgresSSLMode = "verify-ca"

	// PostgresSSLModeVerifyFull represents `verify-full` sslmode. I want my data encrypted, and I accept the overhead.
	// I want to be sure that I connect to a server I trust, and that it's the one I specify.
	PostgresSSLModeVerifyFull PostgresSSLMode = "verify-full"
)
