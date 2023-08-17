package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	v1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllServiceAccounts(namespaceName string) K8sWorkloadResult {
	result := []v1.ServiceAccount{}

	provider := NewKubeProvider()
	rolesList, err := provider.ClientSet.CoreV1().ServiceAccounts(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllServiceAccounts ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, role := range rolesList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, role.ObjectMeta.Namespace) {
			result = append(result, role)
		}
	}
	return WorkloadResult(result, nil)
}

func UpdateK8sServiceAccount(data v1.ServiceAccount) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	roleClient := kubeProvider.ClientSet.CoreV1().ServiceAccounts(data.Namespace)
	_, err := roleClient.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sServiceAccount(data v1.ServiceAccount) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	roleClient := kubeProvider.ClientSet.CoreV1().ServiceAccounts(data.Namespace)
	err := roleClient.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sServiceAccount(namespace string, name string) K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "serviceaccount", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func NewK8sServiceAccount() K8sNewWorkload {
	return NewWorkload(
		RES_SERVICE_ACCOUNT,
		utils.InitServiceAccountExampleYaml(),
		"A ServiceAccount is an object within Kubernetes that provides an identity for processes that run in a Pod. In this example, a service account named 'my-serviceaccount' is created in the 'my-namespace' namespace. If the namespace field is omitted, Kubernetes will assume the default namespace. ServiceAccounts aren't just limited to the following fields. They can have secrets and imagePullSecrets associated with them. The secrets field is used to attach arbitrary secrets to the service account which can then be mounted into pods. The imagePullSecrets field is used to specify Docker registry credentials.")
}
