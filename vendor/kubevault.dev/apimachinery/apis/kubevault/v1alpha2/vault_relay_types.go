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
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	ResourceKindVaultRelay = "VaultRelay"
	ResourceVaultRelay     = "vaultrelay"
	ResourceVaultRelays    = "vaultrelays"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=vaultrelays,singular=vaultrelay,shortName=vr,categories={vault,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Spoke",type="string",JSONPath=".spec.spokeName"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type VaultRelay struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              VaultRelaySpec   `json:"spec,omitempty"`
	Status            VaultRelayStatus `json:"status,omitempty"`
}

type VaultRelaySpec struct {
	// HubVaultRef references the hub VaultServer
	HubVaultRef HubVaultReference `json:"hubVaultRef"`

	// SpokeName is the unique identifier for this spoke cluster
	SpokeName string `json:"spokeName"`

	// IsolateTenants opts this spoke into per-tenant OpenBao namespace isolation against
	// the hub. Hub-managed (placement-driven) spokes inherit this from the hub
	// VaultServer.spec.isolateTenants automatically; a standalone VaultRelay sets it
	// explicitly. It is stamped into the hub AppBinding's VaultServerConfiguration
	// (design/tenant-namespace-design.md §5.2). Default false.
	// +optional
	IsolateTenants bool `json:"isolateTenants,omitempty"`

	// TokenSecretRef references a secret containing vault token for authentication
	// Secret data:
	//  - token: <vault-token>
	// +optional
	TokenSecretRef *core.LocalObjectReference `json:"tokenSecretRef,omitempty"`

	// Image is the spoke-relay container image
	// +optional
	Image string `json:"image,omitempty"`

	// TLS configuration for gRPC connection
	// +optional
	TLS *VaultRelayTLSConfig `json:"tls,omitempty"`

	// Reconnect settings for automatic reconnection
	// +optional
	// +kubebuilder:default={enabled: true, backoffSeconds: 5, maxBackoffSeconds: 300}
	Reconnect *ReconnectConfig `json:"reconnect,omitempty"`

	// PodTemplate is an optional configuration for the spoke-relay pod
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

	// Bootstrap configures the automated `bao relay join` flow. When set, the
	// spoke-relay pod runs a join init container that exchanges the bootstrap
	// token for mTLS client credentials before the long-running relay starts.
	// Exactly one of Bootstrap or TLS.CertSecret (pre-provisioned credentials)
	// should be used.
	// +optional
	Bootstrap *RelayBootstrapConfig `json:"bootstrap,omitempty"`
}

// RelayBootstrapConfig configures the automated `bao relay join` trust bootstrap.
type RelayBootstrapConfig struct {
	// JoinSecretRef references a Secret with the join parameters:
	//  - token:       hub bootstrap token (<id>.<secret>)
	//  - hubCertHash: "sha256:<hex>" SPKI pin of the spoke-CA
	//  - caBundle:    optional PEM CA bundle for the hub Vault API endpoint
	JoinSecretRef core.LocalObjectReference `json:"joinSecretRef"`
}

// HubVaultReference contains information to connect to hub vault
type HubVaultReference struct {
	// Name of the VaultServer in hub cluster
	Name string `json:"name"`

	// Namespace of the VaultServer in hub cluster
	Namespace string `json:"namespace"`

	// Address is the hub vault URL (e.g., http://10.2.0.88:30820)
	Address string `json:"address"`

	// GRPCPort is the port for gRPC proxy communication
	// +optional
	// +kubebuilder:default=50053
	GRPCPort int32 `json:"grpcPort,omitempty"`

	// CABundle is the PEM bundle used to verify the hub Vault API endpoint.
	// Takes precedence over the caBundle key of bootstrap.joinSecretRef.
	// +optional
	CABundle []byte `json:"caBundle,omitempty"`
}

// VaultRelayTLSConfig contains TLS configuration for spoke-relay
type VaultRelayTLSConfig struct {
	// Enabled indicates whether TLS is enabled
	// +optional
	Enabled bool `json:"enabled,omitempty"`

	// CASecret contains the CA certificate
	// Secret data:
	//  - ca.crt: <value>
	// +optional
	CASecret *core.LocalObjectReference `json:"caSecret,omitempty"`

	// CertSecret contains the client certificate and key
	// Secret data:
	//  - tls.crt: <value>
	//  - tls.key: <value>
	// +optional
	CertSecret *core.LocalObjectReference `json:"certSecret,omitempty"`
}

// ReconnectConfig contains automatic reconnection settings
type ReconnectConfig struct {
	// Enabled indicates whether auto-reconnect is enabled
	// +optional
	// +kubebuilder:default=true
	Enabled bool `json:"enabled,omitempty"`

	// BackoffSeconds is the initial backoff duration in seconds
	// +optional
	// +kubebuilder:default=5
	BackoffSeconds int32 `json:"backoffSeconds,omitempty"`

	// MaxBackoffSeconds is the maximum backoff duration in seconds
	// +optional
	// +kubebuilder:default=300
	MaxBackoffSeconds int32 `json:"maxBackoffSeconds,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

type VaultRelayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VaultRelay `json:"items,omitempty"`
}

// +kubebuilder:validation:Enum=Pending;Connected;Disconnected;Error
type VaultRelayPhase string

const (
	VaultRelayPhasePending      VaultRelayPhase = "Pending"
	VaultRelayPhaseConnected    VaultRelayPhase = "Connected"
	VaultRelayPhaseDisconnected VaultRelayPhase = "Disconnected"
	VaultRelayPhaseError        VaultRelayPhase = "Error"
)

type VaultRelayStatus struct {
	// ObservedGeneration is the most recent generation observed
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Phase indicates the current state of the VaultRelay
	// +optional
	Phase VaultRelayPhase `json:"phase,omitempty"`

	// LastHeartbeat is the timestamp of the last successful heartbeat
	// +optional
	LastHeartbeat *metav1.Time `json:"lastHeartbeat,omitempty"`

	// AppBindingRef references the created AppBinding for hub vault
	// +optional
	AppBindingRef *kmapi.ObjectReference `json:"appBindingRef,omitempty"`

	// PodName is the name of the spoke-relay pod
	// +optional
	PodName string `json:"podName,omitempty"`

	// CertExpiry of the current spoke client certificate, if known.
	// +optional
	CertExpiry *metav1.Time `json:"certExpiry,omitempty"`

	// Conditions represent the latest available observations of the VaultRelay's state
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}
