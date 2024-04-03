package kubernetes

import (
	"context"
	"fmt"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	v1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllServiceAccounts(namespaceName string, contextId *string) []v1.ServiceAccount {
	result := []v1.ServiceAccount{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return result
	}
	srvAccList, err := provider.ClientSet.CoreV1().ServiceAccounts(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllServiceAccounts ERROR: %s", err.Error())
		return result
	}

	for _, srvAcc := range srvAccList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, srvAcc.ObjectMeta.Namespace) {
			srvAcc.Kind = "ServiceAccount"
			srvAcc.APIVersion = "v1"
			result = append(result, srvAcc)
		}
	}
	return result
}

func AllK8sServiceAccounts(namespaceName string, contextId *string) utils.K8sWorkloadResult {
	result := []v1.ServiceAccount{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	srvAccList, err := provider.ClientSet.CoreV1().ServiceAccounts(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllServiceAccounts ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, srvAcc := range srvAccList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, srvAcc.ObjectMeta.Namespace) {
			srvAcc.Kind = "ServiceAccount"
			srvAcc.APIVersion = "v1"
			result = append(result, srvAcc)
		}
	}
	return WorkloadResult(result, nil)
}

func GetServiceAccount(namespaceName string, name string, contextId *string) (*v1.ServiceAccount, error) {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return nil, err
	}
	return provider.ClientSet.CoreV1().ServiceAccounts(namespaceName).Get(context.TODO(), name, metav1.GetOptions{})
}

func UpdateK8sServiceAccount(data v1.ServiceAccount, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.CoreV1().ServiceAccounts(data.Namespace)
	res, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func DeleteK8sServiceAccount(data v1.ServiceAccount, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.CoreV1().ServiceAccounts(data.Namespace)
	err = client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sServiceAccountBy(namespace string, name string, contextId *string) error {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return err
	}
	client := provider.ClientSet.CoreV1().ServiceAccounts(namespace)
	return client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sServiceAccount(namespace string, name string, contextId *string) utils.K8sWorkloadResult {
	cmd := utils.RunOnLocalShell(fmt.Sprintf("kubectl describe serviceaccount %s -n %s%s", name, namespace, ContextFlag(contextId)))

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sServiceAccount(data v1.ServiceAccount, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.CoreV1().ServiceAccounts(data.Namespace)
	res, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func NewK8sServiceAccount() K8sNewWorkload {
	return NewWorkload(
		RES_SERVICE_ACCOUNT,
		utils.InitServiceAccountExampleYaml(),
		"A ServiceAccount is an object within Kubernetes that provides an identity for processes that run in a Pod. In this example, a service account named 'my-serviceaccount' is created in the 'my-namespace' namespace. If the namespace field is omitted, Kubernetes will assume the default namespace. ServiceAccounts aren't just limited to the following fields. They can have secrets and imagePullSecrets associated with them. The secrets field is used to attach arbitrary secrets to the service account which can then be mounted into pods. The imagePullSecrets field is used to specify Docker registry credentials.")
}
