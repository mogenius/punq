package kubernetes

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/jedib0t/go-pretty/table"
	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServicePodExistsResult struct {
	PodExists bool `json:"podExists"`
}

func PodStatus(namespace string, name string, statusOnly bool, contextId *string) *v1.Pod {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return nil
	}
	getOptions := metav1.GetOptions{}

	podClient := provider.ClientSet.CoreV1().Pods(namespace)

	pod, err := podClient.Get(context.TODO(), name, getOptions)
	pod.Kind = "Pod"
	pod.APIVersion = "v1"
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

func ServicePodStatus(namespace string, serviceName string, contextId *string) []v1.Pod {
	result := []v1.Pod{}
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return result
	}

	podClient := provider.ClientSet.CoreV1().Pods(namespace)

	pods, err := podClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("ServicePodStatus Error: %s", err.Error())
		return result
	}

	for _, pod := range pods.Items {
		if strings.Contains(pod.Name, serviceName) {
			pod.ManagedFields = nil
			pod.Spec = v1.PodSpec{}
			pod.Kind = "Pod"
			pod.APIVersion = "v1"
			result = append(result, pod)
		}
	}

	return result
}

// labelname should look like app=my-app-name (like you defined your label)
func GetFirstPodForLabelName(namespace string, labelName string, contextId *string) *v1.Pod {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return nil
	}

	pods, err := provider.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelName})

	for _, pod := range pods.Items {
		pod.Kind = "Pod"
		pod.APIVersion = "v1"
		return &pod
	}

	if err != nil {
		logger.Log.Errorf("GetFirstPodForLabelName ERR:", err)
		return nil
	}

	logger.Log.Errorf("No pod labeled '%s/%s' not found", namespace, labelName)
	return nil
}

func GetPod(namespace string, podName string, contextId *string) *v1.Pod {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return nil
	}

	client := provider.ClientSet.CoreV1().Pods(namespace)
	pod, err := client.Get(context.TODO(), podName, metav1.GetOptions{})
	pod.Kind = "Pod"
	pod.APIVersion = "v1"
	if err != nil {
		logger.Log.Errorf("GetPod Error: %s", err.Error())
		return nil
	}
	return pod
}

func GetPodBy(namespace string, podName string, contextId *string) (*v1.Pod, error) {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return nil, err
	}
	client := provider.ClientSet.CoreV1().Pods(namespace)
	pod, err := client.Get(context.TODO(), podName, metav1.GetOptions{})
	pod.Kind = "Pod"
	pod.APIVersion = "v1"

	return pod, err
}

func PodExists(namespace string, name string, contextId *string) ServicePodExistsResult {
	result := ServicePodExistsResult{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return result
	}
	podClient := provider.ClientSet.CoreV1().Pods(namespace)
	pod, err := podClient.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil || pod == nil {
		result.PodExists = false
		return result
	}

	result.PodExists = true
	return result
}

func AllPodsOnNode(nodeName string, contextId *string) []v1.Pod {
	result := []v1.Pod{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return result
	}

	podsList, err := provider.ClientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeName,
	})
	if err != nil {
		logger.Log.Errorf("AllPodsOnNode ERROR: %s", err.Error())
		return result
	}
	for _, pod := range podsList.Items {
		pod.Kind = "Pod"
		pod.APIVersion = "v1"
		result = append(result, pod)
	}

	return result
}

func AllPods(namespaceName string, contextId *string) []v1.Pod {
	result := []v1.Pod{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return result
	}
	podsList, err := provider.ClientSet.CoreV1().Pods(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllPods podMetricsList ERROR: %s", err.Error())
		return result
	}

	for _, pod := range podsList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, pod.ObjectMeta.Namespace) {
			pod.Kind = "Pod"
			pod.APIVersion = "v1"
			result = append(result, pod)
		}
	}
	return result
}

func AllK8sPods(namespaceName string, contextId *string) utils.K8sWorkloadResult {
	result := []v1.Pod{}
	pods := AllPods(namespaceName, contextId)
	for _, pod := range pods {
		pod.Kind = "Pod"
		pod.APIVersion = "v1"
		result = append(result, pod)
	}

	return WorkloadResult(result, nil)
}

func AllPodNames(contextId *string) []string {
	result := []string{}
	allPods := AllPods("", contextId)
	for _, pod := range allPods {
		result = append(result, pod.ObjectMeta.Name)
	}
	return result
}

func AllPodNamesForLabel(namespace string, labelKey string, labelValue string, contextId *string) []string {
	result := []string{}
	allPods := AllPods(namespace, contextId)
	for _, pod := range allPods {
		if pod.Labels[labelKey] == labelValue {
			result = append(result, pod.ObjectMeta.Name)
		}
	}
	return result
}

func PodIdsFor(namespace string, serviceId *string, contextId *string) []string {
	result := []string{}

	provider, err := NewKubeProviderMetrics(contextId)
	if provider == nil || err != nil {
		logger.Log.Errorf(err.Error())
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

func UpdateK8sPod(data v1.Pod, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	podClient := provider.ClientSet.CoreV1().Pods(data.Namespace)
	res, err := podClient.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func DeleteK8sPod(data v1.Pod, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	podClient := provider.ClientSet.CoreV1().Pods(data.Namespace)
	err = podClient.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sPodBy(namespace string, name string, contextId *string) error {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return err
	}
	podClient := provider.ClientSet.CoreV1().Pods(namespace)
	return podClient.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sPod(namespace string, name string, contextId *string) utils.K8sWorkloadResult {
	cmd := utils.RunOnLocalShell(fmt.Sprintf("kubectl describe pod %s -n %s%s", name, namespace, ContextFlag(contextId)))

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sPod(data v1.Pod, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.CoreV1().Pods(data.Namespace)
	res, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
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

func ListPodsTerminal(namespace string, contextId *string) {
	pods := AllPods(namespace, contextId)
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Namespace", "Name", "Ready", "Status", "Restarts", "Age"})
	for index, pod := range pods {
		t.AppendRow(
			table.Row{index + 1, pod.Namespace, pod.Name, pod.Status.ContainerStatuses[0].Ready, pod.Status.Phase, pod.Status.ContainerStatuses[0].RestartCount, utils.JsonStringToHumanDuration(pod.Status.StartTime.Format(time.RFC3339))},
		)
	}
	t.Render()
}
