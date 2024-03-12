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
)

const (
	ResourceKindPKIRole = "PKIRole"
	ResourcePKIRole     = "pkirole"
	ResourcePKIRoles    = "pkiroles"
)

// PKIRole

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=pkiroles,singular=pkirole,categories={vault,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type PKIRole struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PKIRoleSpec `json:"spec,omitempty"`
	Status            RoleStatus  `json:"status,omitempty"`
}

// PKIRoleSpec contains connection information, PKI role info, etc
// More info: https://developer.hashicorp.com/vault/api-docs/secret/pki#create-update-role
type PKIRoleSpec struct {
	// SecretEngineRef is the name of a Secret Engine
	SecretEngineRef   core.LocalObjectReference `json:"secretEngineRef"`
	IssuerRef         string                    `json:"issuerRef,omitempty"`
	TTL               string                    `json:"ttl,omitempty"`
	MaxTTL            string                    `json:"maxTTL,omitempty"`
	AllowedDomains    []string                  `json:"allowedDomains,omitempty"`
	AllowSubdomains   *bool                     `json:"allowSubdomains,omitempty"`
	OU                string                    `json:"ou,omitempty"`
	Organization      string                    `json:"organization,omitempty"`
	Country           string                    `json:"country,omitempty"`
	AdditionalPayload map[string]string         `json:"additionalPayload,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

type PKIRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is a list of PKIRole objects
	Items []PKIRole `json:"items,omitempty"`
}
