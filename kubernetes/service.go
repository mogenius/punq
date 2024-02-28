package kubernetes

import (
	"context"
	"fmt"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func UpdateServiceWith(service *v1.Service, contextId *string) error {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return err
	}
	serviceClient := provider.ClientSet.CoreV1().Services("")
	_, err = serviceClient.Update(context.TODO(), service, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}
func ServiceFor(namespace string, serviceName string, contextId *string) *v1.Service {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return nil
	}
	serviceClient := provider.ClientSet.CoreV1().Services(namespace)
	service, err := serviceClient.Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		logger.Log.Errorf("ServiceFor ERROR: %s", err.Error())
		return nil
	}
	return service
}

func GetService(namespace string, serviceName string, contextId *string) (*v1.Service, error) {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return nil, err
	}
	serviceClient := provider.ClientSet.CoreV1().Services(namespace)
	return serviceClient.Get(context.TODO(), serviceName, metav1.GetOptions{})
}

func AllServices(namespaceName string, contextId *string) []v1.Service {
	result := []v1.Service{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return result
	}
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

func AllK8sServices(namespaceName string, contextId *string) utils.K8sWorkloadResult {
	results := AllServices(namespaceName, contextId)
	return WorkloadResult(results, nil)
}

func UpdateK8sService(data v1.Service, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.CoreV1().Services(data.ObjectMeta.Namespace)
	res, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func DeleteK8sService(data v1.Service, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.CoreV1().Services(data.ObjectMeta.Namespace)
	err = client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sServiceBy(namespace string, name string, contextId *string) error {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return err
	}
	client := provider.ClientSet.CoreV1().Services(namespace)
	return client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sService(namespace string, name string, contextId *string) utils.K8sWorkloadResult {
	cmd := utils.RunOnLocalShell(fmt.Sprintf("kubectl describe service %s -n %s%s", name, namespace, ContextFlag(contextId)))

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sService(data v1.Service, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.CoreV1().Services(data.ObjectMeta.Namespace)
	res, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func NewK8sService() K8sNewWorkload {
	return NewWorkload(
		RES_SERVICE,
		utils.InitServiceExampleYaml(),
		"A Kubernetes Service is an abstraction which defines a logical set of Pods and a policy by which to access them. The set of Pods targeted by a Service is usually determined by a selector. In this example, the service named 'my-service' listens on port 80, and forwards the requests to port 9376 on the pods which have the label app=MyApp.")
}
