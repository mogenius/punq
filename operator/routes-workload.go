package operator

import (
	"net/http"

	v1Cert "github.com/cert-manager/cert-manager/pkg/apis/acme/v1"
	cmapi "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	v6Snap "github.com/kubernetes-csi/external-snapshotter/client/v6/apis/volumesnapshot/v1"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/kubernetes"
	v1Apps "k8s.io/api/apps/v1"
	v2Scale "k8s.io/api/autoscaling/v2"
	v1Job "k8s.io/api/batch/v1"
	v1Coordination "k8s.io/api/coordination/v1"
	v1 "k8s.io/api/core/v1"
	v1Networking "k8s.io/api/networking/v1"
	v1Rbac "k8s.io/api/rbac/v1"
	v1Scheduling "k8s.io/api/scheduling/v1"
	v1Storage "k8s.io/api/storage/v1"
	apiExt "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

func InitWorkloadRoutes(router *gin.Engine) {
	router.GET("/workload/templates", Auth(dtos.USER), allWorkloadTemplates)
	router.GET("/workload/available-resources", Auth(dtos.READER), allKubernetesResources)

	router.GET("/workload/namespace/all", Auth(dtos.USER), allNamespaces)   // QUERY: -
	router.DELETE("/workload/namespace", Auth(dtos.ADMIN), deleteNamespace) // BODY: json-object
	router.POST("/workload/namespace", Auth(dtos.USER), createNamespace)    // BODY: yaml-object

	router.GET("/workload/pod", Auth(dtos.USER), allPods)              // QUERY: namespace
	router.GET("/workload/pod/describe", Auth(dtos.USER), describePod) // QUERY: namespace, name
	router.DELETE("/workload/pod", Auth(dtos.USER), deletePod)         // BODY: json-object
	router.PATCH("/workload/pod", Auth(dtos.USER), patchPod)           // BODY: json-object
	router.POST("/workload/pod", Auth(dtos.USER), createPod)           // BODY: yaml-object

	router.GET("/workload/deployment", Auth(dtos.USER), allDeployments)              // QUERY: namespace
	router.GET("/workload/deployment/describe", Auth(dtos.USER), describeDeployment) // QUERY: namespace, name
	router.DELETE("/workload/deployment", Auth(dtos.USER), deleteDeployment)         // BODY: json-object
	router.PATCH("/workload/deployment", Auth(dtos.USER), patchDeployment)           // BODY: json-object
	router.POST("/workload/deployment", Auth(dtos.USER), createDeployment)           // BODY: yaml-object

	router.GET("/workload/service", Auth(dtos.USER), allServices)              // QUERY: namespace
	router.GET("/workload/service/describe", Auth(dtos.USER), describeService) // QUERY: namespace, name
	router.DELETE("/workload/service", Auth(dtos.USER), deleteService)         // BODY: json-object
	router.PATCH("/workload/service", Auth(dtos.USER), patchService)           // BODY: json-object
	router.POST("/workload/service", Auth(dtos.USER), createService)           // BODY: yaml-object

	router.GET("/workload/ingress", Auth(dtos.USER), allIngresses)             // QUERY: namespace
	router.GET("/workload/ingress/describe", Auth(dtos.USER), describeIngress) // QUERY: namespace, name
	router.DELETE("/workload/ingress", Auth(dtos.USER), deleteIngress)         // BODY: json-object
	router.PATCH("/workload/ingress", Auth(dtos.USER), patchIngress)           // BODY: json-object
	router.POST("/workload/ingress", Auth(dtos.USER), createIngress)           // BODY: yaml-object

	router.GET("/workload/configmap", Auth(dtos.USER), allConfigmaps)              // QUERY: namespace
	router.GET("/workload/configmap/describe", Auth(dtos.USER), describeConfigmap) // QUERY: namespace, name
	router.DELETE("/workload/configmap", Auth(dtos.USER), deleteConfigmap)         // BODY: json-object
	router.PATCH("/workload/configmap", Auth(dtos.USER), patchConfigmap)           // BODY: json-object
	router.POST("/workload/configmap", Auth(dtos.USER), createConfigmap)           // BODY: yaml-object

	router.GET("/workload/secret", Auth(dtos.ADMIN), allSecrets)              // QUERY: namespace
	router.GET("/workload/secret/describe", Auth(dtos.ADMIN), describeSecret) // QUERY: namespace, name
	router.DELETE("/workload/secret", Auth(dtos.ADMIN), deleteSecret)         // BODY: json-object
	router.PATCH("/workload/secret", Auth(dtos.ADMIN), patchSecret)           // BODY: json-object
	router.POST("/workload/secret", Auth(dtos.ADMIN), createSecret)           // BODY: yaml-object

	router.GET("/workload/node", Auth(dtos.USER), allNodes)              // QUERY: -
	router.GET("/workload/node/describe", Auth(dtos.USER), describeNode) // QUERY:  name

	router.GET("/workload/daemon_set", Auth(dtos.USER), allDaemonSets)              // QUERY: namespace
	router.GET("/workload/daemon_set/describe", Auth(dtos.USER), describeDaemonSet) // QUERY: namespace, name
	router.DELETE("/workload/daemon_set", Auth(dtos.USER), deleteDaemonSet)         // BODY: json-object
	router.PATCH("/workload/daemon_set", Auth(dtos.USER), patchDaemonSet)           // BODY: json-object
	router.POST("/workload/daemon_set", Auth(dtos.USER), createDaemonSet)           // BODY: yaml-object

	router.GET("/workload/stateful_set", Auth(dtos.USER), allStatefulSets)              // QUERY: namespace
	router.GET("/workload/stateful_set/describe", Auth(dtos.USER), describeStatefulSet) // QUERY: namespace, name
	router.DELETE("/workload/stateful_set", Auth(dtos.USER), deleteStatefulSet)         // BODY: json-object
	router.PATCH("/workload/stateful_set", Auth(dtos.USER), patchStatefulSet)           // BODY: json-object
	router.POST("/workload/stateful_set", Auth(dtos.USER), createStatefulSet)           // BODY: yaml-object

	router.GET("/workload/job", Auth(dtos.USER), allJobs)              // QUERY: namespace
	router.GET("/workload/job/describe", Auth(dtos.USER), describeJob) // QUERY: namespace, name
	router.DELETE("/workload/job", Auth(dtos.USER), deleteJob)         // BODY: json-object
	router.PATCH("/workload/job", Auth(dtos.USER), patchJob)           // BODY: json-object
	router.POST("/workload/job", Auth(dtos.USER), createJob)           // BODY: yaml-object

	router.GET("/workload/cron_job", Auth(dtos.USER), allCronJobs)              // QUERY: namespace
	router.GET("/workload/cron_job/describe", Auth(dtos.USER), describeCronJob) // QUERY: namespace, name
	router.DELETE("/workload/cron_job", Auth(dtos.USER), deleteCronJob)         // BODY: json-object
	router.PATCH("/workload/cron_job", Auth(dtos.USER), patchCronJob)           // BODY: json-object
	router.POST("/workload/cron_job", Auth(dtos.USER), createCronJob)           // BODY: yaml-object

	router.GET("/workload/replica_set", Auth(dtos.USER), allReplicasets)              // QUERY: namespace
	router.GET("/workload/replica_set/describe", Auth(dtos.USER), describeReplicaset) // QUERY: namespace, name
	router.DELETE("/workload/replica_set", Auth(dtos.USER), deleteReplicaset)         // BODY: json-object
	router.PATCH("/workload/replica_set", Auth(dtos.USER), patchReplicaset)           // BODY: json-object
	router.POST("/workload/replica_set", Auth(dtos.USER), createReplicaset)           // BODY: yaml-object

	router.GET("/workload/persistent_volume", Auth(dtos.ADMIN), allPersistentVolumes)              // QUERY: -
	router.GET("/workload/persistent_volume/describe", Auth(dtos.ADMIN), describePersistentVolume) // QUERY: name
	router.DELETE("/workload/persistent_volume", Auth(dtos.ADMIN), deletePersistentVolume)         // BODY: json-object
	router.PATCH("/workload/persistent_volume", Auth(dtos.ADMIN), patchPersistentVolume)           // BODY: json-object
	router.POST("/workload/persistent_volume", Auth(dtos.ADMIN), createPersistentVolume)           // BODY: yaml-object

	router.GET("/workload/persistent_volume_claim", Auth(dtos.USER), allPersistentVolumeClaims)              // QUERY: namespace
	router.GET("/workload/persistent_volume_claim/describe", Auth(dtos.USER), describePersistentVolumeClaim) // QUERY: namespace, name
	router.DELETE("/workload/persistent_volume_claim", Auth(dtos.ADMIN), deletePersistentVolumeClaim)        // BODY: json-object
	router.PATCH("/workload/persistent_volume_claim", Auth(dtos.ADMIN), patchPersistentVolumeClaim)          // BODY: json-object
	router.POST("/workload/persistent_volume_claim", Auth(dtos.ADMIN), createPersistentVolumeClaim)          // BODY: yaml-object

	router.GET("/workload/horizontal_pod_autoscaler", Auth(dtos.USER), allHpas)              // QUERY: namespace
	router.GET("/workload/horizontal_pod_autoscaler/describe", Auth(dtos.USER), describeHpa) // QUERY: namespace, name
	router.DELETE("/workload/horizontal_pod_autoscaler", Auth(dtos.ADMIN), deleteHpa)        // BODY: json-object
	router.PATCH("/workload/horizontal_pod_autoscaler", Auth(dtos.ADMIN), patchHpa)          // BODY: json-object
	router.POST("/workload/horizontal_pod_autoscaler", Auth(dtos.ADMIN), createHpa)          // BODY: yaml-object

	router.GET("/workload/event", Auth(dtos.USER), allEvents)              // QUERY: namespace
	router.GET("/workload/event/describe", Auth(dtos.USER), describeEvent) // QUERY: namespace, name

	router.GET("/workload/certificate", Auth(dtos.USER), allCertificates)              // QUERY: namespace
	router.GET("/workload/certificate/describe", Auth(dtos.USER), describeCertificate) // QUERY: namespace, name
	router.DELETE("/workload/certificate", Auth(dtos.USER), deleteCertificate)         // BODY: json-object
	router.PATCH("/workload/certificate", Auth(dtos.USER), patchCertificate)           // BODY: json-object
	router.POST("/workload/certificate", Auth(dtos.USER), createCertificate)           // BODY: yaml-object

	router.GET("/workload/certificaterequest", Auth(dtos.USER), allCertificateRequests)              // QUERY: namespace
	router.GET("/workload/certificaterequest/describe", Auth(dtos.USER), describeCertificateRequest) // QUERY: namespace, name
	router.DELETE("/workload/certificaterequest", Auth(dtos.USER), deleteCertificateRequest)         // BODY: json-object
	router.PATCH("/workload/certificaterequest", Auth(dtos.USER), patchCertificateRequest)           // BODY: json-object
	router.POST("/workload/certificaterequest", Auth(dtos.USER), createCertificateRequest)           // BODY: yaml-object

	router.GET("/workload/orders", Auth(dtos.USER), allOrders)              // QUERY: namespace
	router.GET("/workload/orders/describe", Auth(dtos.USER), describeOrder) // QUERY: namespace, name
	router.DELETE("/workload/orders", Auth(dtos.USER), deleteOrder)         // BODY: json-object
	router.PATCH("/workload/orders", Auth(dtos.USER), patchOrder)           // BODY: json-object
	router.POST("/workload/orders", Auth(dtos.USER), createOrder)           // BODY: yaml-object

	router.GET("/workload/issuer", Auth(dtos.USER), allIssuers)              // QUERY: namespace
	router.GET("/workload/issuer/describe", Auth(dtos.USER), describeIssuer) // QUERY: namespace, name
	router.DELETE("/workload/issuer", Auth(dtos.USER), deleteIssuer)         // BODY: json-object
	router.PATCH("/workload/issuer", Auth(dtos.USER), patchIssuer)           // BODY: json-object
	router.POST("/workload/issuer", Auth(dtos.USER), createIssuer)           // BODY: yaml-object

	router.GET("/workload/clusterissuer", Auth(dtos.ADMIN), allClusterIssuers)              // QUERY: -
	router.GET("/workload/clusterissuer/describe", Auth(dtos.ADMIN), describeClusterIssuer) // QUERY: name
	router.DELETE("/workload/clusterissuer", Auth(dtos.ADMIN), deleteClusterIssuer)         // BODY: json-object
	router.PATCH("/workload/clusterissuer", Auth(dtos.ADMIN), patchClusterIssuer)           // BODY: json-object
	router.POST("/workload/clusterissuer", Auth(dtos.ADMIN), createClusterIssuer)           // BODY: yaml-object

	router.GET("/workload/service_account", Auth(dtos.ADMIN), allServiceAccounts)              // QUERY: namespace
	router.GET("/workload/service_account/describe", Auth(dtos.ADMIN), describeServiceAccount) // QUERY: namespace, name
	router.DELETE("/workload/service_account", Auth(dtos.ADMIN), deleteServiceAccount)         // BODY: json-object
	router.PATCH("/workload/service_account", Auth(dtos.ADMIN), patchServiceAccount)           // BODY: json-object
	router.POST("/workload/service_account", Auth(dtos.ADMIN), createServiceAccount)           // BODY: yaml-object

	router.GET("/workload/role", Auth(dtos.USER), allRoles)              // QUERY: namespace
	router.GET("/workload/role/describe", Auth(dtos.USER), describeRole) // QUERY: namespace, name
	router.DELETE("/workload/role", Auth(dtos.ADMIN), deleteRole)        // BODY: json-object
	router.PATCH("/workload/role", Auth(dtos.ADMIN), patchRole)          // BODY: json-object
	router.POST("/workload/role", Auth(dtos.ADMIN), createRole)          // BODY: yaml-object

	router.GET("/workload/role_binding", Auth(dtos.USER), allRoleBindings)              // QUERY: namespace
	router.GET("/workload/role_binding/describe", Auth(dtos.USER), describeRoleBinding) // QUERY: namespace, name
	router.DELETE("/workload/role_binding", Auth(dtos.ADMIN), deleteRoleBinding)        // BODY: json-object
	router.PATCH("/workload/role_binding", Auth(dtos.ADMIN), patchRoleBinding)          // BODY: json-object
	router.POST("/workload/role_binding", Auth(dtos.ADMIN), createRoleBinding)          // BODY: yaml-object

	router.GET("/workload/cluster_role", Auth(dtos.ADMIN), allClusterRoles)              // QUERY: -
	router.GET("/workload/cluster_role/describe", Auth(dtos.ADMIN), describeClusterRole) // QUERY: name
	router.DELETE("/workload/cluster_role", Auth(dtos.ADMIN), deleteClusterRole)         // BODY: json-object
	router.PATCH("/workload/cluster_role", Auth(dtos.ADMIN), patchClusterRole)           // BODY: json-object
	router.POST("/workload/cluster_role", Auth(dtos.ADMIN), createClusterRole)           // BODY: yaml-object

	router.GET("/workload/cluster_role_binding", Auth(dtos.ADMIN), allClusterRoleBindings)              // QUERY: -
	router.GET("/workload/cluster_role_binding/describe", Auth(dtos.ADMIN), describeClusterRoleBinding) // QUERY: name
	router.DELETE("/workload/cluster_role_binding", Auth(dtos.ADMIN), deleteClusterRoleBinding)         // BODY: json-object
	router.PATCH("/workload/cluster_role_binding", Auth(dtos.ADMIN), patchClusterRoleBinding)           // BODY: json-object
	router.POST("/workload/cluster_role_binding", Auth(dtos.ADMIN), createClusterRoleBinding)           // BODY: yaml-object

	router.GET("/workload/volume_attachment", Auth(dtos.ADMIN), allVolumeAttachments)              // QUERY: -
	router.GET("/workload/volume_attachment/describe", Auth(dtos.ADMIN), describeVolumeAttachment) // QUERY: name
	router.DELETE("/workload/volume_attachment", Auth(dtos.ADMIN), deleteVolumeAttachment)         // BODY: json-object
	router.PATCH("/workload/volume_attachment", Auth(dtos.ADMIN), patchVolumeAttachment)           // BODY: json-object
	router.POST("/workload/volume_attachment", Auth(dtos.ADMIN), createVolumeAttachment)           // BODY: yaml-object

	router.GET("/workload/network_policy", Auth(dtos.USER), allNetworkPolicies)             // QUERY: namespace
	router.GET("/workload/network_policy/describe", Auth(dtos.USER), describeNetworkPolicy) // QUERY: namespace, name
	router.DELETE("/workload/network_policy", Auth(dtos.ADMIN), deleteNetworkPolicy)        // BODY: json-object
	router.PATCH("/workload/network_policy", Auth(dtos.ADMIN), patchNetworkPolicy)          // BODY: json-object
	router.POST("/workload/network_policy", Auth(dtos.ADMIN), createNetworkPolicy)          // BODY: yaml-object

	router.GET("/workload/storageclass", Auth(dtos.USER), allStorageClasses)             // QUERY: namespace
	router.GET("/workload/storageclass/describe", Auth(dtos.USER), describeStorageClass) // QUERY: namespace, name
	router.DELETE("/workload/storageclass", Auth(dtos.ADMIN), deleteStorageClass)        // BODY: json-object
	router.PATCH("/workload/storageclass", Auth(dtos.ADMIN), patchStorageClass)          // BODY: json-object
	router.POST("/workload/storageclass", Auth(dtos.ADMIN), createStorageClass)          // BODY: yaml-object

	router.GET("/workload/crds", Auth(dtos.ADMIN), allCrds)              // QUERY: -
	router.GET("/workload/crds/describe", Auth(dtos.ADMIN), describeCrd) // QUERY: name
	router.DELETE("/workload/crds", Auth(dtos.ADMIN), deleteCrd)         // BODY: json-object
	router.PATCH("/workload/crds", Auth(dtos.ADMIN), patchCrd)           // BODY: json-object
	router.POST("/workload/crds", Auth(dtos.ADMIN), createCrd)           // BODY: yaml-object

	router.GET("/workload/endpoints", Auth(dtos.USER), allEndpoints)              // QUERY: namespace
	router.GET("/workload/endpoints/describe", Auth(dtos.USER), describeEndpoint) // QUERY: namespace, name
	router.DELETE("/workload/endpoints", Auth(dtos.USER), deleteEndpoint)         // BODY: json-object
	router.PATCH("/workload/endpoints", Auth(dtos.USER), patchEndpoint)           // BODY: json-object
	router.POST("/workload/endpoints", Auth(dtos.USER), createEndpoint)           // BODY: yaml-object

	router.GET("/workload/leases", Auth(dtos.USER), allLeases)              // QUERY: namespace
	router.GET("/workload/leases/describe", Auth(dtos.USER), describeLease) // QUERY: namespace, name
	router.DELETE("/workload/leases", Auth(dtos.USER), deleteLease)         // BODY: json-object
	router.PATCH("/workload/leases", Auth(dtos.USER), patchLease)           // BODY: json-object
	router.POST("/workload/leases", Auth(dtos.USER), createLease)           // BODY: yaml-object

	router.GET("/workload/priorityclasses", Auth(dtos.ADMIN), allPriorityClasses)             // QUERY: -
	router.GET("/workload/priorityclasses/describe", Auth(dtos.ADMIN), describePriorityClass) // QUERY: name
	router.DELETE("/workload/priorityclasses", Auth(dtos.ADMIN), deletePriorityClass)         // BODY: json-object
	router.PATCH("/workload/priorityclasses", Auth(dtos.ADMIN), patchPriorityClass)           // BODY: json-object
	router.POST("/workload/priorityclasses", Auth(dtos.ADMIN), createPriorityClass)           // BODY: yaml-object

	router.GET("/workload/volumesnapshots", Auth(dtos.USER), allVolumeSnapshots)              // QUERY: namespace
	router.GET("/workload/volumesnapshots/describe", Auth(dtos.USER), describeVolumeSnapshot) // QUERY: namespace, name
	router.DELETE("/workload/volumesnapshots", Auth(dtos.USER), deleteVolumeSnapshot)         // BODY: json-object
	router.PATCH("/workload/volumesnapshots", Auth(dtos.USER), patchVolumeSnapshot)           // BODY: json-object
	router.POST("/workload/volumesnapshots", Auth(dtos.USER), createVolumeSnapshot)           // BODY: yaml-object

	router.GET("/workload/resourcequota", Auth(dtos.ADMIN), allResourceQuotas)              // QUERY: namespace
	router.GET("/workload/resourcequota/describe", Auth(dtos.ADMIN), describeResourceQuota) // QUERY: namespace, name
	router.DELETE("/workload/resourcequota", Auth(dtos.ADMIN), deleteResourceQuota)         // BODY: json-object
	router.PATCH("/workload/resourcequota", Auth(dtos.ADMIN), patchResourceQuota)           // BODY: json-object
	router.POST("/workload/resourcequota", Auth(dtos.ADMIN), createResourceQuota)           // BODY: yaml-object
}

// GENERAL
// @Tags General
// @Produce json
// @Success 200 {array} kubernetes.K8sNewWorkload
// @Router /workload/templates [get]
func allWorkloadTemplates(c *gin.Context) {
	c.JSON(http.StatusOK, kubernetes.ListCreateTemplates())
}

// @Tags General
// @Produce json
// @Success 200 {array} string
// @Router /workload/available-resources [get]
func allKubernetesResources(c *gin.Context) {
	user, err := CheckUserAuthorization(c)
	if err != nil || user == nil {
		MalformedMessage(c, "User not found.")
		return
	}
	c.JSON(http.StatusOK, kubernetes.WorkloadsForAccesslevel(user.AccessLevel))
}

// NAMESPACES
func allNamespaces(c *gin.Context) {
	c.JSON(http.StatusOK, kubernetes.ListAllNamespace())
}
func createNamespace(c *gin.Context) {
	var data v1.Namespace
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, kubernetes.CreateK8sNamespace(data))
}
func deleteNamespace(c *gin.Context) {
	var data v1.Namespace
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, kubernetes.DeleteK8sNamespace(data))
}

// PODS
func allPods(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllK8sPods(namespace))
}
func describePod(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sPod(namespace, name))
}
func deletePod(c *gin.Context) {
	var data v1.Pod
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sPod(data))
}
func patchPod(c *gin.Context) {
	var data v1.Pod
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sPod(data))
}
func createPod(c *gin.Context) {
	var data v1.Pod
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sPod(data))
}

// DEPLOYMENTS
func allDeployments(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllK8sDeployments(namespace))
}
func describeDeployment(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sDeployment(namespace, name))
}
func deleteDeployment(c *gin.Context) {
	var data v1Apps.Deployment
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sDeployment(data))
}
func patchDeployment(c *gin.Context) {
	var data v1Apps.Deployment
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sDeployment(data))
}
func createDeployment(c *gin.Context) {
	var data v1Apps.Deployment
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sDeployment(data))
}

// SERVICES
func allServices(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllK8sServices(namespace))
}
func describeService(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sService(namespace, name))
}
func deleteService(c *gin.Context) {
	var data v1.Service
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sService(data))
}
func patchService(c *gin.Context) {
	var data v1.Service
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sService(data))
}
func createService(c *gin.Context) {
	var data v1.Service
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sService(data))
}

// INGRESSES
func allIngresses(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllK8sIngresses(namespace))
}
func describeIngress(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sIngress(namespace, name))
}
func deleteIngress(c *gin.Context) {
	var data v1Networking.Ingress
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sIngress(data))
}
func patchIngress(c *gin.Context) {
	var data v1Networking.Ingress
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sIngress(data))
}
func createIngress(c *gin.Context) {
	var data v1Networking.Ingress
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sIngress(data))
}

// CONFIGMAPS
func allConfigmaps(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllK8sConfigmaps(namespace))
}
func describeConfigmap(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sConfigmap(namespace, name))
}
func deleteConfigmap(c *gin.Context) {
	var data v1.ConfigMap
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sConfigmap(data))
}
func patchConfigmap(c *gin.Context) {
	var data v1.ConfigMap
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sConfigMap(data))
}
func createConfigmap(c *gin.Context) {
	var data v1.ConfigMap
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sConfigMap(data))
}

// SECRETS
func allSecrets(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllK8sSecrets(namespace))
}
func describeSecret(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sSecret(namespace, name))
}
func deleteSecret(c *gin.Context) {
	var data v1.Secret
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sSecret(data))
}
func patchSecret(c *gin.Context) {
	var data v1.Secret
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sSecret(data))
}
func createSecret(c *gin.Context) {
	var data v1.Secret
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sSecret(data))
}

// NODES
func allNodes(c *gin.Context) {
	RespondForWorkloadResult(c, kubernetes.ListK8sNodes())
}
func describeNode(c *gin.Context) {
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sNode(name))
}

// DAEMONSETS
func allDaemonSets(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllK8sDaemonsets(namespace))
}
func describeDaemonSet(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sDaemonSet(namespace, name))
}
func deleteDaemonSet(c *gin.Context) {
	var data v1Apps.DaemonSet
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sDaemonSet(data))
}
func patchDaemonSet(c *gin.Context) {
	var data v1Apps.DaemonSet
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sDaemonSet(data))
}
func createDaemonSet(c *gin.Context) {
	var data v1Apps.DaemonSet
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sDaemonSet(data))
}

// STATEFULSETS
func allStatefulSets(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllStatefulSets(namespace))
}
func describeStatefulSet(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sStatefulset(namespace, name))
}
func deleteStatefulSet(c *gin.Context) {
	var data v1Apps.StatefulSet
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sStatefulset(data))
}
func patchStatefulSet(c *gin.Context) {
	var data v1Apps.StatefulSet
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sStatefulset(data))
}
func createStatefulSet(c *gin.Context) {
	var data v1Apps.StatefulSet
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sStatefulset(data))
}

// JOBS
func allJobs(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllJobs(namespace))
}
func describeJob(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sJob(namespace, name))
}
func deleteJob(c *gin.Context) {
	var data v1Job.Job
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sJob(data))
}
func patchJob(c *gin.Context) {
	var data v1Job.Job
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sJob(data))
}
func createJob(c *gin.Context) {
	var data v1Job.Job
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sJob(data))
}

// CRONJOBS
func allCronJobs(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllCronjobs(namespace))
}
func describeCronJob(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sCronJob(namespace, name))
}
func deleteCronJob(c *gin.Context) {
	var data v1Job.CronJob
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sCronJob(data))
}
func patchCronJob(c *gin.Context) {
	var data v1Job.CronJob
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sCronJob(data))
}
func createCronJob(c *gin.Context) {
	var data v1Job.CronJob
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sCronJob(data))
}

// REPLICASETS
func allReplicasets(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllK8sReplicasets(namespace))
}
func describeReplicaset(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sReplicaset(namespace, name))
}
func deleteReplicaset(c *gin.Context) {
	var data v1Apps.ReplicaSet
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sReplicaset(data))
}
func patchReplicaset(c *gin.Context) {
	var data v1Apps.ReplicaSet
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sReplicaset(data))
}
func createReplicaset(c *gin.Context) {
	var data v1Apps.ReplicaSet
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sReplicaSet(data))
}

// PERSISTENTVOLUMES
func allPersistentVolumes(c *gin.Context) {
	RespondForWorkloadResult(c, kubernetes.AllPersistentVolumes())
}
func describePersistentVolume(c *gin.Context) {
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sPersistentVolume(name))
}
func deletePersistentVolume(c *gin.Context) {
	var data v1.PersistentVolume
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sPersistentVolume(data))
}
func patchPersistentVolume(c *gin.Context) {
	var data v1.PersistentVolume
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sPersistentVolume(data))
}
func createPersistentVolume(c *gin.Context) {
	var data v1.PersistentVolume
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sPersistentVolume(data))
}

// PERSISTENTVOLUMECLAIMS
func allPersistentVolumeClaims(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllK8sPersistentVolumeClaims(namespace))
}
func describePersistentVolumeClaim(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sPersistentVolumeClaim(namespace, name))
}
func deletePersistentVolumeClaim(c *gin.Context) {
	var data v1.PersistentVolumeClaim
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sPersistentVolumeClaim(data))
}
func patchPersistentVolumeClaim(c *gin.Context) {
	var data v1.PersistentVolumeClaim
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sPersistentVolumeClaim(data))
}
func createPersistentVolumeClaim(c *gin.Context) {
	var data v1.PersistentVolumeClaim
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sPersistentVolumeClaim(data))
}

// HPA
func allHpas(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllHpas(namespace))
}
func describeHpa(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sHpa(namespace, name))
}
func deleteHpa(c *gin.Context) {
	var data v2Scale.HorizontalPodAutoscaler
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sHpa(data))
}
func patchHpa(c *gin.Context) {
	var data v2Scale.HorizontalPodAutoscaler
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sHpa(data))
}
func createHpa(c *gin.Context) {
	var data v2Scale.HorizontalPodAutoscaler
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sHpa(data))
}

// EVENTS
func allEvents(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllEvents(namespace))
}
func describeEvent(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sEvent(namespace, name))
}

// CERTIFICATES
func allCertificates(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllK8sCertificates(namespace))
}
func describeCertificate(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sCertificate(namespace, name))
}
func deleteCertificate(c *gin.Context) {
	var data cmapi.Certificate
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sCertificate(data))
}
func patchCertificate(c *gin.Context) {
	var data cmapi.Certificate
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sCertificate(data))
}
func createCertificate(c *gin.Context) {
	var data cmapi.Certificate
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sCertificate(data))
}

// CERTIFICATEREQUESTS
func allCertificateRequests(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllCertificateSigningRequests(namespace))
}
func describeCertificateRequest(c *gin.Context) {
	name := c.Query("name")
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sCertificateSigningRequest(namespace, name))
}
func deleteCertificateRequest(c *gin.Context) {
	var data cmapi.CertificateRequest
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sCertificateSigningRequest(data))
}
func patchCertificateRequest(c *gin.Context) {
	var data cmapi.CertificateRequest
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sCertificateSigningRequest(data))
}
func createCertificateRequest(c *gin.Context) {
	var data cmapi.CertificateRequest
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sCertificateSigningRequest(data))
}

// ORDERS
func allOrders(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllOrders(namespace))
}
func describeOrder(c *gin.Context) {
	name := c.Query("name")
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sOrder(namespace, name))
}
func deleteOrder(c *gin.Context) {
	var data v1Cert.Order
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sOrder(data))
}
func patchOrder(c *gin.Context) {
	var data v1Cert.Order
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sOrder(data))
}
func createOrder(c *gin.Context) {
	var data v1Cert.Order
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sOrder(data))
}

// ISSUERS
func allIssuers(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllIssuer(namespace))
}
func describeIssuer(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sIssuer(namespace, name))
}
func deleteIssuer(c *gin.Context) {
	var data cmapi.Issuer
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sIssuer(data))
}
func patchIssuer(c *gin.Context) {
	var data cmapi.Issuer
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sIssuer(data))
}
func createIssuer(c *gin.Context) {
	var data cmapi.Issuer
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sIssuer(data))
}

// CLUSTERISSUERS
func allClusterIssuers(c *gin.Context) {
	RespondForWorkloadResult(c, kubernetes.AllClusterIssuers())
}
func describeClusterIssuer(c *gin.Context) {
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sClusterIssuer(name))
}
func deleteClusterIssuer(c *gin.Context) {
	var data cmapi.ClusterIssuer
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sClusterIssuer(data))
}
func patchClusterIssuer(c *gin.Context) {
	var data cmapi.ClusterIssuer
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sClusterIssuer(data))
}
func createClusterIssuer(c *gin.Context) {
	var data cmapi.ClusterIssuer
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sClusterIssuer(data))
}

// SERVICEACCOUNTS
func allServiceAccounts(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllServiceAccounts(namespace))
}
func describeServiceAccount(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sServiceAccount(namespace, name))
}
func deleteServiceAccount(c *gin.Context) {
	var data v1.ServiceAccount
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sServiceAccount(data))
}
func patchServiceAccount(c *gin.Context) {
	var data v1.ServiceAccount
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sServiceAccount(data))
}
func createServiceAccount(c *gin.Context) {
	var data v1.ServiceAccount
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sServiceAccount(data))
}

// ROLES
func allRoles(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllRoles(namespace))
}
func describeRole(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sRole(namespace, name))
}
func deleteRole(c *gin.Context) {
	var data v1Rbac.Role
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sRole(data))
}
func patchRole(c *gin.Context) {
	var data v1Rbac.Role
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sRole(data))
}
func createRole(c *gin.Context) {
	var data v1Rbac.Role
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sRole(data))
}

// ROLEBINDINGS
func allRoleBindings(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllRoleBindings(namespace))
}
func describeRoleBinding(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sRoleBinding(namespace, name))
}
func deleteRoleBinding(c *gin.Context) {
	var data v1Rbac.RoleBinding
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sRoleBinding(data))
}
func patchRoleBinding(c *gin.Context) {
	var data v1Rbac.RoleBinding
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sRoleBinding(data))
}
func createRoleBinding(c *gin.Context) {
	var data v1Rbac.RoleBinding
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sRoleBinding(data))
}

// CLUSTERROLES
func allClusterRoles(c *gin.Context) {
	RespondForWorkloadResult(c, kubernetes.AllClusterRoles())
}
func describeClusterRole(c *gin.Context) {
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sClusterRole(name))
}
func deleteClusterRole(c *gin.Context) {
	var data v1Rbac.ClusterRole
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sClusterRole(data))
}
func patchClusterRole(c *gin.Context) {
	var data v1Rbac.ClusterRole
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sClusterRole(data))
}
func createClusterRole(c *gin.Context) {
	var data v1Rbac.ClusterRole
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sClusterRole(data))
}

// CLUSTERROLEBINDINGS
func allClusterRoleBindings(c *gin.Context) {
	RespondForWorkloadResult(c, kubernetes.AllClusterRoleBindings())
}
func describeClusterRoleBinding(c *gin.Context) {
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sClusterRoleBinding(name))
}
func deleteClusterRoleBinding(c *gin.Context) {
	var data v1Rbac.ClusterRoleBinding
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sClusterRoleBinding(data))
}
func patchClusterRoleBinding(c *gin.Context) {
	var data v1Rbac.ClusterRoleBinding
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sClusterRoleBinding(data))
}
func createClusterRoleBinding(c *gin.Context) {
	var data v1Rbac.ClusterRoleBinding
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sClusterRoleBinding(data))
}

// VOLUMEATTACHMENTS
func allVolumeAttachments(c *gin.Context) {
	RespondForWorkloadResult(c, kubernetes.AllVolumeAttachments())
}
func describeVolumeAttachment(c *gin.Context) {
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sVolumeAttachment(name))
}
func deleteVolumeAttachment(c *gin.Context) {
	var data v1Storage.VolumeAttachment
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sVolumeAttachment(data))
}
func patchVolumeAttachment(c *gin.Context) {
	var data v1Storage.VolumeAttachment
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sVolumeAttachment(data))
}
func createVolumeAttachment(c *gin.Context) {
	var data v1Storage.VolumeAttachment
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sVolumeAttachment(data))
}

// NETWORKPOLICIES
func allNetworkPolicies(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllNetworkPolicies(namespace))
}
func describeNetworkPolicy(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sNetworkPolicy(namespace, name))
}
func deleteNetworkPolicy(c *gin.Context) {
	var data v1Networking.NetworkPolicy
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sNetworkPolicy(data))
}
func patchNetworkPolicy(c *gin.Context) {
	var data v1Networking.NetworkPolicy
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sNetworkPolicy(data))
}
func createNetworkPolicy(c *gin.Context) {
	var data v1Networking.NetworkPolicy
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sNetworkpolicy(data))
}

// STORAGECLASSES
func allStorageClasses(c *gin.Context) {
	RespondForWorkloadResult(c, kubernetes.AllStorageClasses())
}
func describeStorageClass(c *gin.Context) {
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sStorageClass(name))
}
func deleteStorageClass(c *gin.Context) {
	var data v1Storage.StorageClass
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sStorageClass(data))
}
func patchStorageClass(c *gin.Context) {
	var data v1Storage.StorageClass
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sStorageClass(data))
}
func createStorageClass(c *gin.Context) {
	var data v1Storage.StorageClass
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sStorageClass(data))
}

// CRDS
func allCrds(c *gin.Context) {
	RespondForWorkloadResult(c, kubernetes.AllCustomResourceDefinitions())
}
func describeCrd(c *gin.Context) {
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sCustomResourceDefinition(name))
}
func deleteCrd(c *gin.Context) {
	var data apiExt.CustomResourceDefinition
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sCustomResourceDefinition(data))
}
func patchCrd(c *gin.Context) {
	var data apiExt.CustomResourceDefinition
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sCustomResourceDefinition(data))
}
func createCrd(c *gin.Context) {
	var data apiExt.CustomResourceDefinition
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sCustomResourceDefinition(data))
}

// ENDPOINTS
func allEndpoints(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllEndpoints(namespace))
}
func describeEndpoint(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sEndpoint(namespace, name))
}
func deleteEndpoint(c *gin.Context) {
	var data v1.Endpoints
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sEndpoint(data))
}
func patchEndpoint(c *gin.Context) {
	var data v1.Endpoints
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sEndpoint(data))
}
func createEndpoint(c *gin.Context) {
	var data v1.Endpoints
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sEndpoint(data))
}

// LEASES
func allLeases(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllLeases(namespace))
}
func describeLease(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sLease(namespace, name))
}
func deleteLease(c *gin.Context) {
	var data v1Coordination.Lease
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sLease(data))
}
func patchLease(c *gin.Context) {
	var data v1Coordination.Lease
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sLease(data))
}
func createLease(c *gin.Context) {
	var data v1Coordination.Lease
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sLease(data))
}

// PRIORITYCLASSES
func allPriorityClasses(c *gin.Context) {
	RespondForWorkloadResult(c, kubernetes.AllPriorityClasses())
}
func describePriorityClass(c *gin.Context) {
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sPriorityClass(name))
}
func deletePriorityClass(c *gin.Context) {
	var data v1Scheduling.PriorityClass
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sPriorityClass(data))
}
func patchPriorityClass(c *gin.Context) {
	var data v1Scheduling.PriorityClass
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sPriorityClass(data))
}
func createPriorityClass(c *gin.Context) {
	var data v1Scheduling.PriorityClass
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sPriorityClass(data))
}

// VOLUMESNAPSHOTS
func allVolumeSnapshots(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllVolumeSnapshots(namespace))
}
func describeVolumeSnapshot(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sVolumeSnapshot(namespace, name))
}
func deleteVolumeSnapshot(c *gin.Context) {
	var data v6Snap.VolumeSnapshot
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sVolumeSnapshot(data))
}
func patchVolumeSnapshot(c *gin.Context) {
	var data v6Snap.VolumeSnapshot
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sVolumeSnapshot(data))
}
func createVolumeSnapshot(c *gin.Context) {
	var data v6Snap.VolumeSnapshot
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sVolumeSnapshot(data))
}

// RESOURCEQUOTAS
func allResourceQuotas(c *gin.Context) {
	namespace := c.Query("namespace")
	RespondForWorkloadResult(c, kubernetes.AllResourceQuotas(namespace))
}
func describeResourceQuota(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	RespondForWorkloadResult(c, kubernetes.DescribeK8sResourceQuota(namespace, name))
}
func deleteResourceQuota(c *gin.Context) {
	var data v1.ResourceQuota
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.DeleteK8sResourceQuota(data))
}
func patchResourceQuota(c *gin.Context) {
	var data v1.ResourceQuota
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.UpdateK8sResourceQuota(data))
}
func createResourceQuota(c *gin.Context) {
	var data v1.ResourceQuota
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, kubernetes.CreateK8sResourceQuota(data))
}
