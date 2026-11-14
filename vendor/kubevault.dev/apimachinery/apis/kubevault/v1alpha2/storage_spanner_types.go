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

// vault doc: https://openbao.org/docs/configuration/storage/google-cloud-spanner/
//
// SpannerSpec defines configuration to set up Google Cloud Spanner as backend storage in vault
type SpannerSpec struct {
	// Specifies the full name of the Spanner database, in the form
	// projects/<project>/instances/<instance>/databases/<database>.
	Database string `json:"database"`

	// Specifies the name of the table to use for Vault data.
	// +optional
	// +kubebuilder:default:="Vault"
	Table string `json:"table,omitempty"`

	// High Availability Parameter
	// Specifies if high availability mode is enabled.
	// +optional
	HAEnabled bool `json:"haEnabled,omitempty"`

	// Specifies the name of the table used for HA leader election. Defaults to the
	// value of Table suffixed with "HA".
	// +optional
	HATable string `json:"haTable,omitempty"`

	// Specifies the maximum number of concurrent requests to Spanner.
	// +optional
	MaxParallel int64 `json:"maxParallel,omitempty"`

	// Secret containing the Google application credential used to reach Spanner.
	// secret data:
	//  - sa.json:<value>
	// +optional
	CredentialSecretRef *core.LocalObjectReference `json:"credentialSecretRef,omitempty"`
}
