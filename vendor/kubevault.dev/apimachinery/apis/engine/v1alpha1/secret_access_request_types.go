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
	ResourceKindSecretAccessRequest = "SecretAccessRequest"
	ResourceSecretAccessRequest     = "secretaccessrequest"
	ResourceSecretAccessRequests    = "secretaccessrequests"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=secretaccessrequests,singular=secretaccessrequest,categories={vault,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type SecretAccessRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              SecretAccessRequestSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status            SecretAccessRequestStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// SecretAccessRequestSpec contains information to request for database credential
type SecretAccessRequestSpec struct {
	// Contains vault database role info
	RoleRef core.TypedLocalObjectReference `json:"roleRef" protobuf:"bytes,1,opt,name=roleRef"`

	Subjects []rbac.Subject `json:"subjects" protobuf:"bytes,2,rep,name=subjects"`

	// Specifies the TTL for the leases associated with this role.
	// Accepts time suffixed strings ("1h") or an integer number of seconds.
	// Defaults to roles default TTL time
	TTL string `json:"ttl,omitempty" protobuf:"bytes,3,opt,name=ttl"`

	SecretAccessRequestConfiguration `json:",inline" protobuf:"bytes,4,opt,name=secretAccessRequestConfiguration"`
}

// SecretAccessRequestConfiguration contains information to request for database credential
type SecretAccessRequestConfiguration struct {
	// +optional
	AWS *AWSAccessRequestConfiguration `json:"aws,omitempty" protobuf:"bytes,1,opt,name=aws"`
	GCP *GCPAccessRequestConfiguration `json:"gcp,omitempty" protobuf:"bytes,2,opt,name=gcp"`
}

// https://www.vaultproject.io/api/secret/aws/index.html#parameters-6
// AWSAccessKeyRequestSpec contains information to request for vault aws credential
type AWSAccessRequestConfiguration struct {
	// The ARN of the role to assume if credential_type on the Vault role is assumed_role.
	// Must match one of the allowed role ARNs in the Vault role. Optional if the Vault role
	// only allows a single AWS role ARN; required otherwise.
	RoleARN string `json:"roleARN,omitempty" protobuf:"bytes,1,opt,name=roleARN"`

	// If true, '/aws/sts' endpoint will be used to retrieve credential
	// Otherwise, '/aws/creds' endpoint will be used to retrieve credential
	UseSTS bool `json:"useSTS,omitempty" protobuf:"varint,2,opt,name=useSTS"`
}

// Link:
//  - https://www.vaultproject.io/api/secret/gcp/index.html#generate-secret-iam-service-account-creds-oauth2-access-token
//  - https://www.vaultproject.io/api/secret/gcp/index.html#generate-secret-iam-service-account-creds-service-account-key

// GCPAccessRequestConfiguration contains information to request for vault gcp credentials
type GCPAccessRequestConfiguration struct {
	// Specifies the algorithm used to generate key.
	// Defaults to 2k RSA key.
	// Accepted values: KEY_ALG_UNSPECIFIED, KEY_ALG_RSA_1024, KEY_ALG_RSA_2048
	// +optional
	KeyAlgorithm string `json:"keyAlgorithm,omitempty" protobuf:"bytes,1,opt,name=keyAlgorithm"`

	// Specifies the private key type to generate.
	// Defaults to JSON credentials file
	// Accepted values: TYPE_UNSPECIFIED, TYPE_PKCS12_FILE, TYPE_GOOGLE_CREDENTIALS_FILE
	// +optional
	KeyType string `json:"keyType,omitempty" protobuf:"bytes,2,opt,name=keyType"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SecretAccessRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is a list of SecretAccessRequest objects
	Items []SecretAccessRequest `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}

type SecretAccessRequestStatus struct {
	// Specifies the phase of SecretAccessRequest object
	Phase RequestStatusPhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=RequestStatusPhase"`

	// Conditions applied to the request, such as approval or denial.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty" protobuf:"bytes,2,rep,name=conditions"`

	// Contains lease info
	Lease *Lease `json:"lease,omitempty" protobuf:"bytes,3,opt,name=lease"`

	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,4,opt,name=observedGeneration"`

	// Name of the secret containing secret engine role credentials
	Secret *kmapi.ObjectReference `json:"secret,omitempty" protobuf:"bytes,5,opt,name=secret"`
}
