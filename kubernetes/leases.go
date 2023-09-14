package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	v1 "k8s.io/api/coordination/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllLeases(namespaceName string, contextId *string) utils.K8sWorkloadResult {
	result := []v1.Lease{}

	provider := NewKubeProvider(contextId)
	rolesList, err := provider.ClientSet.CoordinationV1().Leases(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllLeases ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, role := range rolesList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, role.ObjectMeta.Namespace) {
			result = append(result, role)
		}
	}
	return WorkloadResult(result, nil)
}

func GetLeas(namespaceName string, name string, contextId *string) (*v1.Lease, error) {
	provider := NewKubeProvider(contextId)
	return provider.ClientSet.CoordinationV1().Leases(namespaceName).Get(context.TODO(), name, metav1.GetOptions{})
}

func UpdateK8sLease(data v1.Lease, contextId *string) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider(contextId)
	client := kubeProvider.ClientSet.CoordinationV1().Leases(data.Namespace)
	_, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sLease(data v1.Lease, contextId *string) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider(contextId)
	client := kubeProvider.ClientSet.CoordinationV1().Leases(data.Namespace)
	err := client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sLeaseBy(namespace string, name string, contextId *string) error {
	kubeProvider := NewKubeProvider(contextId)
	client := kubeProvider.ClientSet.CoordinationV1().Leases(namespace)
	return client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sLease(namespace string, name string, contextId *string) utils.K8sWorkloadResult {
	cmd := exec.Command("kubectl", ContextFlag(contextId), "describe", "lease", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sLease(data v1.Lease, contextId *string) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider(contextId)
	client := kubeProvider.ClientSet.CoordinationV1().Leases(data.Namespace)
	_, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func NewK8sLease() K8sNewWorkload {
	return NewWorkload(
		RES_LEASES,
		utils.InitLeaseYaml(),
		"A Lease is a simple object that allows coordination between different components or processes running in a cluster. In this example, a Lease named 'my-lease' is created in the 'my-namespace' namespace. The Lease is associated with the identity 'my-identity'. It has a duration of 60 seconds and should be renewed at the specified time. Leases are typically used for coordination and to ensure that only one instance of a component or process is active at a given time.")
}
