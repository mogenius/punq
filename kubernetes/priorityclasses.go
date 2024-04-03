package kubernetes

import (
	"context"
	"fmt"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	v1 "k8s.io/api/scheduling/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllPriorityClasses(contextId *string) []v1.PriorityClass {
	result := []v1.PriorityClass{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return result
	}
	pcList, err := provider.ClientSet.SchedulingV1().PriorityClasses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllPriorityClasses ERROR: %s", err.Error())
		return result
	}

	for _, prioClass := range pcList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, prioClass.ObjectMeta.Namespace) {
			prioClass.Kind = "PriorityClass"
			prioClass.APIVersion = "scheduling.k8s.io/v1"
			result = append(result, prioClass)
		}
	}
	return result
}

func AllK8sPriorityClasses(contextId *string) utils.K8sWorkloadResult {
	result := []v1.PriorityClass{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	pcList, err := provider.ClientSet.SchedulingV1().PriorityClasses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllPriorityClasses ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, prioClass := range pcList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, prioClass.ObjectMeta.Namespace) {
			prioClass.Kind = "PriorityClass"
			prioClass.APIVersion = "scheduling.k8s.io/v1"
			result = append(result, prioClass)
		}
	}
	return WorkloadResult(result, nil)
}

func GetPriorityClass(name string, contextId *string) (*v1.PriorityClass, error) {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return nil, err
	}
	return provider.ClientSet.SchedulingV1().PriorityClasses().Get(context.TODO(), name, metav1.GetOptions{})
}

func UpdateK8sPriorityClass(data v1.PriorityClass, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.SchedulingV1().PriorityClasses()
	res, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func DeleteK8sPriorityClass(data v1.PriorityClass, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.SchedulingV1().PriorityClasses()
	err = client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sPriorityClassBy(name string, contextId *string) error {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return err
	}
	client := provider.ClientSet.SchedulingV1().PriorityClasses()
	return client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sPriorityClass(name string, contextId *string) utils.K8sWorkloadResult {
	cmd := utils.RunOnLocalShell(fmt.Sprintf("kubectl describe priorityclasses %s%s", name, ContextFlag(contextId)))

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sPriorityClass(data v1.PriorityClass, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.SchedulingV1().PriorityClasses()
	res, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func NewK8sPriorityClass() K8sNewWorkload {
	return NewWorkload(
		RES_PRIORITY_CLASS,
		utils.InitPriorityClassYaml(),
		"PriorityClass is used to assign priority to Pods and allows the Kubernetes scheduler to make scheduling decisions based on the relative priorities of the Pods. In this example, a PriorityClass named 'high-priority' is created. It has a value of 1000000, indicating high priority. The globalDefault field is set to false, meaning it is not the default priority class for all Pods. The description field provides a brief description of the PriorityClass.")
}
