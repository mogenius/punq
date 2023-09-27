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

func NewKubeProviderMetrics(contextId *string) (*KubeProviderMetrics, error) {
	var provider *KubeProviderMetrics
	var err error
	if RunsInCluster {
		provider, err = newKubeProviderMetricsInCluster(contextId)
	} else {
		provider, err = newKubeProviderMetricsLocal(contextId)
	}

	if err != nil {
		logger.Log.Errorf("ERROR: %s", err.Error())
	}
	return provider, err
}

func newKubeProviderMetricsLocal(contextId *string) (*KubeProviderMetrics, error) {
	config, err := ContextSwitcher(contextId)
	if err != nil {
		return nil, err
	}

	clientSet, errClientSet := metricsv.NewForConfig(config)
	if errClientSet != nil {
		return nil, errClientSet
	}

	return &KubeProviderMetrics{
		ClientSet:    clientSet,
		ClientConfig: *config,
	}, nil
}

func newKubeProviderMetricsInCluster(contextId *string) (*KubeProviderMetrics, error) {
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

	clientset, err := metricsv.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &KubeProviderMetrics{
		ClientSet:    clientset,
		ClientConfig: *config,
	}, nil
}
