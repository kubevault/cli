package cmds

import (
	"flag"

	v "github.com/appscode/go/version"
	dbscheme "github.com/kubedb/apimachinery/client/clientset/versioned/scheme"
	"github.com/kubevault/operator/client/clientset/versioned/scheme"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	cliflag "k8s.io/component-base/cli/flag"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"kmodules.xyz/client-go/logs"
	"kmodules.xyz/client-go/tools/cli"
	appcatscheme "kmodules.xyz/custom-resources/client/clientset/versioned/scheme"
)

var (
	EnableStatusSubresource bool
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "vault [command]",
		Short:             `KubeVault cli by AppsCode`,
		DisableAutoGenTag: true,
		PersistentPreRun: func(c *cobra.Command, args []string) {
			cli.SendAnalytics(c, v.Version.Version)

			scheme.AddToScheme(clientsetscheme.Scheme)
			appcatscheme.AddToScheme(clientsetscheme.Scheme)
			dbscheme.AddToScheme(clientsetscheme.Scheme)
		},
	}

	flags := rootCmd.PersistentFlags()
	// Normalize all flags that are coming from other packages or pre-configurations
	// a.k.a. change all "_" to "-". e.g. glog package
	flags.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)

	kubeConfigFlags := genericclioptions.NewConfigFlags(true)
	kubeConfigFlags.AddFlags(flags)
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)
	matchVersionKubeConfigFlags.AddFlags(flags)

	flags.AddGoFlagSet(flag.CommandLine)
	logs.ParseFlags()
	flags.BoolVar(&cli.EnableAnalytics, "analytics", cli.EnableAnalytics, "Send analytical events to Google Analytics")
	flag.Set("stderrthreshold", "ERROR")
	flags.BoolVar(&EnableStatusSubresource, "enable-status-subresource", true, "If true, uses sub resource for crds.")

	rootCmd.AddCommand(NewCmdApprove(matchVersionKubeConfigFlags))
	rootCmd.AddCommand(NewCmdDeny(matchVersionKubeConfigFlags))
	rootCmd.AddCommand(v.NewCmdVersion())
	return rootCmd
}
