package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	v2 "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllHpas(namespaceName string) K8sWorkloadResult {
	result := []v2.HorizontalPodAutoscaler{}

	provider := NewKubeProvider()
	hpaList, err := provider.ClientSet.AutoscalingV2().HorizontalPodAutoscalers(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllHpas ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, hpa := range hpaList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, hpa.ObjectMeta.Namespace) {
			result = append(result, hpa)
		}
	}
	return WorkloadResult(result, nil)
}

func UpdateK8sHpa(data v2.HorizontalPodAutoscaler) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	hpaClient := kubeProvider.ClientSet.AutoscalingV2().HorizontalPodAutoscalers(data.Namespace)
	_, err := hpaClient.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sHpa(data v2.HorizontalPodAutoscaler) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	hpaClient := kubeProvider.ClientSet.AutoscalingV2().HorizontalPodAutoscalers(data.Namespace)
	err := hpaClient.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sHpa(namespace string, name string) K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "hpa", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func NewK8sHpa() K8sNewWorkload {
	return NewWorkload(
		RES_HORIZONTAL_POD_AUTOSCALER,
		utils.InitHpaYaml(),
		"The Horizontal Pod Autoscaler automatically scales the number of pods in a replication controller, deployment, or replica set based on observed CPU utilization. In this example, an HPA named 'example-hpa' is created. It will automatically scale the number of pods in the deployment named 'my-app-deployment' between 1 and 10, aiming to target an average CPU utilization across all Pods of 50%.")
}
