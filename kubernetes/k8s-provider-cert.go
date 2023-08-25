package kubernetes

import (
	"path/filepath"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	cmclientset "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type KubeProviderCertManager struct {
	ClientSet    *cmclientset.Clientset
	ClientConfig rest.Config
}

func NewKubeProviderCertManager() *KubeProviderCertManager {
	var kubeProvider *KubeProviderCertManager
	var err error
	if utils.CONFIG.Kubernetes.RunInCluster {
		kubeProvider, err = newKubeProviderCertManagerInCluster()
	} else {
		kubeProvider, err = newKubeProviderCertManagerLocal()
	}

	if err != nil {
		logger.Log.Errorf("ERROR: %s", err.Error())
	}
	return kubeProvider
}

func newKubeProviderCertManagerLocal() (*KubeProviderCertManager, error) {
	var kubeconfig string = ""
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	restConfig, errConfig := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if errConfig != nil {
		panic(errConfig.Error())
	}

	cmClientset, err := cmclientset.NewForConfig(restConfig)
	if err != nil {
		logger.Log.Panicf("Failed to create cert-manager clientset: %v\n", err)
	}

	return &KubeProviderCertManager{
		ClientSet:    cmClientset,
		ClientConfig: *restConfig,
	}, nil
}

func newKubeProviderCertManagerInCluster() (*KubeProviderCertManager, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := cmclientset.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &KubeProviderCertManager{
		ClientSet:    clientset,
		ClientConfig: *config,
	}, nil
}
