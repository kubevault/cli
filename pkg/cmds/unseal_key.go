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

type getKeyOptions struct {
	keyId   int
	keyName string
}

type setKeyOptions struct {
	keyId    int
	keyValue string
	keyName  string
}

type delKeyOptions struct {
	keyId   int
	keyName string
}

func newGetKeyOptions() *getKeyOptions {
	return &getKeyOptions{}
}

func newSetKeyOptions() *setKeyOptions {
	return &setKeyOptions{}
}

func newDelKeyOptions() *delKeyOptions {
	return &delKeyOptions{}
}

func (o *getKeyOptions) addGetKeyFlags(fs *pflag.FlagSet) {
	fs.IntVar(&o.keyId, "key-id", o.keyId, "get the latest unseal key with id")
	fs.StringVar(&o.keyName, "key-name", o.keyName, "get unseal key with key-name")
}

func (o *setKeyOptions) addSetKeyFlags(fs *pflag.FlagSet) {
	fs.IntVar(&o.keyId, "key-id", o.keyId, "set the latest unseal key with id")
	fs.StringVar(&o.keyName, "key-name", o.keyName, "set unseal key with key-name")
	fs.StringVar(&o.keyValue, "key-value", o.keyValue, "set unseal key with key-value")
}

func (o *delKeyOptions) addDelKeyFlags(fs *pflag.FlagSet) {
	fs.IntVar(&o.keyId, "key-id", o.keyId, "delete the latest unseal key with id")
	fs.StringVar(&o.keyName, "key-name", o.keyName, "delete unseal key with key-name")
}

func NewCmdUnsealKey(clientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "unseal-key",
		Short:             "unseal-key\" short cmd",
		Long:              "unseal-key [command]",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
			os.Exit(1)
		},
	}

	cmd.AddCommand(NewCmdGetKey(clientGetter))
	cmd.AddCommand(NewCmdSetKey(clientGetter))
	cmd.AddCommand(NewCmdDeleteKey(clientGetter))
	cmd.AddCommand(NewCmdListKey(clientGetter))
	return cmd
}

func NewCmdGetKey(clientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	o := newGetKeyOptions()
	cmd := &cobra.Command{
		Use:               "get",
		Short:             "get unseal-key short cmd",
		Long:              "get unseal-key long cmd",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("get unseal-key")
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

	o.addGetKeyFlags(cmd.Flags())
	return cmd
}

func NewCmdSetKey(clientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	o := newSetKeyOptions()
	cmd := &cobra.Command{
		Use:               "set",
		Short:             "set unseal-key short cmd",
		Long:              "set unseal-key long cmd",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("set unseal-key: ", args)
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

	o.addSetKeyFlags(cmd.Flags())
	return cmd
}

func NewCmdDeleteKey(clientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	o := newDelKeyOptions()
	cmd := &cobra.Command{
		Use:               "delete",
		Short:             "delete unseal-key short cmd",
		Long:              "delete unseal-key long cmd",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("delete unseal-key: ", args)
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

	o.addDelKeyFlags(cmd.Flags())
	return cmd
}

func NewCmdListKey(clientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	o := newGetKeyOptions()
	cmd := &cobra.Command{
		Use:               "list",
		Short:             "list unseal-key short cmd",
		Long:              "list unseal-key long cmd",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("list unseal-key: ", args)
			if len(args) > 0 {
				ResourceName = args[0]
				ObjectNames = args[1:]
			}

			if err := o.list(clientGetter); err != nil {
				Fatal(err)
			}
			os.Exit(0)
		},
	}

	return cmd
}

func (o *getKeyOptions) list(clientGetter genericclioptions.RESTClientGetter) error {
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
			o.listUnsealKey(obj, kubeClient)
		default:
			err2 = errors.New("unknown/unsupported type")
		}
		return err2
	})
	return err
}

func (o *getKeyOptions) listUnsealKey(vs *vaultapi.VaultServer, kubeClient kubernetes.Interface) {
	cnt := vs.Spec.Unsealer.SecretShares
	for i := 0; int64(i) < cnt; i++ {
		o.keyId = i
		err := o.getUnsealKey(vs, kubeClient)
		if err != nil {
			fmt.Printf("vault-unseal-key-%d not found\n", i)
		}
	}
}

func (o *getKeyOptions) get(clientGetter genericclioptions.RESTClientGetter) error {
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
			err2 = o.getUnsealKey(obj, kubeClient)
		default:
			err2 = errors.New("unknown/unsupported type")
		}
		return err2
	})
	return err
}

func (o *getKeyOptions) getUnsealKey(vs *vaultapi.VaultServer, kubeClient kubernetes.Interface) error {
	ti, err := token_key_store.NewTokenKeyInterface(vs, kubeClient)
	if err != nil {
		return err
	}

	defer func() {
		ti.Clean()
	}()

	name, err := ti.NewUnsealKeyName(o.keyId)
	if err != nil {
		return err
	}

	if len(o.keyName) > 0 {
		name = o.keyName
	}

	rToken, err := ti.Get(name)
	if err != nil {
		return err
	}

	o.Print(name, rToken)

	return nil
}

func (o *delKeyOptions) del(clientGetter genericclioptions.RESTClientGetter) error {
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
			err2 = o.deleteUnsealKey(obj, kubeClient)
		default:
			err2 = errors.New("unknown/unsupported type")
		}
		return err2
	})
	return err
}

func (o *delKeyOptions) deleteUnsealKey(vs *vaultapi.VaultServer, kubeClient kubernetes.Interface) error {
	ti, err := token_key_store.NewTokenKeyInterface(vs, kubeClient)
	if err != nil {
		return err
	}

	defer func() {
		ti.Clean()
	}()

	name, err := ti.NewUnsealKeyName(o.keyId)
	if err != nil {
		return err
	}

	if len(o.keyName) > 0 {
		name = o.keyName
	}

	err = ti.Delete(name)
	if err != nil {
		return err
	}

	fmt.Printf("unseal-key with name %s successfully deleted\n", name)

	return nil
}

func (o *setKeyOptions) set(clientGetter genericclioptions.RESTClientGetter) error {
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
			err2 = o.setUnsealKey(obj, kubeClient)
		default:
			err2 = errors.New("unknown/unsupported type")
		}
		return err2
	})
	return err
}

func (o *setKeyOptions) setUnsealKey(vs *vaultapi.VaultServer, kubeClient kubernetes.Interface) error {
	ti, err := token_key_store.NewTokenKeyInterface(vs, kubeClient)
	if err != nil {
		return err
	}

	defer func() {
		ti.Clean()
	}()

	if len(o.keyValue) == 0 {
		return errors.New("unseal key value is empty")
	}

	name, err := ti.NewUnsealKeyName(o.keyId)
	if err != nil {
		return err
	}

	if len(o.keyName) > 0 {
		name = o.keyName
	}

	err = ti.Set(name, o.keyValue)
	if err != nil {
		return err
	}

	fmt.Printf("unseal-key with name %s, value %s successfully set\n", name, o.keyValue)

	return nil
}

func (o *getKeyOptions) Print(name, token string) {
	fmt.Printf("%s: %s\n", name, token)
}
