package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ConfigMapFor(namespace string, configMapName string) *v1.ConfigMap {
	kubeProvider := NewKubeProvider()
	configMapClient := kubeProvider.ClientSet.CoreV1().ConfigMaps(namespace)
	configMap, err := configMapClient.Get(context.TODO(), configMapName, metav1.GetOptions{})
	if err != nil {
		logger.Log.Errorf("ConfigMapFor ERROR: %s", err.Error())
		return nil
	}
	return configMap
}

func AllConfigmaps(namespaceName string) []v1.ConfigMap {
	result := []v1.ConfigMap{}

	provider := NewKubeProvider()
	configmapList, err := provider.ClientSet.CoreV1().ConfigMaps(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllConfigmaps ERROR: %s", err.Error())
		return result
	}

	for _, configmap := range configmapList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, configmap.ObjectMeta.Namespace) {
			result = append(result, configmap)
		}
	}
	return result
}

func AllK8sConfigmaps(namespaceName string) utils.K8sWorkloadResult {
	result := []v1.ConfigMap{}

	provider := NewKubeProvider()
	configmapList, err := provider.ClientSet.CoreV1().ConfigMaps(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllConfigmaps ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, configmap := range configmapList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, configmap.ObjectMeta.Namespace) {
			result = append(result, configmap)
		}
	}
	return WorkloadResult(result, nil)
}

func GetK8sConfigmap(namespaceName string, name string) (*v1.ConfigMap, error) {
	provider := NewKubeProvider()
	return provider.ClientSet.CoreV1().ConfigMaps(namespaceName).Get(context.TODO(), name, metav1.GetOptions{})
}

func UpdateK8sConfigMap(data v1.ConfigMap) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.CoreV1().ConfigMaps(data.Namespace)
	_, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		logger.Log.Errorf("UpdateK8sConfigMap ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sConfigmap(data v1.ConfigMap) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.CoreV1().ConfigMaps(data.Namespace)
	err := client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		logger.Log.Errorf("DeleteK8sConfigmap ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sConfigmapBy(data v1.ConfigMap) error {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.CoreV1().ConfigMaps(data.Namespace)
	return client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
}

func DescribeK8sConfigmap(namespace string, name string) utils.K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "configmap", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sConfigMap(data v1.ConfigMap) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.CoreV1().ConfigMaps(data.Namespace)
	_, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func NewK8sConfigmap() K8sNewWorkload {
	return NewWorkload(
		RES_CONFIGMAP,
		utils.InitConfigMapYaml(),
		"ConfigMaps allow you to decouple configuration artifacts from image content to keep containerized applications portable. In this example, a ConfigMap named 'my-configmap' is created with two key-value pairs: my-key and my-value, another-key and another-value. ConfigMap data can be referenced in many ways depending on where you need the data to be used. For example, you could use a ConfigMap to set environment variables for a Pod, or mount a ConfigMap as a volume in a Pod.")
}
