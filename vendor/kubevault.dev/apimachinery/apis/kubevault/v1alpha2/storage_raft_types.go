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

// RaftSpec defines the configuration for the Raft integrated storage.
// https://www.vaultproject.io/docs/configuration/storage/raft
type RaftSpec struct {
	// An integer multiplier used by servers to scale key Raft timing parameters.
	// Tuning this affects the time it takes Vault to detect leader failures and to perform leader elections,
	// at the expense of requiring more network and CPU resources for better performance.
	// default: 0
	// +optional
	PerformanceMultiplier int64 `json:"performanceMultiplier,omitempty"`

	// This controls how many log entries are left in the log store on disk after a snapshot is made.
	// default: 10000
	// +optional
	TrailingLogs *int64 `json:"trailingLogs,omitempty"`

	// This controls the minimum number of raft commit entries between snapshots that are saved to disk.
	// default: 8192
	// +optional
	SnapshotThreshold *int64 `json:"snapshotThreshold,omitempty"`

	// This configures the maximum number of bytes for a raft entry. It applies to both Put operations and transactions.
	// default: 1048576
	// +optional
	MaxEntrySize *int64 `json:"maxEntrySize,omitempty"`

	// This is the interval after which autopilot will pick up any state changes.
	// default: ""
	// +optional
	AutopilotReconcileInterval string `json:"autopilotReconcileInterval,omitempty"`

	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`
}
