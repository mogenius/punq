package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	storage "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllStorageClasses(contextId *string) utils.K8sWorkloadResult {
	result := []storage.StorageClass{}

	provider := NewKubeProvider(contextId)
	scList, err := provider.ClientSet.StorageV1().StorageClasses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllStorageClasses ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, pv := range scList.Items {
		result = append(result, pv)
	}
	return WorkloadResult(result, nil)
}

func GetStorageClass(name string, contextId *string) (*storage.StorageClass, error) {
	provider := NewKubeProvider(contextId)
	return provider.ClientSet.StorageV1().StorageClasses().Get(context.TODO(), name, metav1.GetOptions{})
}

func UpdateK8sStorageClass(data storage.StorageClass, contextId *string) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider(contextId)
	client := kubeProvider.ClientSet.StorageV1().StorageClasses()
	_, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sStorageClass(data storage.StorageClass, contextId *string) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider(contextId)
	client := kubeProvider.ClientSet.StorageV1().StorageClasses()
	err := client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sStorageClassBy(name string, contextId *string) error {
	kubeProvider := NewKubeProvider(contextId)
	client := kubeProvider.ClientSet.StorageV1().StorageClasses()
	return client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sStorageClass(name string, contextId *string) utils.K8sWorkloadResult {
	cmd := exec.Command("kubectl", ContextFlag(contextId), "describe", "storageclass", name)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sStorageClass(data storage.StorageClass, contextId *string) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider(contextId)
	client := kubeProvider.ClientSet.StorageV1().StorageClasses()
	_, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func NewK8sStorageClass() K8sNewWorkload {
	return NewWorkload(
		RES_STORAGECLASS,
		utils.InitStorageClassYaml(),
		"A StorageClass provides a way for administrators to describe the 'classes' of storage they offer. Different classes might map to quality-of-service levels, backup policies, or arbitrary policies determined by the cluster administrators. Please note, the above example uses kubernetes.io/aws-ebs as the provisioner which means this StorageClass is specific to AWS EBS volumes. The parameters may vary based on the provisioner.")
}
