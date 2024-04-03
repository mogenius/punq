package kubernetes

import (
	"context"
	"fmt"

	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllEndpoints(namespaceName string, contextId *string) []corev1.Endpoints {
	result := []corev1.Endpoints{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return result
	}
	endpointList, err := provider.ClientSet.CoreV1().Endpoints(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllEndpoints ERROR: %s", err.Error())
		return result
	}

	for _, endpoint := range endpointList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, endpoint.ObjectMeta.Namespace) {
			endpoint.Kind = "Endpoints"
			result = append(result, endpoint)
		}
	}
	return result
}

func AllK8sEndpoints(namespaceName string, contextId *string) utils.K8sWorkloadResult {
	result := []corev1.Endpoints{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	endpointList, err := provider.ClientSet.CoreV1().Endpoints(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllEndpoints ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, endpoint := range endpointList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, endpoint.ObjectMeta.Namespace) {
			endpoint.Kind = "Endpoints"
			result = append(result, endpoint)
		}
	}
	return WorkloadResult(result, nil)
}

func GetEndpoint(namespaceName string, name string, contextId *string) (*corev1.Endpoints, error) {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return nil, err
	}
	return provider.ClientSet.CoreV1().Endpoints(namespaceName).Get(context.TODO(), name, metav1.GetOptions{})
}

func UpdateK8sEndpoint(data corev1.Endpoints, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.CoreV1().Endpoints(data.Namespace)
	res, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func DeleteK8sEndpoint(data corev1.Endpoints, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.CoreV1().Endpoints(data.Namespace)
	err = client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sEndpointBy(namespace string, name string, contextId *string) error {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return err
	}
	client := provider.ClientSet.CoreV1().Endpoints(namespace)
	return client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sEndpoint(namespace string, name string, contextId *string) utils.K8sWorkloadResult {
	cmd := utils.RunOnLocalShell(fmt.Sprintf("kubectl describe endpoint %s -n %s%s", name, namespace, ContextFlag(contextId)))

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sEndpoint(data corev1.Endpoints, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.CoreV1().Endpoints(data.Namespace)
	res, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func NewK8sEndpoint() K8sNewWorkload {
	return NewWorkload(
		RES_ENDPOINT,
		utils.InitEndPointYaml(),
		"The Endpoints resource represents a set of network addresses for a service and allows the service to be accessed internally by other resources in the cluster. In this example, an Endpoints resource named 'my-service' is created. It specifies two IP addresses, '10.0.0.1' and '10.0.0.2', as the network endpoints for the service. The service is accessible on port 8080. Endpoints are typically managed automatically by Kubernetes controllers based on the availability and readiness of the corresponding Pods.")
}
