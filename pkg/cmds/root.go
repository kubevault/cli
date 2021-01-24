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
	"flag"

	"kubevault.dev/apimachinery/client/clientset/versioned/scheme"

	"github.com/spf13/cobra"
	v "gomodules.xyz/x/version"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	cliflag "k8s.io/component-base/cli/flag"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"kmodules.xyz/client-go/logs"
	"kmodules.xyz/client-go/tools/cli"
	appcatscheme "kmodules.xyz/custom-resources/client/clientset/versioned/scheme"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "vault [command]",
		Short:             `KubeVault cli by AppsCode`,
		DisableAutoGenTag: true,
		PersistentPreRun: func(c *cobra.Command, args []string) {
			cli.SendAnalytics(c, v.Version.Version)

			utilruntime.Must(scheme.AddToScheme(clientsetscheme.Scheme))
			utilruntime.Must(appcatscheme.AddToScheme(clientsetscheme.Scheme))
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
	utilruntime.Must(flag.Set("stderrthreshold", "ERROR"))

	rootCmd.AddCommand(NewCmdApprove(matchVersionKubeConfigFlags))
	rootCmd.AddCommand(NewCmdDeny(matchVersionKubeConfigFlags))
	rootCmd.AddCommand(NewCmdCompletion())
	rootCmd.AddCommand(v.NewCmdVersion())
	return rootCmd
}
