package kubernetes

import (
	"context"
	"fmt"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	v2 "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllHpas(namespaceName string, contextId *string) []v2.HorizontalPodAutoscaler {
	result := []v2.HorizontalPodAutoscaler{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return result
	}
	hpaList, err := provider.ClientSet.AutoscalingV2().HorizontalPodAutoscalers(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllHpas ERROR: %s", err.Error())
		return result
	}

	for _, hpa := range hpaList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, hpa.ObjectMeta.Namespace) {
			hpa.Kind = "HorizontalPodAutoscaler"
			hpa.APIVersion = "autoscaling/v2"
			result = append(result, hpa)
		}
	}
	return result
}

func AllK8sHpas(namespaceName string, contextId *string) utils.K8sWorkloadResult {
	result := []v2.HorizontalPodAutoscaler{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	hpaList, err := provider.ClientSet.AutoscalingV2().HorizontalPodAutoscalers(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllHpas ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, hpa := range hpaList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, hpa.ObjectMeta.Namespace) {
			hpa.Kind = "HorizontalPodAutoscaler"
			hpa.APIVersion = "autoscaling/v2"
			result = append(result, hpa)
		}
	}
	return WorkloadResult(result, nil)
}

func GetHpa(namespaceName string, name string, contextId *string) (*v2.HorizontalPodAutoscaler, error) {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return nil, err
	}
	return provider.ClientSet.AutoscalingV2().HorizontalPodAutoscalers(namespaceName).Get(context.TODO(), name, metav1.GetOptions{})
}

func UpdateK8sHpa(data v2.HorizontalPodAutoscaler, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.AutoscalingV2().HorizontalPodAutoscalers(data.Namespace)
	res, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func DeleteK8sHpa(data v2.HorizontalPodAutoscaler, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.AutoscalingV2().HorizontalPodAutoscalers(data.Namespace)
	err = client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sHpaBy(namespace string, name string, contextId *string) error {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return err
	}
	client := provider.ClientSet.AutoscalingV2().HorizontalPodAutoscalers(namespace)
	return client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sHpa(namespace string, name string, contextId *string) utils.K8sWorkloadResult {
	cmd := utils.RunOnLocalShell(fmt.Sprintf("kubectl describe hpa %s -n %s%s", name, namespace, ContextFlag(contextId)))

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sHpa(data v2.HorizontalPodAutoscaler, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.AutoscalingV2().HorizontalPodAutoscalers(data.Namespace)
	res, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func NewK8sHpa() K8sNewWorkload {
	return NewWorkload(
		RES_HORIZONTAL_POD_AUTOSCALER,
		utils.InitHpaYaml(),
		"The Horizontal Pod Autoscaler automatically scales the number of pods in a replication controller, deployment, or replica set based on observed CPU utilization. In this example, an HPA named 'example-hpa' is created. It will automatically scale the number of pods in the deployment named 'my-app-deployment' between 1 and 10, aiming to target an average CPU utilization across all Pods of 50%.")
}
