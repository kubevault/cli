package cmds

import (
	"fmt"
	"os"

	engineapi "github.com/kubevault/operator/apis/engine/v1alpha1"
	enginecs "github.com/kubevault/operator/client/clientset/versioned/typed/engine/v1alpha1"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

var (
	awsDeniedCond = engineapi.AWSAccessKeyRequestCondition{
		Type:           engineapi.AccessDenied,
		Reason:         "KubectlDeny",
		Message:        "This was denied by kubectl vault deny awsaccesskeyrequest",
		LastUpdateTime: metav1.Now(),
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

			cfg, err := clientGetter.ToRESTConfig()
			if err != nil {
				Fatal(errors.Wrap(err, "failed to read kubeconfig"))
			}

			namespace, _, err := clientGetter.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				Fatal(err)
			}

			builder := cmdutil.NewFactory(clientGetter).NewBuilder()

			c, err := enginecs.NewForConfig(cfg)
			if err != nil {
				Fatal(err)
			}

			if err := modifyAWSAccessKeyCondition(builder, c, namespace, awsDeniedCond); err != nil {
				Fatal(err)
			}

			fmt.Printf("%s is denied", ResourceName)
			os.Exit(0)
		},
	}

	cmdutil.AddFilenameOptionFlags(cmd, &FilenameOptions, "identifying the resource to update")
	return cmd
}
