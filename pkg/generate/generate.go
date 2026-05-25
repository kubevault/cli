/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Community License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Community-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package generate

import (
	engineapi "kubevault.dev/apimachinery/apis/engine/v1alpha1"
	"kubevault.dev/cli/pkg/generate/api"
	"kubevault.dev/cli/pkg/generate/aws"
	"kubevault.dev/cli/pkg/generate/azure"
	es "kubevault.dev/cli/pkg/generate/database/elasticsearch"
	maria "kubevault.dev/cli/pkg/generate/database/mariadb"
	mongo "kubevault.dev/cli/pkg/generate/database/mongodb"
	sql "kubevault.dev/cli/pkg/generate/database/mysql"
	pg "kubevault.dev/cli/pkg/generate/database/postgres"
	rd "kubevault.dev/cli/pkg/generate/database/redis"
	"kubevault.dev/cli/pkg/generate/gcp"

	"github.com/go-errors/errors"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewGenerator(role []string, srb *engineapi.SecretRoleBinding, keys map[string]string, kc client.Client, kubeClient *kubernetes.Clientset) (api.GeneratorInterface, error) {
	switch role[0] {
	case engineapi.ResourceKindGCPRole:
		return gcp.NewGCPGenerator(role, srb, keys, kc, kubeClient)
	case engineapi.ResourceKindAWSRole:
		return aws.NewAWSGenerator(role, srb, keys, kc, kubeClient)
	case engineapi.ResourceKindAzureRole:
		return azure.NewAzureGenerator(role, srb, keys, kc, kubeClient)
	case engineapi.ResourceKindMongoDBRole:
		return mongo.NewMongoGenerator(role, srb, keys, kc, kubeClient)
	case engineapi.ResourceKindElasticsearchRole:
		return es.NewElasticsearchGenerator(role, srb, keys, kc, kubeClient)
	case engineapi.ResourceKindMySQLRole:
		return sql.NewMySQLGenerator(role, srb, keys, kc, kubeClient)
	case engineapi.ResourceKindMariaDBRole:
		return maria.NewMariaDBGenerator(role, srb, keys, kc, kubeClient)
	case engineapi.ResourceKindRedisRole:
		return rd.NewRedisGenerator(role, srb, keys, kc, kubeClient)
	case engineapi.ResourceKindPostgresRole:
		return pg.NewPostgresGenerator(role, srb, keys, kc, kubeClient)
	default:
		return nil, errors.New("unknown role")
	}
}
