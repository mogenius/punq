package kubernetes

import (
	"context"
	"fmt"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllPersistentVolumeClaims(namespaceName string, contextId *string) []core.PersistentVolumeClaim {
	result := []core.PersistentVolumeClaim{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return result
	}
	pvList, err := provider.ClientSet.CoreV1().PersistentVolumeClaims(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllPersistentVolumeClaims ERROR: %s", err.Error())
		return result
	}

	for _, v := range pvList.Items {
		v.Kind = "PersistentVolumeClaim"
		v.APIVersion = "v1"
		result = append(result, v)
	}

	return result
}

func GetPersistentVolumeClaim(namespaceName string, name string, contextId *string) (*core.PersistentVolumeClaim, error) {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return nil, err
	}
	pvc, err := provider.ClientSet.CoreV1().PersistentVolumeClaims(namespaceName).Get(context.TODO(), name, metav1.GetOptions{})
	pvc.Kind = "PersistentVolumeClaim"
	pvc.APIVersion = "v1"

	return pvc, err
}

func AllK8sPersistentVolumeClaims(namespaceName string, contextId *string) utils.K8sWorkloadResult {
	result := []core.PersistentVolumeClaim{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	pvList, err := provider.ClientSet.CoreV1().PersistentVolumeClaims(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllPersistentVolumeClaims ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, v := range pvList.Items {
		v.Kind = "PersistentVolumeClaim"
		v.APIVersion = "v1"
		result = append(result, v)
	}
	return WorkloadResult(result, nil)
}

func UpdateK8sPersistentVolumeClaim(data core.PersistentVolumeClaim, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.CoreV1().PersistentVolumeClaims(data.Namespace)
	res, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func DeleteK8sPersistentVolumeClaim(data core.PersistentVolumeClaim, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.CoreV1().PersistentVolumeClaims(data.Namespace)
	err = client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sPersistentVolumeClaimBy(namespace string, name string, contextId *string) error {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return err
	}
	client := provider.ClientSet.CoreV1().PersistentVolumeClaims(namespace)
	return client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sPersistentVolumeClaim(namespace string, name string, contextId *string) utils.K8sWorkloadResult {
	cmd := utils.RunOnLocalShell(fmt.Sprintf("kubectl describe persistentvolumeclaim %s -n %s%s", name, namespace, ContextFlag(contextId)))

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sPersistentVolumeClaim(data core.PersistentVolumeClaim, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.CoreV1().PersistentVolumeClaims(data.Namespace)
	res, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func NewK8sPersistentVolumeClaim() K8sNewWorkload {
	return NewWorkload(
		RES_PERSISTENT_VOLUME_CLAIM,
		utils.InitPersistentVolumeClaimYaml(),
		"A PersistentVolumeClaim (PVC) is a request for storage by a user. It is similar to a Pod. Pods consume node resources, and PVCs consume PersistentVolume resources.")
}
