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
	RemotePostgresDatabasePlugin = "remote-postgres-plugin"
	RemoteMySQLDatabasePlugin    = "remote-mysql-plugin"
	RemoteRedisDatabasePlugin    = "remote-redis-plugin"
	RemoteValkeyDatabasePlugin   = "remote-valkey-plugin"
)

// RemoteDatabasePlugin maps a database engine kind to the hub-side proxy
// plugin used when the Vault AppBinding is of deployment type RemoteRelay.
// Engine kinds without a remote plugin (mongodb, elasticsearch) return an
// error so misconfigurations fail loudly instead of silently writing a
// local plugin name into a hub mount.
func RemoteDatabasePlugin(engineKind string) (string, error) {
	switch engineKind {
	case "postgres":
		return RemotePostgresDatabasePlugin, nil
	case "mysql", "mariadb":
		// MariaDB uses the MySQL plugin locally and remotely.
		return RemoteMySQLDatabasePlugin, nil
	case "redis":
		return RemoteRedisDatabasePlugin, nil
	case "valkey":
		return RemoteValkeyDatabasePlugin, nil
	default:
		return "", fmt.Errorf("database engine %q is not supported through the OpenBao spoke relay; supported: postgres, mysql, mariadb, redis, valkey", engineKind)
	}
}
