/*
Copyright The KubeVault Authors.

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

package cmds

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	core "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	kmapi "kmodules.xyz/client-go/api/v1"
)

var (
	awsDeniedCond = kmapi.Condition{
		Type:    kmapi.ConditionRequestDenied,
		Status:  core.ConditionTrue,
		Reason:  "KubectlDeny",
		Message: "This was denied by: kubectl vault deny awsaccesskeyrequest",
	}

	dbDeniedCond = kmapi.Condition{
		Type:    kmapi.ConditionRequestDenied,
		Status:  core.ConditionTrue,
		Reason:  "KubectlDeny",
		Message: "This was denied by: kubectl vault deny databaseaccessrequest",
	}

	gcpDeniedCond = kmapi.Condition{
		Type:    kmapi.ConditionRequestDenied,
		Status:  core.ConditionTrue,
		Reason:  "KubectlDeny",
		Message: "This was denied by: kubectl vault deny gcpaccesskeyrequest",
	}

	azureDeniedCond = kmapi.Condition{
		Type:    kmapi.ConditionRequestDenied,
		Status:  core.ConditionTrue,
		Reason:  "KubectlDeny",
		Message: "This was denied by: kubectl vault deny azureaccesskeyrequest",
	}
)

func NewCmdDeny(clientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "deny",
		Short:             "Deny request",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				ResourceName = args[0]
				ObjectNames = args[1:]
			}

			if err := modifyStatusCondition(clientGetter, false); err != nil {
				Fatal(err)
			} else {
				fmt.Println("Denied")
			}
			os.Exit(0)
		},
	}

	cmdutil.AddFilenameOptionFlags(cmd, &FilenameOptions, "identifying the resource to update")
	return cmd
}
