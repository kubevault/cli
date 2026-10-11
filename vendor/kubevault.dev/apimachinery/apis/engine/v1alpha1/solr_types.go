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
	ResourceKindSolrRole = "SolrRole"
	ResourceSolrRole     = "solrrole"
	ResourceSolrRoles    = "solrroles"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=solrroles,singular=solrrole,categories={vault,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type SolrRole struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SolrRoleSpec `json:"spec,omitempty"`
	Status            RoleStatus   `json:"status,omitempty"`
}

// SolrRoleSpec describes a dynamic-role binding against the Apache
// Solr database secret engine. The OpenBao `solr-database-plugin`
// (sigilr/openbao#11) provisions credentials via Solr's Security
// Plugin API (Basic Auth Plugin user + Rule-Based Authorization role
// bindings), so `creationStatements` is a single-element JSON role
// document of the form `{"roles":["admin","read"]}`. The referenced
// roles must already exist on the configured Solr authorizer.
type SolrRoleSpec struct {
	// SecretEngineRef is the name of a Secret Engine
	SecretEngineRef core.LocalObjectReference `json:"secretEngineRef"`

	// CreationStatements is a JSON role document listing pre-existing
	// Solr roles to bind via the `set-user-role` API, of the form
	// `{"roles":["admin","read"]}`. Roles must already exist on the
	// configured Solr authorizer.
	CreationStatements []string `json:"creationStatements"`

	// Specifies the TTL for the leases associated with this role.
	// Accepts time suffixed strings ("1h") or an integer number of seconds.
	// Defaults to system/engine default TTL time
	// +optional
	DefaultTTL string `json:"defaultTTL,omitempty"`

	// Specifies the maximum TTL for the leases associated with this role.
	// Accepts time suffixed strings ("1h") or an integer number of seconds.
	// Defaults to system/engine default TTL time.
	// +optional
	MaxTTL string `json:"maxTTL,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SolrRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is a list of SolrRole objects
	Items []SolrRole `json:"items,omitempty"`
}
