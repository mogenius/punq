package kubernetes

import (
	"encoding/base64"
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
		kubeProvider, err = newKubeProviderLocal()
	}

	if err != nil {
		logger.Log.Errorf("ERROR: %s", err.Error())
	}
	return kubeProvider
}

func newKubeProviderLocal() (*KubeProvider, error) {
	var kubeconfig string = getKubeConfig()

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

func newKubeProviderInCluster(contextId *string) (*KubeProvider, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// CONTEXT SWITCHER
	if contextId != nil {
		config, err = ContextConfigLoader(contextId)
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

func ContextConfigLoader(contextId *string) (*rest.Config, error) {
	// get current context
	ctx := ContextForId(*contextId)
	if ctx == nil {
		return nil, fmt.Errorf("Context not found for id: %s", *contextId)
	}

	// decode base64
	decodedBytes, err := base64.StdEncoding.DecodeString(ctx.ContextBase64)
	if err != nil {
		logger.Log.Errorf("Error decoding string:", err.Error())
		return nil, err
	}

	configFromString, err := clientcmd.NewClientConfigFromBytes(decodedBytes)
	if err != nil {
		logger.Log.Errorf("Error creating client config from string:", err.Error())
		return nil, err
	}

	config, err := configFromString.ClientConfig()
	return config, err
}
