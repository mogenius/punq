package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	v1 "k8s.io/api/coordination/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllLeases(namespaceName string) K8sWorkloadResult {
	result := []v1.Lease{}

	provider := NewKubeProvider()
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

func UpdateK8sLease(data v1.Lease) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	leasesClient := kubeProvider.ClientSet.CoordinationV1().Leases(data.Namespace)
	_, err := leasesClient.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sLease(data v1.Lease) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	leasesClient := kubeProvider.ClientSet.CoordinationV1().Leases(data.Namespace)
	err := leasesClient.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sLease(namespace string, name string) K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "lease", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func NewK8sLease() K8sNewWorkload {
	return NewWorkload(
		RES_LEASES,
		utils.InitLeaseYaml(),
		"A Lease is a simple object that allows coordination between different components or processes running in a cluster. In this example, a Lease named 'my-lease' is created in the 'my-namespace' namespace. The Lease is associated with the identity 'my-identity'. It has a duration of 60 seconds and should be renewed at the specified time. Leases are typically used for coordination and to ensure that only one instance of a component or process is active at a given time.")
}
