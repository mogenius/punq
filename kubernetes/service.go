package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func UpdateServiceWith(service *v1.Service) error {
	kubeProvider := NewKubeProvider()
	serviceClient := kubeProvider.ClientSet.CoreV1().Services("")
	_, err := serviceClient.Update(context.TODO(), service, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}
func ServiceFor(namespace string, serviceName string) *v1.Service {
	kubeProvider := NewKubeProvider()
	serviceClient := kubeProvider.ClientSet.CoreV1().Services(namespace)
	service, err := serviceClient.Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		logger.Log.Errorf("ServiceFor ERROR: %s", err.Error())
		return nil
	}
	return service
}

func AllServices(namespaceName string) []v1.Service {
	result := []v1.Service{}

	provider := NewKubeProvider()
	serviceList, err := provider.ClientSet.CoreV1().Services(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllServices ERROR: %s", err.Error())
		return result
	}

	for _, service := range serviceList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, service.ObjectMeta.Namespace) {
			result = append(result, service)
		}
	}
	return result
}

func AllK8sServices(namespaceName string) utils.K8sWorkloadResult {
	results := AllServices(namespaceName)
	return WorkloadResult(results, nil)
}

func UpdateK8sService(data v1.Service) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.CoreV1().Services(data.ObjectMeta.Namespace)
	_, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sService(data v1.Service) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.CoreV1().Services(data.ObjectMeta.Namespace)
	err := client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sService(namespace string, name string) utils.K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "service", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sService(data v1.Service) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.CoreV1().Services(data.ObjectMeta.Namespace)
	_, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func NewK8sService() K8sNewWorkload {
	return NewWorkload(
		RES_SERVICE,
		utils.InitServiceExampleYaml(),
		"A Kubernetes Service is an abstraction which defines a logical set of Pods and a policy by which to access them. The set of Pods targeted by a Service is usually determined by a selector. In this example, the service named 'my-service' listens on port 80, and forwards the requests to port 9376 on the pods which have the label app=MyApp.")
}
