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
)

// vault doc: https://www.vaultproject.io/docs/configuration/storage/swift.html
//
// SwiftSpec defines configuration to set up Swift Storage as backend storage in vault
type SwiftSpec struct {
	// Specifies the OpenStack authentication endpoint.
	AuthURL string `json:"authURL"`

	// Specifies the name of the Swift container.
	Container string `json:"container"`

	// Specifies the name of the secret containing the OpenStack account/username and password
	// Specifies secret containing auth token from alternate authentication.
	// secret data:
	//  - username=<value>
	//  - password=<value>
	//  - auth_token=<value>
	// +optional
	CredentialSecretRef *core.LocalObjectReference `json:"credentialSecretRef,omitempty"`

	// Specifies the name of the tenant. If left blank, this will default to the default tenant of the username.
	// +optional
	Tenant string `json:"tenant,omitempty"`

	// Specifies the name of the region.
	// +optional
	Region string `json:"region,omitempty"`

	// Specifies the id of the tenant.
	// +optional
	TenantID string `json:"tenantID,omitempty"`

	// Specifies the name of the user domain.
	// +optional
	Domain string `json:"domain,omitempty"`

	// Specifies the name of the project's domain.
	// +optional
	ProjectDomain string `json:"projectDomain,omitempty"`

	// Specifies the id of the trust.
	// +optional
	TrustID string `json:"trustID,omitempty"`

	// Specifies storage URL from alternate authentication.
	// +optional
	StorageURL string `json:"storageURL,omitempty"`

	//  Specifies the maximum number of concurrent requests to take place.
	// +optional
	MaxParallel int64 `json:"maxParallel,omitempty"`
}
