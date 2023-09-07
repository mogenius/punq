package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func UpdateK8sDeployment(data v1.Deployment) utils.HttpResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.AppsV1().Deployments(data.Namespace)
	_, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sDeployment(data v1.Deployment) utils.HttpResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.AppsV1().Deployments(data.Namespace)
	err := client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func AllDeployments(namespaceName string) []v1.Deployment {
	result := []v1.Deployment{}

	provider := NewKubeProvider()
	deploymentList, err := provider.ClientSet.AppsV1().Deployments(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllDeployments ERROR: %s", err.Error())
		return result
	}

	for _, deployment := range deploymentList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, deployment.ObjectMeta.Namespace) {
			result = append(result, deployment)
		}
	}
	return result
}

func AllK8sDeployments(namespaceName string) utils.HttpResult {
	result := []v1.Deployment{}

	provider := NewKubeProvider()
	deploymentList, err := provider.ClientSet.AppsV1().Deployments(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("AllDeployments ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, deployment := range deploymentList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, deployment.ObjectMeta.Namespace) {
			result = append(result, deployment)
		}
	}
	return WorkloadResult(result, nil)
}

func DescribeK8sDeployment(namespace string, name string) utils.HttpResult {
	cmd := exec.Command("kubectl", "describe", "deployment", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sDeployment(data v1.Deployment) utils.HttpResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.AppsV1().Deployments(data.Namespace)
	_, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func NewK8sDeployment() K8sNewWorkload {
	return NewWorkload(
		RES_DEPLOYMENT,
		utils.InitDeploymentYaml(),
		"A Deployment provides declarative updates for Pods and ReplicaSets. You describe a desired state in a Deployment, and the Deployment controller changes the actual state to the desired state at a controlled rate. In this example, a Deployment named 'my-app-deployment' is created. It will create 3 replicas of the pod, each running a single container from the 'my-app-image:1.0.0' image and exposing port 8080.")
}
