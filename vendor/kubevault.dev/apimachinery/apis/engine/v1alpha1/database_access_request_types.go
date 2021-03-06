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
	ResourceKindDatabaseAccessRequest = "DatabaseAccessRequest"
	ResourceDatabaseAccessRequest     = "databaseaccessrequest"
	ResourceDatabaseAccessRequests    = "databaseaccessrequests"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=databaseaccessrequests,singular=databaseaccessrequest,categories={vault,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type DatabaseAccessRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              DatabaseAccessRequestSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status            DatabaseAccessRequestStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// DatabaseAccessRequestSpec contains information to request for database credential
type DatabaseAccessRequestSpec struct {
	// Contains vault database role info
	RoleRef RoleRef `json:"roleRef" protobuf:"bytes,1,opt,name=roleRef"`

	Subjects []rbac.Subject `json:"subjects" protobuf:"bytes,2,rep,name=subjects"`

	// Specifies the TTL for the leases associated with this role.
	// Accepts time suffixed strings ("1h") or an integer number of seconds.
	// Defaults to roles default TTL time
	TTL string `json:"ttl,omitempty" protobuf:"bytes,3,opt,name=ttl"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type DatabaseAccessRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is a list of DatabaseAccessRequest objects
	Items []DatabaseAccessRequest `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}

type DatabaseAccessRequestStatus struct {
	// Specifies the phase of DatabaseAccessRequest object
	Phase RequestStatusPhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=RequestStatusPhase"`

	// Conditions applied to the request, such as approval or denial.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty" protobuf:"bytes,2,rep,name=conditions"`

	// Name of the secret containing database credentials
	Secret *core.LocalObjectReference `json:"secret,omitempty" protobuf:"bytes,3,opt,name=secret"`

	// Contains lease info
	Lease *Lease `json:"lease,omitempty" protobuf:"bytes,4,opt,name=lease"`

	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,5,opt,name=observedGeneration"`
}
