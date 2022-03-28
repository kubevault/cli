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
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
)

const (
	ResourceKindSecretRoleBinding = "SecretRoleBinding"
	ResourceSecretRoleBinding     = "secretrolebinding"
	ResourceSecretRoleBindings    = "secretrolebindings"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=secretrolebindings,singular=secretrolebinding,categories={vault,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type SecretRoleBinding struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SecretRoleBindingSpec   `json:"spec,omitempty"`
	Status            SecretRoleBindingStatus `json:"status,omitempty"`
}

// SecretRoleBindingSpec contains information to request for database credential
type SecretRoleBindingSpec struct {
	Roles []core.TypedLocalObjectReference `json:"roles"`

	Subjects []rbac.Subject `json:"subjects"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SecretRoleBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is a list of SecretRoleBinding objects
	Items []SecretRoleBinding `json:"items,omitempty"`
}

type SecretRoleBindingStatus struct {
	// Specifies the phase of SecretRoleBinding object
	Phase RequestStatusPhase `json:"phase,omitempty"`

	// Conditions applied to the request, such as approval or denial.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`

	// Contains lease info
	Lease *Lease `json:"lease,omitempty"`

	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	PolicyRef *kmapi.ObjectReference `json:"policyRef,omitempty"`

	PolicyBindingRef *kmapi.ObjectReference `json:"policyBindingRef,omitempty"`
}
