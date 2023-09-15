package kubernetes

import (
	"context"
	"os/exec"
	"strings"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListAllNamespaceNames(contextId *string) []string {
	result := []string{}

	kubeProvider := NewKubeProvider(contextId)
	namespaceClient := kubeProvider.ClientSet.CoreV1().Namespaces()

	namespaceList, err := namespaceClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("ListAll ERROR: %s", err.Error())
		return result
	}

	for _, ns := range namespaceList.Items {
		result = append(result, ns.Name)
	}

	return result
}

func ListAllNamespace(contextId *string) []v1.Namespace {
	result := []v1.Namespace{}

	kubeProvider := NewKubeProvider(contextId)
	namespaceClient := kubeProvider.ClientSet.CoreV1().Namespaces()

	namespaceList, err := namespaceClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("ListAllNamespace ERROR: %s", err.Error())
		return result
	}

	result = append(result, namespaceList.Items...)

	return result
}

func GetNamespace(name string, contextId *string) (*v1.Namespace, error) {
	kubeProvider := NewKubeProvider(contextId)
	namespaceClient := kubeProvider.ClientSet.CoreV1().Namespaces()
	return namespaceClient.Get(context.TODO(), name, metav1.GetOptions{})
}

func ListK8sNamespaces(namespaceName string, contextId *string) utils.K8sWorkloadResult {
	result := []v1.Namespace{}

	kubeProvider := NewKubeProvider(contextId)
	namespaceClient := kubeProvider.ClientSet.CoreV1().Namespaces()

	namespaceList, err := namespaceClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("ListAllNamespace ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, ns := range namespaceList.Items {
		if namespaceName == "" {
			result = append(result, ns)
		} else {
			if strings.HasPrefix(ns.Name, namespaceName) {
				result = append(result, ns)
			}
		}
	}

	return WorkloadResult(result, nil)
}

func DeleteK8sNamespace(data v1.Namespace, contextId *string) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider(contextId)
	client := kubeProvider.ClientSet.CoreV1().Namespaces()
	err := client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sNamespaceBy(name string, contextId *string) error {
	kubeProvider := NewKubeProvider(contextId)
	client := kubeProvider.ClientSet.CoreV1().Namespaces()
	return client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sNamespace(name string, contextId *string) utils.K8sWorkloadResult {
	cmd := exec.Command("kubectl", ContextFlag(contextId), "describe", "namespace", name)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func NamespaceExists(namespaceName string, contextId *string) (bool, error) {
	kubeProvider := NewKubeProvider(contextId)
	namespaceClient := kubeProvider.ClientSet.CoreV1().Namespaces()
	ns, err := namespaceClient.Get(context.TODO(), namespaceName, metav1.GetOptions{})
	return (ns != nil && err == nil), err
}

func CreateK8sNamespace(data v1.Namespace, contextId *string) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider(contextId)
	client := kubeProvider.ClientSet.CoreV1().Namespaces()
	_, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func UpdateK8sNamespace(data v1.Namespace, contextId *string) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider(contextId)
	client := kubeProvider.ClientSet.CoreV1().Namespaces()
	_, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func NewK8sNamespace() K8sNewWorkload {
	return NewWorkload(
		RES_NAMESPACE,
		utils.InitNamespaceYaml(),
		"A Namespace is a way to divide cluster resources between multiple users. They are intended for use in environments with many users spread across multiple teams, or projects. In this example, a Namespace named 'my-namespace' is created. Namespaces provide a scope for names. Names of resources need to be unique within a namespace but not across namespaces. Namespaces can not be nested inside one another and each Kubernetes resource can only be in one namespace.")
}
