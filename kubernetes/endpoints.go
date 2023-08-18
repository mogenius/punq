package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllEndpoints(namespaceName string) K8sWorkloadResult {
	result := []corev1.Endpoints{}

	provider := NewKubeProvider()
	hpaList, err := provider.ClientSet.CoreV1().Endpoints(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
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

func UpdateK8sEndpoint(data corev1.Endpoints) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	hpaClient := kubeProvider.ClientSet.CoreV1().Endpoints(data.Namespace)
	_, err := hpaClient.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sEndpoint(data corev1.Endpoints) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	hpaClient := kubeProvider.ClientSet.CoreV1().Endpoints(data.Namespace)
	err := hpaClient.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sEndpoint(namespace string, name string) K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "endpoint", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func NewK8sEndpoint() K8sNewWorkload {
	return NewWorkload(
		RES_ENDPOINTS,
		utils.InitEndPointYaml(),
		"The Endpoints resource represents a set of network addresses for a service and allows the service to be accessed internally by other resources in the cluster. In this example, an Endpoints resource named 'my-service' is created. It specifies two IP addresses, '10.0.0.1' and '10.0.0.2', as the network endpoints for the service. The service is accessible on port 8080. Endpoints are typically managed automatically by Kubernetes controllers based on the availability and readiness of the corresponding Pods.")
}