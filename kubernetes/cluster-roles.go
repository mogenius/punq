package kubernetes

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	v1 "k8s.io/api/rbac/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllClusterRoles(contextId *string) utils.K8sWorkloadResult {
	result := []v1.ClusterRole{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
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

func GetClusterRole(name string, contextId *string) (*v1.ClusterRole, error) {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return nil, err
	}
	return provider.ClientSet.RbacV1().ClusterRoles().Get(context.TODO(), name, metav1.GetOptions{})
}

func UpdateK8sClusterRole(data v1.ClusterRole, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.RbacV1().ClusterRoles()
	res, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func DeleteK8sClusterRole(data v1.ClusterRole, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.RbacV1().ClusterRoles()
	err = client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sClusterRoleBy(name string, contextId *string) error {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return err
	}
	client := provider.ClientSet.RbacV1().ClusterRoles()
	return client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sClusterRole(name string, contextId *string) utils.K8sWorkloadResult {
	cmd := exec.Command("/bin/ash", "-c", fmt.Sprintf("/usr/local/bin/kubectl describe clusterrole %s%s", name, ContextFlag(contextId)))

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sClusterRole(data v1.ClusterRole, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.RbacV1().ClusterRoles()
	res, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func NewK8sClusterRole() K8sNewWorkload {
	return NewWorkload(
		RES_CLUSTER_ROLE,
		utils.InitClusterRoleYaml(),
		"A ClusterRole is a non-namespaced resource that defines permissions for accessing cluster-level resources in Kubernetes. In this example, a ClusterRole named 'my-cluster-role' is created. It grants permissions to access and perform various actions on pods and deployments including get, list, watch, create, update, and delete. It also grants permissions to access roles and rolebindings (from the 'rbac.authorization.k8s.io' API group) including get, list, and watch.")
}
