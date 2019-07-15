package cmds

import (
	"fmt"
	"os"

	dbapi "github.com/kubedb/apimachinery/apis/authorization/v1alpha1"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	engineapi "kubevault.dev/operator/apis/engine/v1alpha1"
)

var (
	awsDeniedCond = engineapi.AWSAccessKeyRequestCondition{
		Type:    engineapi.AccessDenied,
		Reason:  "KubectlDeny",
		Message: "This was denied by kubectl vault deny awsaccesskeyrequest",
	}

	dbDeniedCond = dbapi.DatabaseAccessRequestCondition{
		Type:    dbapi.AccessDenied,
		Reason:  "KubectlDeny",
		Message: "This was denied by kubectl vault deny databaseaccessrequest",
	}

	gcpDeniedCond = engineapi.GCPAccessKeyRequestCondition{
		Type:    engineapi.AccessDenied,
		Reason:  "KubectlDeny",
		Message: "This was denied by kubectl vault deny gcpaccesskeyrequest",
	}

	azureDeniedCond = engineapi.AzureAccessKeyRequestCondition{
		Type:           engineapi.AccessDenied,
		Reason:         "KubectlDeny",
		Message:        "This was denied by kubectl vault deny azureaccesskeyrequest",
		LastUpdateTime: v1.Time{},
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
			if EnableStatusSubresource {
				EnableStatusSubresource = GetDefaultValueForStatusSubresource(clientGetter)
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
