package kubernetes

import (
	"context"
	"os/exec"

	"punq/logger"
	"punq/utils"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllReplicasets(namespaceName string) []v1.ReplicaSet {
	result := []v1.ReplicaSet{}

	provider := NewKubeProvider()
	replicaSetList, err := provider.ClientSet.AppsV1().ReplicaSets(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllReplicasets ERROR: %s", err.Error())
		return result
	}

	for _, replicaSet := range replicaSetList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, replicaSet.ObjectMeta.Namespace) {
			result = append(result, replicaSet)
		}
	}
	return result
}

func AllK8sReplicasets(namespaceName string) K8sWorkloadResult {
	result := []v1.ReplicaSet{}

	provider := NewKubeProvider()
	replicaSetList, err := provider.ClientSet.AppsV1().ReplicaSets(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllReplicasets ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, replicaSet := range replicaSetList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, replicaSet.ObjectMeta.Namespace) {
			result = append(result, replicaSet)
		}
	}
	return WorkloadResult(result, nil)
}

func UpdateK8sReplicaset(data v1.ReplicaSet) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	replicasetClient := kubeProvider.ClientSet.AppsV1().ReplicaSets(data.Namespace)
	_, err := replicasetClient.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sReplicaset(data v1.ReplicaSet) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	replicasetClient := kubeProvider.ClientSet.AppsV1().ReplicaSets(data.Namespace)
	err := replicasetClient.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sReplicaset(namespace string, name string) K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "replicaset", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func NewK8sReplicaSet() K8sNewWorkload {
	return NewWorkload(
		RES_REPLICA_SET,
		utils.InitReplicaSetYaml(),
		"A ReplicaSet's purpose is to maintain a stable set of replica Pods running at any given time. It's often used to guarantee the availability of a specified number of identical Pods. In this example, a ReplicaSet named 'my-replicaset' is created to ensure that exactly three Pods with labels app=myapp and tier=frontend are running at all times. Please note, although ReplicaSets are a powerful tool for maintaining sets of pods, Deployments are a higher-level concept that manage ReplicaSets and provide declarative updates to Pods along with a lot of other useful features. Hence, it's recommended to use Deployments instead of directly using ReplicaSets, unless you require custom update orchestration or don't require updates at all.")
}
