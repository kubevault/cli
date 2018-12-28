package cmds

import (
	"fmt"
	"os"

	engineapi "github.com/kubevault/operator/apis/engine/v1alpha1"
	enginescheme "github.com/kubevault/operator/client/clientset/versioned/scheme"
	enginecs "github.com/kubevault/operator/client/clientset/versioned/typed/engine/v1alpha1"
	engineutil "github.com/kubevault/operator/client/clientset/versioned/typed/engine/v1alpha1/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericclioptions/resource"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

var (
	FilenameOptions resource.FilenameOptions
	ObjectNames     []string
	ResourceName    string
)

var (
	awsApprovedCond = engineapi.AWSAccessKeyRequestCondition{
		Type:           engineapi.AccessApproved,
		Reason:         "KubectlApprove",
		Message:        "This was approved by kubectl vault approve awsaccesskeyrequest",
		LastUpdateTime: metav1.Now(),
	}
)

func NewCmdApprove(clientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "approve",
		Short:             "Approve request",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				ResourceName = args[0]
				ObjectNames = args[1:]
			}

			cfg, err := clientGetter.ToRESTConfig()
			if err != nil {
				Fatal(errors.Wrap(err, "failed to read kubeconfig"))
			}

			namespace, _, err := clientGetter.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				Fatal(err)
			}

			builder := cmdutil.NewFactory(clientGetter).NewBuilder()

			c, err := enginecs.NewForConfig(cfg)
			if err != nil {
				Fatal(err)
			}

			if err := modifyAWSAccessKeyCondition(builder, c, namespace, awsApprovedCond); err != nil {
				Fatal(err)
			}

			fmt.Printf("%s is approved", ResourceName)
			os.Exit(0)
		},
	}

	cmdutil.AddFilenameOptionFlags(cmd, &FilenameOptions, "identifying the resource to update")
	return cmd
}

func modifyAWSAccessKeyCondition(builder *resource.Builder, c enginecs.EngineV1alpha1Interface, namespace string, cond engineapi.AWSAccessKeyRequestCondition) error {
	var found int
	r := builder.
		WithScheme(enginescheme.Scheme, enginescheme.Scheme.PrioritizedVersionsAllGroups()...).
		ContinueOnError().
		FilenameParam(false, &FilenameOptions).
		ResourceNames(engineapi.ResourceAWSAccessKeyRequest, ObjectNames...).
		RequireObject(true).
		Flatten().
		Latest().
		NamespaceParam(namespace).
		Do()
	err := r.Visit(func(info *resource.Info, err error) error {
		if err != nil {
			return err
		}
		obj := info.Object.(*engineapi.AWSAccessKeyRequest)
		_, err2 := engineutil.UpdateAWSAccessKeyRequestStatus(c, obj, func(in *engineapi.AWSAccessKeyRequestStatus) *engineapi.AWSAccessKeyRequestStatus {
			for _, cond := range in.Conditions {
				if cond.Type == cond.Type {
					return in
				}
			}
			in.Conditions = append(in.Conditions, cond)
			return in
		}, EnableStatusSubresource)
		found++

		return err2
	})
	if found == 0 {
		fmt.Println("No resources found")
	}
	return err
}
