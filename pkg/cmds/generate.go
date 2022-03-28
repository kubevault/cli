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
	"encoding/json"
	"fmt"
	"os"
	"strings"

	enginecs "kubevault.dev/apimachinery/client/clientset/versioned/typed/engine/v1alpha1"
	vaultcs "kubevault.dev/apimachinery/client/clientset/versioned/typed/kubevault/v1alpha2"
	policycs "kubevault.dev/apimachinery/client/clientset/versioned/typed/policy/v1alpha1"
	"kubevault.dev/cli/pkg/generate"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	secretsstore "sigs.k8s.io/secrets-store-csi-driver/apis/v1alpha1"
	"sigs.k8s.io/yaml"
)

const (
	ResourceKindSecretProviderClass = "secretproviderclass"
)

type generateOption struct {
	secretRoleBinding string
	vaultRole         string
	keys              map[string]string
	output            string
	vaultCACertPath   string
}

func NewOptions() *generateOption {
	return &generateOption{}
}

type SecretProviderClassOptions struct {
	options    *generateOption
	apiVersion string
	kind       string
	provider   secretsstore.Provider
	namespace  string
	name       string
	vsURL      string
	roleName   string
}

func NewSecretProviderClassOptions(op *generateOption, namespace, name string) *SecretProviderClassOptions {
	return &SecretProviderClassOptions{
		options:    op,
		apiVersion: fmt.Sprintf("%s/v1alpha1", secretsstore.GroupName),
		kind:       "SecretProviderClass",
		provider:   "vault",
		namespace:  namespace,
		name:       name,
	}
}

func (o *generateOption) AddSecretProviderClassFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.vaultRole, "vaultrole", "r", o.vaultRole, "vault role. RoleKind/name")
	fs.StringVarP(&o.secretRoleBinding, "secretrolebinding", "b", o.secretRoleBinding, "secret role binding. namespace/name")
	fs.StringToStringVar(&o.keys, "keys", o.keys, "Key/Value map used to store the keys to read and their mapping keys. secretKey=objectName")
	fs.StringVarP(&o.output, "output", "o", o.output, "output format yaml/json. default to yaml")
	fs.StringVarP(&o.vaultCACertPath, "vault-ca-cert-path", "p", o.vaultCACertPath, "vault CA cert path in secret provider, default to Insecure mode.")
}

func NewCmdGenerate(clientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	o := NewOptions()

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate secretproviderclass",
		Long: "Generate secretproviderclass from secretrolebinding. Provide flags secretrolebinding, role and keys to mount.\n\n" +
			"See more about Secrets-Store-CSI-Driver and the usage of SecretProviderClass:\n Link: https://secrets-store-csi-driver.sigs.k8s.io/concepts.html#secretproviderclass \n\n" +
			"secretrolebinding needs to be created and successful beforehand.\nProvided roles must be in the seretrolebinding and " +
			"provided keys must be available for the RoleKind.\n" +
			"Output format can be yaml or json, defaults to yaml\n\n" +
			"Examples:\n" +
			" # Generate secretproviderclass with name <name1> and namespace <ns1>\n" +
			" # secretrolebinding with namespace <ns2> and name <name2>\n" +
			" # vaultrole kind MongoDBRole and name <name3>\n" +
			" # keys to mount <secretKey> and it's mapping name <objectName> \n" +
			"\n $ kubectl vault generate secretproviderclass <name1> -n <ns1> \\\n  " +
			"--secretrolebinding=<ns2>/<name2> \\\n  --vaultrole=MongoDBRole/<name3> \\\n  --keys <secretKey>=<objectName> -o yaml\n\n" +
			" # Generate secretproviderclass for the MongoDB username and password\n" +
			" $ kubectl vault generate secretproviderclass mongo-secret-provider -n test     " +
			" \\\n  --secretrolebinding=dev/secret-r-binding \\\n  --vaultrole=MongoDBRole/mongo-role \\\n  --keys username=mongo-user --keys password=mongo-pass -o yaml\n",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				ResourceName = args[0]
				ObjectNames = args[1:]
			}

			if err := o.generate(clientGetter); err != nil {
				Fatal(err)
			}

			os.Exit(0)
		},
	}
	cmdutil.AddFilenameOptionFlags(cmd, &FilenameOptions, "identifying the resource to update")
	o.AddSecretProviderClassFlags(cmd.Flags())
	return cmd
}

func (o *generateOption) generate(clientGetter genericclioptions.RESTClientGetter) error {
	var resourceName string
	switch strings.ToLower(ResourceName) {
	case ResourceKindSecretProviderClass:
		resourceName = ResourceKindSecretProviderClass
	default:
		return errors.New(fmt.Sprintf("unknown/unsupported resource %s", ResourceName))
	}

	if len(resourceName) == 0 {
		return errors.New("resourceName empty")
	}
	if len(ObjectNames) == 0 {
		return errors.New("secretproviderclass name not provided")
	}

	namespace, _, err := clientGetter.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	cfg, err := clientGetter.ToRESTConfig()
	if err != nil {
		return errors.Wrap(err, "failed to read kubeconfig")
	}

	engineClient, vaultClient, policyClient, kubeClient, err := initClients(cfg)
	if err != nil {
		return err
	}

	spc := NewSecretProviderClassOptions(o, namespace, ObjectNames[0])

	objectsList, err := spc.generateSecretObjects(engineClient, vaultClient, policyClient, kubeClient)
	if err != nil {
		return err
	}

	if err = spc.generateSecretProviderClass(objectsList); err != nil {
		return err
	}

	return nil
}

func (s *SecretProviderClassOptions) generateSecretObjects(engineClient *enginecs.EngineV1alpha1Client, vaultClient *vaultcs.KubevaultV1alpha2Client, policyClient *policycs.PolicyV1alpha1Client, kubeClient *kubernetes.Clientset) (string, error) {
	if engineClient == nil || vaultClient == nil || policyClient == nil || kubeClient == nil {
		return "", errors.New("engineClient/vaultClient/policyClient/kubeClient is nil")
	}

	var srbNs, srbName string
	srb := strings.Split(s.options.secretRoleBinding, "/")
	if len(srb) != 2 {
		srbNs = metav1.NamespaceDefault
		srbName = srb[0]
	} else {
		srbNs = srb[0]
		srbName = srb[1]
	}

	srbObj, err := engineClient.SecretRoleBindings(srbNs).Get(context.TODO(), srbName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	role := strings.Split(s.options.vaultRole, "/")
	if len(role) != 2 {
		return "", errors.New("vault role/name not provided")
	}

	roleAvailable := false
	for _, srbRole := range srbObj.Spec.Roles {
		if srbRole.Kind == role[0] && srbRole.Name == role[1] {
			roleAvailable = true
		}
	}

	if !roleAvailable {
		return "", errors.Errorf("%s/%s not found in secretrolebinding", role[0], role[1])
	}

	gen, err := generate.NewGenerator(role, srbObj, s.options.keys, engineClient, vaultClient, policyClient, kubeClient)
	if err != nil {
		return "", err
	}

	address, err := gen.GetVaultServerURL()
	if err != nil {
		return "", err
	}
	s.vsURL = address

	vaultRoleName, err := gen.GetVaultRoleName()
	if err != nil {
		return "", err
	}
	s.roleName = vaultRoleName

	return gen.Generate()
}

func (s *SecretProviderClassOptions) generateSecretProviderClass(objectsList string) error {
	spc := &secretsstore.SecretProviderClass{
		TypeMeta: metav1.TypeMeta{
			Kind:       s.kind,
			APIVersion: s.apiVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.name,
			Namespace: s.namespace,
		},
		Spec: secretsstore.SecretProviderClassSpec{
			Provider: s.provider,
			Parameters: map[string]string{
				"vaultAddress": s.vsURL,
				"roleName":     s.roleName,
				"objects":      objectsList,
			},
		},
	}

	if len(s.options.vaultCACertPath) != 0 {
		if !strings.HasPrefix(s.vsURL, "https:") {
			return errors.New("VaultServer isn't secure with SSL, vaultCACertPath isn't supported")
		}
		spc.Spec.Parameters["vaultCACertPath"] = s.options.vaultCACertPath
		spc.Spec.Parameters["vaultSkipTLSVerify"] = "false"
	} else {
		spc.Spec.Parameters["vaultSkipTLSVerify"] = "true"
	}

	jsonData, err := json.MarshalIndent(&spc, "", "\t")
	if err != nil {
		return errors.Errorf("Error while Marshaling to yaml with %s", err.Error())
	}

	if s.options.output == "json" {
		fmt.Println(string(jsonData))
		return nil
	}

	yamlData, err := yaml.JSONToYAML(jsonData)
	if err != nil {
		return errors.Errorf("Error while Marshaling to yaml with %s", err.Error())
	}

	fmt.Println(string(yamlData))

	return nil
}

func initClients(cfg *rest.Config) (*enginecs.EngineV1alpha1Client, *vaultcs.KubevaultV1alpha2Client, *policycs.PolicyV1alpha1Client, *kubernetes.Clientset, error) {
	engineClient, err := enginecs.NewForConfig(cfg)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	vaultClient, err := vaultcs.NewForConfig(cfg)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	policyClient, err := policycs.NewForConfig(cfg)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return engineClient, vaultClient, policyClient, kubeClient, nil
}
