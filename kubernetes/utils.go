package kubernetes

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"time"

	version2 "k8s.io/apimachinery/pkg/version"

	"github.com/jedib0t/go-pretty/table"
	"github.com/mogenius/punq/version"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/structs"

	"github.com/mogenius/punq/logger"

	"github.com/mogenius/punq/dtos"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var RunsInCluster bool = false

const (
	RES_NAMESPACE                  string = "Namespace"
	RES_POD                        string = "Pod"
	RES_DEPLOYMENT                 string = "Deployment"
	RES_SERVICE                    string = "Service"
	RES_INGRESS                    string = "Ingress"
	RES_CONFIG_MAP                 string = "ConfigMap"
	RES_SECRET                     string = "Secret"
	RES_NODE                       string = "Node"
	RES_DAEMON_SET                 string = "DaemonSet"
	RES_STATEFUL_SET               string = "StatefulSet"
	RES_JOB                        string = "Job"
	RES_CRON_JOB                   string = "CronJob"
	RES_REPLICA_SET                string = "ReplicaSet"
	RES_PERSISTENT_VOLUME          string = "PersistentVolume"
	RES_PERSISTENT_VOLUME_CLAIM    string = "PersistentVolumeClaim"
	RES_HORIZONTAL_POD_AUTOSCALER  string = "HorizontalPodAutoscaler"
	RES_EVENT                      string = "Event"
	RES_CERTIFICATE                string = "Certificate"
	RES_CERTIFICATE_REQUEST        string = "CertificateRequest"
	RES_ORDER                      string = "Order"
	RES_ISSUER                     string = "Issuer"
	RES_CLUSTER_ISSUER             string = "ClusterIssuer"
	RES_SERVICE_ACCOUNT            string = "ServiceAccount"
	RES_ROLE                       string = "Role"
	RES_ROLE_BINDING               string = "RoleBinding"
	RES_CLUSTER_ROLE               string = "ClusterRole"
	RES_CLUSTER_ROLE_BINDING       string = "ClusterRoleBinding"
	RES_VOLUME_ATTACHMENT          string = "VolumeAttachment"
	RES_NETWORK_POLICY             string = "NetworkPolicy"
	RES_STORAGE_CLASS              string = "StorageClass"
	RES_CUSTOM_RESOURCE_DEFINITION string = "CustomResourceDefinition"
	RES_ENDPOINT                   string = "Endpoint"
	RES_LEASE                      string = "Lease"
	RES_PRIORITY_CLASS             string = "PriorityClass"
	RES_VOLUME_SNAPSHOT            string = "VolumeSnapshot"
	RES_RESOURCE_QUOTA             string = "ResourceQuota"
)

var ALL_RESOURCES []string = []string{
	RES_NAMESPACE,
	RES_POD,
	RES_DEPLOYMENT,
	RES_SERVICE,
	RES_INGRESS,
	RES_CONFIG_MAP,
	RES_SECRET,
	RES_NODE,
	RES_DAEMON_SET,
	RES_STATEFUL_SET,
	RES_JOB,
	RES_CRON_JOB,
	RES_REPLICA_SET,
	RES_PERSISTENT_VOLUME,
	RES_PERSISTENT_VOLUME_CLAIM,
	RES_HORIZONTAL_POD_AUTOSCALER,
	RES_EVENT,
	RES_CERTIFICATE,
	RES_CERTIFICATE_REQUEST,
	RES_ORDER,
	RES_ISSUER,
	RES_CLUSTER_ISSUER,
	RES_SERVICE_ACCOUNT,
	RES_ROLE,
	RES_ROLE_BINDING,
	RES_CLUSTER_ROLE,
	RES_CLUSTER_ROLE_BINDING,
	RES_VOLUME_ATTACHMENT,
	RES_NETWORK_POLICY,
	RES_STORAGE_CLASS,
	RES_CUSTOM_RESOURCE_DEFINITION,
	RES_ENDPOINT,
	RES_LEASE,
	RES_PRIORITY_CLASS,
	RES_VOLUME_SNAPSHOT,
	RES_RESOURCE_QUOTA,
}

var ALL_RESOURCES_USER []string = []string{
	RES_NAMESPACE,
	RES_POD,
	RES_DEPLOYMENT,
	RES_SERVICE,
	RES_INGRESS,
	RES_CONFIG_MAP,
	RES_NODE,
	RES_DAEMON_SET,
	RES_STATEFUL_SET,
	RES_JOB,
	RES_CRON_JOB,
	RES_REPLICA_SET,
	RES_PERSISTENT_VOLUME_CLAIM,
	RES_EVENT,
	RES_NETWORK_POLICY,
	RES_ENDPOINT,
}

var ALL_RESOURCES_READER []string = []string{
	RES_NAMESPACE,
	RES_POD,
	RES_DEPLOYMENT,
	RES_SERVICE,
	RES_INGRESS,
	RES_CONFIG_MAP,
	RES_NODE,
	RES_DAEMON_SET,
	RES_STATEFUL_SET,
	RES_JOB,
	RES_CRON_JOB,
	RES_REPLICA_SET,
	RES_PERSISTENT_VOLUME_CLAIM,
	RES_EVENT,
	RES_NETWORK_POLICY,
	RES_ENDPOINT,
}

var (
	SERVICEACCOUNTNAME     = fmt.Sprintf("%s-service-account-app", version.Name)
	CLUSTERROLENAME        = fmt.Sprintf("%s--cluster-role-app", version.Name)
	CLUSTERROLEBINDINGNAME = fmt.Sprintf("%s--cluster-role-binding-app", version.Name)
	RBACRESOURCES          = []string{"*"}
	SERVICENAME            = fmt.Sprintf("%s-service", version.Name)
	INGRESSNAME            = fmt.Sprintf("%s-ingress", version.Name)
)

type K8sNewWorkload struct {
	Name        string `json:"name"`
	YamlString  string `json:"yamlString"`
	Description string `json:"description"`
}

type MogeniusNfsInstallationStatus struct {
	Error       string `json:"error,omitempty"`
	IsInstalled bool   `json:"isInstalled"`
}

func Init(runsInCluster bool) {
	RunsInCluster = runsInCluster
}

func ListWorkloadsOnTerminal(access dtos.AccessLevel) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Name"})

	resources := WorkloadsForAccesslevel(access)
	for index, resource := range resources {
		t.AppendRow(
			table.Row{index + 1, resource},
		)
	}
	t.Render()
}

func WorkloadsForAccesslevel(access dtos.AccessLevel) []string {
	resources := []string{}
	switch access {
	case dtos.READER:
		resources = ALL_RESOURCES_READER
	case dtos.USER:
		resources = ALL_RESOURCES_USER
	case dtos.ADMIN:
		resources = ALL_RESOURCES
	}
	return resources
}

func WorkloadResult(result interface{}, err interface{}) utils.K8sWorkloadResult {
	if fmt.Sprint(reflect.TypeOf(err)) == "*errors.errorString" {
		err = err.(error).Error()
	}
	return utils.K8sWorkloadResult{
		Result: result,
		Error:  err,
	}
}

func WorkloadResultError(error string) utils.K8sWorkloadResult {
	return utils.K8sWorkloadResult{
		Result: nil,
		Error:  error,
	}
}

func NewWorkload(name string, yaml string, description string) K8sNewWorkload {
	return K8sNewWorkload{
		Name:        name,
		YamlString:  yaml,
		Description: description,
	}
}

func CurrentContextName() string {
	var kubeconfig string = utils.GetDefaultKubeConfig()

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{
			CurrentContext: "",
		}).RawConfig()

	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	return config.CurrentContext
}

func Hostname(contextId *string) string {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		logger.Log.Error("Hostname:", err)
		return ""
	}
	return provider.ClientConfig.Host
}

func ClusterStatus(contextId *string) dtos.ClusterStatusDto {
	var currentPods = make(map[string]v1.Pod)
	pods := listAllPods(contextId)
	for _, pod := range pods {
		currentPods[pod.Name] = pod
	}

	result, err := podStats(currentPods, contextId)
	if err != nil {
		logger.Log.Error("podStats:", err)
	}

	var cpu int64 = 0
	var cpuLimit int64 = 0
	var memory int64 = 0
	var memoryLimit int64 = 0
	var ephemeralStorageLimit int64 = 0
	for _, pod := range result {
		cpu += pod.Cpu
		cpuLimit += pod.CpuLimit
		memory += pod.Memory
		memoryLimit += pod.MemoryLimit
		ephemeralStorageLimit += pod.EphemeralStorageLimit
	}

	kubernetesVersion := ""
	platform := ""

	info := KubernetesVersion(contextId)
	if info != nil {
		kubernetesVersion = info.String()
		platform = info.Platform
	}

	return dtos.ClusterStatusDto{
		ClusterName:                  utils.CONFIG.Kubernetes.ClusterName,
		Pods:                         len(result),
		CpuInMilliCores:              int(cpu),
		CpuLimitInMilliCores:         int(cpuLimit),
		MemoryInBytes:                memory,
		MemoryLimitInBytes:           memoryLimit,
		EphemeralStorageLimitInBytes: ephemeralStorageLimit,
		KubernetesVersion:            kubernetesVersion,
		Platform:                     platform,
	}
}

func KubernetesVersion(contextId *string) *version2.Info {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return nil
	}
	info, err := provider.ClientSet.Discovery().ServerVersion()
	if err != nil {
		logger.Log.Error("Error KubernetesVersion:", err)
		return nil
	}
	return info
}

func ClusterInfo(contextId *string) dtos.ClusterInfoDto {
	result := dtos.ClusterInfoDto{
		ClusterStatus: ClusterStatus(contextId),
		NodeStats:     GetNodeStats(contextId),
	}
	return result
}

func listAllPods(contextId *string) []v1.Pod {
	var result []v1.Pod

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return result
	}
	pods, err := provider.ClientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system,metadata.namespace!=default"})

	if err != nil {
		logger.Log.Error("Error listAllPods:", err)
		return result
	}
	return pods.Items
}

func ListNodes(contextId *string) []v1.Node {
	provider, err := NewKubeProvider(contextId)
	if provider == nil || err != nil {
		logger.Log.Errorf("ListNodes ERROR: %s", err.Error())
		return []v1.Node{}
	}

	nodeMetricsList, err := provider.ClientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("ListNodeMetrics ERROR: %s", err.Error())
		return []v1.Node{}
	}
	return nodeMetricsList.Items
}

func podStats(pods map[string]v1.Pod, contextId *string) ([]structs.Stats, error) {
	provider, err := NewKubeProviderMetrics(contextId)
	if provider == nil || err != nil {
		logger.Log.Errorf(err.Error())
		return []structs.Stats{}, err
	}

	podMetricsList, err := provider.ClientSet.MetricsV1beta1().PodMetricses("").List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system,metadata.namespace!=default"})
	if err != nil {
		return nil, err
	}

	var result []structs.Stats
	// I HATE THIS BUT I DONT SEE ANY OTHER SOLUTION! SPEND HOURS (to find something better) ON THIS UGGLY SHIT!!!!

	for _, podMetrics := range podMetricsList.Items {
		var pod = pods[podMetrics.Name]

		var entry = structs.Stats{}
		entry.Cluster = utils.CONFIG.Kubernetes.ClusterName
		entry.Namespace = podMetrics.Namespace
		entry.PodName = podMetrics.Name
		entry.StartTime = pod.Status.StartTime.Format(time.RFC3339)
		for _, container := range pod.Spec.Containers {
			entry.CpuLimit += container.Resources.Limits.Cpu().MilliValue()
			entry.MemoryLimit += container.Resources.Limits.Memory().Value()
			entry.EphemeralStorageLimit += container.Resources.Limits.StorageEphemeral().Value()
		}
		for _, containerMetric := range podMetrics.Containers {
			entry.Cpu += containerMetric.Usage.Cpu().MilliValue()
			entry.Memory += containerMetric.Usage.Memory().Value()
		}

		result = append(result, entry)
	}

	return result, nil
}

// TAKEN FROM Kubernetes apimachineryv0.25.1
func HumanDuration(d time.Duration) string {
	// Allow deviation no more than 2 seconds(excluded) to tolerate machine time
	// inconsistence, it can be considered as almost now.
	if seconds := int(d.Seconds()); seconds < -1 {
		return "<invalid>"
	} else if seconds < 0 {
		return "0s"
	} else if seconds < 60*2 {
		return fmt.Sprintf("%ds", seconds)
	}
	minutes := int(d / time.Minute)
	if minutes < 10 {
		s := int(d/time.Second) % 60
		if s == 0 {
			return fmt.Sprintf("%dm", minutes)
		}
		return fmt.Sprintf("%dm%ds", minutes, s)
	} else if minutes < 60*3 {
		return fmt.Sprintf("%dm", minutes)
	}
	hours := int(d / time.Hour)
	if hours < 8 {
		m := int(d/time.Minute) % 60
		if m == 0 {
			return fmt.Sprintf("%dh", hours)
		}
		return fmt.Sprintf("%dh%dm", hours, m)
	} else if hours < 48 {
		return fmt.Sprintf("%dh", hours)
	} else if hours < 24*8 {
		h := hours % 24
		if h == 0 {
			return fmt.Sprintf("%dd", hours/24)
		}
		return fmt.Sprintf("%dd%dh", hours/24, h)
	} else if hours < 24*365*2 {
		return fmt.Sprintf("%dd", hours/24)
	} else if hours < 24*365*8 {
		dy := int(hours/24) % 365
		if dy == 0 {
			return fmt.Sprintf("%dy", hours/24/365)
		}
		return fmt.Sprintf("%dy%dd", hours/24/365, dy)
	}
	return fmt.Sprintf("%dy", int(hours/24/365))
}

func MoCreateOptions() metav1.CreateOptions {
	return metav1.CreateOptions{
		FieldManager: version.Name,
	}
}

func MoUpdateOptions() metav1.UpdateOptions {
	return metav1.UpdateOptions{
		FieldManager: version.Name,
	}
}

func ListCreateTemplates() []K8sNewWorkload {
	result := []K8sNewWorkload{}

	result = append(result,
		NewK8sCertificate(),
		NewK8sCertificateSigningRequest(),
		NewK8sClusterIssuer(),
		NewK8sClusterRole(),
		NewK8sClusterRoleBinding(),
		NewK8sConfigmap(),
		NewK8sCronJob(),
		NewK8sDaemonSet(),
		NewK8sDeployment(),
		NewK8sHpa(),
		NewK8sIngress(),
		NewK8sIssuer(),
		NewK8sJob(),
		NewK8sNamespace(),
		NewK8sNetPol(),
		NewK8sOrder(),
		NewK8sPersistentVolumeClaim(),
		NewK8sPod(),
		NewK8sReplicaSet(),
		NewK8sRole(),
		NewK8sRoleBinding(),
		NewK8sSecret(),
		NewK8sService(),
		NewK8sServiceAccount(),
		NewK8sStatefulset(),
		NewK8sStorageClass(),
		NewK8sVolume(),
		NewK8sVolumeAttachment())

	return result
}

func ListTemplatesTerminal() {
	for _, template := range ListCreateTemplates() {
		structs.PrettyPrint(template)
	}
}
