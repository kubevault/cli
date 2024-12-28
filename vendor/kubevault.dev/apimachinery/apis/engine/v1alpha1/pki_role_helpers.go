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

	"kubevault.dev/apimachinery/crds"

	"kmodules.xyz/client-go/apiextensions"
	clustermeta "kmodules.xyz/client-go/cluster"
)

func (_ PKIRole) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePKIRoles))
}

func (p PKIRole) RoleName() string {
	cluster := "-"
	if clustermeta.ClusterName() != "" {
		cluster = clustermeta.ClusterName()
	}
	return fmt.Sprintf("k8s.%s.%s.%s", cluster, p.Namespace, p.Name)
}

func (p PKIRole) IsValid() error {
	return nil
}