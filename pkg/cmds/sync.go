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
	"kubevault.dev/cli/pkg/token-keys-store/api"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func NewCmdSync(clientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "sync",
		Short:             "sync short cmd",
		Long:              "sync long cmd",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				ResourceName = args[0]
				ObjectNames = args[1:]
			}

			if err := sync(clientGetter); err != nil {
				Fatal(err)
			}

			os.Exit(0)
		},
	}
	return cmd
}

func sync(clientGetter genericclioptions.RESTClientGetter) error {
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
			err2 = syncTokenKeys(obj, kubeClient)
		default:
			err2 = errors.New("unknown/unsupported type")
		}
		return err2
	})
	return err
}

func syncTokenKeys(vs *vaultapi.VaultServer, kubeClient kubernetes.Interface) error {
	if err := syncRootToken(vs, kubeClient); err != nil {
		return err
	}

	if err := syncUnsealKeys(vs, kubeClient); err != nil {
		return err
	}

	fmt.Println("vault root-token & unseal-keys successfully synced")
	return nil
}

func syncRootToken(vs *vaultapi.VaultServer, kubeClient kubernetes.Interface) error {
	ti, err := token_key_store.NewTokenKeyInterface(vs, kubeClient)
	if err != nil {
		return err
	}

	defer func() {
		err = ti.Clean()
	}()

	newKey := ti.NewTokenName()
	_, err = ti.Get(newKey)
	if err == nil {
		return err
	}

	oldKey := ti.OldTokenName()
	value, err := ti.Get(oldKey)
	if err == nil {
		err = ti.Set(newKey, value)
		if err != nil {
			return err
		}

		err = ti.Delete(oldKey)
		return err
	}

	return err
}

func syncUnsealKeys(vs *vaultapi.VaultServer, kubeClient kubernetes.Interface) error {
	ti, err := token_key_store.NewTokenKeyInterface(vs, kubeClient)
	if err != nil {
		return err
	}

	defer func() {
		err = ti.Clean()
	}()

	for i := 0; int64(i) < vs.Spec.Unsealer.SecretShares; i++ {
		fmt.Println("trying for vault-unseal-key: ", i)
		err := syncUnsealKey(i, ti)
		if err != nil {
			return err
		}
	}

	return nil
}

func syncUnsealKey(id int, ti api.TokenKeyInterface) error {
	newKey := ti.NewUnsealKeyName(id)
	_, err := ti.Get(newKey)
	if err == nil {
		return err
	}

	oldKey := ti.OldUnsealKeyName(id)
	value, err := ti.Get(oldKey)
	if err == nil {
		err = ti.Set(newKey, value)
		if err != nil {
			return err
		}

		err = ti.Delete(oldKey)
		return err
	}

	return err
}
