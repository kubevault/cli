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
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
)

const (
	ResourceKindNeo4jRole = "Neo4jRole"
	ResourceNeo4jRole     = "neo4jrole"
	ResourceNeo4jRoles    = "neo4jroles"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=neo4jroles,singular=neo4jrole,categories={vault,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Neo4jRole struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              Neo4jRoleSpec `json:"spec,omitempty"`
	Status            RoleStatus    `json:"status,omitempty"`
}

// Neo4jRoleSpec describes a dynamic-role binding against the Neo4j
// database secret engine. The OpenBao `neo4j-database-plugin`
// (sigilr/openbao#10) provisions credentials as native Neo4j users
// created with Cypher's `CREATE USER` against the `system` database
// (via the neo4j-go-driver/v5), so `creationStatements` is a
// single-element JSON role document of the form
// `{"roles":["role1","role2"]}`. The referenced roles must already
// exist on the target Neo4j cluster.
type Neo4jRoleSpec struct {
	// SecretEngineRef is the name of a Secret Engine
	SecretEngineRef core.LocalObjectReference `json:"secretEngineRef"`

	// CreationStatements is a JSON role document listing pre-existing Neo4j
	// roles to grant. Example: ['{"roles":["reader","publisher"]}']
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

type Neo4jRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is a list of Neo4jRole objects
	Items []Neo4jRole `json:"items,omitempty"`
}

// Neo4jConfiguration defines a Neo4j app configuration. The OpenBao
// `neo4j-database-plugin` (sigilr/openbao#10) provisions credentials
// as native Neo4j users created with Cypher's `CREATE USER` against
// the `system` database, so the connection payload uses `uri` (a Bolt
// URI such as `bolt://host:7687` or `neo4j://host:7687`) rather than
// `connection_url`. Authentication is HTTP Basic Auth (username +
// password from the AppBinding secret). Neo4j is dynamic:
// NewUser/UpdateUser/DeleteUser are all supported. Revocation is the
// plugin's default `DROP USER ... IF EXISTS` so no
// `revocation_statements` field is exposed here.
// https://neo4j.com/docs/operations-manual/current/authentication-authorization/
type Neo4jConfiguration struct {
	// Specifies the Neo4j database appbinding reference. The AppBinding
	// URL points at the Neo4j Bolt endpoint (e.g.
	// `bolt://neo4j.demo.svc:7687` or `neo4j://neo4j.demo.svc:7687`);
	// the secret contributes username/password used to authenticate
	// against Neo4j when the plugin runs `CREATE USER` Cypher
	// statements on the `system` database.
	DatabaseRef appcat.AppReference `json:"databaseRef"`

	// Specifies the name of the plugin to use for this connection.
	// Default plugin:
	//  - for neo4j: neo4j-database-plugin
	// +optional
	PluginName string `json:"pluginName,omitempty"`

	// List of the roles allowed to use this connection.
	// Defaults to empty (no roles), if contains a "*" any role can use this connection.
	// +optional
	AllowedRoles []string `json:"allowedRoles,omitempty"`
}
