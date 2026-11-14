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

// vault doc: https://www.vaultproject.io/docs/configuration/storage/azure.html
//
// AzureSpec defines configuration to set up Google Cloud Storage as backend storage in vault
type AzureSpec struct {
	// Specifies the Azure Storage account name.
	AccountName string `json:"accountName"`

	// Specifies the secret containing Azure Storage account key.
	// secret data:
	//  - account_key:<value>
	// +optional
	CredentialSecretRef *core.LocalObjectReference `json:"credentialSecretRef,omitempty"`

	// Specifies the Azure Storage Blob container name.
	Container string `json:"container"`

	//  Specifies the maximum number of concurrent operations to take place.
	// +optional
	MaxParallel int64 `json:"maxParallel,omitempty"`
}
