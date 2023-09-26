package kubernetes

import (
	"github.com/mogenius/punq/logger"
	"k8s.io/client-go/rest"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

type KubeProviderMetrics struct {
	ClientSet    *metricsv.Clientset
	ClientConfig rest.Config
}

func NewKubeProviderMetrics(contextId *string) *KubeProviderMetrics {
	var kubeProvider *KubeProviderMetrics
	var err error
	if RunsInCluster {
		kubeProvider, err = newKubeProviderMetricsInCluster(contextId)
	} else {
		kubeProvider, err = newKubeProviderMetricsLocal(contextId)
	}

	if err != nil {
		logger.Log.Fatalf("ERROR: %s", err.Error())
	}
	return kubeProvider
}

func newKubeProviderMetricsLocal(contextId *string) (*KubeProviderMetrics, error) {
	config := ContextSwitcher(contextId)

	clientSet, errClientSet := metricsv.NewForConfig(config)
	if errClientSet != nil {
		panic(errClientSet.Error())
	}

	return &KubeProviderMetrics{
		ClientSet:    clientSet,
		ClientConfig: *config,
	}, nil
}

func newKubeProviderMetricsInCluster(contextId *string) (*KubeProviderMetrics, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	config = ContextSwitcher(contextId)

	clientset, err := metricsv.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &KubeProviderMetrics{
		ClientSet:    clientset,
		ClientConfig: *config,
	}, nil
}
