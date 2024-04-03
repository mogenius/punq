package kubernetes

import (
	"context"
	"fmt"

	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"

	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	INGRESS_PREFIX = "ingress"
)

func AllIngresses(namespaceName string, contextId *string) []v1.Ingress {
	result := []v1.Ingress{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return result
	}
	ingressList, err := provider.ClientSet.NetworkingV1().Ingresses(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllIngresses ERROR: %s", err.Error())
		return result
	}

	for _, ingress := range ingressList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, ingress.ObjectMeta.Namespace) {
			ingress.Kind = "Ingress"
			result = append(result, ingress)
		}
	}
	return result
}

func AllK8sIngresses(namespaceName string, contextId *string) utils.K8sWorkloadResult {
	result := []v1.Ingress{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	ingressList, err := provider.ClientSet.NetworkingV1().Ingresses(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllIngresses ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, ingress := range ingressList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, ingress.ObjectMeta.Namespace) {
			ingress.Kind = "Ingress"
			result = append(result, ingress)
		}
	}
	return WorkloadResult(result, nil)
}

func GetK8sIngress(namespaceName string, name string, contextId *string) (*v1.Ingress, error) {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return nil, err
	}
	return provider.ClientSet.NetworkingV1().Ingresses(namespaceName).Get(context.TODO(), name, metav1.GetOptions{})
}

func UpdateK8sIngress(data v1.Ingress, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.NetworkingV1().Ingresses(data.Namespace)
	res, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func DeleteK8sIngress(data v1.Ingress, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.NetworkingV1().Ingresses(data.Namespace)
	err = client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sIngressBy(namespace string, name string, contextId *string) error {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return err
	}
	client := provider.ClientSet.NetworkingV1().Ingresses(namespace)
	return client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sIngress(namespace string, name string, contextId *string) utils.K8sWorkloadResult {
	cmd := utils.RunOnLocalShell(fmt.Sprintf("kubectl describe ingress %s -n %s%s", name, namespace, ContextFlag(contextId)))

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sIngress(data v1.Ingress, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.NetworkingV1().Ingresses(data.Namespace)
	res, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func NewK8sIngress() K8sNewWorkload {
	return NewWorkload(
		RES_INGRESS,
		utils.InitIngresYaml(),
		"An Ingress is a collection of rules that allow inbound connections to reach the cluster services. In this example, an Ingress named 'example-ingress' is created. It will route traffic that comes in on 'myapp.mydomain.com' with a URL path that starts with '/testpath' to the service named 'test' on port 80.")
}
