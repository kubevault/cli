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
)

// Lease contains lease info
type Lease struct {
	// lease id
	ID string `json:"id,omitempty" protobuf:"bytes,1,opt,name=id"`

	// lease duration
	Duration metav1.Duration `json:"duration,omitempty" protobuf:"bytes,2,opt,name=duration"`

	// Specifies whether this lease is renewable
	Renewable bool `json:"renewable,omitempty" protobuf:"varint,3,opt,name=renewable"`
}

// Specifies the access key request phase
type RequestStatusPhase string

var (
	RequestStatusPhaseWaitingForApproval RequestStatusPhase = "WaitingForApproval"
	RequestStatusPhaseApproved           RequestStatusPhase = "Approved"
	RequestStatusPhaseDenied             RequestStatusPhase = "Denied"
)

type RolePhase string

const (
	// RolePhase constants
	RolePhaseSuccess    RolePhase = "Success"
	RolePhaseProcessing RolePhase = "Processing"
)

type RoleStatus struct {
	Phase RolePhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=RolePhase"`

	// ObservedGeneration is the most recent generation observed for this MySQLRole. It corresponds to the
	// MySQLRole's generation, which is updated on mutation by the API Server.
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,2,opt,name=observedGeneration"`

	// Represents the latest available observations of a MySQLRole current state.
	Conditions []kmapi.Condition `json:"conditions,omitempty" protobuf:"bytes,3,rep,name=conditions"`

	PolicyRef *kmapi.ObjectReference `json:"policyRef,omitempty" protobuf:"bytes,4,opt,name=policyRef"`
}

const (
	SecretRoleBindingAnnotationName      = "secretrolebindings.engine.kubevault.com/name"
	SecretRoleBindingAnnotationNamespace = "secretrolebindings.engine.kubevault.com/namespace"
)

// SecretRoleBinding Phases

const (
	SecretRoleBindingPhaseSuccess    RequestStatusPhase = "Success"
	SecretRoleBindingPhaseProcessing RequestStatusPhase = "Processing"
	SecretRoleBindingPhaseFailed     RequestStatusPhase = "Failed"
)

// SecretRoleBinding Conditions

const (
	VaultPolicySuccess        = "VaultPolicySuccess"
	VaultPolicyBindingSuccess = "VaultPolicyBindingSuccess"
	SecretRoleBindingSuccess  = "SecretRoleBindingSuccess"
)
