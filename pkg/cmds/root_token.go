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
	token_key_store "kubevault.dev/cli/pkg/token-keys-store"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

type getTokenOptions struct {
	valueOnly bool
}

type setTokenOptions struct {
	tokenValue string
}

func newGetTokenOptions() *getTokenOptions {
	return &getTokenOptions{}
}

func newSetTokenOptions() *setTokenOptions {
	return &setTokenOptions{}
}

func (o *getTokenOptions) AddGetTokenFlag(fs *pflag.FlagSet) {
	fs.BoolVar(&o.valueOnly, "value-only", o.valueOnly, "prints only the value if flag value-only is true.")
}

func (o *setTokenOptions) AddSetTokenFlag(fs *pflag.FlagSet) {
	fs.StringVar(&o.tokenValue, "token-value", o.tokenValue, "token value to be set.")
}

func NewCmdRootToken(clientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "root-token",
		Short:             "root-token short cmd",
		Long:              "root-token [command]",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
			os.Exit(1)
		},
	}

	cmd.AddCommand(NewCmdGetToken(clientGetter))
	cmd.AddCommand(NewCmdSetToken(clientGetter))
	cmd.AddCommand(NewCmdDeleteToken(clientGetter))
	return cmd
}

func NewCmdGetToken(clientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	o := newGetTokenOptions()
	cmd := &cobra.Command{
		Use:               "get",
		Short:             "get root-token short cmd",
		Long:              "get root-token long cmd",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("get root-token")
			if len(args) > 0 {
				ResourceName = args[0]
				ObjectNames = args[1:]
			}

			if err := o.get(clientGetter); err != nil {
				Fatal(err)
			}
			os.Exit(0)
		},
	}

	o.AddGetTokenFlag(cmd.Flags())
	return cmd
}

func NewCmdSetToken(clientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	o := newSetTokenOptions()
	cmd := &cobra.Command{
		Use:               "set",
		Short:             "set root-token short cmd",
		Long:              "set root-token long cmd",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("set root-token: ", args)
			if len(args) > 0 {
				ResourceName = args[0]
				ObjectNames = args[1:]
			}

			if err := o.set(clientGetter); err != nil {
				Fatal(err)
			}
			os.Exit(0)
		},
	}

	o.AddSetTokenFlag(cmd.Flags())
	return cmd
}

func NewCmdDeleteToken(clientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete",
		Short:             "delete root-token short cmd",
		Long:              "delete root-token long cmd",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("delete root-token: ", args)
			if len(args) > 0 {
				ResourceName = args[0]
				ObjectNames = args[1:]
			}

			if err := del(clientGetter); err != nil {
				Fatal(err)
			}
			os.Exit(0)
		},
	}

	return cmd
}

func (o *getTokenOptions) get(clientGetter genericclioptions.RESTClientGetter) error {
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
			err2 = o.getRootToken(obj, kubeClient)
		default:
			err2 = errors.New("unknown/unsupported type")
		}
		return err2
	})
	return err
}

func (o *getTokenOptions) getRootToken(vs *vaultapi.VaultServer, kubeClient kubernetes.Interface) error {
	ti, err := token_key_store.NewTokenKeyInterface(vs, kubeClient)
	if err != nil {
		return err
	}

	defer func() {
		err = ti.Clean()
	}()

	name := ti.NewTokenName()
	rToken, err := ti.Get(name)
	if err == nil {
		o.Print(name, rToken)
		return nil
	}

	name = ti.OldTokenName()
	rToken, err = ti.Get(name)
	if err == nil {
		o.Print(name, rToken)
		return nil
	}

	return err
}

func del(clientGetter genericclioptions.RESTClientGetter) error {
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
			err2 = deleteRootToken(obj, kubeClient)
		default:
			err2 = errors.New("unknown/unsupported type")
		}
		return err2
	})
	return err
}

func deleteRootToken(vs *vaultapi.VaultServer, kubeClient kubernetes.Interface) error {
	ti, err := token_key_store.NewTokenKeyInterface(vs, kubeClient)
	if err != nil {
		return err
	}

	defer func() {
		err = ti.Clean()
	}()

	name := ti.NewTokenName()
	err = ti.Delete(name)
	if err == nil {
		fmt.Printf("root-token with name %s successfully deleted\n", name)
		return nil
	}

	name = ti.OldTokenName()
	err = ti.Delete(name)
	if err == nil {
		fmt.Printf("root-token with name %s successfully deleted\n", name)
		return nil
	}

	return err
}

func (o *setTokenOptions) set(clientGetter genericclioptions.RESTClientGetter) error {
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
			err2 = o.setRootToken(obj, kubeClient)
		default:
			err2 = errors.New("unknown/unsupported type")
		}
		return err2
	})
	return err
}

func (o *setTokenOptions) setRootToken(vs *vaultapi.VaultServer, kubeClient kubernetes.Interface) error {
	ti, err := token_key_store.NewTokenKeyInterface(vs, kubeClient)
	if err != nil {
		return err
	}

	defer func() {
		err = ti.Clean()
	}()

	name := ti.NewTokenName()
	if err = ti.Set(name, o.tokenValue); err != nil {
		return err
	}

	fmt.Printf("root-token with name %s, value %s successfully set\n", name, o.tokenValue)
	return err
}

func (o *getTokenOptions) Print(key, value string) {
	if o.valueOnly {
		fmt.Printf("%s\n", value)
	} else {
		fmt.Printf("%s: %s\n", key, value)
	}
}
