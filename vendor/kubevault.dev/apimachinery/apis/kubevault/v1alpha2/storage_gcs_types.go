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

// vault doc: https://www.vaultproject.io/docs/configuration/storage/google-cloud-storage.html
//
// GcsSpec defines configuration to set up Google Cloud Storage as backend storage in vault
type GcsSpec struct {
	// Specifies the name of the bucket to use for storage.
	Bucket string `json:"bucket"`

	// Specifies the maximum size (in kilobytes) to send in a single request. If set to 0,
	// it will attempt to send the whole object at once, but will not retry any failures.
	// +optional
	ChunkSize string `json:"chunkSize,omitempty"`

	//  Specifies the maximum number of parallel operations to take place.
	// +optional
	MaxParallel int64 `json:"maxParallel,omitempty"`

	// Specifies if high availability mode is enabled.
	// +optional
	HAEnabled bool `json:"haEnabled,omitempty"`

	// Secret containing Google application credential
	// secret data:
	//  - sa.json:<value>
	// +optional
	CredentialSecretRef *core.LocalObjectReference `json:"credentialSecretRef,omitempty"`
}
