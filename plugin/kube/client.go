package kube

import (
	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/plugin"
	clientset "github.com/baetyl/baetyl-cloud/plugin/kube/client/clientset/versioned"
	"github.com/baetyl/baetyl-go/log"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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
	plugin.RegisterFactory("kubernetes", New)
}

// New New
func New() (plugin.Plugin, error) {
	var cfg CloudConfig
	if err := common.LoadConfig(&cfg); err != nil {
		return nil, err
	}

	kubeConfig, err := func() (*rest.Config, error) {
		if cfg.Kubernetes.InCluster {
			return rest.InClusterConfig()
		}
		return clientcmd.BuildConfigFromFlags(
			"", cfg.Kubernetes.ConfigPath)
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
		aesKey:       []byte(cfg.Kubernetes.AES.Key),
		log:          log.With(log.Any("plugin", "kube")),
	}, nil
}
