package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"

	v1 "k8s.io/api/rbac/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllRoles(namespaceName string) utils.K8sWorkloadResult {
	result := []v1.Role{}

	provider := NewKubeProvider()
	rolesList, err := provider.ClientSet.RbacV1().Roles(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllRoles ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, role := range rolesList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, role.ObjectMeta.Namespace) {
			result = append(result, role)
		}
	}
	return WorkloadResult(result, nil)
}

func GetRole(namespaceName string, name string) (*v1.Role, error) {
	provider := NewKubeProvider()
	return provider.ClientSet.RbacV1().Roles(namespaceName).Get(context.TODO(), name, metav1.GetOptions{})
}

func UpdateK8sRole(data v1.Role) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.RbacV1().Roles(data.Namespace)
	_, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sRole(data v1.Role) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.RbacV1().Roles(data.Namespace)
	err := client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sRoleBy(namespace string, name string) error {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.RbacV1().Roles(namespace)
	return client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sRole(namespace string, name string) utils.K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "role", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sRole(data v1.Role) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.RbacV1().Roles(data.Namespace)
	_, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func NewK8sRole() K8sNewWorkload {
	return NewWorkload(
		RES_ROLE,
		utils.InitRoleYaml(),
		"Roles in Kubernetes provide a mechanism to define authorizations within a particular namespace. In this example, a Role named 'pod-reader' is created in the 'default' namespace. This Role has permissions to 'get', 'watch', and 'list' Pods. Please note, Roles define permissions within a specific namespace. If you want to define permissions cluster-wide, you would use a ClusterRole instead of a Role.")
}
