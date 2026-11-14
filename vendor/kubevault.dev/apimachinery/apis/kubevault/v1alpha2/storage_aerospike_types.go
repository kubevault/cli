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

// vault doc: https://openbao.org/docs/configuration/storage/aerospike/
//
// AerospikeSpec defines configuration to set up Aerospike as backend storage in vault
type AerospikeSpec struct {
	// Specifies the hostname of the Aerospike server to connect to, when HostList is
	// not set.
	// +optional
	// +kubebuilder:default:="127.0.0.1"
	Hostname string `json:"hostname,omitempty"`

	// Specifies the port of the Aerospike server to connect to, when HostList is not
	// set.
	// +optional
	// +kubebuilder:default:="3000"
	Port string `json:"port,omitempty"`

	// A comma-separated list of host[:port] entries describing the Aerospike cluster's
	// seed nodes. Takes precedence over Hostname/Port when set.
	// +optional
	HostList string `json:"hostList,omitempty"`

	// Specifies the Aerospike namespace to store data in. The namespace must already
	// exist on the server; Vault does not create it.
	// +optional
	// +kubebuilder:default:="test"
	Namespace string `json:"namespace,omitempty"`

	// Specifies the Aerospike set to store data in, within Namespace.
	// +optional
	Set string `json:"set,omitempty"`

	// Specifies the secret name containing username and password to authenticate to
	// the Aerospike cluster with, if it has security enabled.
	// secret data:
	//  - username=<value>
	//  - password=<value>
	// +optional
	CredentialSecretRef *core.LocalObjectReference `json:"credentialSecretRef,omitempty"`

	// Either INTERNAL (Aerospike's own credential store) or EXTERNAL (e.g. LDAP).
	// +optional
	// +kubebuilder:default:="INTERNAL"
	AuthMode string `json:"authMode,omitempty"`

	// If set, the client verifies it is connected to a cluster with this name.
	// +optional
	ClusterName string `json:"clusterName,omitempty"`

	// Socket connection timeout, in milliseconds.
	// +optional
	Timeout int64 `json:"timeout,omitempty"`

	// Idle connection timeout, in milliseconds. 0 disables idle connection trimming.
	// +optional
	IdleTimeout int64 `json:"idleTimeout,omitempty"`

	// DatabaseRef contains the info of a KubeDB managed Aerospike.
	// If set, Hostname will be generated from the referenced AppBinding.
	// +optional
	DatabaseRef *appcat.AppReference `json:"databaseRef,omitempty"`
}
