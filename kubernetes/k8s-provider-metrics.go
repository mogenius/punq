package kubernetes

import (
	"github.com/mogenius/punq/logger"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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
		if contextId == nil || *contextId == "" {
			kubeProvider, err = newKubeProviderMetricsLocal()
		} else {
			kubeProvider, err = newKubeProviderMetricsInCluster(contextId)
		}
	}

	if err != nil {
		logger.Log.Errorf("ERROR: %s", err.Error())
	}
	return kubeProvider
}

func newKubeProviderMetricsLocal() (*KubeProviderMetrics, error) {
	kubeconfig := getKubeConfig()

	restConfig, errConfig := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if errConfig != nil {
		panic(errConfig.Error())
	}

	clientSet, errClientSet := metricsv.NewForConfig(restConfig)
	if errClientSet != nil {
		panic(errClientSet.Error())
	}

	return &KubeProviderMetrics{
		ClientSet:    clientSet,
		ClientConfig: *restConfig,
	}, nil
}

func newKubeProviderMetricsInCluster(contextId *string) (*KubeProviderMetrics, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// CONTEXT SWITCHER
	if contextId != nil {
		config, err = ContextConfigLoader(contextId)
		if err != nil || config == nil {
			return nil, err
		}
	}

	clientset, err := metricsv.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &KubeProviderMetrics{
		ClientSet:    clientset,
		ClientConfig: *config,
	}, nil
}
