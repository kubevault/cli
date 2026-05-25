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
	engineapi "kubevault.dev/apimachinery/apis/engine/v1alpha1"
	vaultv1a2 "kubevault.dev/apimachinery/apis/kubevault/v1alpha2"
	policyapi "kubevault.dev/apimachinery/apis/policy/v1alpha1"

	"github.com/spf13/cobra"
	v "gomodules.xyz/x/version"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	appcatscheme "kmodules.xyz/custom-resources/client/clientset/versioned/scheme"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "vault [command]",
		Short:             `KubeVault cli by AppsCode`,
		DisableAutoGenTag: true,
		PersistentPreRun: func(c *cobra.Command, args []string) {
			utilruntime.Must(engineapi.AddToScheme(clientsetscheme.Scheme))
			utilruntime.Must(vaultv1a2.AddToScheme(clientsetscheme.Scheme))
			utilruntime.Must(policyapi.AddToScheme(clientsetscheme.Scheme))
			utilruntime.Must(appcatscheme.AddToScheme(clientsetscheme.Scheme))
		},
	}

	flags := rootCmd.PersistentFlags()

	kubeConfigFlags := genericclioptions.NewConfigFlags(true)
	kubeConfigFlags.AddFlags(flags)
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)
	matchVersionKubeConfigFlags.AddFlags(flags)

	rootCmd.AddCommand(NewCmdCompletion())
	rootCmd.AddCommand(v.NewCmdVersion())

	rootCmd.AddCommand(NewCmdApprove(matchVersionKubeConfigFlags))
	rootCmd.AddCommand(NewCmdDeny(matchVersionKubeConfigFlags))
	rootCmd.AddCommand(NewCmdRevoke(matchVersionKubeConfigFlags))
	rootCmd.AddCommand(NewCmdGenerate(matchVersionKubeConfigFlags))
	rootCmd.AddCommand(NewCmdRootToken(matchVersionKubeConfigFlags))
	rootCmd.AddCommand(NewCmdUnsealKey(matchVersionKubeConfigFlags))
	rootCmd.AddCommand(NewCmdMergeSecrets(matchVersionKubeConfigFlags))
	return rootCmd
}
