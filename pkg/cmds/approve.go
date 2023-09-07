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

package cmds

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	kmapi "kmodules.xyz/client-go/api/v1"
	condutil "kmodules.xyz/client-go/conditions"
)

var (
	FilenameOptions resource.FilenameOptions
	ObjectNames     []string
	ResourceName    string
)

var secretAccessApprovedCond = kmapi.Condition{
	Type:    condutil.ConditionRequestApproved,
	Status:  metav1.ConditionTrue,
	Reason:  "KubectlApprove",
	Message: "This was approved by: kubectl vault approve secretaccessrequest",
}

func NewCmdApprove(clientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "approve",
		Short:             "Approve request",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				ResourceName = args[0]
				ObjectNames = args[1:]
			}

			if err := modifyStatusCondition(clientGetter, secretAccessApprovedCond); err != nil {
				Fatal(err)
			} else {
				fmt.Printf("secretaccessrequests %s approved\n", strings.Join(ObjectNames, ", "))
			}
			os.Exit(0)
		},
	}

	cmdutil.AddFilenameOptionFlags(cmd, &FilenameOptions, "identifying the resource to update")
	return cmd
}
