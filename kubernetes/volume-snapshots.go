package kubernetes

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"

	snap "github.com/kubernetes-csi/external-snapshotter/client/v6/apis/volumesnapshot/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllVolumeSnapshots(namespace string, contextId *string) utils.K8sWorkloadResult {
	result := []snap.VolumeSnapshot{}

	provider, err := NewKubeProviderSnapshot(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	volSnapshotsList, err := provider.ClientSet.SnapshotV1().VolumeSnapshots(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllVolumeSnapshots ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	result = append(result, volSnapshotsList.Items...)
	return WorkloadResult(result, nil)
}

func GetVolumeSnapshot(namespace string, name string, contextId *string) (*snap.VolumeSnapshot, error) {
	provider, err := NewKubeProviderSnapshot(contextId)
	if err != nil {
		return nil, err
	}
	return provider.ClientSet.SnapshotV1().VolumeSnapshots(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func UpdateK8sVolumeSnapshot(data snap.VolumeSnapshot) utils.K8sWorkloadResult {
	return WorkloadResult(nil, fmt.Errorf("UPDATE not available in VolumeSnapshot"))
}

func DeleteK8sVolumeSnapshot(data snap.VolumeSnapshot, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProviderSnapshot(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.SnapshotV1().VolumeSnapshots(data.Namespace)
	err = client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sVolumeSnapshotBy(namespace string, name string, contextId *string) error {
	provider, err := NewKubeProviderSnapshot(contextId)
	if err != nil {
		return err
	}
	client := provider.ClientSet.SnapshotV1().VolumeSnapshots(namespace)
	return client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sVolumeSnapshot(namespace string, name string, contextId *string) utils.K8sWorkloadResult {
	cmd := exec.Command("kubectl", ContextFlag(contextId), "describe", "volumesnapshots", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sVolumeSnapshot(data snap.VolumeSnapshot, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProviderSnapshot(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.SnapshotV1().VolumeSnapshots(data.Namespace)
	res, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func NewK8sVolumeSnapshots() K8sNewWorkload {
	return NewWorkload(
		RES_VOLUME_SNAPSHOT,
		utils.InitVolumeSnapshotYaml(),
		"A VolumeSnapshot in Kubernetes is a representation of a storage volume at a particular point in time. It's part of the Kubernetes storage system and is used for creating backups of data.	This YAML file will create a VolumeSnapshot named 'snapshot-test' from the PersistentVolumeClaim named 'pvc-test'. The snapshot will be taken using the VolumeSnapshotClass named 'snapshot-class'. The VolumeSnapshotClass would typically be defined by your storage provider and would specify the underlying snapshotting technology to use.")
}
