package kubernetes

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	version2 "k8s.io/apimachinery/pkg/version"
	v1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"

	"github.com/jedib0t/go-pretty/table"
	"github.com/mogenius/punq/version"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/structs"

	"github.com/mogenius/punq/logger"

	"github.com/mogenius/punq/dtos"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
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
	RES_INGRESS_CLASS              string = "IngressClass"
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
	RES_INGRESS_CLASS,
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

type IngressType int

const (
	NGINX IngressType = iota
	TRAEFIK
	MULTIPLE
	NONE
	UNKNOWN
)

func (i IngressType) String() string {
	return [...]string{"NGINX", "TRAEFIK", "MULTIPLE", "NONE", "UNKNOWN"}[i]
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

func InitKubernetes(runsInCluster bool) {
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
		PodCpuUsageInMilliCores:      int(cpu),
		PodCpuLimitInMilliCores:      int(cpuLimit),
		PodMemoryUsageInBytes:        memory,
		PodMemoryLimitInBytes:        memoryLimit,
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

func ListNodeMetricss(contextId *string) []v1beta1.NodeMetrics {
	provider, err := NewKubeProviderMetrics(contextId)
	if provider == nil || err != nil {
		logger.Log.Errorf("ListNodeMetricss ERROR: %s", err.Error())
		return []v1beta1.NodeMetrics{}
	}

	nodeMetricsList, err := provider.ClientSet.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("ListNodeMetrics ERROR: %s", err.Error())
		return []v1beta1.NodeMetrics{}
	}
	return nodeMetricsList.Items
}

func podStats(pods map[string]v1.Pod, contextId *string) ([]structs.Stats, error) {
	provider, err := NewKubeProviderMetrics(contextId)
	if provider == nil || err != nil {
		logger.Log.Errorf(err.Error())
		return []structs.Stats{}, err
	}

	podMetricsList, err := provider.ClientSet.MetricsV1beta1().PodMetricses("").List(context.TODO(), metav1.ListOptions{})
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
		if pod.Status.StartTime != nil {
			entry.StartTime = pod.Status.StartTime.Format(time.RFC3339)
		}
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

func GetCurrentOperatorVersion() (string, error) {
	ownDeployment, err := GetK8sDeployment(utils.CONFIG.Kubernetes.OwnNamespace, version.Name, nil)
	if err != nil {
		logger.Log.Error("GetCurrentOperatorVersion:", err)
		return "", err
	}

	return strings.Split(ownDeployment.Spec.Template.Spec.Containers[0].Image, ":")[1], nil
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

func DetermineIngressControllerType(contextId *string) (IngressType, error) {
	ingressClasses := AllIngressClasses(contextId)

	if len(ingressClasses) > 1 {
		return MULTIPLE, fmt.Errorf("multiple ingress controllers found")
	}

	if len(ingressClasses) == 0 {
		return NONE, fmt.Errorf("no ingress controller found")
	}

	unknownController := ""
	for _, ingressClass := range ingressClasses {
		if ingressClass.Spec.Controller == "k8s.io/ingress-nginx" {
			return NGINX, nil
		} else if ingressClass.Spec.Controller == "traefik.io/ingress-controller" {
			return TRAEFIK, nil
		} else {
			unknownController = ingressClass.Spec.Controller
		}
	}

	return UNKNOWN, fmt.Errorf("unknown ingress controller: %s", unknownController)
}

func IsMetricsServerAvailable(contextId *string) (bool, string, error) {
	// kube-system would be the right namespace but if somebody installed it in another namespace we want to find it
	deployments := AllDeploymentsIncludeIgnored("", contextId)

	for _, deployment := range deployments {
		for key, label := range deployment.Labels {
			if key == "k8s-app" && label == "metrics-server" {
				if deployment.Status.UnavailableReplicas > 0 {
					return false, "", fmt.Errorf("metrics-server installed but not running")
				}
				return true, deployment.Spec.Template.Spec.Containers[0].Image, nil
			}
		}
	}

	return false, "", fmt.Errorf("no metrics-server found")
}

func ApiVersions(contextId *string) ([]string, error) {
	result := []string{}

	provider, err := NewKubeProvider(contextId)
	if provider == nil || err != nil {
		return result, err
	}

	groupResources, err := provider.ClientSet.DiscoveryClient.ServerPreferredResources()
	if err != nil {
		fmt.Printf("Error fetching API GroupResources: %v\n", err)
		return result, err
	}

	for _, groupList := range groupResources {
		result = append(result, groupList.GroupVersion)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})

	return result, nil
}

func GuessClusterProvider(contextId *string) (dtos.KubernetesProvider, error) {
	provider, err := NewKubeProvider(contextId)
	if provider == nil || err != nil {
		return dtos.SELF_HOSTED, err
	}

	nodes, err := provider.ClientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return dtos.SELF_HOSTED, err
	}

	return GuessCluserProviderFromNodeList(nodes)
}

func GuessCluserProviderFromNodeList(nodes *v1.NodeList) (dtos.KubernetesProvider, error) {

	for _, node := range nodes.Items {
		labels := node.GetLabels()

		if LabelsContain(labels, "eks.amazonaws.com/") {
			return dtos.EKS, nil
		} else if LabelsContain(labels, "docker-desktop") {
			return dtos.DOCKER_DESKTOP, nil
		} else if LabelsContain(labels, "kubernetes.azure.com/role") {
			return dtos.AKS, nil
		} else if LabelsContain(labels, "cloud.google.com/gke-nodepool") {
			return dtos.GKE, nil
		} else if LabelsContain(labels, "k3s.io/hostname") {
			return dtos.K3S, nil
		} else if LabelsContain(labels, "ibm-cloud.kubernetes.io/worker-version") {
			return dtos.IBM, nil
		} else if LabelsContain(labels, "doks.digitalocean.com/node-id") {
			return dtos.DOKS, nil
		} else if LabelsContain(labels, "oke.oraclecloud.com/node-pool") {
			return dtos.OKE, nil
		} else if LabelsContain(labels, "ack.aliyun.com") {
			return dtos.ACK, nil
		} else if LabelsContain(labels, "node-role.kubernetes.io/master") && LabelsContain(labels, "node.openshift.io/os_id") {
			return dtos.OPEN_SHIFT, nil
		} else if LabelsContain(labels, "vmware-system-vmware.io/role") {
			return dtos.VMWARE, nil
		} else if LabelsContain(labels, "io.rancher.os/hostname") {
			return dtos.RKE, nil
		} else if LabelsContain(labels, "linode-lke/") {
			return dtos.LINODE, nil
		} else if LabelsContain(labels, "scaleway-kapsule/") {
			return dtos.SCALEWAY, nil
		} else if LabelsContain(labels, "microk8s.io/cluster") {
			return dtos.MICROK8S, nil
		} else if strings.ToLower(node.Name) == "minikube" {
			return dtos.MINIKUBE, nil
		} else if LabelsContain(labels, "io.k8s.sigs.kind/role") {
			return dtos.KIND, nil
		} else if LabelsContain(labels, "civo/") {
			return dtos.CIVO, nil
		} else if LabelsContain(labels, "giantswarm.io/") {
			return dtos.GIANTSWARM, nil
		} else if LabelsContain(labels, "ovhcloud/") {
			return dtos.OVHCLOUD, nil
		} else if LabelsContain(labels, "gardener.cloud/role") {
			return dtos.GARDENER, nil
		} else if LabelsContain(labels, "cce.huawei.com") {
			return dtos.HUAWEI, nil
		} else if LabelsContain(labels, "nirmata.io") {
			return dtos.NIRMATA, nil
		} else if LabelsContain(labels, "platform9.com/role") {
			return dtos.PF9, nil
		} else if LabelsContain(labels, "nks.netapp.io") {
			return dtos.NKS, nil
		} else if LabelsContain(labels, "appscode.com") {
			return dtos.APPSCODE, nil
		} else if LabelsContain(labels, "loft.sh") {
			return dtos.LOFT, nil
		} else if LabelsContain(labels, "spectrocloud.com") {
			return dtos.SPECTROCLOUD, nil
		} else if LabelsContain(labels, "diamanti.com") {
			return dtos.DIAMANTI, nil
		} else if strings.HasPrefix(strings.ToLower(node.Name), "k3d-") {
			return dtos.K3D, nil
		} else if LabelsContain(labels, "cloud.google.com/gke-on-prem") {
			return dtos.GKE_ON_PREM, nil
		} else if LabelsContain(labels, "rke.cattle.io") {
			return dtos.RKE, nil
		} else {
			fmt.Println("This cluster's provider is unknown or it might be self-managed.")
			return dtos.UNKNOWN, nil
		}
	}
	return dtos.UNKNOWN, nil
}

func LabelsContain(labels map[string]string, str string) bool {
	// Keys EQUAL
	if _, ok := labels[strings.ToLower(str)]; ok {
		return true
	}

	// Values
	for key, label := range labels {
		if strings.EqualFold(label, str) {
			return true
		}
		// KEY CONTAINS
		if strings.Contains(key, str) {
			return true
		}
	}
	return false
}

func AllResourcesFrom(namespace string, contextId *string) ([]interface{}, error) {
	ignoredResources := []string{
		"events.k8s.io/v1",
		"events.k8s.io/v1beta1",
		"metrics.k8s.io/v1beta1",
		"discovery.k8s.io/v1",
	}

	result := []interface{}{}

	provider, err := NewKubeProvider(contextId)
	if provider == nil || err != nil {
		return result, err
	}

	// Get a list of all resource types in the cluster
	resourceList, err := provider.ClientSet.Discovery().ServerPreferredResources()
	if err != nil {
		return result, err
	}

	// Iterate over each resource type and backup all resources in the namespace
	for _, resource := range resourceList {
		if utils.Contains(ignoredResources, resource.GroupVersion) {
			continue
		}
		gv, _ := schema.ParseGroupVersion(resource.GroupVersion)
		if len(resource.APIResources) <= 0 {
			continue
		}

		for _, aApiResource := range resource.APIResources {
			if !aApiResource.Namespaced {
				continue
			}

			resourceId := schema.GroupVersionResource{
				Group:    gv.Group,
				Version:  gv.Version,
				Resource: aApiResource.Name,
			}
			// Get the REST client for this resource type
			restClient := dynamic.New(provider.ClientSet.RESTClient()).Resource(resourceId).Namespace(namespace)

			// Get a list of all resources of this type in the namespace
			list, err := restClient.List(context.Background(), metav1.ListOptions{})
			if err != nil {
				logger.Log.Error("%s: %s", resourceId.Resource, err.Error())
				continue
			}

			// Iterate over each resource and write it to a file
			for _, obj := range list.Items {
				logger.Log.Noticef("(SUCCESS) %s: %s/%s", resourceId.Resource, obj.GetNamespace(), obj.GetName())
				obj.SetManagedFields(nil)
				delete(obj.Object, "status")
				obj.SetUID("")
				obj.SetResourceVersion("")
				obj.SetCreationTimestamp(metav1.Time{})

				result = append(result, obj.Object)
			}
		}
	}
	return result, nil
}
