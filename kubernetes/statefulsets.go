package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllStatefulSets(namespaceName string) K8sWorkloadResult {
	result := []v1.StatefulSet{}

	provider := NewKubeProvider()
	statefulSetList, err := provider.ClientSet.AppsV1().StatefulSets(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllStatefulSets ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, statefulSet := range statefulSetList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, statefulSet.ObjectMeta.Namespace) {
			result = append(result, statefulSet)
		}
	}
	return WorkloadResult(result, nil)
}

func UpdateK8sStatefulset(data v1.StatefulSet) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	statefulsetClient := kubeProvider.ClientSet.AppsV1().StatefulSets(data.Namespace)
	_, err := statefulsetClient.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sStatefulset(data v1.StatefulSet) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	statefulsetClient := kubeProvider.ClientSet.AppsV1().StatefulSets(data.Namespace)
	err := statefulsetClient.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sStatefulset(namespace string, name string) K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "statefulset", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func NewK8sStatefulset() K8sNewWorkload {
	return NewWorkload(
		RES_STATEFUL_SET,
		utils.InitStatefulsetYaml(),
		"StatefulSets are intended to be used with stateful applications and distributed systems. They manage the deployment and scaling of a set of Pods and provide guarantees about the ordering and uniqueness of these Pods. In this example, a StatefulSet named 'web' is created, which runs 3 replicas of the nginx container. Each pod is available through the service named 'nginx', and each has a single PersistentVolumeClaim, 'www', which requests a 1Gi storage volume.")
}
