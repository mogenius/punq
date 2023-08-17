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

func ListAllNamespaceNames() []string {
	result := []string{}

	kubeProvider := NewKubeProvider()
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

func ListAllNamespace() []v1.Namespace {
	result := []v1.Namespace{}

	kubeProvider := NewKubeProvider()
	namespaceClient := kubeProvider.ClientSet.CoreV1().Namespaces()

	namespaceList, err := namespaceClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("ListAllNamespace ERROR: %s", err.Error())
		return result
	}

	result = append(result, namespaceList.Items...)

	return result
}

func ListK8sNamespaces(namespaceName string) K8sWorkloadResult {
	result := []v1.Namespace{}

	kubeProvider := NewKubeProvider()
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

func DeleteK8sNamespace(data v1.Namespace) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	namespaceClient := kubeProvider.ClientSet.CoreV1().Namespaces()
	err := namespaceClient.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sNamespace(name string) K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "namespace", name)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func NamespaceExists(namespaceName string) (bool, error) {
	kubeProvider := NewKubeProvider()
	namespaceClient := kubeProvider.ClientSet.CoreV1().Namespaces()
	ns, err := namespaceClient.Get(context.TODO(), namespaceName, metav1.GetOptions{})
	return (ns != nil && err == nil), err
}

func NewK8sNamespace() K8sNewWorkload {
	return NewWorkload(
		RES_NAMESPACE,
		utils.InitNamespaceYaml(),
		"A Namespace is a way to divide cluster resources between multiple users. They are intended for use in environments with many users spread across multiple teams, or projects. In this example, a Namespace named 'my-namespace' is created. Namespaces provide a scope for names. Names of resources need to be unique within a namespace but not across namespaces. Namespaces can not be nested inside one another and each Kubernetes resource can only be in one namespace.")
}
