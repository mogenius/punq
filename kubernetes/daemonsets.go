package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllDaemonsets(namespaceName string) []v1.DaemonSet {
	result := []v1.DaemonSet{}

	provider := NewKubeProvider()
	daemonsetList, err := provider.ClientSet.AppsV1().DaemonSets(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllDaemonsets ERROR: %s", err.Error())
		return result
	}

	for _, daemonset := range daemonsetList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, daemonset.ObjectMeta.Namespace) {
			result = append(result, daemonset)
		}
	}
	return result
}

func AllK8sDaemonsets(namespaceName string) K8sWorkloadResult {
	result := []v1.DaemonSet{}

	provider := NewKubeProvider()
	daemonsetList, err := provider.ClientSet.AppsV1().DaemonSets(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllDaemonsets ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, daemonset := range daemonsetList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, daemonset.ObjectMeta.Namespace) {
			result = append(result, daemonset)
		}
	}
	return WorkloadResult(result, nil)
}

func UpdateK8sDaemonSet(data v1.DaemonSet) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	daemonSetClient := kubeProvider.ClientSet.AppsV1().DaemonSets(data.Namespace)
	_, err := daemonSetClient.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sDaemonSet(data v1.DaemonSet) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	daemonSetClient := kubeProvider.ClientSet.AppsV1().DaemonSets(data.Namespace)
	err := daemonSetClient.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sDaemonSet(namespace string, name string) K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "daemonset", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func NewK8sDaemonSet() K8sNewWorkload {
	return NewWorkload(
		RES_DAEMON_SET,
		utils.InitDaemonsetYaml(),
		"A DaemonSet ensures that all (or some) nodes run a copy of a Pod. As nodes are added to the cluster, Pods are added to them. As nodes are removed from the cluster, those Pods are garbage collected. In this example, a DaemonSet named 'my-daemonset' is created. It ensures that each node in the cluster runs a Pod with a single container from the 'my-daemonset-image' image and exposing port 8080.")
}
