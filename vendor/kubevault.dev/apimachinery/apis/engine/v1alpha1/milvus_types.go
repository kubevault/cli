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
	ResourceKindMilvusRole = "MilvusRole"
	ResourceMilvusRole     = "milvusrole"
	ResourceMilvusRoles    = "milvusroles"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=milvusroles,singular=milvusrole,categories={vault,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MilvusRole struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MilvusRoleSpec `json:"spec,omitempty"`
	Status            RoleStatus     `json:"status,omitempty"`
}

// MilvusRoleSpec describes a dynamic-role binding against the Milvus
// database secret engine. The OpenBao `milvus-database-plugin`
// (sigilr/openbao#13) provisions credentials via Milvus's HTTP RESTful
// API v2 user-management endpoints (Create User + Grant Role), so
// `creationStatements` is a single-element JSON role document of the
// form `{"roles":["role1","role2"]}`. The referenced roles must already
// exist on the target Milvus cluster.
type MilvusRoleSpec struct {
	// SecretEngineRef is the name of a Secret Engine
	SecretEngineRef core.LocalObjectReference `json:"secretEngineRef"`

	// CreationStatements is a JSON role document listing pre-existing Milvus
	// roles to grant. Example: ['{"roles":["dba","readonly"]}']
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

type MilvusRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is a list of MilvusRole objects
	Items []MilvusRole `json:"items,omitempty"`
}

// MilvusConfiguration defines a Milvus app configuration. The OpenBao
// `milvus-database-plugin` (sigilr/openbao#13) provisions credentials
// via Milvus's HTTP RESTful API v2 user-management endpoints (Create
// User + Grant Role), so the connection payload uses `url` (Milvus
// HTTP endpoint) rather than a connection_url. Authentication can be
// HTTP Basic Auth (username/password from the AppBinding secret) OR a
// Bearer token from a `token` key on the secret (Zilliz Cloud style).
// Milvus is dynamic: NewUser/UpdateUser/DeleteUser are all supported.
// https://milvus.io/docs/users_and_roles.md
type MilvusConfiguration struct {
	// Specifies the Milvus database appbinding reference. The AppBinding
	// URL points at the Milvus HTTP endpoint (e.g.
	// `http://milvus.demo.svc:19530`); the secret contributes either
	// username/password (HTTP Basic) or a `token` key (Bearer / Zilliz
	// Cloud) used to authenticate against Milvus's user-management API.
	DatabaseRef appcat.AppReference `json:"databaseRef"`

	// Specifies the name of the plugin to use for this connection.
	// Default plugin:
	//  - for milvus: milvus-database-plugin
	// +optional
	PluginName string `json:"pluginName,omitempty"`

	// List of the roles allowed to use this connection.
	// Defaults to empty (no roles), if contains a "*" any role can use this connection.
	// +optional
	AllowedRoles []string `json:"allowedRoles,omitempty"`

	// DBName is forwarded as the dbName request header on every API call.
	// +optional
	DBName string `json:"dbName,omitempty"`
}
