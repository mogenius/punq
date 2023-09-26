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

func NewKubeProviderSnapshot(contextId *string) *KubeProviderSnapshot {
	var kubeProvider *KubeProviderSnapshot
	var err error
	if RunsInCluster {
		kubeProvider, err = newKubeProviderCsiInCluster(contextId)
	} else {
		kubeProvider, err = newKubeProviderCsiLocal(contextId)
	}

	if err != nil {
		logger.Log.Fatalf("ERROR: %s", err.Error())
	}
	return kubeProvider
}

func newKubeProviderCsiLocal(contextId *string) (*KubeProviderSnapshot, error) {
	config := ContextSwitcher(contextId)

	clientSet, errClientSet := snapClientset.NewForConfig(config)
	if errClientSet != nil {
		panic(errClientSet.Error())
	}

	return &KubeProviderSnapshot{
		ClientSet:    clientSet,
		ClientConfig: *config,
	}, nil
}

func newKubeProviderCsiInCluster(contextId *string) (*KubeProviderSnapshot, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	if contextId != nil {
		config = ContextSwitcher(contextId)
	}

	clientset, err := snapClientset.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &KubeProviderSnapshot{
		ClientSet:    clientset,
		ClientConfig: *config,
	}, nil
}
