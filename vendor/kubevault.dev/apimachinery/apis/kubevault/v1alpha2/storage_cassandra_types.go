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

// vault doc: https://openbao.org/docs/configuration/storage/cassandra/
//
// CassandraSpec defines configuration to set up Cassandra as backend storage in vault
type CassandraSpec struct {
	// Specifies a comma-separated list of Cassandra hosts to connect to. All hosts must
	// listen on the same port; include the port in each host as "<host>:<port>" if it is
	// not the CQL native protocol default.
	// +optional
	Hosts string `json:"hosts,omitempty"`

	// Specifies the keyspace that is used for storing the Vault data. The keyspace must
	// already exist, be reachable, and writable.
	// +optional
	// +kubebuilder:default:="vault"
	Keyspace string `json:"keyspace,omitempty"`

	// Specifies the table in the keyspace that is used for storing the Vault data. The
	// table must already exist.
	// +optional
	// +kubebuilder:default:="entries"
	Table string `json:"table,omitempty"`

	// Specifies the consistency level for read and write operations. Must be one of ANY,
	// ONE, TWO, THREE, QUORUM, ALL, LOCAL_QUORUM, EACH_QUORUM, or LOCAL_ONE.
	// +optional
	Consistency string `json:"consistency,omitempty"`

	// Specifies the CQL protocol version to use. Set to "3" or higher to use
	// username/password authentication with a proto version that requires it.
	// +optional
	ProtocolVersion string `json:"protocolVersion,omitempty"`

	// Specifies the secret name containing username and password to use for
	// authentication (PasswordAuthenticator). Requires protocolVersion of "2" or higher.
	// secret data:
	//  - username=<value>
	//  - password=<value>
	// +optional
	CredentialSecretRef *core.LocalObjectReference `json:"credentialSecretRef,omitempty"`

	// When set to a positive integer, enables Cassandra's SimpleRetryPolicy with the
	// given number of retries for queries that time out or fail.
	// +optional
	SimpleRetryPolicyRetries string `json:"simpleRetryPolicyRetries,omitempty"`

	// Specifies the timeout, in seconds, for the initial connection to the cluster.
	// +optional
	InitialConnectionTimeout string `json:"initialConnectionTimeout,omitempty"`

	// Specifies the timeout, in seconds, for individual queries against the cluster.
	// +optional
	ConnectionTimeout string `json:"connectionTimeout,omitempty"`

	// Set to enable a TLS connection to Cassandra.
	// +optional
	TLSEnabled bool `json:"tlsEnabled,omitempty"`

	// Specifies the secret name that contains a PEM-encoded certificate bundle (and,
	// optionally, a private key) used for the TLS connection to Cassandra.
	// secret data:
	//  - pem_bundle=<value>
	// +optional
	TLSSecretRef *core.LocalObjectReference `json:"tlsSecretRef,omitempty"`

	// Disables verification of the Cassandra server's certificate chain and host name.
	// Not recommended for production use.
	// +optional
	TLSSkipVerify bool `json:"tlsSkipVerify,omitempty"`

	// Specifies the minimum acceptable TLS version. One of tls10, tls11, tls12, or tls13.
	// +optional
	TLSMinVersion string `json:"tlsMinVersion,omitempty"`

	// Specifies the name to verify the certificate of every Cassandra node against. By
	// default each node is checked against the host name it was reached by, which, for
	// nodes discovered from the cluster, is their IP address. Set it when the nodes share
	// a certificate that lists DNS names but not the node IPs. If DatabaseRef is set and
	// the TLS material comes from the AppBinding, it defaults to the AppBinding's host.
	// +optional
	TLSServerName string `json:"tlsServerName,omitempty"`

	// DatabaseRef contains the info of a KubeDB managed Cassandra.
	// If set, Hosts will be generated from the referenced AppBinding.
	// +optional
	DatabaseRef *appcat.AppReference `json:"databaseRef,omitempty"`
}
