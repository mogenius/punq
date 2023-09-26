package kubernetes

import (
	"github.com/mogenius/punq/logger"

	cmclientset "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	"k8s.io/client-go/rest"
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
		kubeProvider, err = newKubeProviderCertManagerLocal(contextId)
	}

	if err != nil {
		logger.Log.Fatalf("ERROR: %s", err.Error())
	}
	return kubeProvider
}

func newKubeProviderCertManagerLocal(contextId *string) (*KubeProviderCertManager, error) {
	config := ContextSwitcher(contextId)

	cmClientset, err := cmclientset.NewForConfig(config)
	if err != nil {
		logger.Log.Panicf("Failed to create cert-manager clientset: %v\n", err)
	}

	return &KubeProviderCertManager{
		ClientSet:    cmClientset,
		ClientConfig: *config,
	}, nil
}

func newKubeProviderCertManagerInCluster(contextId *string) (*KubeProviderCertManager, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	if contextId != nil {
		config = ContextSwitcher(contextId)
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
