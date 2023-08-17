package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"

	storage "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllVolumeSnapshots() K8sWorkloadResult {
	result := []storage.VolumeAttachment{}

	provider := NewKubeProvider()
	volAttachList, err := provider.ClientSet.StorageV1().VolumeAttachments().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllVolumeSnapshots ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	result = append(result, volAttachList.Items...)
	return WorkloadResult(result, nil)
}

func UpdateK8sVolumeSnapshot(data storage.VolumeAttachment) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	volAttachClient := kubeProvider.ClientSet.StorageV1().VolumeAttachments()
	_, err := volAttachClient.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sVolumeSnapshot(data storage.VolumeAttachment) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	volAttachClient := kubeProvider.ClientSet.StorageV1().VolumeAttachments()
	err := volAttachClient.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sVolumeSnapshot(namespace string, name string) K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "volumesnapshots", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func NewK8sVolumeSnapshots() K8sNewWorkload {
	return NewWorkload(
		RES_VOLUMESNAPSHOTS,
		utils.InitVolumeSnapshotYaml(),
		"A VolumeSnapshot in Kubernetes is a representation of a storage volume at a particular point in time. It's part of the Kubernetes storage system and is used for creating backups of data.	This YAML file will create a VolumeSnapshot named 'snapshot-test' from the PersistentVolumeClaim named 'pvc-test'. The snapshot will be taken using the VolumeSnapshotClass named 'snapshot-class'. The VolumeSnapshotClass would typically be defined by your storage provider and would specify the underlying snapshotting technology to use.")
}
