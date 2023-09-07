package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"

	v1Core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllEvents(namespaceName string) utils.K8sWorkloadResult {
	result := []v1Core.Event{}

	provider := NewKubeProvider()
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

func DescribeK8sEvent(namespace string, name string) utils.K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "event", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}
