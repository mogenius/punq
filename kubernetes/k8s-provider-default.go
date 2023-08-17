package kubernetes

import (
	"path/filepath"
	"punq/logger"
	"punq/utils"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type KubeProvider struct {
	ClientSet    *kubernetes.Clientset
	ClientConfig rest.Config
}

func NewKubeProvider() *KubeProvider {
	var kubeProvider *KubeProvider
	var err error
	if utils.CONFIG.Kubernetes.RunInCluster {
		kubeProvider, err = NewKubeProviderInCluster()
	} else {
		kubeProvider, err = NewKubeProviderLocal()
	}

	if err != nil {
		logger.Log.Errorf("ERROR: %s", err.Error())
	}
	return kubeProvider
}

func NewKubeProviderLocal() (*KubeProvider, error) {
	var kubeconfig string = ""
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	restConfig, errConfig := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if errConfig != nil {
		panic(errConfig.Error())
	}

	clientSet, errClientSet := kubernetes.NewForConfig(restConfig)
	if errClientSet != nil {
		panic(errClientSet.Error())
	}

	return &KubeProvider{
		ClientSet:    clientSet,
		ClientConfig: *restConfig,
	}, nil
}

func NewKubeProviderInCluster() (*KubeProvider, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &KubeProvider{
		ClientSet:    clientset,
		ClientConfig: *config,
	}, nil
}
