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

// vault doc: https://openbao.org/docs/configuration/storage/zookeeper/
//
// ZookeeperSpec defines configuration to set up ZooKeeper as backend storage in vault
type ZookeeperSpec struct {
	// Specifies the addresses of the ZooKeeper instances as a comma-separated list.
	// +optional
	// +kubebuilder:default:="localhost:2181"
	Address string `json:"address,omitempty"`

	// Specifies the path in ZooKeeper's tree where Vault data will be stored.
	// +optional
	// +kubebuilder:default:="vault/"
	Path string `json:"path,omitempty"`

	// Specifies the ACL scheme:id applied to every znode Vault creates. Defaults to
	// "world:anyone", i.e. unrestricted access.
	// +optional
	ZnodeOwner string `json:"znodeOwner,omitempty"`

	// Specifies the secret name containing a scheme:auth pair passed to ZooKeeper's
	// AddAuth API immediately after connecting, so the client authenticates as a
	// specific principal.
	// secret data:
	//  - authInfo=<scheme:auth>
	// +optional
	AuthInfoSecretRef *core.LocalObjectReference `json:"authInfoSecretRef,omitempty"`

	// Enables a TLS connection to ZooKeeper.
	// +optional
	TLSEnabled bool `json:"tlsEnabled,omitempty"`

	// Specifies the secret name that contains ca.crt, tls.crt and tls.key for ZooKeeper
	// communication.
	// secret data:
	//  - ca.crt
	//  - tls.crt
	//  - tls.key
	// +optional
	TLSSecretRef *core.LocalObjectReference `json:"tlsSecretRef,omitempty"`

	// Disables verification of the ZooKeeper server's certificate chain and host name.
	// Not recommended for production use.
	// +optional
	TLSSkipVerify bool `json:"tlsSkipVerify,omitempty"`

	// When set, verifies the server's IP address against the certificate instead of its
	// DNS name. Only consulted when TLSSkipVerify is false.
	// +optional
	TLSVerifyIP bool `json:"tlsVerifyIP,omitempty"`

	// Specifies the minimum acceptable TLS version. One of tls10, tls11, tls12, or tls13.
	// +optional
	TLSMinVersion string `json:"tlsMinVersion,omitempty"`

	// DatabaseRef contains the info of a KubeDB managed ZooKeeper.
	// If set, Address will be generated from the referenced AppBinding.
	// +optional
	DatabaseRef *appcat.AppReference `json:"databaseRef,omitempty"`
}
