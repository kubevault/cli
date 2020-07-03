/*
Copyright The KubeVault Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

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

	engineapi "kubevault.dev/operator/apis/engine/v1alpha1"
	enginecs "kubevault.dev/operator/client/clientset/versioned/typed/engine/v1alpha1"
	engineutil "kubevault.dev/operator/client/clientset/versioned/typed/engine/v1alpha1/util"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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
	awsApprovedCond = kmapi.Condition{
		Type:    kmapi.ConditionRequestApproved,
		Status:  kmapi.ConditionTrue,
		Reason:  "KubectlApprove",
		Message: "This was approved by: kubectl vault approve awsaccesskeyrequest",
	}

	dbApprovedCond = kmapi.Condition{
		Type:    kmapi.ConditionRequestApproved,
		Status:  kmapi.ConditionTrue,
		Reason:  "KubectlApprove",
		Message: "This was approved by: kubectl vault approve databaseaccessrequest",
	}

	gcpApprovedCond = kmapi.Condition{
		Type:    kmapi.ConditionRequestApproved,
		Status:  kmapi.ConditionTrue,
		Reason:  "KubectlApprove",
		Message: "This was approved by: kubectl vault approve gcpaccesskeyrequest",
	}

	azureApprovedCond = kmapi.Condition{
		Type:    kmapi.ConditionRequestApproved,
		Status:  kmapi.ConditionTrue,
		Reason:  "KubectlApprove",
		Message: "This was approved by: kubectl vault approve azureaccesskeyrequest",
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
				fmt.Println("approved")
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
	case engineapi.ResourceAWSAccessKeyRequest, engineapi.ResourceAWSAccessKeyRequests:
		resourceName = engineapi.ResourceAWSAccessKeyRequest
	case engineapi.ResourceDatabaseAccessRequest, engineapi.ResourceDatabaseAccessRequests:
		resourceName = engineapi.ResourceDatabaseAccessRequest
	case engineapi.ResourceGCPAccessKeyRequest, engineapi.ResourceGCPAccessKeyRequests:
		resourceName = engineapi.ResourceGCPAccessKeyRequest
	case engineapi.ResourceAzureAccessKeyRequest, engineapi.ResourceAzureAccessKeyRequests:
		resourceName = engineapi.ResourceAzureAccessKeyRequest
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
		case *engineapi.AWSAccessKeyRequest:
			obj := info.Object.(*engineapi.AWSAccessKeyRequest)
			cond := awsDeniedCond
			if isApproveReq {
				cond = awsApprovedCond
			}
			err2 = UpdateAWSAccessKeyRequestCondition(engineClient, obj.ObjectMeta, cond, isApproveReq)
		case *engineapi.DatabaseAccessRequest:
			obj := info.Object.(*engineapi.DatabaseAccessRequest)
			cond := dbDeniedCond
			if isApproveReq {
				cond = dbApprovedCond
			}
			err2 = UpdateDBAccessRequestCondition(engineClient, obj.ObjectMeta, cond, isApproveReq)
		case *engineapi.GCPAccessKeyRequest:
			obj := info.Object.(*engineapi.GCPAccessKeyRequest)
			cond := gcpDeniedCond
			if isApproveReq {
				cond = gcpApprovedCond
			}
			err2 = UpdateGCPAccessKeyRequest(engineClient, obj.ObjectMeta, cond, isApproveReq)
		case *engineapi.AzureAccessKeyRequest:
			obj := info.Object.(*engineapi.AzureAccessKeyRequest)
			cond := azureDeniedCond
			if isApproveReq {
				cond = azureApprovedCond
			}
			err2 = UpdateAzureAccessKeyRequest(engineClient, obj.ObjectMeta, cond, isApproveReq)
		default:
			err2 = errors.New("unknown/unsupported type")
		}
		found++
		return err2
	})
	if found == 0 {
		fmt.Println("No resources found")
	}
	return err
}

func UpdateAWSAccessKeyRequestCondition(c enginecs.EngineV1alpha1Interface, awsAKR metav1.ObjectMeta, cond kmapi.Condition, isApproveReq bool) error {
	_, err := engineutil.UpdateAWSAccessKeyRequestStatus(context.TODO(), c, awsAKR, func(in *engineapi.AWSAccessKeyRequestStatus) *engineapi.AWSAccessKeyRequestStatus {
		cond.LastTransitionTime = metav1.Now()
		in.Conditions = kmapi.SetCondition(in.Conditions, cond)
		if isApproveReq {
			in.Phase = engineapi.RequestStatusPhaseApproved
		} else {
			in.Phase = engineapi.RequestStatusPhaseDenied
		}
		return in
	}, metav1.UpdateOptions{})
	return err
}

func UpdateDBAccessRequestCondition(c enginecs.EngineV1alpha1Interface, dbAR metav1.ObjectMeta, cond kmapi.Condition, isApproveReq bool) error {
	_, err := engineutil.UpdateDatabaseAccessRequestStatus(context.TODO(), c, dbAR, func(in *engineapi.DatabaseAccessRequestStatus) *engineapi.DatabaseAccessRequestStatus {
		cond.LastTransitionTime = metav1.Now()
		in.Conditions = kmapi.SetCondition(in.Conditions, cond)
		if isApproveReq {
			in.Phase = engineapi.RequestStatusPhaseApproved
		} else {
			in.Phase = engineapi.RequestStatusPhaseDenied
		}
		return in
	}, metav1.UpdateOptions{})
	return err
}

func UpdateGCPAccessKeyRequest(c enginecs.EngineV1alpha1Interface, gcpAKR metav1.ObjectMeta, cond kmapi.Condition, isApproveReq bool) error {
	_, err := engineutil.UpdateGCPAccessKeyRequestStatus(context.TODO(), c, gcpAKR, func(in *engineapi.GCPAccessKeyRequestStatus) *engineapi.GCPAccessKeyRequestStatus {
		cond.LastTransitionTime = metav1.Now()
		in.Conditions = kmapi.SetCondition(in.Conditions, cond)
		if isApproveReq {
			in.Phase = engineapi.RequestStatusPhaseApproved
		} else {
			in.Phase = engineapi.RequestStatusPhaseDenied
		}
		return in
	}, metav1.UpdateOptions{})
	return err
}

func UpdateAzureAccessKeyRequest(c enginecs.EngineV1alpha1Interface, azureAKR metav1.ObjectMeta, cond kmapi.Condition, isApproveReq bool) error {
	_, err := engineutil.UpdateAzureAccessKeyRequestStatus(context.TODO(), c, azureAKR, func(in *engineapi.AzureAccessKeyRequestStatus) *engineapi.AzureAccessKeyRequestStatus {
		cond.LastTransitionTime = metav1.Now()
		in.Conditions = kmapi.SetCondition(in.Conditions, cond)
		if isApproveReq {
			in.Phase = engineapi.RequestStatusPhaseApproved
		} else {
			in.Phase = engineapi.RequestStatusPhaseDenied
		}
		return in
	}, metav1.UpdateOptions{})
	return err
}
