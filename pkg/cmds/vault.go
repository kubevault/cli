package cmds

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	EnableStatusSubresource bool
)

func NewCmdVault(clientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "vault",
		Short:             `manage vault`,
		DisableAutoGenTag: true,
	}

	cmd.PersistentFlags().BoolVar(&EnableStatusSubresource, "enable-status-subresource", GetDefaultValueForStatusSubresource(clientGetter), "If true, uses sub resource for crds.")

	cmd.AddCommand(NewCmdApprove(clientGetter))
	cmd.AddCommand(NewCmdDeny(clientGetter))
	return cmd
}
