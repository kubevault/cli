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

	vaultapi "kubevault.dev/apimachinery/apis/kubevault/v1alpha2"
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
	tokenName string
}

type setTokenOptions struct {
	tokenValue string
	tokenName  string
}

type delTokenOptions struct {
	tokenName string
}

func newGetTokenOptions() *getTokenOptions {
	return &getTokenOptions{}
}

func newSetTokenOptions() *setTokenOptions {
	return &setTokenOptions{}
}

func newDelTokenOptions() *delTokenOptions {
	return &delTokenOptions{}
}

func (o *getTokenOptions) AddGetTokenFlag(fs *pflag.FlagSet) {
	fs.BoolVar(&o.valueOnly, "value-only", o.valueOnly, "prints only the value if flag value-only is true.")
	fs.StringVar(&o.tokenName, "token-name", o.tokenName, "get root-token with token-name.")
}

func (o *setTokenOptions) AddSetTokenFlag(fs *pflag.FlagSet) {
	fs.StringVar(&o.tokenValue, "token-value", o.tokenValue, "set latest token-name with token-value")
	fs.StringVar(&o.tokenName, "token-name", o.tokenName, "set token value root-token with token-name.")
}

func (o *delTokenOptions) AddDelTokenFlag(fs *pflag.FlagSet) {
	fs.StringVar(&o.tokenName, "token-name", o.tokenName, "delete root-token with token-name. delete the latest root-token otherwise.")
}

func NewCmdRootToken(clientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "root-token",
		Short: "get, set, delete and sync root-token",
		Long: `
$ kubectl vault root-token [command] [flags] to get, set, delete or sync vault root-token

Examples:
 $ kubectl vault root-token get [flags]
 $ kubectl vault root-token set [flags]
 $ kubectl vault root-token delete [flags]
 $ kubectl vault root-token sync [flags]
`,

		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
			os.Exit(1)
		},
	}

	cmd.AddCommand(NewCmdGetToken(clientGetter))
	cmd.AddCommand(NewCmdSetToken(clientGetter))
	cmd.AddCommand(NewCmdDeleteToken(clientGetter))
	cmd.AddCommand(NewCmdSyncToken(clientGetter))
	return cmd
}

func NewCmdGetToken(clientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	o := newGetTokenOptions()
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get vault root-token",
		Long: `
$ kubectl vault root-token get vaultserver <name> -n <namespace> [flags]

Examples:
 # get the decrypted root-token of a vaultserver with name vault in demo namespace
 $ kubectl vault root-token get vaultserver vault -n demo

 # pass the --value-only flag to get only the decrypted value
 $ kubectl vault root-token get vaultserver vault -n demo --value-only

 # pass the --token-name flag to get only the decrypted root-token value with a specific token name
 $ kubectl vault root-token get vaultserver vault -n demo --token-name <token-name> --value-only
`,
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
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
		Use:   "set",
		Short: "set vault root-token",
		Long: `
$ kubectl vault root-token set vaultserver <name> -n <namespace> [flags]

Examples:
 # set the root-token with name --token-name flag & value --token-value flag
 $ kubectl vault root-token set vaultserver vault -n demo --token-name <name> --token-value <value>

 # default name for root-token will be used if --token-name flag is not provided
 # default root-token naming format: k8s.{cluster-name or UID}.{vault-namespace}.{vault-name}-root-token
 $ kubectl vault root-token set vaultserver vault -n demo --token-value <value>
`,
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
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
	o := newDelTokenOptions()
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete vault root-token",
		Long: `
$ kubectl vault root-token delete vaultserver <name> -n <namespace> [flags]

Examples:
 # delete the root-token with name set by --token-name flag
 $ kubectl vault root-token delete vaultserver vault -n demo --token-name <name>

 # default name for root-token will be used if --token-name flag is not provided
 # default root-token naming format: k8s.{cluster-name or UID}.{vault-namespace}.{vault-name}-root-token
 $ kubectl vault root-token delete vaultserver vault -n demo
`,
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				ResourceName = args[0]
				ObjectNames = args[1:]
			}

			if err := o.del(clientGetter); err != nil {
				Fatal(err)
			}
			os.Exit(0)
		},
	}

	o.AddDelTokenFlag(cmd.Flags())
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
		ti.Clean()
	}()

	// if --token-name if provided, get token with this name
	if len(o.tokenName) > 0 {
		rToken, err := ti.Get(o.tokenName)
		if err != nil {
			return err
		}
		o.Print(o.tokenName, rToken)
		return nil
	}

	// --token-name isn't provided, look for the token with the latest naming format
	name := ti.NewTokenName()
	rToken, err := ti.Get(name)
	if err != nil {
		return err
	}

	o.Print(name, rToken)

	return nil
}

func (o *delTokenOptions) del(clientGetter genericclioptions.RESTClientGetter) error {
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
			err2 = o.deleteRootToken(obj, kubeClient)
		default:
			err2 = errors.New("unknown/unsupported type")
		}
		return err2
	})
	return err
}

func (o *delTokenOptions) deleteRootToken(vs *vaultapi.VaultServer, kubeClient kubernetes.Interface) error {
	ti, err := token_key_store.NewTokenKeyInterface(vs, kubeClient)
	if err != nil {
		return err
	}

	defer func() {
		ti.Clean()
	}()

	name := ti.NewTokenName()
	if len(o.tokenName) > 0 {
		name = o.tokenName
	}

	err = ti.Delete(name)
	if err != nil {
		return err
	}

	fmt.Printf("root-token with name %s successfully deleted\n", name)

	return nil
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

func NewCmdSyncToken(clientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "sync vault root-token",
		Long: `
$ kubectl vault root-token sync vaultserver <name> -n <namespace> [flags]

Examples:
 # sync the vaultserver root-token 
 # old naming conventions: vault-root-token
 # new naming convention for root-token: k8s.{cluster-name or UID}.{vault-namespace}.{vault-name}-root-token
 $ kubectl vault root-token sync vaultserver vault -n demo
`,
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				ResourceName = args[0]
				ObjectNames = args[1:]
			}

			if err := syncRootToken(clientGetter); err != nil {
				Fatal(err)
			}
			os.Exit(0)
		},
	}

	return cmd
}

func syncRootToken(clientGetter genericclioptions.RESTClientGetter) error {
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
			err2 = syncToken(obj, kubeClient)
		default:
			err2 = errors.New("unknown/unsupported type")
		}
		return err2
	})
	return err
}

func syncToken(vs *vaultapi.VaultServer, kubeClient kubernetes.Interface) error {
	ti, err := token_key_store.NewTokenKeyInterface(vs, kubeClient)
	if err != nil {
		return err
	}

	defer func() {
		ti.Clean()
	}()

	// if new key already exists just return
	newKey := ti.NewTokenName()
	if _, err = ti.Get(newKey); err == nil {
		fmt.Printf("%s already up-to-date\n", newKey)
		fmt.Println("successfully synced root-token")
		return nil
	}

	// new key doesn't exist, check for old key
	oldKey := ti.OldTokenName()
	value, err := ti.Get(oldKey)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// old key exist, set the value to new key
	if err = ti.Set(newKey, value); err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("%s successfully synced\n", newKey)
	fmt.Println("successfully synced root-token")
	return nil
}

func (o *setTokenOptions) setRootToken(vs *vaultapi.VaultServer, kubeClient kubernetes.Interface) error {
	ti, err := token_key_store.NewTokenKeyInterface(vs, kubeClient)
	if err != nil {
		return err
	}

	defer func() {
		ti.Clean()
	}()

	if len(o.tokenValue) == 0 {
		return errors.New("token value is empty")
	}

	name := ti.NewTokenName()
	if len(o.tokenName) > 0 {
		name = o.tokenName
	}

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
