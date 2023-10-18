package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"

	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllIngressClasses(contextId *string) []v1.IngressClass {
	result := []v1.IngressClass{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return result
	}
	ingressList, err := provider.ClientSet.NetworkingV1().IngressClasses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllIngressClasses ERROR: %s", err.Error())
		return result
	}

	result = append(result, ingressList.Items...)

	return result
}

func AllK8sIngressClasses(contextId *string) utils.K8sWorkloadResult {
	result := []v1.IngressClass{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	ingressList, err := provider.ClientSet.NetworkingV1().IngressClasses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllK8sIngressClasses ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	result = append(result, ingressList.Items...)

	return WorkloadResult(result, nil)
}

func GetK8sIngressClass(name string, contextId *string) (*v1.IngressClass, error) {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return nil, err
	}
	return provider.ClientSet.NetworkingV1().IngressClasses().Get(context.TODO(), name, metav1.GetOptions{})
}

func UpdateK8sIngressClass(data v1.IngressClass, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.NetworkingV1().IngressClasses()
	res, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func DeleteK8sIngressClass(data v1.IngressClass, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.NetworkingV1().IngressClasses()
	err = client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sIngressClassBy(name string, contextId *string) error {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return err
	}
	client := provider.ClientSet.NetworkingV1().IngressClasses()
	return client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sIngressClass(name string, contextId *string) utils.K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "ingressclass", name, ContextFlag(contextId))

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sIngressClass(data v1.IngressClass, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.NetworkingV1().IngressClasses()
	res, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func NewK8sIngressClass() K8sNewWorkload {
	return NewWorkload(
		RES_INGRESS_CLASS,
		utils.InitIngresClassYaml(),
		"In Kubernetes, an IngressClass is a way to determine which controller should handle Ingress resources. It's important to create an IngressClass in environments that have multiple Ingress controllers or in cases where you're using specific configurations for your Ingress rules.")
}
