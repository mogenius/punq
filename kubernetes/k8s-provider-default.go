package kubernetes

import (
	"fmt"

	"github.com/mogenius/punq/logger"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeProvider struct {
	ClientSet    *kubernetes.Clientset
	ClientConfig rest.Config
}

func NewKubeProvider(contextId *string) *KubeProvider {
	var kubeProvider *KubeProvider
	var err error
	if RunsInCluster {
		kubeProvider, err = newKubeProviderInCluster(contextId)
	} else {
		kubeProvider, err = newKubeProviderLocal(contextId)
	}

	if err != nil {
		logger.Log.Fatalf("ERROR: %s", err.Error())
	}
	return kubeProvider
}

func newKubeProviderLocal(contextId *string) (*KubeProvider, error) {
	config := ContextSwitcher(contextId)

	clientSet, errClientSet := kubernetes.NewForConfig(config)
	if errClientSet != nil {
		panic(errClientSet.Error())
	}

	return &KubeProvider{
		ClientSet:    clientSet,
		ClientConfig: *config,
	}, nil
}

func newKubeProviderInCluster(contextId *string) (*KubeProvider, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	config = ContextSwitcher(contextId)

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &KubeProvider{
		ClientSet:    clientset,
		ClientConfig: *config,
	}, nil
}

func ContextConfigLoader(contextId *string) (*rest.Config, error) {
	// get current context
	ctx := ContextForId(*contextId)
	if ctx == nil {
		return nil, fmt.Errorf("Context not found for id: %s", *contextId)
	}

	configFromString, err := clientcmd.NewClientConfigFromBytes([]byte(ctx.Context))
	if err != nil {
		logger.Log.Errorf("Error creating client config from string:", err.Error())
		return nil, err
	}

	config, err := configFromString.ClientConfig()
	return config, err
}

func ContextSwitcher(contextId *string) (restConfig *rest.Config) {
	var kubeconfig string = getKubeConfig()

	// CONTEXT SWITCHER
	if contextId != nil && *contextId != "" {
		restConfig, err := ContextConfigLoader(contextId)
		if err != nil || restConfig == nil {
			logger.Log.Fatal(err.Error())
		}
		return restConfig
	} else {
		restConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil || restConfig == nil {
			logger.Log.Fatal(err.Error())
		}
		return restConfig
	}
}
