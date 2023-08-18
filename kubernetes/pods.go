package kubernetes

import (
	"bytes"
	"context"
	"os/exec"
	"sort"
	"strings"
	"text/template"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServicePodExistsResult struct {
	PodExists bool `json:"podExists"`
}

func PodStatus(namespace string, name string, statusOnly bool) *v1.Pod {
	kubeProvider := NewKubeProvider()
	getOptions := metav1.GetOptions{}

	podClient := kubeProvider.ClientSet.CoreV1().Pods(namespace)

	pod, err := podClient.Get(context.TODO(), name, getOptions)
	if err != nil {
		logger.Log.Errorf("PodStatus Error: %s", err.Error())
		return nil
	}

	if statusOnly {
		filterStatus(pod)
	}

	return pod
}

func LastTerminatedStateIfAny(pod *v1.Pod) *v1.ContainerStateTerminated {
	if pod != nil {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			state := containerStatus.LastTerminationState

			if state.Terminated != nil {
				return state.Terminated
			}
		}
	}

	return nil
}

func LastTerminatedStateToString(terminatedState *v1.ContainerStateTerminated) string {
	if terminatedState == nil {
		return "Last State:	   nil\n"
	}

	tpl, err := template.New("state").Parse(
		"Last State:    Terminated\n" +
			"  Reason:      {{.Reason}}\n" +
			"  Message:     {{.Message}}\n" +
			"  Exit Code:   {{.ExitCode}}\n" +
			"  Started:     {{.StartedAt}}\n" +
			"  Finished:    {{.FinishedAt}}\n")
	if err != nil {
		logger.Log.Error(err.Error())
		return ""
	}

	buf := bytes.Buffer{}
	err = tpl.Execute(&buf, terminatedState)
	if err != nil {
		logger.Log.Error(err.Error())
		return ""
	}

	return buf.String()
}

func ServicePodStatus(namespace string, serviceName string) []v1.Pod {
	result := []v1.Pod{}
	kubeProvider := NewKubeProvider()

	podClient := kubeProvider.ClientSet.CoreV1().Pods(namespace)

	pods, err := podClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Error("ServicePodStatus Error: %s", err.Error())
		return result
	}

	for _, pod := range pods.Items {
		if strings.Contains(pod.Name, serviceName) {
			pod.ManagedFields = nil
			pod.Spec = v1.PodSpec{}
			result = append(result, pod)
		}
	}

	return result
}

func PodExists(namespace string, name string) ServicePodExistsResult {
	result := ServicePodExistsResult{}

	kubeProvider := NewKubeProvider()
	podClient := kubeProvider.ClientSet.CoreV1().Pods(namespace)
	pod, err := podClient.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil || pod == nil {
		result.PodExists = false
		return result
	}

	result.PodExists = true
	return result
}

func AllPods(namespaceName string) []v1.Pod {
	result := []v1.Pod{}

	provider := NewKubeProvider()
	podsList, err := provider.ClientSet.CoreV1().Pods(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllPods podMetricsList ERROR: %s", err.Error())
		return result
	}

	for _, pod := range podsList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, pod.ObjectMeta.Namespace) {
			result = append(result, pod)
		}
	}
	return result
}

func AllK8sPods(namespaceName string) K8sWorkloadResult {
	result := AllPods(namespaceName)
	return WorkloadResult(result, nil)
}

func AllPodNames() []string {
	result := []string{}
	allPods := AllPods("")
	for _, pod := range allPods {
		result = append(result, pod.ObjectMeta.Name)
	}
	return result
}

func AllPodNamesForLabel(namespace string, labelKey string, labelValue string) []string {
	result := []string{}
	allPods := AllPods(namespace)
	for _, pod := range allPods {
		if pod.Labels[labelKey] == labelValue {
			result = append(result, pod.ObjectMeta.Name)
		}
	}
	return result
}

func PodIdsFor(namespace string, serviceId *string) []string {
	result := []string{}

	var provider *KubeProviderMetrics
	var err error
	if !utils.CONFIG.Kubernetes.RunInCluster {
		provider, err = NewKubeProviderMetricsLocal()
	} else {
		provider, err = NewKubeProviderMetricsInCluster()
	}
	if err != nil {
		logger.Log.Errorf("PodIdsForServiceId ERROR: %s", err.Error())
		return result
	}

	podMetricsList, err := provider.ClientSet.MetricsV1beta1().PodMetricses(namespace).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("PodIdsForServiceId podMetricsList ERROR: %s", err.Error())
		return result
	}

	for _, podMetrics := range podMetricsList.Items {
		if serviceId != nil {
			if strings.Contains(podMetrics.ObjectMeta.Name, *serviceId) {
				result = append(result, podMetrics.ObjectMeta.Name)
			}
		} else {
			result = append(result, podMetrics.ObjectMeta.Name)
		}
	}
	// SORT TO HAVE A DETERMINISTIC ORDERING
	sort.Strings(result)

	return result
}

func UpdateK8sPod(data v1.Pod) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	podClient := kubeProvider.ClientSet.CoreV1().Pods(data.Namespace)
	_, err := podClient.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sPod(data v1.Pod) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	podClient := kubeProvider.ClientSet.CoreV1().Pods(data.Namespace)
	err := podClient.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sPod(namespace string, name string) K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "pod", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func NewK8sPod() K8sNewWorkload {
	return NewWorkload(
		RES_POD,
		utils.InitPodYaml(),
		"A Pod is the smallest and simplest unit in the Kubernetes object model that you create or deploy. It represents a single instance of a running process in a cluster and can contain one or more containers. In this example, a pod named 'my-pod' is created with a single container running the 'busybox' image. When the container starts, it runs the command sh -c 'echo Hello, Kubernetes! && sleep 3600', which prints 'Hello, Mogenius!' and then sleeps for 1 hour.")
}

func filterStatus(pod *v1.Pod) {
	pod.ManagedFields = nil
	pod.ObjectMeta = metav1.ObjectMeta{}
	pod.Spec = v1.PodSpec{}
}