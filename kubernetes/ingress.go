package kubernetes

import (
	"context"
	"os/exec"
	"punq/logger"
	"punq/utils"

	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	INGRESS_PREFIX = "ingress"
)

func AllIngresses(namespaceName string) []v1.Ingress {
	result := []v1.Ingress{}

	provider := NewKubeProvider()
	ingressList, err := provider.ClientSet.NetworkingV1().Ingresses(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllIngresses ERROR: %s", err.Error())
		return result
	}

	for _, ingress := range ingressList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, ingress.ObjectMeta.Namespace) {
			result = append(result, ingress)
		}
	}
	return result
}

func AllK8sIngresses(namespaceName string) K8sWorkloadResult {
	result := []v1.Ingress{}

	provider := NewKubeProvider()
	ingressList, err := provider.ClientSet.NetworkingV1().Ingresses(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllIngresses ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, ingress := range ingressList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, ingress.ObjectMeta.Namespace) {
			result = append(result, ingress)
		}
	}
	return WorkloadResult(result, nil)
}

func UpdateK8sIngress(data v1.Ingress) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	ingressClient := kubeProvider.ClientSet.NetworkingV1().Ingresses(data.Namespace)
	_, err := ingressClient.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sIngress(data v1.Ingress) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	ingressClient := kubeProvider.ClientSet.NetworkingV1().Ingresses(data.Namespace)
	err := ingressClient.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sIngress(namespace string, name string) K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "ingress", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func NewK8sIngress() K8sNewWorkload {
	return NewWorkload(
		RES_INGRESS,
		utils.InitIngresYaml(),
		"An Ingress is a collection of rules that allow inbound connections to reach the cluster services. In this example, an Ingress named 'example-ingress' is created. It will route traffic that comes in on 'myapp.mydomain.com' with a URL path that starts with '/testpath' to the service named 'test' on port 80.")
}
