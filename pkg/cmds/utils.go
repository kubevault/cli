package cmds

import (
	"fmt"
	"os"

	"github.com/golang/glog"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"kmodules.xyz/client-go/discovery"
)

func Fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func GetDefaultValueForStatusSubresource(clientGetter genericclioptions.RESTClientGetter) bool {
	cfg, err := clientGetter.ToRESTConfig()
	if err != nil {
		glog.Errorln("failed to get rest config", err)
		return false
	}

	kc, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Errorln("failed to create Kubernetes client from rest config", err)
		return false
	}

	resp, err := discovery.CheckAPIVersion(kc.Discovery(), ">=1.11.0")
	if err != nil {
		glog.Errorln("failed to check Kubernetes api version", err)
		return false
	}
	return resp
}
