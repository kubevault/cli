module github.com/kubevault/cli

go 1.12

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/MakeNowJust/heredoc v0.0.0-20170808103936-bb23615498cd // indirect
	github.com/appscode/go v0.0.0-20190424183524-60025f1135c9
	github.com/emicklei/go-restful v2.9.5+incompatible // indirect
	github.com/exponent-io/jsonpath v0.0.0-20151013193312-d6023ce2651d // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/kubedb/apimachinery v0.0.0-20190506191700-871d6b5d30ee
	github.com/kubevault/operator v0.0.0-20190522102911-da03839fcb8c
	github.com/mattn/go-isatty v0.0.8 // indirect
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/spf13/cobra v0.0.4
	k8s.io/apimachinery v0.0.0-20190508063446-a3da69d3723c
	k8s.io/cli-runtime v0.0.0-20190508184404-b26560c459bd
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/component-base v0.0.0-20190509023737-8de8845fb642
	k8s.io/kubernetes v1.14.2
	kmodules.xyz/client-go v0.0.0-20190508091620-0d215c04352f
	kmodules.xyz/custom-resources v0.0.0-20190508103408-464e8324c3ec
)

replace (
	github.com/graymeta/stow => github.com/appscode/stow v0.0.0-20190506085026-ca5baa008ea3
	gopkg.in/robfig/cron.v2 => github.com/appscode/cron v0.0.0-20170717094345-ca60c6d796d4
	k8s.io/api => k8s.io/api v0.0.0-20190313235455-40a48860b5ab
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190315093550-53c4693659ed
	k8s.io/apimachinery => github.com/kmodules/apimachinery v0.0.0-20190508045248-a52a97a7a2bf
	k8s.io/apiserver => github.com/kmodules/apiserver v0.0.0-20190508082252-8397d761d4b5
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20190314001948-2899ed30580f
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20190314002645-c892ea32361a
	k8s.io/component-base => k8s.io/component-base v0.0.0-20190314000054-4a91899592f4
	k8s.io/klog => k8s.io/klog v0.3.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20190314000639-da8327669ac5
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20190228160746-b3a7cee44a30
	k8s.io/metrics => k8s.io/metrics v0.0.0-20190314001731-1bd6a4002213
	k8s.io/utils => k8s.io/utils v0.0.0-20190221042446-c2654d5206da
)
