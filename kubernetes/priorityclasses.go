package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	v1 "k8s.io/api/scheduling/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllPriorityClasses(namespaceName string) K8sWorkloadResult {
	result := []v1.PriorityClass{}

	provider := NewKubeProvider()
	pcList, err := provider.ClientSet.SchedulingV1().PriorityClasses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllPriorityClasses ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, roleBinding := range pcList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, roleBinding.ObjectMeta.Namespace) {
			result = append(result, roleBinding)
		}
	}
	return WorkloadResult(result, nil)
}

func UpdateK8sPriorityClass(data v1.PriorityClass) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	pcClient := kubeProvider.ClientSet.SchedulingV1().PriorityClasses()
	_, err := pcClient.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sPriorityClass(data v1.PriorityClass) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	pcClient := kubeProvider.ClientSet.SchedulingV1().PriorityClasses()
	err := pcClient.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sPriorityClass(name string) K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "priorityclasses", name)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func NewK8sPriorityClass() K8sNewWorkload {
	return NewWorkload(
		RES_PRIORITYCLASSES,
		utils.InitPriorityClassYaml(),
		"PriorityClass is used to assign priority to Pods and allows the Kubernetes scheduler to make scheduling decisions based on the relative priorities of the Pods. In this example, a PriorityClass named 'high-priority' is created. It has a value of 1000000, indicating high priority. The globalDefault field is set to false, meaning it is not the default priority class for all Pods. The description field provides a brief description of the PriorityClass.")
}