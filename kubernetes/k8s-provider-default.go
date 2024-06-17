package kubernetes

import (
	"fmt"

	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeProvider struct {
	ClientSet    *kubernetes.Clientset
	ClientConfig rest.Config
}

func NewKubeProvider(contextId *string) (*KubeProvider, error) {
	var provider *KubeProvider
	var err error
	if RunsInCluster {
		provider, err = newKubeProviderInCluster(contextId)
	} else {
		provider, err = newKubeProviderLocal(contextId)
	}

	if err != nil {
		logger.Log.Errorf("ERROR: %s", err.Error())
	}
	return provider, err
}

func newKubeProviderLocal(contextId *string) (*KubeProvider, error) {
	config, err := ContextSwitcher(contextId)
	if err != nil {
		return nil, err
	}

	clientSet, errClientSet := kubernetes.NewForConfig(config)
	if errClientSet != nil {
		return nil, errClientSet
	}

	return &KubeProvider{
		ClientSet:    clientSet,
		ClientConfig: *config,
	}, nil
}

func newKubeProviderInCluster(contextId *string) (*KubeProvider, error) {
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

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &KubeProvider{
		ClientSet:    clientset,
		ClientConfig: *config,
	}, nil
}

func ContextConfigLoader(contextId *string) (*rest.Config, error) {
	ctx := ContextForId(*contextId)
	if ctx == nil {
		return nil, fmt.Errorf("context not found for id: %s", *contextId)
	}

	configFromString, err := clientcmd.NewClientConfigFromBytes([]byte(ctx.Context))
	if err != nil {
		logger.Log.Errorf("Error creating client config from string:", err.Error())
		return nil, err
	}

	config, err := configFromString.ClientConfig()
	return config, err
}

func ContextSwitcher(contextId *string) (*rest.Config, error) {
	if contextId != nil && *contextId != "" {
		return ContextConfigLoader(contextId)
	} else {
		var kubeconfigs []string = utils.GetDefaultKubeConfig()
		var config *rest.Config
		var err error
		for _, kubeconfig := range kubeconfigs {
			config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
			if err == nil {
				return config, nil
			}
		}

		return nil, fmt.Errorf("Error loading kubeconfig: %s", err.Error())
	}
}
