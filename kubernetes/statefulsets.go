package kubernetes

import (
	"context"
	"fmt"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllK8sStatefulSets(namespaceName string, contextId *string) []v1.StatefulSet {
	result := []v1.StatefulSet{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return result
	}
	statefulSetList, err := provider.ClientSet.AppsV1().StatefulSets(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllStatefulSets ERROR: %s", err.Error())
		return result
	}

	for _, statefulSet := range statefulSetList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, statefulSet.ObjectMeta.Namespace) {
			statefulSet.Kind = "StatefulSet"
			result = append(result, statefulSet)
		}
	}
	return result
}

func AllStatefulSets(namespaceName string, contextId *string) utils.K8sWorkloadResult {
	result := []v1.StatefulSet{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	statefulSetList, err := provider.ClientSet.AppsV1().StatefulSets(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllStatefulSets ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, statefulSet := range statefulSetList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, statefulSet.ObjectMeta.Namespace) {
			statefulSet.Kind = "StatefulSet"
			result = append(result, statefulSet)
		}
	}
	return WorkloadResult(result, nil)
}

func GetStatefulSet(namespaceName string, name string, contextId *string) (*v1.StatefulSet, error) {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return nil, err
	}
	return provider.ClientSet.AppsV1().StatefulSets(namespaceName).Get(context.TODO(), name, metav1.GetOptions{})
}

func UpdateK8sStatefulset(data v1.StatefulSet, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.AppsV1().StatefulSets(data.Namespace)
	res, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func DeleteK8sStatefulset(data v1.StatefulSet, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.AppsV1().StatefulSets(data.Namespace)
	err = client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sStatefulsetBy(namespace string, name string, contextId *string) error {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return err
	}
	client := provider.ClientSet.AppsV1().StatefulSets(namespace)
	return client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sStatefulset(namespace string, name string, contextId *string) utils.K8sWorkloadResult {
	cmd := utils.RunOnLocalShell(fmt.Sprintf("kubectl describe statefulset %s -n %s%s", name, namespace, ContextFlag(contextId)))

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sStatefulset(data v1.StatefulSet, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.AppsV1().StatefulSets(data.Namespace)
	res, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func NewK8sStatefulset() K8sNewWorkload {
	return NewWorkload(
		RES_STATEFUL_SET,
		utils.InitStatefulsetYaml(),
		"StatefulSets are intended to be used with stateful applications and distributed systems. They manage the deployment and scaling of a set of Pods and provide guarantees about the ordering and uniqueness of these Pods. In this example, a StatefulSet named 'web' is created, which runs 3 replicas of the nginx container. Each pod is available through the service named 'nginx', and each has a single PersistentVolumeClaim, 'www', which requests a 1Gi storage volume.")
}
