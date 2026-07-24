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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
)

const (
	ResourceKindNamespaceSlice = "NamespaceSlice"
	ResourceNamespaceSlice     = "namespaceslice"
	ResourceNamespaceSlices    = "namespaceslices"
)

// NamespaceSlice is a shard of the OpenBao namespaces a KubeVault operator needs
// provisioned, following the Kubernetes EndpointSlice pattern: a large set is split
// across multiple NamespaceSlice objects, each grouped back to the owning VaultServer
// through the kubevault.com/vaultserver-name + kubevault.com/vaultserver-namespace
// labels (cross-namespace owner references are not allowed, so labels associate them).
//
// In the hub-spoke model, the KubeVault operator on a managed (spoke) cluster records
// here the OpenBao org namespaces its client-org databases require on the hub — the
// deduplicated, validated set, sharded if large. The hub reads the slice(s) and
// idempotently creates each namespace; the spoke never creates hub namespaces itself
// (design/tenant-namespace-design.md §7.2).

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=namespaceslices,singular=namespaceslice,shortName=nss,categories={vault,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="NamespaceCount",type="integer",JSONPath=".status.namespaceCount"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type NamespaceSlice struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              NamespaceSliceSpec   `json:"spec,omitempty"`
	Status            NamespaceSliceStatus `json:"status,omitempty"`
}

type NamespaceSliceSpec struct {
	// HubVaultRef identifies the hub VaultServer whose namespaces this slice tracks —
	// the same VaultServer the kubevault.com/vaultserver-name + -namespace labels group
	// this slice to, so the ref and the labels always agree. It is set by the hub (which
	// owns those labels); the spoke fills only spec.namespaces. Optional for non-spoke uses.
	// +optional
	HubVaultRef kmapi.ObjectReference `json:"hubVaultRef,omitempty"`

	// Namespaces is this slice's shard of required OpenBao namespaces, keyed on name.
	// +optional
	// +listType=map
	// +listMapKey=name
	Namespaces []NamespaceSliceEntry `json:"namespaces,omitempty"`
}

// NamespaceSliceEntry is one required OpenBao namespace — the analogue of a single
// Endpoint in an EndpointSlice.
type NamespaceSliceEntry struct {
	// Name is the effective OpenBao namespace to provision (the org-id, in the current
	// tenant-isolation model where the namespace is keyed on the org identity).
	Name string `json:"name"`

	// ExternalID is the external identity this namespace maps to — the KubeDB Platform
	// org-id, in string form. It lets the effective namespace name diverge from the
	// org-id in the future without changing the association.
	// +optional
	ExternalID string `json:"externalID,omitempty"`

	// Conditions carries the readiness of this single namespace entry, mirroring the
	// EndpointConditions shape on EndpointSlice.
	// +optional
	Conditions NamespaceSliceConditions `json:"conditions,omitempty"`
}

// NamespaceSliceConditions is the readiness of one NamespaceSliceEntry.
type NamespaceSliceConditions struct {
	// Ready is true when the namespace is validated and required (as reported by the
	// authoring spoke operator). A nil value means the condition is unknown.
	// +optional
	Ready *bool `json:"ready,omitempty"`
}

type NamespaceSliceStatus struct {
	// ObservedGeneration is the most recent generation observed for this resource.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// NamespaceCount mirrors len(spec.namespaces) for a cheap print column and a
	// bounded status-feedback scalar.
	// +optional
	NamespaceCount int32 `json:"namespaceCount,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

type NamespaceSliceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NamespaceSlice `json:"items,omitempty"`
}
