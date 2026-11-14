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

package v1alpha2

import (
	core "k8s.io/api/core/v1"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
)

// vault doc: https://openbao.org/docs/configuration/storage/couchdb/
//
// CouchDBSpec defines configuration to set up CouchDB as backend storage in vault
type CouchDBSpec struct {
	// Specifies the full URL to the CouchDB database to use, including the database
	// name. The database must already exist; Vault does not create it.
	// +optional
	Endpoint string `json:"endpoint,omitempty"`

	// Specifies the secret name containing the CouchDB username and password to
	// connect with.
	// secret data:
	//  - username=<value>
	//  - password=<value>
	// +optional
	CredentialSecretRef *core.LocalObjectReference `json:"credentialSecretRef,omitempty"`

	// Specifies the maximum number of concurrent requests to CouchDB.
	// +optional
	MaxParallel int64 `json:"maxParallel,omitempty"`
	// Specifies the name of the CouchDB database to store data in. Only used with
	// DatabaseRef, where the endpoint is generated as <appbinding url>/<database>.
	// +optional
	// +kubebuilder:default:="vault"
	Database string `json:"database,omitempty"`

	// DatabaseRef contains the info of a CouchDB database exposed through an AppBinding.
	// If set, Endpoint will be generated from the referenced AppBinding.
	// +optional
	DatabaseRef *appcat.AppReference `json:"databaseRef,omitempty"`
}
