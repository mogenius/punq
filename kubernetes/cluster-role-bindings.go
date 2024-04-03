package kubernetes

import (
	"context"
	"fmt"

	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllClusterRoleBindings(contextId *string) []v1.ClusterRoleBinding {
	result := []v1.ClusterRoleBinding{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return result
	}
	rolesList, err := provider.ClientSet.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllClusterRoleBindings ERROR: %s", err.Error())
		return result
	}

	for _, role := range rolesList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, role.ObjectMeta.Namespace) {
			role.Kind = "ClusterRoleBinding"
			result = append(result, role)
		}
	}
	return result
}

func AllK8sClusterRoleBindings(contextId *string) utils.K8sWorkloadResult {
	result := []v1.ClusterRoleBinding{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	rolesList, err := provider.ClientSet.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllClusterRoleBindings ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, role := range rolesList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, role.ObjectMeta.Namespace) {
			role.Kind = "ClusterRoleBinding"
			result = append(result, role)
		}
	}
	return WorkloadResult(result, nil)
}

func GetClusterRoleBinding(name string, contextId *string) (*v1.ClusterRoleBinding, error) {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return nil, err
	}
	return provider.ClientSet.RbacV1().ClusterRoleBindings().Get(context.TODO(), name, metav1.GetOptions{})
}

func UpdateK8sClusterRoleBinding(data v1.ClusterRoleBinding, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.RbacV1().ClusterRoleBindings()
	res, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func DeleteK8sClusterRoleBinding(data v1.ClusterRoleBinding, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.RbacV1().ClusterRoleBindings()
	err = client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sClusterRoleBindingBy(name string, contextId *string) error {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return err
	}
	client := provider.ClientSet.RbacV1().ClusterRoleBindings()
	return client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sClusterRoleBinding(name string, contextId *string) utils.K8sWorkloadResult {
	cmd := utils.RunOnLocalShell(fmt.Sprintf("kubectl describe clusterrolebinding %s%s", name, ContextFlag(contextId)))

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sClusterRoleBinding(data v1.ClusterRoleBinding, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.RbacV1().ClusterRoleBindings()
	res, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func NewK8sClusterRoleBinding() K8sNewWorkload {
	return NewWorkload(
		RES_CLUSTER_ROLE_BINDING,
		utils.InitClusterRoleBindingYaml(),
		"A ClusterRoleBinding binds a ClusterRole to a group of subjects, granting them the permissions defined by the ClusterRole at the cluster level. In this example, a ClusterRoleBinding named 'my-cluster-role-binding' is created. It binds the ClusterRole named 'my-cluster-role' to the group named 'my-group' in the 'rbac.authorization.k8s.io' API group.")
}
