package kubernetes

import (
	"context"
	"fmt"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	v1job "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllJobs(namespaceName string, contextId *string) []v1job.Job {
	result := []v1job.Job{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return result
	}
	jobList, err := provider.ClientSet.BatchV1().Jobs(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllJobs ERROR: %s", err.Error())
		return result
	}

	for _, job := range jobList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, job.ObjectMeta.Namespace) {
			job.Kind = "Job"
			result = append(result, job)
		}
	}
	return result
}

func AllK8sJobs(namespaceName string, contextId *string) utils.K8sWorkloadResult {
	result := []v1job.Job{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	jobList, err := provider.ClientSet.BatchV1().Jobs(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllJobs ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, job := range jobList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, job.ObjectMeta.Namespace) {
			job.Kind = "Job"
			result = append(result, job)
		}
	}
	return WorkloadResult(result, nil)
}

func GetJob(namespaceName string, name string, contextId *string) (*v1job.Job, error) {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return nil, err
	}
	return provider.ClientSet.BatchV1().Jobs(namespaceName).Get(context.TODO(), name, metav1.GetOptions{})
}

func UpdateK8sJob(data v1job.Job, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.BatchV1().Jobs(data.Namespace)
	res, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func DeleteK8sJob(data v1job.Job, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.BatchV1().Jobs(data.Namespace)
	err = client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sJobBy(namespace string, name string, contextId *string) error {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return err
	}
	client := provider.ClientSet.BatchV1().Jobs(namespace)
	return client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sJob(namespace string, name string, contextId *string) utils.K8sWorkloadResult {
	cmd := utils.RunOnLocalShell(fmt.Sprintf("kubectl describe job %s -n %s%s", name, namespace, ContextFlag(contextId)))

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sJob(data v1job.Job, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.BatchV1().Jobs(data.Namespace)
	res, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func NewK8sJob() K8sNewWorkload {
	return NewWorkload(
		RES_JOB,
		utils.InitJobYaml(),
		"A Job creates one or more Pods and ensures that a specified number of them successfully terminate. As pods successfully complete, the Job tracks the successful completions. In this example, a Job named 'my-job' is created. It will create a pod that runs a single container using the 'busybox' image. When the container starts, it will run the command sh -c 'echo Hello, mogenius! && sleep 30'. If the job fails, Kubernetes will try to restart it up to 4 times.")
}
