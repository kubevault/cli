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

	vaultapi "kubevault.dev/apimachinery/apis/kubevault/v1alpha1"
	token "kubevault.dev/cli/pkg/get-root-token"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

type options struct {
	valueOnly bool
}

func NewOptions() *options {
	return &options{}
}

func (o *options) AddTokenFlag(fs *pflag.FlagSet) {
	fs.BoolVar(&o.valueOnly, "value-only", o.valueOnly, "prints only the value if flag value-only is true.")
}

func NewCmdGetRootToken(clientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	o := NewOptions()
	cmd := &cobra.Command{
		Use:   "get-root-token",
		Short: "Get root token for vault server",
		Long: "Get the decrypted root token for a vault server. You can provide flags vaultserver name and namespace.\n\n" +
			"Examples: \n # Get the decrypted root-token for resource vaultserver with name vault and namespace demo.\n # Prints the root-token name and the value.\n " +
			"$ kubectl vault get-root-token vaultserver -n demo vault\n\n" +
			" # Provide flag --value-only if you want to print only the root-token value.\n" +
			" $ kubectl vault get-root-token vaultserver -n demo vault --value-only",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				ResourceName = args[0]
				ObjectNames = args[1:]
			}

			if err := o.getRootToken(clientGetter); err != nil {
				Fatal(err)
			}
			os.Exit(0)
		},
	}

	cmdutil.AddFilenameOptionFlags(cmd, &FilenameOptions, "identifying the resource to update")
	o.AddTokenFlag(cmd.Flags())
	return cmd
}

func (o *options) getRootToken(clientGetter genericclioptions.RESTClientGetter) error {
	var resourceName string
	switch ResourceName {
	case strings.ToLower(vaultapi.ResourceVaultServer), strings.ToLower(vaultapi.ResourceVaultServers):
		resourceName = vaultapi.ResourceVaultServer
	default:
		return errors.New(fmt.Sprintf("unknown/unsupported resource %s", ResourceName))
	}

	namespace, _, err := clientGetter.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	cfg, err := clientGetter.ToRESTConfig()
	if err != nil {
		return errors.Wrap(err, "failed to read kubeconfig")
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return err
	}

	builder := cmdutil.NewFactory(clientGetter).NewBuilder()
	r := builder.
		WithScheme(clientsetscheme.Scheme, clientsetscheme.Scheme.PrioritizedVersionsAllGroups()...).
		ContinueOnError().
		NamespaceParam(namespace).DefaultNamespace().
		FilenameParam(false, &FilenameOptions).
		ResourceNames(resourceName, ObjectNames...).
		RequireObject(true).
		Flatten().
		Latest().
		Do()

	err = r.Visit(func(info *resource.Info, err error) error {
		if err != nil {
			return err
		}

		var err2 error
		switch info.Object.(type) {
		case *vaultapi.VaultServer:
			obj := info.Object.(*vaultapi.VaultServer)
			err2 = o.printRootToken(obj, kubeClient)
		default:
			err2 = errors.New("unknown/unsupported type")
		}
		return err2
	})
	return err
}

func (o *options) printRootToken(vs *vaultapi.VaultServer, kubeClient kubernetes.Interface) error {
	ti, err := token.NewTokenInterface(vs, kubeClient)
	if err != nil {
		return err
	}
	rToken, err := ti.Token()
	if err != nil {
		return err
	}
	if o.valueOnly {
		fmt.Printf("%s", rToken)
	} else {
		fmt.Printf("%s: %s\n", ti.TokenName(), rToken)
	}
	return nil
}
