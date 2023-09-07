package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	v1 "k8s.io/api/rbac/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllClusterRoleBindings() utils.HttpResult {
	result := []v1.ClusterRoleBinding{}

	provider := NewKubeProvider()
	rolesList, err := provider.ClientSet.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllClusterRoleBindings ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, role := range rolesList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, role.ObjectMeta.Namespace) {
			result = append(result, role)
		}
	}
	return WorkloadResult(result, nil)
}

func UpdateK8sClusterRoleBinding(data v1.ClusterRoleBinding) utils.HttpResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.RbacV1().ClusterRoleBindings()
	_, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sClusterRoleBinding(data v1.ClusterRoleBinding) utils.HttpResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.RbacV1().ClusterRoleBindings()
	err := client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sClusterRoleBinding(name string) utils.HttpResult {
	cmd := exec.Command("kubectl", "describe", "clusterrolebinding", name)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sClusterRoleBinding(data v1.ClusterRoleBinding) utils.HttpResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.RbacV1().ClusterRoleBindings()
	_, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func NewK8sClusterRoleBinding() K8sNewWorkload {
	return NewWorkload(
		RES_CLUSTER_ROLE_BINDING,
		utils.InitClusterRoleBindingYaml(),
		"A ClusterRoleBinding binds a ClusterRole to a group of subjects, granting them the permissions defined by the ClusterRole at the cluster level. In this example, a ClusterRoleBinding named 'my-cluster-role-binding' is created. It binds the ClusterRole named 'my-cluster-role' to the group named 'my-group' in the 'rbac.authorization.k8s.io' API group.")
}
