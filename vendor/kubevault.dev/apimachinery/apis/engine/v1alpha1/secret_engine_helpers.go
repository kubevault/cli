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

	api "kubevault.dev/apimachinery/apis/kubevault/v1alpha2"
	"kubevault.dev/apimachinery/crds"

	"kmodules.xyz/client-go/apiextensions"
	"kmodules.xyz/client-go/tools/clusterid"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
)

func (_ SecretEngine) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourceSecretEngines))
}

func (e SecretEngine) IsValid() error {
	return nil
}

// Generates the policy name which contains
// required permission for this secret engine
func (e SecretEngine) GetPolicyName() string {
	cluster := "-"
	if clusterid.ClusterName() != "" {
		cluster = clusterid.ClusterName()
	}
	return fmt.Sprintf("k8s.%s.%s.%s", cluster, e.Namespace, e.Name)
}

// Generates unique database name from database appbinding reference
func GetDBNameFromAppBindingRef(dbAppRef *appcat.AppReference) string {
	cluster := "-"
	if clusterid.ClusterName() != "" {
		cluster = clusterid.ClusterName()
	}
	return fmt.Sprintf("k8s.%s.%s.%s", cluster, dbAppRef.Namespace, dbAppRef.Name)
}

func (se SecretEngine) GetSecretEnginePath() string {
	// Todo: update SecretEngine path
	//  - k8s.{cluster-name or -}.{se-type}.se-ns.se-name
	cluster := "-"
	if clusterid.ClusterName() != "" {
		cluster = clusterid.ClusterName()
	}
	return fmt.Sprintf("k8s.%s.%s.%s.%s", cluster, se.GetSecretEngineType(), se.Namespace, se.Name)
}

func (se SecretEngine) GetSecretEngineType() api.SecretEngineType {
	switch seSpec := se.Spec; {
	case seSpec.GCP != nil:
		return api.SecretEngineTypeGCP
	case seSpec.AWS != nil:
		return api.SecretEngineTypeAWS
	case seSpec.Azure != nil:
		return api.SecretEngineTypeAzure
	case seSpec.Elasticsearch != nil:
		return api.SecretEngineTypeElasticsearch
	case seSpec.MongoDB != nil:
		return api.SecretEngineTypeMongoDB
	case seSpec.MySQL != nil:
		return api.SecretEngineTypeMySQL
	case seSpec.Postgres != nil:
		return api.SecretEngineTypePostgres
	case seSpec.MariaDB != nil:
		return api.SecretEngineTypeMariaDB
	default:
		return ""
	}
}
