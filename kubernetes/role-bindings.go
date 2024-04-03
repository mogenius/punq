package kubernetes

import (
	"context"
	"fmt"

	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"

	v1 "k8s.io/api/rbac/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllRoleBindings(namespaceName string, contextId *string) []v1.RoleBinding {
	result := []v1.RoleBinding{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return result
	}
	rolesList, err := provider.ClientSet.RbacV1().RoleBindings(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllBindings ERROR: %s", err.Error())
		return result
	}

	for _, roleBinding := range rolesList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, roleBinding.ObjectMeta.Namespace) {
			roleBinding.Kind = "RoleBinding"
			roleBinding.APIVersion = "rbac.authorization.k8s.io/v1"
			result = append(result, roleBinding)
		}
	}
	return result
}

func AllK8sRoleBindings(namespaceName string, contextId *string) utils.K8sWorkloadResult {
	result := []v1.RoleBinding{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	rolesList, err := provider.ClientSet.RbacV1().RoleBindings(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllBindings ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, roleBinding := range rolesList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, roleBinding.ObjectMeta.Namespace) {
			roleBinding.Kind = "RoleBinding"
			roleBinding.APIVersion = "rbac.authorization.k8s.io/v1"
			result = append(result, roleBinding)
		}
	}
	return WorkloadResult(result, nil)
}

func GetRoleBinding(namespaceName string, name string, contextId *string) (*v1.RoleBinding, error) {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return nil, err
	}
	return provider.ClientSet.RbacV1().RoleBindings(namespaceName).Get(context.TODO(), name, metav1.GetOptions{})
}

func UpdateK8sRoleBinding(data v1.RoleBinding, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.RbacV1().RoleBindings(data.Namespace)
	res, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func DeleteK8sRoleBinding(data v1.RoleBinding, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.RbacV1().RoleBindings(data.Namespace)
	err = client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sRoleBindingBy(namespace string, name string, contextId *string) error {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return err
	}
	client := provider.ClientSet.RbacV1().RoleBindings(namespace)
	return client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sRoleBinding(namespace string, name string, contextId *string) utils.K8sWorkloadResult {
	cmd := utils.RunOnLocalShell(fmt.Sprintf("kubectl describe rolebinding %s -n %s%s", name, namespace, ContextFlag(contextId)))

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sRoleBinding(data v1.RoleBinding, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.RbacV1().RoleBindings(data.Namespace)
	res, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func NewK8sRoleBinding() K8sNewWorkload {
	return NewWorkload(
		RES_ROLE_BINDING,
		utils.InitRoleBindingYaml(),
		"RoleBindings in Kubernetes provide a way to bind a Role or ClusterRole to a set of subjects (which can be Users, Groups, or ServiceAccounts) within a particular namespace. In this example, a RoleBinding named 'read-pods' is created in the 'default' namespace. This RoleBinding binds the Role 'pod-reader' to the user 'jane'.")
}
