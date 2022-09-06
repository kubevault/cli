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

import (
	"fmt"

	"kubevault.dev/apimachinery/apis/ops"
	"kubevault.dev/apimachinery/crds"

	"kmodules.xyz/client-go/apiextensions"
)

func (_ VaultOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralVaultOpsRequest))
}

func (e VaultOpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralVaultOpsRequest, ops.GroupName)
}

func (e VaultOpsRequest) ResourceShortCode() string {
	return ResourceCodeVaultOpsRequest
}

func (e VaultOpsRequest) ResourceKind() string {
	return ResourceKindVaultOpsRequest
}

func (e VaultOpsRequest) ResourceSingular() string {
	return ResourceSingularVaultOpsRequest
}

func (e VaultOpsRequest) ResourcePlural() string {
	return ResourcePluralVaultOpsRequest
}

func (e VaultOpsRequest) ValidateSpecs() error {
	return nil
}
