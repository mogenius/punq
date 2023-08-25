package kubernetes

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/jedib0t/go-pretty/table"
	"github.com/mogenius/punq/version"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/structs"

	"github.com/mogenius/punq/logger"

	"github.com/mogenius/punq/dtos"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	RES_NAMESPACE                   string = "namespace"
	RES_POD                         string = "pod"
	RES_DEPLOYMENT                  string = "deployment"
	RES_SERVICE                     string = "service"
	RES_INGRESS                     string = "ingress"
	RES_CONFIGMAP                   string = "configmap"
	RES_SECRET                      string = "secret"
	RES_NODE                        string = "node"
	RES_DAEMON_SET                  string = "daemon_set"
	RES_STATEFUL_SET                string = "stateful_set"
	RES_JOB                         string = "job"
	RES_CRON_JOB                    string = "cron_job"
	RES_REPLICA_SET                 string = "replica_set"
	RES_PERSISTENT_VOLUME           string = "persistent_volume"
	RES_PERSISTENT_VOLUME_CLAIM     string = "persistent_volume_claim"
	RES_HORIZONTAL_POD_AUTOSCALER   string = "horizontal_pod_autoscaler"
	RES_EVENT                       string = "event"
	RES_CERTIFICATE                 string = "certificate"
	RES_CERTIFICATE_REQUEST         string = "certificaterequest"
	RES_ORDER                       string = "orders"
	RES_ISSUER                      string = "issuer"
	RES_CLUSTER_ISSUER              string = "clusterissuer"
	RES_SERVICE_ACCOUNT             string = "service_account"
	RES_ROLE                        string = "role"
	RES_ROLE_BINDING                string = "role_binding"
	RES_CLUSTER_ROLE                string = "cluster_role"
	RES_CLUSTER_ROLE_BINDING        string = "cluster_role_binding"
	RES_VOLUME_ATTACHMENT           string = "volume_attachment"
	RES_NETWORK_POLICY              string = "network_policy"
	RES_STORAGECLASS                string = "storageclass"
	RES_CUSTOM_RESOURCE_DEFINITIONS string = "crds"
	RES_ENDPOINTS                   string = "endpoints"
	RES_LEASES                      string = "leases"
	RES_PRIORITYCLASSES             string = "priorityclasses"
	RES_VOLUMESNAPSHOTS             string = "volumesnapshots"
	RES_RESOURCEQUOTAS              string = "resourcequotas"
)

var ALL_RESOURCES []string = []string{
	RES_NAMESPACE,
	RES_POD,
	RES_DEPLOYMENT,
	RES_SERVICE,
	RES_INGRESS,
	RES_CONFIGMAP,
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
	RES_STORAGECLASS,
	RES_CUSTOM_RESOURCE_DEFINITIONS,
	RES_ENDPOINTS,
	RES_LEASES,
	RES_PRIORITYCLASSES,
	RES_VOLUMESNAPSHOTS,
	RES_RESOURCEQUOTAS,
}

var (
	DEPLOYMENTIMAGE = fmt.Sprintf("ghcr.io/mogenius/%s:v%s", version.Name, version.Ver)

	SERVICEACCOUNTNAME     = fmt.Sprintf("%s-service-account-app", version.Name)
	CLUSTERROLENAME        = fmt.Sprintf("%s--cluster-role-app", version.Name)
	CLUSTERROLEBINDINGNAME = fmt.Sprintf("%s--cluster-role-binding-app", version.Name)
	RBACRESOURCES          = []string{"pods", "services", "endpoints", "secrets"}
	SERVICENAME            = fmt.Sprintf("%s-service", version.Name)
	INGRESSNAME            = fmt.Sprintf("%s-ingress", version.Name)
)

type K8sWorkloadResult struct {
	Result interface{} `json:"result,omitempty"`
	Error  interface{} `json:"error,omitempty"`
}

type K8sNewWorkload struct {
	Name        string `json:"name"`
	YamlString  string `json:"yamlString"`
	Description string `json:"description"`
}

type MogeniusNfsInstallationStatus struct {
	Error       string `json:"error,omitempty"`
	IsInstalled bool   `json:"isInstalled"`
}

func ListWorkloads() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Name"})
	for index, resource := range ALL_RESOURCES {
		t.AppendRow(
			table.Row{index + 1, resource},
		)
	}
	t.Render()
}

func WorkloadResult(result interface{}, err interface{}) K8sWorkloadResult {
	fmt.Println(reflect.TypeOf(err))
	if fmt.Sprint(reflect.TypeOf(err)) == "*errors.errorString" {
		err = err.(error).Error()
	}
	return K8sWorkloadResult{
		Result: result,
		Error:  err,
	}
}

func WorkloadResultError(error string) K8sWorkloadResult {
	return K8sWorkloadResult{
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
	var kubeconfig string = ""
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

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

func Hostname() string {
	provider := NewKubeProvider()
	return provider.ClientConfig.Host
}

func ClusterStatus() dtos.ClusterStatusDto {
	var currentPods = make(map[string]v1.Pod)
	pods := listAllPods()
	for _, pod := range pods {
		currentPods[pod.Name] = pod
	}

	result, err := podStats(currentPods)
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

	return dtos.ClusterStatusDto{
		ClusterName:           utils.CONFIG.Kubernetes.ClusterName,
		Pods:                  len(result),
		CpuInMilliCores:       int(cpu),
		CpuLimitInMilliCores:  int(cpuLimit),
		Memory:                utils.BytesToHumanReadable(memory),
		MemoryLimit:           utils.BytesToHumanReadable(memoryLimit),
		EphemeralStorageLimit: utils.BytesToHumanReadable(ephemeralStorageLimit),
	}
}

func listAllPods() []v1.Pod {
	var result []v1.Pod

	kubeProvider := NewKubeProvider()
	pods, err := kubeProvider.ClientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system,metadata.namespace!=default"})

	if err != nil {
		logger.Log.Error("Error listAllPods:", err)
		return result
	}
	return pods.Items
}

func ListNodes() []v1.Node {
	var provider *KubeProvider = NewKubeProvider()
	if provider == nil {
		logger.Log.Errorf("Failed to load kubeprovider.")
		return []v1.Node{}
	}

	nodeMetricsList, err := provider.ClientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("ListNodeMetrics ERROR: %s", err.Error())
		return []v1.Node{}
	}
	return nodeMetricsList.Items
}

func podStats(pods map[string]v1.Pod) ([]structs.Stats, error) {
	var provider *KubeProviderMetrics = NewKubeProviderMetrics()
	if provider == nil {
		err := fmt.Errorf("Failed to load kubeprovider.")
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

func getKubeConfig() string {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		kubeconfig = ""
	}
	return kubeconfig
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
