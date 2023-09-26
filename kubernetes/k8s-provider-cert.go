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

func NewKubeProviderCertManager(contextId *string) (*KubeProviderCertManager, error) {
	var provider *KubeProviderCertManager
	var err error
	if RunsInCluster {
		provider, err = newKubeProviderCertManagerInCluster(contextId)
	} else {
		provider, err = newKubeProviderCertManagerLocal(contextId)
	}

	if err != nil {
		logger.Log.Errorf("ERROR: %s", err.Error())
	}
	return provider, err
}

func newKubeProviderCertManagerLocal(contextId *string) (*KubeProviderCertManager, error) {
	config, err := ContextSwitcher(contextId)
	if err != nil {
		return nil, err
	}

	cmClientset, err := cmclientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &KubeProviderCertManager{
		ClientSet:    cmClientset,
		ClientConfig: *config,
	}, nil
}

func newKubeProviderCertManagerInCluster(contextId *string) (*KubeProviderCertManager, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	if contextId != nil {
		config, err = ContextSwitcher(contextId)
		if err != nil {
			return nil, err
		}
	}

	clientset, err := cmclientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &KubeProviderCertManager{
		ClientSet:    clientset,
		ClientConfig: *config,
	}, nil
}
