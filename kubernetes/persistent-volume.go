package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllPersistentVolumes() K8sWorkloadResult {
	result := []core.PersistentVolume{}

	provider := NewKubeProvider()
	pvList, err := provider.ClientSet.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllPersistentVolumes ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, pv := range pvList.Items {
		result = append(result, pv)
	}
	return WorkloadResult(result, nil)
}

func UpdateK8sPersistentVolume(data core.PersistentVolume) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.CoreV1().PersistentVolumes()
	_, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sPersistentVolume(data core.PersistentVolume) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.CoreV1().PersistentVolumes()
	err := client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sPersistentVolume(name string) K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "persistentvolume", name)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sPersistentVolume(data core.PersistentVolume) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.CoreV1().PersistentVolumes()
	_, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func NewK8sVolume() K8sNewWorkload {
	return NewWorkload(
		RES_PERSISTENT_VOLUME,
		utils.InitPersistentVolumeYaml(),
		"A PersistentVolume (PV) is a piece of storage in the cluster that has been provisioned by an administrator or dynamically provisioned using Storage Classes. It is a resource in the cluster just like a node is a cluster resource. In this example, a PersistentVolume named 'my-pv' is created with a capacity of 10Gi, it uses the local directory /data/my-pv on the host for storage. Please note, that hostPath is a simple type of storage and useful for development and testing. For production usage, you might want to use a more robust solution like networked storage (NFS, iSCSI, GlusterFS, etc) or cloud-provided storage (AWS EBS, GCE Persistent Disk, Azure Disk, etc).")
}
