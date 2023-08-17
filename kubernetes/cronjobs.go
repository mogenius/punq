package kubernetes

import (
	"context"
	"os/exec"

	"punq/logger"
	"punq/utils"

	v1 "k8s.io/api/batch/v1"
	v1job "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllCronjobs(namespaceName string) K8sWorkloadResult {
	result := []v1job.CronJob{}

	provider := NewKubeProvider()
	cronJobList, err := provider.ClientSet.BatchV1().CronJobs(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllCronjobs ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, cronJob := range cronJobList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, cronJob.ObjectMeta.Namespace) {
			result = append(result, cronJob)
		}
	}
	return WorkloadResult(result, nil)
}

func UpdateK8sCronJob(data v1.CronJob) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	cronJobClient := kubeProvider.ClientSet.BatchV1().CronJobs(data.Namespace)
	_, err := cronJobClient.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sCronJob(data v1job.CronJob) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	jobClient := kubeProvider.ClientSet.BatchV1().CronJobs(data.Namespace)
	err := jobClient.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sCronJob(namespace string, name string) K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "cronjob", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func NewK8sCronJob() K8sNewWorkload {
	return NewWorkload(
		RES_CRON_JOB,
		utils.InitCronJobYaml(),
		"A CronJob creates Jobs on a repeating schedule, like the cron utility in Unix-like systems. In this example, a CronJob named 'my-cronjob' is created. It runs a Job every minute. Each Job creates a Pod with a single container from the 'my-cronjob-image' image.")
}
