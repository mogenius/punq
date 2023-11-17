package kubernetes

import (
	"context"
	"fmt"

	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"

	v1Core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllEvents(namespaceName string, contextId *string) utils.K8sWorkloadResult {
	result := []v1Core.Event{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	eventList, err := provider.ClientSet.CoreV1().Events(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllEvents ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, event := range eventList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, event.ObjectMeta.Namespace) {
			result = append(result, event)
		}
	}
	return WorkloadResult(result, nil)
}

func GetEvent(namespaceName string, name string, contextId *string) (*v1Core.Event, error) {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return nil, err
	}
	return provider.ClientSet.CoreV1().Events(namespaceName).Get(context.TODO(), name, metav1.GetOptions{})
}

func DescribeK8sEvent(namespace string, name string, contextId *string) utils.K8sWorkloadResult {
	cmd := utils.RunOnLocalShell(fmt.Sprintf("/usr/local/bin/kubectl describe event %s -n %s%s", name, namespace, ContextFlag(contextId)))

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}
