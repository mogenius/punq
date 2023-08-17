package kubernetes

import (
	"context"
	"os/exec"

	"punq/logger"
	"punq/utils"

	v1 "k8s.io/api/rbac/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllClusterRoles(namespaceName string) K8sWorkloadResult {
	result := []v1.ClusterRole{}

	provider := NewKubeProvider()
	rolesList, err := provider.ClientSet.RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllClusterRoles ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, role := range rolesList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, role.ObjectMeta.Namespace) {
			result = append(result, role)
		}
	}
	return WorkloadResult(result, nil)
}

func UpdateK8sClusterRole(data v1.ClusterRole) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	roleClient := kubeProvider.ClientSet.RbacV1().ClusterRoles()
	_, err := roleClient.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sClusterRole(data v1.ClusterRole) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	roleClient := kubeProvider.ClientSet.RbacV1().ClusterRoles()
	err := roleClient.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sClusterRole(name string) K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "clusterrole", name)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func NewK8sClusterRole() K8sNewWorkload {
	return NewWorkload(
		RES_CLUSTER_ROLE,
		utils.InitClusterRoleYaml(),
		"A ClusterRole is a non-namespaced resource that defines permissions for accessing cluster-level resources in Kubernetes. In this example, a ClusterRole named 'my-cluster-role' is created. It grants permissions to access and perform various actions on pods and deployments including get, list, watch, create, update, and delete. It also grants permissions to access roles and rolebindings (from the 'rbac.authorization.k8s.io' API group) including get, list, and watch.")
}
