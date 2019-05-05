package cmds

import (
	"fmt"
	"os"

	dbapi "github.com/kubedb/apimachinery/apis/authorization/v1alpha1"
	dbcs "github.com/kubedb/apimachinery/client/clientset/versioned/typed/authorization/v1alpha1"
	dbutil "github.com/kubedb/apimachinery/client/clientset/versioned/typed/authorization/v1alpha1/util"
	engineapi "github.com/kubevault/operator/apis/engine/v1alpha1"
	enginecs "github.com/kubevault/operator/client/clientset/versioned/typed/engine/v1alpha1"
	engineutil "github.com/kubevault/operator/client/clientset/versioned/typed/engine/v1alpha1/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericclioptions/resource"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

var (
	FilenameOptions resource.FilenameOptions
	ObjectNames     []string
	ResourceName    string
)

var (
	awsApprovedCond = engineapi.AWSAccessKeyRequestCondition{
		Type:    engineapi.AccessApproved,
		Reason:  "KubectlApprove",
		Message: "This was approved by kubectl vault approve awsaccesskeyrequest",
	}

	dbApprovedCond = dbapi.DatabaseAccessRequestCondition{
		Type:    dbapi.AccessApproved,
		Reason:  "KubectlApprove",
		Message: "This was approved by kubectl vault approve databaseaccessrequest",
	}

	gcpApprovedCond = engineapi.GCPAccessKeyRequestCondition{
		Type:    engineapi.AccessApproved,
		Reason:  "KubectlApprove",
		Message: "This was approved by kubectl vault approve gcpaccesskeyrequest",
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
	case dbapi.ResourceDatabaseAccessRequest, dbapi.ResourceDatabaseAccessRequests:
		resourceName = dbapi.ResourceDatabaseAccessRequest
	case engineapi.ResourceGCPAccessKeyRequest, engineapi.ResourceGCPAccessKeyRequests:
		resourceName = engineapi.ResourceGCPAccessKeyRequest
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

	dbClient, err := dbcs.NewForConfig(cfg)
	if err != nil {
		return err
	}

	var found int
	r := builder.
		WithScheme(clientsetscheme.Scheme, clientsetscheme.Scheme.PrioritizedVersionsAllGroups()...).
		ContinueOnError().
		FilenameParam(false, &FilenameOptions).
		ResourceNames(resourceName, ObjectNames...).
		RequireObject(true).
		Flatten().
		Latest().
		NamespaceParam(namespace).
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
			err2 = UpdateAWSAccessKeyRequestCondition(engineClient, obj, cond)
		case *dbapi.DatabaseAccessRequest:
			obj := info.Object.(*dbapi.DatabaseAccessRequest)
			cond := dbDeniedCond
			if isApproveReq {
				cond = dbApprovedCond
			}
			err2 = UpdateDBAccessRequestCondition(dbClient, obj, cond)
		case *engineapi.GCPAccessKeyRequest:
			obj := info.Object.(*engineapi.GCPAccessKeyRequest)
			cond := gcpDeniedCond
			if isApproveReq {
				cond = gcpApprovedCond
			}
			err2 = UpdateGCPAccessKeyRequest(engineClient, obj, cond)
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

func UpdateAWSAccessKeyRequestCondition(c enginecs.EngineV1alpha1Interface, awsAKR *engineapi.AWSAccessKeyRequest, cond engineapi.AWSAccessKeyRequestCondition) error {
	_, err := engineutil.UpdateAWSAccessKeyRequestStatus(c, awsAKR, func(in *engineapi.AWSAccessKeyRequestStatus) *engineapi.AWSAccessKeyRequestStatus {
		for _, cond := range in.Conditions {
			if cond.Type == cond.Type {
				return in
			}
		}
		cond.LastUpdateTime = metav1.Now()
		in.Conditions = append(in.Conditions, cond)
		return in
	}, EnableStatusSubresource)
	return err
}

func UpdateDBAccessRequestCondition(c dbcs.AuthorizationV1alpha1Interface, dbAR *dbapi.DatabaseAccessRequest, cond dbapi.DatabaseAccessRequestCondition) error {
	_, err := dbutil.UpdateDatabaseAccessRequestStatus(c, dbAR, func(in *dbapi.DatabaseAccessRequestStatus) *dbapi.DatabaseAccessRequestStatus {
		for _, cond := range in.Conditions {
			if cond.Type == cond.Type {
				return in
			}
		}
		cond.LastUpdateTime = metav1.Now()
		in.Conditions = append(in.Conditions, cond)
		return in
	}, EnableStatusSubresource)
	return err
}

func UpdateGCPAccessKeyRequest(c enginecs.EngineV1alpha1Interface, gcpAKR *engineapi.GCPAccessKeyRequest, cond engineapi.GCPAccessKeyRequestCondition) error {
	_, err := engineutil.UpdateGCPAccessKeyRequestStatus(c, gcpAKR, func(in *engineapi.GCPAccessKeyRequestStatus) *engineapi.GCPAccessKeyRequestStatus {
		for _, cond := range in.Conditions {
			if cond.Type == cond.Type {
				return in
			}
		}
		cond.LastUpdateTime = metav1.Now()
		in.Conditions = append(in.Conditions, cond)
		return in
	}, EnableStatusSubresource)
	return err
}
