package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	v1job "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllJobs(namespaceName string) K8sWorkloadResult {
	result := []v1job.Job{}

	provider := NewKubeProvider()
	jobList, err := provider.ClientSet.BatchV1().Jobs(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllJobs ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, job := range jobList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, job.ObjectMeta.Namespace) {
			result = append(result, job)
		}
	}
	return WorkloadResult(result, nil)
}

func UpdateK8sJob(data v1job.Job) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.BatchV1().Jobs(data.Namespace)
	_, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sJob(data v1job.Job) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.BatchV1().Jobs(data.Namespace)
	err := client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sJob(namespace string, name string) K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "job", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sJob(data v1job.Job) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	client := kubeProvider.ClientSet.BatchV1().Jobs(data.Namespace)
	_, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func NewK8sJob() K8sNewWorkload {
	return NewWorkload(
		RES_JOB,
		utils.InitJobYaml(),
		"A Job creates one or more Pods and ensures that a specified number of them successfully terminate. As pods successfully complete, the Job tracks the successful completions. In this example, a Job named 'my-job' is created. It will create a pod that runs a single container using the 'busybox' image. When the container starts, it will run the command sh -c 'echo Hello, mogenius! && sleep 30'. If the job fails, Kubernetes will try to restart it up to 4 times.")
}
