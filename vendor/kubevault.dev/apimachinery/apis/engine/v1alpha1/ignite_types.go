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
	ResourceKindIgniteRole = "IgniteRole"
	ResourceIgniteRole     = "igniterole"
	ResourceIgniteRoles    = "igniteroles"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=igniteroles,singular=igniterole,categories={vault,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type IgniteRole struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              IgniteRoleSpec `json:"spec,omitempty"`
	Status            RoleStatus     `json:"status,omitempty"`
}

// IgniteRoleSpec describes a dynamic-role binding against the Apache Ignite
// database secret engine. The OpenBao `ignite-database-plugin`
// (sigilr/openbao#14) drives Ignite's REST API (`cmd=qryfldexe`) to execute
// `CREATE USER` / `ALTER USER` / `DROP USER` SQL statements: each entry in
// `creationStatements` is an Ignite SQL DDL statement that is executed in
// sequence to mint a user.
type IgniteRoleSpec struct {
	// SecretEngineRef is the name of a Secret Engine
	SecretEngineRef core.LocalObjectReference `json:"secretEngineRef"`

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

	// CreationStatements is the list of SQL statements run against Ignite's
	// REST API to create a dynamic user. Placeholders: {{name}}, {{password}}.
	// Example: ["CREATE USER \"{{name}}\" WITH PASSWORD '{{password}}';"]
	CreationStatements []string `json:"creationStatements"`

	// Specifies the database statements to be executed to revoke a user.
	// +optional
	RevocationStatements []string `json:"revocationStatements,omitempty"`

	// Specifies the database statements to be executed rollback a create operation in the event of an error.
	// +optional
	RollbackStatements []string `json:"rollbackStatements,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type IgniteRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is a list of IgniteRole objects
	Items []IgniteRole `json:"items,omitempty"`
}
