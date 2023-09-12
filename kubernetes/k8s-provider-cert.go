package kubernetes

import (
	"github.com/mogenius/punq/logger"

	cmclientset "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeProviderCertManager struct {
	ClientSet    *cmclientset.Clientset
	ClientConfig rest.Config
}

func NewKubeProviderCertManager(contextId *string) *KubeProviderCertManager {
	var kubeProvider *KubeProviderCertManager
	var err error
	if RunsInCluster {
		kubeProvider, err = newKubeProviderCertManagerInCluster(contextId)
	} else {
		kubeProvider, err = newKubeProviderCertManagerLocal()
	}

	if err != nil {
		logger.Log.Errorf("ERROR: %s", err.Error())
	}
	return kubeProvider
}

func newKubeProviderCertManagerLocal() (*KubeProviderCertManager, error) {
	var kubeconfig string = getKubeConfig()

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

func newKubeProviderCertManagerInCluster(contextId *string) (*KubeProviderCertManager, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// CONTEXT SWITCHER
	if contextId != nil {
		config, err = ContextConfigLoader(contextId)
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
