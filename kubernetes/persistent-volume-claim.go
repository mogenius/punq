package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllPersistentVolumeClaims(namespaceName string) []core.PersistentVolumeClaim {
	result := []core.PersistentVolumeClaim{}

	provider := NewKubeProvider()
	pvList, err := provider.ClientSet.CoreV1().PersistentVolumeClaims(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllPersistentVolumeClaims ERROR: %s", err.Error())
		return result
	}
	result = append(result, pvList.Items...)

	return result
}

func AllK8sPersistentVolumeClaims(namespaceName string) utils.HttpResult {
	result := []core.PersistentVolumeClaim{}

	provider := NewKubeProvider()
	pvList, err := provider.ClientSet.CoreV1().PersistentVolumeClaims(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllPersistentVolumeClaims ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, pv := range pvList.Items {
		result = append(result, pv)
	}
	return WorkloadResult(result, nil)
}

func UpdateK8sPersistentVolumeClaim(data core.PersistentVolumeClaim) utils.HttpResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.CoreV1().PersistentVolumeClaims(data.Namespace)
	_, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sPersistentVolumeClaim(data core.PersistentVolumeClaim) utils.HttpResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.CoreV1().PersistentVolumeClaims(data.Namespace)
	err := client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sPersistentVolumeClaim(namespace string, name string) utils.HttpResult {
	cmd := exec.Command("kubectl", "describe", "persistentvolumeclaim", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sPersistentVolumeClaim(data core.PersistentVolumeClaim) utils.HttpResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.CoreV1().PersistentVolumeClaims(data.Namespace)
	_, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func NewK8sPersistentVolumeClaim() K8sNewWorkload {
	return NewWorkload(
		RES_PERSISTENT_VOLUME_CLAIM,
		utils.InitPersistentVolumeClaimYaml(),
		"A PersistentVolumeClaim (PVC) is a request for storage by a user. It is similar to a Pod. Pods consume node resources, and PVCs consume PersistentVolume resources.")
}
