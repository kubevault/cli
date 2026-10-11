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
	ResourceKindMSSQLServerRole = "MSSQLServerRole"
	ResourceMSSQLServerRole     = "mssqlserverrole"
	ResourceMSSQLServerRoles    = "mssqlserverroles"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=mssqlserverroles,singular=mssqlserverrole,categories={vault,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MSSQLServerRole struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MSSQLServerRoleSpec `json:"spec,omitempty"`
	Status            RoleStatus          `json:"status,omitempty"`
}

// MSSQLServerRoleSpec describes a dynamic-role binding against the
// Microsoft SQL Server database secret engine. The OpenBao
// `mssql-database-plugin` (sigilr/openbao#5) was ported from pre-BUSL
// HashiCorp Vault and is a SQL-statement based dynamic plugin
// (postgres-style): each entry in `creationStatements` is a T-SQL
// statement executed in sequence to mint a SQL Server login.
type MSSQLServerRoleSpec struct {
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

	// CreationStatements is the list of T-SQL statements run to create a
	// dynamic SQL Server login. Placeholders: {{name}}, {{password}}.
	CreationStatements []string `json:"creationStatements"`

	// Specifies the database statements to be executed to revoke a user.
	// +optional
	RevocationStatements []string `json:"revocationStatements,omitempty"`

	// Specifies the database statements to be executed rollback a create operation in the event of an error.
	// +optional
	RollbackStatements []string `json:"rollbackStatements,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type MSSQLServerRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is a list of MSSQLServerRole objects
	Items []MSSQLServerRole `json:"items,omitempty"`
}
