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
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

type mergeSecretsOptions struct {
	overwrite bool
	src       string
	dst       string
}

func newMergeSecretsOptions() *mergeSecretsOptions {
	return &mergeSecretsOptions{}
}

func (o *mergeSecretsOptions) addMergeSecretsFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&o.overwrite, "overwrite-keys", o.overwrite, "will overwrite keys in destination if set to true.")
	fs.StringVar(&o.src, "src", o.src, "source secret.")
	fs.StringVar(&o.dst, "dst", o.dst, "destination secret.")
}

func NewCmdMergeSecrets(clientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	o := newMergeSecretsOptions()

	cmd := &cobra.Command{
		Use:   "merge-secrets",
		Short: "merge-secrets merges secrets",
		Long: `merge-secrets command merges two given secrets. Both the src & dst secrets must exist for successful merge operation.

Example: 
 # merge two secret name1 & name2 from ns1 & ns2 namespaces respectively
 $ kubectl vault merge-secrets --src=<ns1>/<name1> --dst=<ns2>/<name2>

 # --overwrite-keys flag will overwrite keys in destination if set to true.
 $ kubectl vault merge-secrets --src=<ns1>/<name1> --dst=<ns2>/<name2> --overwrite-keys=true
`,

		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				ResourceName = args[0]
				ObjectNames = args[1:]
			}

			if err := o.merge(clientGetter); err != nil {
				Fatal(err)
			}
			klog.Infoln("secrets successfully merged")
			os.Exit(0)
		},
	}

	o.addMergeSecretsFlags(cmd.Flags())
	return cmd
}

func (o *mergeSecretsOptions) merge(clientGetter genericclioptions.RESTClientGetter) error {
	cfg, err := clientGetter.ToRESTConfig()
	if err != nil {
		return err
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return err
	}

	srcNs, srcName := split(o.src)
	dstNs, dstName := split(o.dst)

	srcSecret, err := kubeClient.CoreV1().Secrets(srcNs).Get(context.TODO(), srcName, metav1.GetOptions{})
	if err != nil {
		fmt.Println("src secret doesn't exist")
		return err
	}

	dstSecret, err := kubeClient.CoreV1().Secrets(dstNs).Get(context.TODO(), dstName, metav1.GetOptions{})
	if err != nil {
		fmt.Println("dst secret doesn't exist")
		return err
	}

	for key, value := range srcSecret.Data {
		if _, ok := dstSecret.Data[key]; !ok || (ok && o.overwrite) {
			dstSecret.Data[key] = value
		}
	}

	_, err = kubeClient.CoreV1().Secrets(dstNs).Update(context.TODO(), dstSecret, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func split(s string) (string, string) {
	var namespace, name string
	res := strings.Split(s, "/")
	if len(res) != 2 {
		namespace = metav1.NamespaceDefault
		name = res[0]
	} else {
		namespace = res[0]
		name = res[1]
	}

	return namespace, name
}
