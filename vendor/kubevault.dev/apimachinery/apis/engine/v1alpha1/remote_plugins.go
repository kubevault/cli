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

import "fmt"

// Hub-side proxy plugin names registered by the OpenBao remote-db-plugin
// branch (helper/builtinplugins/registry.go). Each brokers credential
// operations over mTLS gRPC to the in-process built-in plugin running on
// the named spoke.
const (
	RemotePostgresDatabasePlugin    = "remote-postgres-plugin"
	RemoteMySQLDatabasePlugin       = "remote-mysql-plugin"
	RemoteRedisDatabasePlugin       = "remote-redis-plugin"
	RemoteValkeyDatabasePlugin      = "remote-valkey-plugin"
	RemoteDruidDatabasePlugin       = "remote-druid-plugin"
	RemoteEtcdDatabasePlugin        = "remote-etcd-plugin"
	RemoteHanaDBDatabasePlugin      = "remote-hana-plugin"
	RemoteIgniteDatabasePlugin      = "remote-ignite-plugin"
	RemoteKafkaDatabasePlugin       = "remote-kafka-plugin"
	RemoteMilvusDatabasePlugin      = "remote-milvus-plugin"
	RemoteMSSQLServerDatabasePlugin = "remote-mssql-plugin"
	RemoteNeo4jDatabasePlugin       = "remote-neo4j-plugin"
	RemoteOracleDatabasePlugin      = "remote-oracle-plugin"
	RemoteQdrantDatabasePlugin      = "remote-qdrant-plugin"
	RemoteRabbitMQDatabasePlugin    = "remote-rabbitmq-plugin"
	RemoteSolrDatabasePlugin        = "remote-solr-plugin"
	RemoteWeaviateDatabasePlugin    = "remote-weaviate-plugin"
)

// RemoteDatabasePlugin maps a database engine kind to the hub-side proxy
// plugin used when the Vault AppBinding is of deployment type RemoteRelay.
// Engine kinds without a remote plugin (mongodb, elasticsearch) return an
// error so misconfigurations fail loudly instead of silently writing a
// local plugin name into a hub mount.
func RemoteDatabasePlugin(engineKind string) (string, error) {
	switch engineKind {
	case "druid":
		return RemoteDruidDatabasePlugin, nil
	case "etcd":
		return RemoteEtcdDatabasePlugin, nil
	case "hanadb":
		return RemoteHanaDBDatabasePlugin, nil
	case "ignite":
		return RemoteIgniteDatabasePlugin, nil
	case "kafka":
		return RemoteKafkaDatabasePlugin, nil
	case "mariadb", "mysql":
		// MariaDB uses the MySQL plugin locally and remotely.
		return RemoteMySQLDatabasePlugin, nil
	case "milvus":
		return RemoteMilvusDatabasePlugin, nil
	case "mssqlserver":
		return RemoteMSSQLServerDatabasePlugin, nil
	case "neo4j":
		return RemoteNeo4jDatabasePlugin, nil
	case "oracle":
		return RemoteOracleDatabasePlugin, nil
	case "postgres":
		return RemotePostgresDatabasePlugin, nil
	case "qdrant":
		return RemoteQdrantDatabasePlugin, nil
	case "rabbitmq":
		return RemoteRabbitMQDatabasePlugin, nil
	case "redis":
		return RemoteRedisDatabasePlugin, nil
	case "solr":
		return RemoteSolrDatabasePlugin, nil
	case "valkey":
		return RemoteValkeyDatabasePlugin, nil
	case "weaviate":
		return RemoteWeaviateDatabasePlugin, nil
	default:
		return "", fmt.Errorf("database engine %q is not supported through the OpenBao spoke relay; supported: druid, etcd, hanadb, ignite, kafka, mariadb, milvus, mssqlserver, mysql, neo4j, oracle, postgres, qdrant, rabbitmq, redis, solr, valkey, weaviate", engineKind)
	}
}
