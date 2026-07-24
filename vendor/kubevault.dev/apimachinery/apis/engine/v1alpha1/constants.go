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

package v1alpha1

const (
	SecretEnginePhaseSuccess    SecretEnginePhase = "Success"
	SecretEnginePhaseProcessing SecretEnginePhase = "Processing"

	ConditionCertificateSigned = "CertificateSigned"

	// Tenant-isolation condition types on SecretEngine (see design/tenant-namespace-design.md).
	// TenantNamespaceUnresolved is set when the database's namespace carries the client-org
	// label but no (or an invalid) org-id annotation; the engine is not mounted and the
	// request is requeued so org data never lands in the shared root tree.
	TenantNamespaceUnresolved = "TenantNamespaceUnresolved"
	// TenantMigrationPending is set when the desired OpenBao namespace differs from the one
	// already recorded in status.effectiveNamespace; the operator does not move a live mount
	// on its own and waits for an admin-authorized migration.
	TenantMigrationPending = "TenantMigrationPending"

	// TenantNamespacePendingHub is set on a spoke (RemoteRelay) SecretEngine whose derived
	// org namespace does not yet exist on the hub. The spoke cannot create hub namespaces;
	// it requeues until the hub creates sys/namespaces/<org-id>
	// (design/tenant-namespace-design.md §7.2).
	TenantNamespacePendingHub = "TenantNamespacePendingHub"
)
