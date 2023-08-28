package kubernetes

import (
	snapClientset "github.com/kubernetes-csi/external-snapshotter/client/v6/clientset/versioned"
	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeProviderSnapshot struct {
	ClientSet    *snapClientset.Clientset
	ClientConfig rest.Config
}

func NewKubeProviderSnapshot() *KubeProviderSnapshot {
	var kubeProvider *KubeProviderSnapshot
	var err error
	if utils.CONFIG.Kubernetes.RunInCluster {
		kubeProvider, err = newKubeProviderCsiInCluster()
	} else {
		kubeProvider, err = newKubeProviderCsiLocal()
	}

	if err != nil {
		logger.Log.Errorf("ERROR: %s", err.Error())
	}
	return kubeProvider
}

func newKubeProviderCsiLocal() (*KubeProviderSnapshot, error) {
	kubeconfig := getKubeConfig()

	restConfig, errConfig := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if errConfig != nil {
		panic(errConfig.Error())
	}

	clientSet, errClientSet := snapClientset.NewForConfig(restConfig)
	if errClientSet != nil {
		panic(errClientSet.Error())
	}

	return &KubeProviderSnapshot{
		ClientSet:    clientSet,
		ClientConfig: *restConfig,
	}, nil
}

func newKubeProviderCsiInCluster() (*KubeProviderSnapshot, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
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