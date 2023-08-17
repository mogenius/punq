package kubernetes

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

func NewKubeProviderMetricsLocal() (*KubeProviderMetrics, error) {
	kubeconfig := getKubeConfig()

	restConfig, errConfig := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if errConfig != nil {
		panic(errConfig.Error())
	}

	clientSet, errClientSet := metricsv.NewForConfig(restConfig)
	if errClientSet != nil {
		panic(errClientSet.Error())
	}

	//logger.Log.Debugf("K8s client config (init with .kube/config), host: %s", restConfig.Host)

	return &KubeProviderMetrics{
		ClientSet:    clientSet,
		ClientConfig: *restConfig,
	}, nil
}

func NewKubeProviderMetricsInCluster() (*KubeProviderMetrics, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := metricsv.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	//logger.Log.Debugf("K8s client config (init InCluster), host: %s", config.Host)

	return &KubeProviderMetrics{
		ClientSet:    clientset,
		ClientConfig: *config,
	}, nil
}
