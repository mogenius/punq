package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"

	storage "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllVolumeAttachments() K8sWorkloadResult {
	result := []storage.VolumeAttachment{}

	provider := NewKubeProvider()
	volAttachList, err := provider.ClientSet.StorageV1().VolumeAttachments().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllCertificateSigningRequests ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	result = append(result, volAttachList.Items...)
	return WorkloadResult(result, nil)
}

func UpdateK8sVolumeAttachment(data storage.VolumeAttachment) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	volAttachClient := kubeProvider.ClientSet.StorageV1().VolumeAttachments()
	_, err := volAttachClient.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sVolumeAttachment(data storage.VolumeAttachment) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	volAttachClient := kubeProvider.ClientSet.StorageV1().VolumeAttachments()
	err := volAttachClient.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sVolumeAttachment(name string) K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "volumeattachment", name)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func NewK8sVolumeAttachment() K8sNewWorkload {
	return NewWorkload(
		RES_VOLUME_ATTACHMENT,
		utils.InitVolumeAttachmentYaml(),
		"The VolumeAttachment kind in Kubernetes provides a mechanism for attaching external volumes to a node. It's typically used by the Container Storage Interface (CSI) to allow for the dynamic provisioning of volumes, but it can be used in more general scenarios as well. However, please note that this is a lower-level construct, and it's usually better to use higher-level abstractions like PersistentVolumeClaim or StorageClass unless you have a very specific reason to directly create VolumeAttachment objects.")
}
