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

type keyOptions struct {
	keyId int
}

func newKeyOptions() *keyOptions {
	return &keyOptions{}
}

func (o *keyOptions) addKeyFlags(fs *pflag.FlagSet) {
	fs.IntVar(&o.keyId, "key-id", o.keyId, "unseal key id")
}

type keyValueOptions struct {
	keyId    int
	keyValue string
}

func newKeyValueOptions() *keyValueOptions {
	return &keyValueOptions{}
}

func (o *keyValueOptions) addKeyValueOptions(fs *pflag.FlagSet) {
	fs.IntVar(&o.keyId, "key-id", o.keyId, "unseal key id")
	fs.StringVar(&o.keyValue, "key-value", o.keyValue, "unseal key value")
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
	o := newKeyOptions()
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

	o.addKeyFlags(cmd.Flags())
	return cmd
}

func NewCmdSetKey(clientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	o := newKeyValueOptions()
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

	o.addKeyValueOptions(cmd.Flags())
	return cmd
}

func NewCmdDeleteKey(clientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	o := newKeyOptions()
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

	o.addKeyFlags(cmd.Flags())
	return cmd
}

func NewCmdListKey(clientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	o := newKeyOptions()
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

func (o *keyOptions) list(clientGetter genericclioptions.RESTClientGetter) error {
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
			err2 = o.listUnsealKey(obj, kubeClient)
		default:
			err2 = errors.New("unknown/unsupported type")
		}
		return err2
	})
	return err
}

func (o *keyOptions) listUnsealKey(vs *vaultapi.VaultServer, kubeClient kubernetes.Interface) error {
	cnt := vs.Spec.Unsealer.SecretShares
	for i := 0; int64(i) < cnt; i++ {
		o.keyId = i
		_ = o.printUnsealKey(vs, kubeClient)
	}

	return nil
}

func (o *keyOptions) get(clientGetter genericclioptions.RESTClientGetter) error {
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
			err2 = o.printUnsealKey(obj, kubeClient)
		default:
			err2 = errors.New("unknown/unsupported type")
		}
		return err2
	})
	return err
}

func (o *keyOptions) printUnsealKey(vs *vaultapi.VaultServer, kubeClient kubernetes.Interface) error {
	ti, err := token_key_store.NewTokenKeyInterface(vs, kubeClient)
	if err != nil {
		return err
	}

	name := ti.NewUnsealKeyName(o.keyId)
	rToken, err := ti.Get(name)
	if err == nil {
		o.Print(name, rToken)
		return nil
	}

	name = ti.OldUnsealKeyName(o.keyId)
	rToken, err = ti.Get(name)
	if err == nil {
		o.Print(name, rToken)
		return nil
	}

	return err
}

func (o *keyOptions) del(clientGetter genericclioptions.RESTClientGetter) error {
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

func (o *keyOptions) deleteUnsealKey(vs *vaultapi.VaultServer, kubeClient kubernetes.Interface) error {
	ti, err := token_key_store.NewTokenKeyInterface(vs, kubeClient)
	if err != nil {
		return err
	}

	name := ti.NewUnsealKeyName(o.keyId)
	err = ti.Delete(name)
	if err == nil {
		fmt.Printf("unseal-key with name %s successfully deleted\n", name)
		return nil
	}

	name = ti.OldUnsealKeyName(o.keyId)
	err = ti.Delete(name)
	if err == nil {
		fmt.Printf("unseal-key with name %s successfully deleted\n", name)
		return nil
	}

	return err
}

func (o *keyValueOptions) set(clientGetter genericclioptions.RESTClientGetter) error {
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

func (o *keyValueOptions) setUnsealKey(vs *vaultapi.VaultServer, kubeClient kubernetes.Interface) error {
	ti, err := token_key_store.NewTokenKeyInterface(vs, kubeClient)
	if err != nil {
		return err
	}

	name := ti.NewUnsealKeyName(o.keyId)
	err = ti.Set(name, o.keyValue)
	if err == nil {
		fmt.Printf("unseal-key with name %s, value %s successfully set\n", name, o.keyValue)
		return nil
	}

	name = ti.OldUnsealKeyName(o.keyId)
	err = ti.Set(name, o.keyValue)
	if err == nil {
		fmt.Printf("unseal-key with name %s, value %s successfully set\n", name, o.keyValue)
		return nil
	}

	return err
}

func (o *keyOptions) Print(name, token string) {
	fmt.Printf("%s: %s\n", name, token)
}
