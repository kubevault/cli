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

// ref: https://www.vaultproject.io/docs/configuration/storage/consul.html
//
// ConsulSpec defines the configuration to set up consul as backend storage in vault
type ConsulSpec struct {
	// Specifies the address of the Consul agent to communicate with.
	// This can be an IP address, DNS record, or unix socket.
	// +optional
	Address string `json:"address,omitempty"`

	// Specifies the check interval used to send health check information
	// back to Consul.
	// This is specified using a label suffix like "30s" or "1h".
	// +optional
	CheckTimeout string `json:"checkTimeout,omitempty"`

	// Specifies the Consul consistency mode.
	// Possible values are "default" or "strong".
	// +optional
	ConsistencyMode string `json:"consistencyMode,omitempty"`

	// Specifies whether Vault should register itself with Consul.
	// Possible values are "true" or "false"
	// +optional
	DisableRegistration string `json:"disableRegistration,omitempty"`

	// Specifies the maximum number of concurrent requests to Consul.
	// +optional
	MaxParallel string `json:"maxParallel,omitempty"`

	// Specifies the path in Consul's key-value store
	// where Vault data will be stored.
	// +optional
	Path string `json:"path,omitempty"`

	// Specifies the scheme to use when communicating with Consul.
	// This can be set to "http" or "https".
	// +optional
	Scheme string `json:"scheme,omitempty"`

	// Specifies the name of the service to register in Consul.
	// +optional
	Service string `json:"service,omitempty"`

	// Specifies a comma-separated list of tags
	// to attach to the service registration in Consul.
	// +optional
	ServiceTags string `json:"serviceTags,omitempty"`

	// Specifies a service-specific address to set on the service registration
	// in Consul.
	// If unset, Vault will use what it knows to be the HA redirect address
	// - which is usually desirable.
	// Setting this parameter to "" will tell Consul to leverage the configuration
	// of the node the service is registered on dynamically.
	// +optional
	ServiceAddress string `json:"serviceAddress,omitempty"`

	// Specifies the secret name that contains ACL token with permission
	// to read and write from the path in Consul's key-value store.
	// secret data:
	//  - aclToken:<value>
	// +optional
	ACLTokenSecretRef *core.LocalObjectReference `json:"aclTokenSecretRef,omitempty"`

	// Specifies the minimum allowed session TTL.
	// Consul server has a lower limit of 10s on the session TTL by default.
	// +optional
	SessionTTL string `json:"sessionTTL,omitempty"`

	// Specifies the wait time before a lock lock acquisition is made.
	// This affects the minimum time it takes to cancel a lock acquisition.
	// +optional
	LockWaitTime string `json:"lockWaitTime,omitempty"`

	// Specifies the secret name that contains tls_ca_file, tls_cert_file and tls_key_file
	// for consul communication
	// Secret data:
	//  - ca.crt
	//  - tls.crt
	//  - tls.key
	// +optional
	TLSSecretRef *core.LocalObjectReference `json:"tlsSecretRef,omitempty"`

	// Specifies the minimum TLS version to use.
	// Accepted values are "tls10", "tls11" or "tls12".
	// +optional
	TLSMinVersion string `json:"tlsMinVersion,omitempty"`

	// Specifies if the TLS host verification should be disabled.
	// It is highly discouraged that you disable this option.
	// +optional
	TLSSkipVerify bool `json:"tlsSkipVerify,omitempty"`
}
