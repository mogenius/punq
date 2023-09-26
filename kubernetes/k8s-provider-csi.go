package kubernetes

import (
	snapClientset "github.com/kubernetes-csi/external-snapshotter/client/v6/clientset/versioned"
	"github.com/mogenius/punq/logger"
	"k8s.io/client-go/rest"
)

type KubeProviderSnapshot struct {
	ClientSet    *snapClientset.Clientset
	ClientConfig rest.Config
}

func NewKubeProviderSnapshot(contextId *string) (*KubeProviderSnapshot, error) {
	var provider *KubeProviderSnapshot
	var err error
	if RunsInCluster {
		provider, err = newKubeProviderCsiInCluster(contextId)
	} else {
		provider, err = newKubeProviderCsiLocal(contextId)
	}

	if err != nil {
		logger.Log.Errorf("ERROR: %s", err.Error())
	}
	return provider, err
}

func newKubeProviderCsiLocal(contextId *string) (*KubeProviderSnapshot, error) {
	config, err := ContextSwitcher(contextId)
	if err != nil {
		return nil, err
	}

	clientSet, errClientSet := snapClientset.NewForConfig(config)
	if errClientSet != nil {
		return nil, errClientSet
	}

	return &KubeProviderSnapshot{
		ClientSet:    clientSet,
		ClientConfig: *config,
	}, nil
}

func newKubeProviderCsiInCluster(contextId *string) (*KubeProviderSnapshot, error) {
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

	clientset, err := snapClientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &KubeProviderSnapshot{
		ClientSet:    clientset,
		ClientConfig: *config,
	}, nil
}
