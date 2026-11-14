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
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	VaultContainerName         = "vault"
	VaultUnsealerContainerName = "vault-unsealer"
	VaultInitContainerName     = "vault-config"
	VaultExporterContainerName = "vault-exporter"
	VaultRelayContainerName    = "spoke-relay"
)

// +kubebuilder:validation:Enum=Initializing;Unsealing;Sealed;Ready;NotReady;Critical
type VaultServerPhase string

const (
	// used for VaultServer that are Initializing
	VaultServerPhaseInitializing VaultServerPhase = "Initializing"
	// used for VaultServer that are Unsealing
	VaultServerPhaseUnsealing VaultServerPhase = "Unsealing"
	// used for VaultServer that are sealed
	VaultServerPhaseSealed VaultServerPhase = "Sealed"
	// used for VaultServer that are Ready
	VaultServerPhaseReady VaultServerPhase = "Ready"
	// used for VaultServer that are NotReady
	VaultServerPhaseNotReady VaultServerPhase = "NotReady"
	// used for VaultServer that are Critical
	VaultServerPhaseCritical VaultServerPhase = "Critical"
)

// +kubebuilder:validation:Enum=Halt;Delete;WipeOut;DoNotTerminate
type TerminationPolicy string

const (
	// Deletes VaultServer pods, service but leave the PVCs and stash backup data intact.
	TerminationPolicyHalt TerminationPolicy = "Halt"
	// Deletes VaultServer pods, service, pvcs but leave the stash backup data intact.
	TerminationPolicyDelete TerminationPolicy = "Delete"
	// Deletes VaultServer pods, service, pvcs and stash backup data.
	TerminationPolicyWipeOut TerminationPolicy = "WipeOut"
	// Rejects attempt to delete VaultServer using ValidationWebhook.
	TerminationPolicyDoNotTerminate TerminationPolicy = "DoNotTerminate"
)

// +kubebuilder:validation:Enum=internal;vault;stats;primary
type ServiceAlias string

const (
	VaultServerServiceInternal ServiceAlias = "internal"
	VaultServerServiceVault    ServiceAlias = "vault"
	VaultServerServiceStats    ServiceAlias = "stats"
	// VaultServerServicePrimary is the client Service that selects only the active
	// (leader) node. The operator creates it only when spec.exposePrimary is true;
	// consumers that need strict read-after-write bind to it, while the vault
	// Service keeps selecting all nodes.
	VaultServerServicePrimary ServiceAlias = "primary"
)

type NamedServiceTemplateSpec struct {
	// Alias represents the identifier of the service.
	Alias ServiceAlias `json:"alias"`

	// ServiceTemplate is an optional configuration for a service used to expose VaultServer
	// +optional
	ofst.ServiceTemplateSpec `json:",inline,omitempty"`
}

// +kubebuilder:validation:Enum=ca;server;client;storage
type VaultCertificateAlias string

const (
	VaultCACert      VaultCertificateAlias = "ca"
	VaultServerCert  VaultCertificateAlias = "server"
	VaultClientCert  VaultCertificateAlias = "client"
	VaultStorageCert VaultCertificateAlias = "storage"
)

// +kubebuilder:validation:Enum=aerospike;alicloudoss;azure;cassandra;cockroachdb;consul;couchdb;dynamodb;etcd;file;gcs;inmem;mssql;mysql;oci;postgresql;raft;s3;spanner;swift;zookeeper
type VaultServerBackend string

const (
	VaultServerAerospike   VaultServerBackend = "aerospike"
	VaultServerAlicloudOSS VaultServerBackend = "alicloudoss"
	VaultServerAzure       VaultServerBackend = "azure"
	VaultServerCassandra   VaultServerBackend = "cassandra"
	VaultServerCockroachDB VaultServerBackend = "cockroachdb"
	VaultServerConsul      VaultServerBackend = "consul"
	VaultServerCouchDB     VaultServerBackend = "couchdb"
	VaultServerDynamoDB    VaultServerBackend = "dynamodb"
	VaultServerEtcd        VaultServerBackend = "etcd"
	VaultServerFile        VaultServerBackend = "file"
	VaultServerGcs         VaultServerBackend = "gcs"
	VaultServerInmem       VaultServerBackend = "inmem"
	VaultServerMSSQL       VaultServerBackend = "mssql"
	VaultServerMySQL       VaultServerBackend = "mysql"
	VaultServerOCI         VaultServerBackend = "oci"
	VaultServerPostgreSQL  VaultServerBackend = "postgresql"
	VaultServerRaft        VaultServerBackend = "raft"
	VaultServerS3          VaultServerBackend = "s3"
	VaultServerSpanner     VaultServerBackend = "spanner"
	VaultServerSwift       VaultServerBackend = "swift"
	VaultServerZookeeper   VaultServerBackend = "zookeeper"
)
