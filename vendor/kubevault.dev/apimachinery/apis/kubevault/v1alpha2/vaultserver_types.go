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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	ResourceKindVaultServer = "VaultServer"
	ResourceVaultServer     = "vaultserver"
	ResourceVaultServers    = "vaultservers"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:path=vaultservers,singular=vaultserver,shortName=vs,categories={vault,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Replicas",type="string",JSONPath=".spec.replicas"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type VaultServer struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              VaultServerSpec   `json:"spec,omitempty"`
	Status            VaultServerStatus `json:"status,omitempty"`
}

type VaultServerSpec struct {
	// Version of VaultServer to be deployed.
	Version string `json:"version"`

	// Number of instances to deploy for a VaultServer.
	Replicas *int32 `json:"replicas,omitempty"`

	// ConfigSecret is an optional field to provide extra configuration for vault.
	// This secret contain extra config for vault
	// File name should be 'vault.hcl'.
	// If specified, this file will be appended to the controller configuration file.
	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`

	// DataSources is a list of Configmaps/Secrets in the same namespace as the VaultServer
	// object, which shall be mounted into the VaultServer Pods.
	// The data are mounted into /etc/vault/data/<name>.
	// The first data will be named as "data-0", second one will be named as "data-1" and so on.
	// +optional
	DataSources []core.VolumeSource `json:"dataSources,omitempty"`

	// TLS policy of vault nodes
	// +optional
	TLS *kmapi.TLSConfig `json:"tls,omitempty"`

	// backend storage configuration for vault
	Backend BackendStorageSpec `json:"backend"`

	// Unsealer configuration for vault
	// +optional
	Unsealer *UnsealerSpec `json:"unsealer,omitempty"`

	// Specifies the list of auth methods to enable
	// +optional
	AuthMethods []AuthMethod `json:"authMethods,omitempty"`

	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`

	// PodTemplate is an optional configuration for pods used to run vault
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

	// ServiceTemplates is an optional configuration for services used to expose database
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty"`

	// Indicates that the vault server is halted and all offshoot Kubernetes resources except PVCs are deleted.
	// +optional
	Halted bool `json:"halted,omitempty"`

	// TerminationPolicy controls the delete operation for vault server
	// +optional
	TerminationPolicy TerminationPolicy `json:"terminationPolicy,omitempty"`

	// AllowedSecretEngines defines the types of Secret Engines that MAY be attached to a
	// Listener and the trusted namespaces where those Route resources MAY be
	// present.
	//
	// Although a client request may match multiple route rules, only one rule
	// may ultimately receive the request. Matching precedence MUST be
	// determined in order of the following criteria:
	//
	// * The most specific match as defined by the Route type.
	// * The oldest Route based on creation timestamp. For example, a Route with
	//   a creation timestamp of "2020-09-08 01:02:03" is given precedence over
	//   a Route with a creation timestamp of "2020-09-08 01:02:04".
	// * If everything else is equivalent, the Route appearing first in
	//   alphabetical order (namespace/name) should be given precedence. For
	//   example, foo/bar is given precedence over foo/baz.
	//
	// All valid rules within a Route attached to this Listener should be
	// implemented. Invalid Route rules can be ignored (sometimes that will mean
	// the full Route). If a Route rule transitions from valid to invalid,
	// support for that Route rule should be dropped to ensure consistency. For
	// example, even if a filter specified by a Route rule is invalid, the rest
	// of the rules within that Route should still be supported.
	//
	// Support: Core
	// +kubebuilder:default={namespaces:{from: Same}}
	// +optional
	AllowedSecretEngines *AllowedSecretEngines `json:"allowedSecretEngines,omitempty"`

	// HealthChecker defines attributes of the health checker
	// +optional
	// +kubebuilder:default={periodSeconds: 10, timeoutSeconds: 10, failureThreshold: 1}
	HealthChecker kmapi.HealthCheckSpec `json:"healthChecker"`

	// ExposePrimary, when true, makes the operator create an additional
	// <vault-name>-primary Service whose selector narrows to the active (leader)
	// node, alongside the always-all-nodes <vault-name> Service. A client that
	// requires strict read-after-write consistency, and cannot tolerate reading from
	// a lagging standby, binds to the primary Service; everything else, including the
	// default AppBinding, keeps using the all-nodes Service and is unaffected.
	//
	// Requires an HA-capable storage backend and a supported distribution; the
	// admission webhook rejects it otherwise. The primary Service has no endpoints
	// during a leader election (the brief window with no active node), which is the
	// cost of the guarantee. Defaults to false. See design/primary-service-routing.md.
	// +optional
	ExposePrimary bool `json:"exposePrimary,omitempty"`

	// IsolateTenants is the master opt-in for per-tenant OpenBao namespace isolation.
	// When true (and the backend is namespace-capable, i.e. OpenBao / Vault Enterprise),
	// SecretEngines on this server whose referenced database lives in a client-org
	// namespace are provisioned in a matching OpenBao namespace keyed on the org-id, and
	// an explicit SecretEngine.spec.namespace is also honored. When false (the default),
	// everything stays in the root namespace and any SecretEngine.spec.namespace is rejected.
	// Enabling this on a running server is safe: existing mounts are never moved on their own
	// (migration is admin-authorized). See design/tenant-namespace-design.md.
	// +optional
	IsolateTenants bool `json:"isolateTenants,omitempty"`

	// RelayPlacementRef points to an OCM Placement object (cluster.open-cluster-management.io/v1beta1)
	// in the same namespace as the VaultServer. When set, the operator deploys a VaultRelay to
	// every managed cluster selected by the Placement, using one ManifestWork per cluster.
	// Requires the OCM hub CRDs (Placement, PlacementDecision, ManifestWork) to be installed;
	// the field is ignored with a warning condition otherwise.
	// +optional
	RelayPlacementRef *core.LocalObjectReference `json:"relayPlacementRef,omitempty"`

	// RelayTemplate customizes the VaultRelays stamped out for clusters selected by
	// RelayPlacementRef. Per-cluster fields (spokeName, hubVaultRef, join material)
	// are filled in by the operator.
	// +optional
	RelayTemplate *VaultRelayTemplate `json:"relayTemplate,omitempty"`
}

// VaultRelayTemplate is the subset of VaultRelaySpec a hub admin may pre-set for
// placement-driven spoke relay deployments.
type VaultRelayTemplate struct {
	// Namespace on the managed cluster where the VaultRelay and its companion
	// resources are created. Defaults to the VaultServer's namespace.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// PodTemplate is an optional configuration for the spoke-relay pods.
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

	// BootstrapTokenTTL controls the TTL of hub bootstrap tokens minted per spoke
	// (and therefore their rotation period). Default 24h, minimum 1h.
	// +optional
	BootstrapTokenTTL *metav1.Duration `json:"bootstrapTokenTTL,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

type VaultServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VaultServer `json:"items,omitempty"`
}

type VaultServerStatus struct {
	// ObservedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Phase indicates the state this Vault server jumps in.
	// +optional
	Phase VaultServerPhase `json:"phase,omitempty"`

	// Initialized indicates if the Vault service is initialized.
	// +optional
	Initialized bool `json:"initialized,omitempty"`

	// ServiceName is the LB service for accessing vault nodes.
	// +optional
	ServiceName string `json:"serviceName,omitempty"`

	// ClientPort is the port for vault client to access.
	// It's the same on client LB service and vault nodes.
	// +optional
	ClientPort int64 `json:"clientPort,omitempty"`

	// VaultStatus is the set of Vault node specific statuses: Active, Standby, and Sealed
	// +optional
	VaultStatus VaultStatus `json:"vaultStatus,omitempty"`

	// PodNames of updated Vault nodes. Updated means the Vault container image version
	// matches the spec's version.
	// +optional
	UpdatedNodes []string `json:"updatedNodes,omitempty"`

	// Represents the latest available observations of a VaultServer current state.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`

	// Status of the vault auth methods
	// +optional
	AuthMethodStatus []AuthMethodStatus `json:"authMethodStatus,omitempty"`

	// RelayPlacement summarizes spoke relay rollout when spec.relayPlacementRef is set.
	// +optional
	RelayPlacement *RelayPlacementStatus `json:"relayPlacement,omitempty"`
}

// RelayPlacementStatus summarizes the rollout of spoke relays to managed clusters.
type RelayPlacementStatus struct {
	// Placement is the resolved Placement name.
	// +optional
	Placement string `json:"placement,omitempty"`

	// Selected is the number of clusters currently listed in the PlacementDecisions.
	// +optional
	Selected int32 `json:"selected,omitempty"`

	// Applied is the number of clusters whose ManifestWork has condition Applied=True.
	// +optional
	Applied int32 `json:"applied,omitempty"`

	// Ready is the number of clusters whose VaultRelay reports phase Connected
	// (scraped via ManifestWork status feedback).
	// +optional
	Ready int32 `json:"ready,omitempty"`

	// Clusters holds per-cluster detail.
	// +optional
	Clusters []SpokeClusterStatus `json:"clusters,omitempty"`
}

// SpokeClusterStatus is the per managed cluster rollout state.
type SpokeClusterStatus struct {
	// ClusterName is the ManagedCluster name.
	ClusterName string `json:"clusterName"`

	// Phase mirrors the spoke VaultRelay phase (Pending|Connected|Disconnected|Error)
	// plus hub-side values (WorkApplied, WorkProgressing, WorkDegraded).
	// +optional
	Phase string `json:"phase,omitempty"`

	// TokenExpiry is when the current bootstrap token for this spoke expires.
	// +optional
	TokenExpiry *metav1.Time `json:"tokenExpiry,omitempty"`

	// CertExpiry is when this spoke's mTLS client certificate expires, as
	// observed by the hub relay backend (relay/spokes). Nil when unknown: the
	// spoke is not connected, or the hub captured no verified peer cert.
	// +optional
	CertExpiry *metav1.Time `json:"certExpiry,omitempty"`
}

// AllowedSecretEngines defines which Secret Engines may be attached to this Listener.
type AllowedSecretEngines struct {
	// Namespaces indicates namespaces from which Secret Engines may be attached to this
	// Listener. This is restricted to the namespace of this VaultServer by default.
	//
	// +optional
	// +kubebuilder:default={from: Same}
	Namespaces *SecretEngineNamespaces `json:"namespaces,omitempty"`

	// SecretEngines specifies the types of Secret Engines that are allowed to bind
	// to this VaultServer. When unspecified or empty, all types of Secret Engines
	// are allowed.
	//
	// +optional
	SecretEngines []SecretEngineType `json:"secretEngines,omitempty"`
}

// +kubebuilder:validation:Enum=kv;pki;aws;azure;gcp;postgres;mongodb;mysql;mariadb;elasticsearch;redis;druid;etcd;hanadb;ignite;kafka;milvus;mssqlserver;neo4j;oracle;qdrant;rabbitmq;solr;weaviate
type SecretEngineType string

const (
	SecretEngineTypeKV            SecretEngineType = "kv"
	SecretEngineTypePKI           SecretEngineType = "pki"
	SecretEngineTypeAWS           SecretEngineType = "aws"
	SecretEngineTypeAzure         SecretEngineType = "azure"
	SecretEngineTypeGCP           SecretEngineType = "gcp"
	SecretEngineTypePostgres      SecretEngineType = "postgres"
	SecretEngineTypeMongoDB       SecretEngineType = "mongodb"
	SecretEngineTypeMySQL         SecretEngineType = "mysql"
	SecretEngineTypeMariaDB       SecretEngineType = "mariadb"
	SecretEngineTypeElasticsearch SecretEngineType = "elasticsearch"
	SecretEngineTypeRedis         SecretEngineType = "redis"
	SecretEngineTypeDruid         SecretEngineType = "druid"
	SecretEngineTypeEtcd          SecretEngineType = "etcd"
	SecretEngineTypeHanaDB        SecretEngineType = "hanadb"
	SecretEngineTypeIgnite        SecretEngineType = "ignite"
	SecretEngineTypeKafka         SecretEngineType = "kafka"
	SecretEngineTypeMilvus        SecretEngineType = "milvus"
	SecretEngineTypeMSSQLServer   SecretEngineType = "mssqlserver"
	SecretEngineTypeNeo4j         SecretEngineType = "neo4j"
	SecretEngineTypeOracle        SecretEngineType = "oracle"
	SecretEngineTypeQdrant        SecretEngineType = "qdrant"
	SecretEngineTypeRabbitMQ      SecretEngineType = "rabbitmq"
	SecretEngineTypeSolr          SecretEngineType = "solr"
	SecretEngineTypeWeaviate      SecretEngineType = "weaviate"
)

// FromNamespaces specifies namespace from which Secret Engines may be attached to a
// VaultServer.
//
// +kubebuilder:validation:Enum=All;Selector;Same
type FromNamespaces string

const (
	// Secret Engines in all namespaces may be attached to this VaultServer.
	NamespacesFromAll FromNamespaces = "All"
	// Only Secret Engines in namespaces selected by the selector may be attached to
	// this VaultServer.
	NamespacesFromSelector FromNamespaces = "Selector"
	// Only Secret Engines in the same namespace as the VaultServer may be attached to this
	// VaultServer.
	NamespacesFromSame FromNamespaces = "Same"
)

// SecretEngineNamespaces indicate which namespaces Secret Engines should be selected from.
type SecretEngineNamespaces struct {
	// From indicates where Secret Engines will be selected for this VaultServer. Possible
	// values are:
	// * All: Secret Engines in all namespaces may be used by this VaultServer.
	// * Selector: Secret Engines in namespaces selected by the selector may be used by
	//   this VaultServer.
	// * Same: Only Secret Engines in the same namespace may be used by this VaultServer.
	//
	// +optional
	// +kubebuilder:default=Same
	From *FromNamespaces `json:"from,omitempty"`

	// Selector must be specified when From is set to "Selector". In that case,
	// only Secret Engines in Namespaces matching this Selector will be selected by this
	// VaultServer. This field is ignored for other values of "From".
	//
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`
}

type VaultStatus struct {
	// PodName of the active (leader) Vault node. The active node is unsealed and
	// every write is committed there. Standby nodes serve reads and forward
	// writes to it. The <vault-name>-primary Service, created when
	// spec.exposePrimary is true, points at this node alone; the <vault-name>
	// Service always points at all nodes.
	// Empty while no node is active, for example during a leader election.
	// +optional
	Active string `json:"active,omitempty"`

	// PodNames of the standby Vault nodes. Standby nodes are unsealed.
	// Standby nodes do not process requests, and instead redirect to the active Vault.
	// +optional
	Standby []string `json:"standby,omitempty"`

	// PodNames of Sealed Vault nodes. Sealed nodes MUST be unsealed to
	// become standby or leader.
	// +optional
	Sealed []string `json:"sealed,omitempty"`

	// PodNames of Unsealed Vault nodes.
	// +optional
	Unsealed []string `json:"unsealed,omitempty"`
}

// TLSPolicy defines the TLS policy of the vault nodes
// If this is not set, operator will auto-gen TLS assets and secrets.
type TLSPolicy struct {
	// TLSSecret is the secret containing TLS certs used by each vault node
	// for the communication between the vault server and its clients.
	// The secret should contain three files:
	// 	- tls.crt
	// 	- tls.key
	//
	// The server certificate must allow the following wildcard domains:
	// 	- localhost
	// 	- *.<namespace>.pod
	// 	- <vaultServer-name>.<namespace>.svc
	TLSSecret string `json:"tlsSecret"`

	// CABundle is a PEM encoded CA bundle which will be used to validate the serving certificate.
	// +optional
	CABundle []byte `json:"caBundle,omitempty"`
}
