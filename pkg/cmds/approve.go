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

	engineapi "kubevault.dev/apimachinery/apis/engine/v1alpha1"
	enginecs "kubevault.dev/apimachinery/client/clientset/versioned/typed/engine/v1alpha1"
	engineutil "kubevault.dev/apimachinery/client/clientset/versioned/typed/engine/v1alpha1/util"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	kmapi "kmodules.xyz/client-go/api/v1"
)

var (
	FilenameOptions resource.FilenameOptions
	ObjectNames     []string
	ResourceName    string
)

var (
	secretAccessApprovedCond = kmapi.Condition{
		Type:    kmapi.ConditionRequestApproved,
		Status:  core.ConditionTrue,
		Reason:  "KubectlApprove",
		Message: "This was approved by: kubectl vault approve secretaccessrequest",
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

			if err := modifyStatusCondition(clientGetter, true); err != nil {
				Fatal(err)
			} else {
				fmt.Println("Approved")
			}
			os.Exit(0)
		},
	}

	cmdutil.AddFilenameOptionFlags(cmd, &FilenameOptions, "identifying the resource to update")
	return cmd
}

func modifyStatusCondition(clientGetter genericclioptions.RESTClientGetter, isApproveReq bool) error {
	var resourceName string
	switch ResourceName {
	case engineapi.ResourceSecretAccessRequest, engineapi.ResourceSecretAccessRequests:
		resourceName = engineapi.ResourceSecretAccessRequest
	case "":
		resourceName = ""
	default:
		return errors.New("unknown/unsupported resource")
	}

	cfg, err := clientGetter.ToRESTConfig()
	if err != nil {
		return errors.Wrap(err, "failed to read kubeconfig")
	}

	namespace, _, err := clientGetter.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	builder := cmdutil.NewFactory(clientGetter).NewBuilder()

	engineClient, err := enginecs.NewForConfig(cfg)
	if err != nil {
		return err
	}

	var found int
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
		case *engineapi.SecretAccessRequest:
			obj := info.Object.(*engineapi.SecretAccessRequest)
			found++
			cond := secretAccessDeniedCond
			if isApproveReq {
				cond = secretAccessApprovedCond
			}
			if cond == secretAccessDeniedCond && kmapi.IsConditionTrue(obj.Status.Conditions, kmapi.ConditionRequestApproved) {
				return errors.New("Failed to update request status to 'Deny'\nRequest already Approved.")
			}
			err2 = UpdateSecretAccessRequestCondition(engineClient, obj.ObjectMeta, cond)
		default:
			err2 = errors.New("unknown/unsupported type")
		}
		return err2
	})
	if found == 0 {
		fmt.Println("No resources found")
	}
	return err
}

func UpdateSecretAccessRequestCondition(c enginecs.EngineV1alpha1Interface, req metav1.ObjectMeta, cond kmapi.Condition) error {
	_, err := engineutil.UpdateSecretAccessRequestStatus(context.TODO(), c, req, func(in *engineapi.SecretAccessRequestStatus) *engineapi.SecretAccessRequestStatus {
		cond.LastTransitionTime = metav1.Now()
		in.Conditions = kmapi.SetCondition(in.Conditions, cond)
		return in
	}, metav1.UpdateOptions{})
	return err
}
