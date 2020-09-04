package kube

import (
	"github.com/baetyl/baetyl-go/v2/log"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	clientset "github.com/baetyl/baetyl-cloud/v2/plugin/kube/client/clientset/versioned"
)

type client struct {
	coreV1       corev1.CoreV1Interface
	customClient clientset.Interface
	cfg          CloudConfig
	aesKey       []byte
	log          *log.Logger
}

// Close Close
func (c *client) Close() error {
	return nil
}

func init() {
	plugin.RegisterFactory("kube", New)
	plugin.RegisterFactory("kubernetes", New)
}

// New New
func New() (plugin.Plugin, error) {
	var cfg CloudConfig
	if err := common.LoadConfig(&cfg); err != nil {
		return nil, err
	}

	kubeConfig, err := func() (*rest.Config, error) {
		if !cfg.Kube.OutCluster {
			return rest.InClusterConfig()
		}
		return clientcmd.BuildConfigFromFlags(
			"", cfg.Kube.ConfigPath)
	}()

	if err != nil {
		return nil, err
	}
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	kubeClient.CoreV1()
	customClient, err := clientset.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	return &client{
		coreV1:       kubeClient.CoreV1(),
		customClient: customClient,
		cfg:          cfg,
		aesKey:       []byte(cfg.Kube.AES.Key),
		log:          log.With(log.Any("plugin", "kube")),
	}, nil
}
