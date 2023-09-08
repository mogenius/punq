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

func AllK8sDaemonsets(namespaceName string) utils.K8sWorkloadResult {
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

func GetK8sDaemonset(namespaceName string, name string) (*v1.DaemonSet, error) {
	provider := NewKubeProvider()
	return provider.ClientSet.AppsV1().DaemonSets(namespaceName).Get(context.TODO(), name, metav1.GetOptions{})
}

func UpdateK8sDaemonSet(data v1.DaemonSet) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.AppsV1().DaemonSets(data.Namespace)
	_, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sDaemonSet(data v1.DaemonSet) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.AppsV1().DaemonSets(data.Namespace)
	err := client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sDaemonSetBy(namespace string, name string) error {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.AppsV1().DaemonSets(namespace)
	return client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sDaemonSet(namespace string, name string) utils.K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "daemonset", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sDaemonSet(data v1.DaemonSet) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.AppsV1().DaemonSets(data.Namespace)
	_, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func NewK8sDaemonSet() K8sNewWorkload {
	return NewWorkload(
		RES_DAEMON_SET,
		utils.InitDaemonsetYaml(),
		"A DaemonSet ensures that all (or some) nodes run a copy of a Pod. As nodes are added to the cluster, Pods are added to them. As nodes are removed from the cluster, those Pods are garbage collected. In this example, a DaemonSet named 'my-daemonset' is created. It ensures that each node in the cluster runs a Pod with a single container from the 'my-daemonset-image' image and exposing port 8080.")
}
