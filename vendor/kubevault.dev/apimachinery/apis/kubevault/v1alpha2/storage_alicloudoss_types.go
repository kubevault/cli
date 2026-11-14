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

// vault doc: https://openbao.org/docs/configuration/storage/alicloudoss/
//
// AlicloudOSSSpec defines configuration to set up Alicloud OSS as backend storage in vault
type AlicloudOSSSpec struct {
	// Specifies the OSS endpoint to connect to, e.g. "http://oss-us-east-1.aliyuncs.com".
	Endpoint string `json:"endpoint"`

	// Specifies the name of the OSS bucket to store data in. It must already exist;
	// Vault does not create it.
	Bucket string `json:"bucket"`

	// Specifies the secret name containing the Alicloud access key ID and secret to
	// connect with. Required: Vault does not start without them.
	// secret data:
	//  - access_key=<value>
	//  - secret_key=<value>
	CredentialSecretRef *core.LocalObjectReference `json:"credentialSecretRef"`

	// Specifies the maximum number of concurrent requests to OSS.
	// +optional
	MaxParallel int64 `json:"maxParallel,omitempty"`
}
