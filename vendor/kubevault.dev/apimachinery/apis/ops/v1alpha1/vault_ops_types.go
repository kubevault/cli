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
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
)

const (
	ResourceCodeVaultOpsRequest     = "vsops"
	ResourceKindVaultOpsRequest     = "VaultOpsRequest"
	ResourceSingularVaultOpsRequest = "vaultopsrequest"
	ResourcePluralVaultOpsRequest   = "vaultopsrequests"
)

// VaultOpsRequest defines a VaultServer operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=vaultopsrequests,singular=vaultopsrequest,shortName=vsops,categories={security,kubevault,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type VaultOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              VaultOpsRequestSpec   `json:"spec,omitempty"`
	Status            VaultOpsRequestStatus `json:"status,omitempty"`
}

// VaultOpsRequestSpec is the spec for VaultOpsRequest
type VaultOpsRequestSpec struct {
	// Specifies the Vault reference
	VaultRef core.LocalObjectReference `json:"vaultRef"`

	// Specifies the ops request type: ReconfigureTLS, Upgrade, etc.
	Type OpsRequestType `json:"type"`

	// Specifies information necessary for configuring TLS
	TLS *TLSSpec `json:"tls,omitempty"`

	// Specifies information necessary for restarting VaultServer
	Restart *RestartSpec `json:"restart,omitempty"`

	// Timeout for each step of the ops request in second. If a step doesn't finish within the specified timeout, the ops request will result in failure.
	Timeout *metav1.Duration `json:"timeout,omitempty"`
}

type RestartSpec struct{}

type TLSSpec struct {
	// TLSConfig contains updated tls configurations for client and server.
	// +optional
	kmapi.TLSConfig `json:",inline,omitempty"`

	// RotateCertificates tells operator to initiate certificate rotation
	// +optional
	RotateCertificates bool `json:"rotateCertificates,omitempty"`

	// Remove tells operator to remove TLS configuration
	// +optional
	Remove bool `json:"remove,omitempty"`
}

// VaultOpsRequestStatus is the status for VaultOpsRequest
type VaultOpsRequestStatus struct {
	// Specifies the current phase of the ops request
	// +optional
	Phase OpsRequestPhase `json:"phase,omitempty"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Conditions applied to the request, such as approval or denial.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VaultOpsRequestList is a list of VaultOpsRequests
type VaultOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of VaultOpsRequest CRD objects
	Items []VaultOpsRequest `json:"items,omitempty"`
}

// +kubebuilder:validation:Enum=Pending;Progressing;Successful;WaitingForApproval;Failed;Approved;Denied
type OpsRequestPhase string

const (
	// used for ops requests that are currently in queue
	OpsRequestPhasePending OpsRequestPhase = "Pending"
	// used for ops requests that are currently Progressing
	OpsRequestPhaseProgressing OpsRequestPhase = "Progressing"
	// used for ops requests that are executed successfully
	OpsRequestPhaseSuccessful OpsRequestPhase = "Successful"
	// used for ops requests that are waiting for approval
	OpsRequestPhaseWaitingForApproval OpsRequestPhase = "WaitingForApproval"
	// used for ops requests that are failed
	OpsRequestPhaseFailed OpsRequestPhase = "Failed"
	// used for ops requests that are approved
	OpsRequestApproved OpsRequestPhase = "Approved"
	// used for ops requests that are denied
	OpsRequestDenied OpsRequestPhase = "Denied"
)

// +kubebuilder:validation:Enum=ReconfigureTLS;Restart
type OpsRequestType string

const (
	// used for Restart operation
	OpsRequestTypeRestart OpsRequestType = "Restart"

	// used for ReconfigureTLS operation
	OpsRequestTypeReconfigureTLSs OpsRequestType = "ReconfigureTLS"
)
