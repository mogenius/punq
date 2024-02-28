package kubernetes

import (
	"context"
	"fmt"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllPersistentVolumesRaw(contextId *string) []core.PersistentVolume {
	result := []core.PersistentVolume{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return result
	}
	pvList, err := provider.ClientSet.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllPersistentVolumesRaw ERROR: %s", err.Error())
		return result
	}
	result = append(result, pvList.Items...)

	return result
}

func AllPersistentVolumes(contextId *string) utils.K8sWorkloadResult {
	result := []core.PersistentVolume{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	pvList, err := provider.ClientSet.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllPersistentVolumes ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	result = append(result, pvList.Items...)
	return WorkloadResult(result, nil)
}

func GetPersistentVolume(name string, contextId *string) (*core.PersistentVolume, error) {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return nil, err
	}
	return provider.ClientSet.CoreV1().PersistentVolumes().Get(context.TODO(), name, metav1.GetOptions{})
}

func UpdateK8sPersistentVolume(data core.PersistentVolume, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.CoreV1().PersistentVolumes()
	res, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func DeleteK8sPersistentVolume(data core.PersistentVolume, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.CoreV1().PersistentVolumes()
	err = client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sPersistentVolumeBy(name string, contextId *string) error {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return err
	}
	client := provider.ClientSet.CoreV1().PersistentVolumes()
	return client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sPersistentVolume(name string, contextId *string) utils.K8sWorkloadResult {
	cmd := utils.RunOnLocalShell(fmt.Sprintf("kubectl describe persistentvolume %s%s", name, ContextFlag(contextId)))

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sPersistentVolume(data core.PersistentVolume, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.CoreV1().PersistentVolumes()
	res, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func NewK8sVolume() K8sNewWorkload {
	return NewWorkload(
		RES_PERSISTENT_VOLUME,
		utils.InitPersistentVolumeYaml(),
		"A PersistentVolume (PV) is a piece of storage in the cluster that has been provisioned by an administrator or dynamically provisioned using Storage Classes. It is a resource in the cluster just like a node is a cluster resource. In this example, a PersistentVolume named 'my-pv' is created with a capacity of 10Gi, it uses the local directory /data/my-pv on the host for storage. Please note, that hostPath is a simple type of storage and useful for development and testing. For production usage, you might want to use a more robust solution like networked storage (NFS, iSCSI, GlusterFS, etc) or cloud-provided storage (AWS EBS, GCE Persistent Disk, Azure Disk, etc).")
}
