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

// vault doc: https://www.vaultproject.io/docs/configuration/storage/s3.html
//
// S3Spec defines configuration to set up Amazon S3 Storage as backend storage in vault
type S3Spec struct {
	// Specifies the name of the bucket to use for storage.
	Bucket string `json:"bucket"`

	// Specifies an alternative, AWS compatible, S3 endpoint.
	// +optional
	Endpoint string `json:"endpoint,omitempty"`

	// Specifies the AWS region
	// +optional
	Region string `json:"region,omitempty"`

	// Specifies the secret name containing AWS session token, AWS access key and AWS secret key
	// secret data:
	//  - access_key=<value>
	//  - secret_key=<value>
	//  - session_token=<value>
	// +optional
	CredentialSecretRef *core.LocalObjectReference `json:"credentialSecretRef,omitempty"`

	// Specifies the maximum number of parallel operations to take place.
	// +optional
	MaxParallel int64 `json:"maxParallel,omitempty"`

	// Specifies whether to use host bucket style domains with the configured endpoint.
	// +optional
	ForcePathStyle bool `json:"forcePathStyle,omitempty"`

	// Specifies if SSL should be used for the endpoint connection
	// +optional
	DisableSSL bool `json:"disableSSL,omitempty"`
}
