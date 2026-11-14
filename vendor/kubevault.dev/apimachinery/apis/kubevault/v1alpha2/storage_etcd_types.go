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

// TODO : set defaults and validation
// vault doc: https://www.vaultproject.io/docs/configuration/storage/etcd.html
//
// EtcdSpec defines configuration to set up etcd as backend storage in vault
type EtcdSpec struct {
	// Specifies the addresses of the etcd instances
	// +optional
	Address string `json:"address,omitempty"`

	// Specifies the version of the API to communicate with etcd
	// +optional
	EtcdApi string `json:"etcdApi,omitempty"`

	// Specifies if high availability should be enabled
	// +optional
	HAEnabled bool `json:"haEnabled,omitempty"`

	// Specifies the path in etcd where vault data will be stored
	// +optional
	Path string `json:"path,omitempty"`

	// Specifies whether to sync list of available etcd services on startup
	// +optional
	Sync bool `json:"sync,omitempty"`

	// Specifies the domain name to query for SRV records describing cluster endpoints
	// +optional
	DiscoverySrv string `json:"discoverySrv,omitempty"`

	// Specifies the secret name that contain username and password to use when authenticating with the etcd server
	// secret data:
	//  - username:<value>
	//  - password:<value>
	// +optional
	CredentialSecretRef *core.LocalObjectReference `json:"credentialSecretRef,omitempty"`

	// Specifies the secret name that contains tls_ca_file, tls_cert_file and tls_key_file for etcd communication
	// secret data:
	//  - ca.crt
	//  - tls.crt
	//  - tls.key
	// +optional
	TLSSecretRef *core.LocalObjectReference `json:"tlsSecretRef,omitempty"`

	// DatabaseRef contains the info of a KubeDB managed Etcd.
	// If set, Address will be generated from the referenced AppBinding.
	// +optional
	DatabaseRef *appcat.AppReference `json:"databaseRef,omitempty"`
}
