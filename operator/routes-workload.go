package operator

import (
	"net/http"

	"github.com/mogenius/punq/utils"

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

	workloadRoutes := router.Group("/workload")
	{
		workloadRoutes.GET("/templates", Auth(dtos.USER), allWorkloadTemplates)
		workloadRoutes.GET("/available-resources", Auth(dtos.READER), allKubernetesResources)

		// namespace
		namespaceWorkloadRoutes := router.Group("/namespace", Auth(dtos.USER))
		{
			namespaceWorkloadRoutes.GET("/", allNamespaces)                                                    // PARAM: -
			namespaceWorkloadRoutes.GET("/:name", validateParam("name"), describeNamespaces)                   // PARAM: name
			namespaceWorkloadRoutes.DELETE("/:name", Auth(dtos.ADMIN), validateParam("name"), deleteNamespace) // PARAM: name
			namespaceWorkloadRoutes.POST("", createNamespace)                                                  // BODY: yaml-object
		}

		// pod
		podWorkloadRoutes := router.Group("/pod", Auth(dtos.USER))
		{
			podWorkloadRoutes.GET("/", allPods)
			podWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describePod) // PARAM: namespace
			podWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deletePod)         // PARAM: namespace, name
			podWorkloadRoutes.PATCH("/", patchPod)                                                               // BODY: json-object
			podWorkloadRoutes.POST("/", createPod)                                                               // BODY: yaml-object
		}

		// deployment
		deploymentWorkloadRoutes := router.Group("/deployment", Auth(dtos.USER))
		{
			deploymentWorkloadRoutes.GET("/", allDeployments)                                                                  // PARAM: namespace
			deploymentWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeDeployment) // PARAM: namespace, name
			deploymentWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteDeployment)         // PARAM: namespace, name
			deploymentWorkloadRoutes.PATCH("/", patchDeployment)                                                               // BODY: json-object
			deploymentWorkloadRoutes.POST("/´", createDeployment)                                                              // BODY: yaml-object
		}

		// service
		serviceWorkloadRoutes := router.Group("/service", Auth(dtos.USER))
		{
			serviceWorkloadRoutes.GET("/", allServices)                                                                  // PARAM: namespace
			serviceWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeService) // PARAM: namespace, name
			serviceWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteService)         // PARAM: namespace, name
			serviceWorkloadRoutes.PATCH("/", patchService)                                                               // BODY: json-object
			serviceWorkloadRoutes.POST("/", createService)                                                               // BODY: yaml-object
		}

		// ingress
		ingressWorkloadRoutes := router.Group("/ingress", Auth(dtos.USER))
		{
			ingressWorkloadRoutes.GET("/", allIngresses)                                                                 // PARAM: namespace
			ingressWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeIngress) // PARAM: namespace, name
			ingressWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteIngress)         // PARAM: json-object
			ingressWorkloadRoutes.PATCH("/", patchIngress)                                                               // BODY: json-object
			ingressWorkloadRoutes.POST("/", createIngress)                                                               // BODY: yaml-object
		}

		// configmap
		configmapWorkloadRoutes := router.Group("/configmap", Auth(dtos.USER))
		{
			configmapWorkloadRoutes.GET("/", allConfigmaps)                                                                  // PARAM: namespace
			configmapWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeConfigmap) // PARAM: namespace, name
			configmapWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteConfigmap)         // PARAM: namespace, name
			configmapWorkloadRoutes.PATCH("/", patchConfigmap)                                                               // BODY: json-object
			workloadRoutes.POST("/´", createConfigmap)                                                                       // BODY: yaml-object
		}

		// secret
		secretWorkloadRoutes := router.Group("/secret", Auth(dtos.ADMIN))
		{
			secretWorkloadRoutes.GET("/", allSecrets)                                                                  // PARAM: namespace
			secretWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeSecret) // PARAM: namespace, name
			secretWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteSecret)         // PARAM: namespace, name
			secretWorkloadRoutes.PATCH("/", patchSecret)                                                               // BODY: json-object
			secretWorkloadRoutes.POST("/", createSecret)                                                               // BODY: yaml-object
		}

		// node
		nodeWorkloadRoutes := router.Group("/node", Auth(dtos.USER))
		{
			nodeWorkloadRoutes.GET("/", allNodes)                                          // -
			nodeWorkloadRoutes.GET("/describe/:name", validateParam("name"), describeNode) // PARAM: namespace
		}

		// daemon-set
		daemonSetWorkloadRoutes := router.Group("/daemon-set", Auth(dtos.USER))
		{
			daemonSetWorkloadRoutes.GET("/", allDaemonSets)                                                                  // PARAM: namespace
			daemonSetWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeDaemonSet) // PARAM: namespace, name
			daemonSetWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteDaemonSet)         // PARAM: namespace, name
			daemonSetWorkloadRoutes.PATCH("/", patchDaemonSet)                                                               // BODY: json-object
			daemonSetWorkloadRoutes.POST("/", createDaemonSet)                                                               // BODY: yaml-object

		}

		// stateful-set
		statefulSetWorkloadRoutes := router.Group("/stateful-set", Auth(dtos.USER))
		{
			statefulSetWorkloadRoutes.GET("/", allStatefulSets)                                                                  // PARAM: namespace
			statefulSetWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeStatefulSet) // PARAM: namespace, name
			statefulSetWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteStatefulSet)         // PARAM: namespace, name
			statefulSetWorkloadRoutes.PATCH("/", patchStatefulSet)                                                               // BODY: json-object
			statefulSetWorkloadRoutes.POST("/", createStatefulSet)                                                               // BODY: yaml-object
		}

		// job
		jobWorkloadRoutes := router.Group("/job", Auth(dtos.USER))
		{
			jobWorkloadRoutes.GET("/", allJobs)                                                                  // PARAM: namespace
			jobWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeJob) // PARAM: namespace, name
			jobWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteJob)         // PARAM: namespace, name
			jobWorkloadRoutes.PATCH("/", patchJob)                                                               // BODY: json-object
			jobWorkloadRoutes.POST("/", createJob)                                                               // BODY: yaml-object
		}

		// cron-job
		cronJobWorkloadRoutes := router.Group("/cron-job", Auth(dtos.USER))
		{
			cronJobWorkloadRoutes.GET("/", allCronJobs)                                                                  // PARAM: namespace
			cronJobWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeCronJob) // PARAM: namespace, name
			cronJobWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteCronJob)         // PARAM: namespace, name
			cronJobWorkloadRoutes.PATCH("/", patchCronJob)                                                               // BODY: json-object
			cronJobWorkloadRoutes.POST("/", createCronJob)                                                               // BODY: yaml-object
		}

		// replicaset
		replicaSetWorkloadRoutes := router.Group("/replica-set", Auth(dtos.USER))
		{
			replicaSetWorkloadRoutes.GET("/", allReplicasets)                                                                  // PARAM: namespace
			replicaSetWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeReplicaset) // PARAM: namespace, name
			replicaSetWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteReplicaset)         // PARAM: namespace, name
			replicaSetWorkloadRoutes.PATCH("(", patchReplicaset)                                                               // BODY: json-object
			replicaSetWorkloadRoutes.POST("/", createReplicaset)                                                               // BODY: yaml-object
		}

		// persistent-volume
		persistentVolumeWorkloadRoutes := router.Group("/persistent-volume", Auth(dtos.ADMIN))
		{
			persistentVolumeWorkloadRoutes.GET("/", allPersistentVolumes)                                          // PARAM: -
			persistentVolumeWorkloadRoutes.GET("/describe/:name", validateParam("name"), describePersistentVolume) // PARAM: name
			persistentVolumeWorkloadRoutes.DELETE("/:name", validateParam("name"), deletePersistentVolume)         // PARAM: name
			persistentVolumeWorkloadRoutes.PATCH("/", patchPersistentVolume)                                       // BODY: json-object
			persistentVolumeWorkloadRoutes.POST("/", createPersistentVolume)                                       // BODY: yaml-object
		}

		// persistent-volume-claim
		persistentVolumeClaimWorkloadRoutes := router.Group("/persistent-volume-claim")
		{
			persistentVolumeClaimWorkloadRoutes.GET("/", Auth(dtos.USER), validateParam("namespace"), allPersistentVolumeClaims)                                      // PARAM: namespace
			persistentVolumeClaimWorkloadRoutes.GET("/describe/:namespace/:name", Auth(dtos.USER), validateParam("namespace", "name"), describePersistentVolumeClaim) // PARAM: namespace, name
			persistentVolumeClaimWorkloadRoutes.DELETE("/:namespace/:name", Auth(dtos.ADMIN), validateParam("namespace", "name"), deletePersistentVolumeClaim)        // PARAM: namespace, name
			persistentVolumeClaimWorkloadRoutes.PATCH("/", Auth(dtos.ADMIN), patchPersistentVolumeClaim)                                                              // BODY: json-object
			persistentVolumeClaimWorkloadRoutes.POST("/", Auth(dtos.ADMIN), createPersistentVolumeClaim)                                                              // BODY: yaml-object
		}

		// horizontal-pod-autoscaler
		horizontalPodAutoscalerWorkloadRoutes := router.Group("/horizontal-pod-autoscaler")
		{
			horizontalPodAutoscalerWorkloadRoutes.GET("/", Auth(dtos.USER), validateParam("namespace"), allHpas)                                      // PARAM: namespace
			horizontalPodAutoscalerWorkloadRoutes.GET("/describe/:namespace/:name", Auth(dtos.USER), validateParam("namespace", "name"), describeHpa) // PARAM: namespace, name
			horizontalPodAutoscalerWorkloadRoutes.DELETE("/:namespace/:name", Auth(dtos.ADMIN), validateParam("namespace", "name"), deleteHpa)        // PARAM: namespace, name
			horizontalPodAutoscalerWorkloadRoutes.PATCH("/", Auth(dtos.ADMIN), patchHpa)                                                              // BODY: json-object
			horizontalPodAutoscalerWorkloadRoutes.POST("/", Auth(dtos.ADMIN), createHpa)                                                              // BODY: yaml-object
		}

		// event
		eventWorkloadRoutes := router.Group("/event", Auth(dtos.USER))
		{
			eventWorkloadRoutes.GET("/", allEvents)                                                                  // PARAM: namespace
			eventWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeEvent) // PARAM: namespace, name
		}

		// certificate
		certificateWorkloadRoutes := router.Group("/certificate", Auth(dtos.USER))
		{
			certificateWorkloadRoutes.GET("/", allCertificates)                                                                  // PARAM: namespace
			certificateWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeCertificate) // PARAM: namespace, name
			certificateWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteCertificate)         // PARAM: namespace, name
			certificateWorkloadRoutes.PATCH("/", patchCertificate)                                                               // BODY: json-object
			certificateWorkloadRoutes.POST("/", createCertificate)                                                               // BODY: yaml-object
		}

		// certificate-request
		certificateRequestWorkloadRoutes := router.Group("/certificate-request", Auth(dtos.USER))
		{
			certificateRequestWorkloadRoutes.GET("/", validateParam("name"), allCertificateRequests)                                           // PARAM: namespace
			certificateRequestWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeCertificateRequest) // PARAM: namespace, name
			certificateRequestWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteCertificateRequest)         // PARAM: namespace, name
			certificateRequestWorkloadRoutes.PATCH("/", patchCertificateRequest)                                                               // BODY: json-object
			certificateRequestWorkloadRoutes.POST("/", createCertificateRequest)                                                               // BODY: yaml-object
		}

		// orders
		ordersWorkloadRoutes := router.Group("/orders", Auth(dtos.USER))
		{
			ordersWorkloadRoutes.GET("/", allOrders)                                                                  // PARAM: namespace
			ordersWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeOrder) // PARAM: namespace, name
			ordersWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteOrder)         // PARAM: namespace, name
			ordersWorkloadRoutes.PATCH("/", patchOrder)                                                               // BODY: json-object
			ordersWorkloadRoutes.POST("/", createOrder)                                                               // BODY: yaml-object
		}

		// issuer
		issuerWorkloadRoutes := router.Group("/issuer", Auth(dtos.USER))
		{
			issuerWorkloadRoutes.GET("/", allIssuers)                                                                  // PARAM: namespace
			issuerWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeIssuer) // PARAM: namespace, name
			issuerWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteIssuer)         // PARAM: namespace, name
			issuerWorkloadRoutes.PATCH("/", patchIssuer)                                                               // BODY: json-object
			issuerWorkloadRoutes.POST("/", createIssuer)                                                               // BODY: yaml-object
		}

		// cluster-issuer
		clusterIssuerWorkloadRoutes := router.Group("/cluster-issuer", Auth(dtos.ADMIN))
		{
			clusterIssuerWorkloadRoutes.GET("/", allClusterIssuers)                                          // PARAM: -
			clusterIssuerWorkloadRoutes.GET("/describe/:name", validateParam("name"), describeClusterIssuer) // PARAM: name
			clusterIssuerWorkloadRoutes.DELETE("/:name", validateParam("name"), deleteClusterIssuer)         // PARAM: name
			clusterIssuerWorkloadRoutes.PATCH("/", patchClusterIssuer)                                       // BODY: json-object
			clusterIssuerWorkloadRoutes.POST("/", createClusterIssuer)                                       // BODY: yaml-object
		}

		// service-account
		serviceAccountWorkloadRoutes := router.Group("/service-account", Auth(dtos.ADMIN))
		{
			serviceAccountWorkloadRoutes.GET("/", allServiceAccounts)                                                                  // PARAM: namespace
			serviceAccountWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeServiceAccount) // PARAM: namespace, name
			serviceAccountWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteServiceAccount)         // PARAM: namespace, name
			serviceAccountWorkloadRoutes.PATCH("/", patchServiceAccount)                                                               // BODY: json-object
			serviceAccountWorkloadRoutes.POST("/", createServiceAccount)                                                               // BODY: yaml-object
		}

		// role
		roleWorkloadRoutes := router.Group("/role")
		{
			roleWorkloadRoutes.GET("/", Auth(dtos.USER), validateParam("namespace"), allRoles)                                      // PARAM: namespace
			roleWorkloadRoutes.GET("/describe/:namespace/:name", Auth(dtos.USER), validateParam("namespace", "name"), describeRole) // PARAM: namespace, name
			roleWorkloadRoutes.DELETE("/:namespace/:name", Auth(dtos.ADMIN), validateParam("namespace", "name"), deleteRole)        // PARAM: namespace, name
			roleWorkloadRoutes.PATCH("/", Auth(dtos.ADMIN), patchRole)                                                              // BODY: json-object
			roleWorkloadRoutes.POST("/", Auth(dtos.ADMIN), createRole)                                                              // BODY: yaml-object
		}

		// role-binding
		roleBindingWorkloadRoutes := router.Group("/role-binding")
		{
			roleBindingWorkloadRoutes.GET("/", Auth(dtos.USER), validateParam("namespace"), allRoleBindings)                                      // PARAM: namespace
			roleBindingWorkloadRoutes.GET("/describe/:namespace/:name", Auth(dtos.USER), validateParam("namespace", "name"), describeRoleBinding) // PARAM: namespace, name
			roleBindingWorkloadRoutes.DELETE("/:namespace/:name", Auth(dtos.ADMIN), validateParam("namespace", "name"), deleteRoleBinding)        // PARAM: namespace, name
			roleBindingWorkloadRoutes.PATCH("/", Auth(dtos.ADMIN), patchRoleBinding)                                                              // BODY: json-object
			roleBindingWorkloadRoutes.POST("/", Auth(dtos.ADMIN), createRoleBinding)                                                              // BODY: yaml-object
		}

		// cluster-role
		clusterRoleWorkloadRoutes := router.Group("/cluster-role", Auth(dtos.ADMIN))
		{
			clusterRoleWorkloadRoutes.GET("/", allClusterRoles)                                          // PARAM: -
			clusterRoleWorkloadRoutes.GET("/describe/:name", validateParam("name"), describeClusterRole) // PARAM: name
			clusterRoleWorkloadRoutes.DELETE("/:name", validateParam("name"), deleteClusterRole)         // PARAM: name
			clusterRoleWorkloadRoutes.PATCH("/cluster_role", patchClusterRole)                           // BODY: json-object
			clusterRoleWorkloadRoutes.POST("/cluster_role", createClusterRole)                           // BODY: yaml-object
		}

		// cluster-role-binding
		clusterRoleBindingWorkloadRoutes := router.Group("/cluster-role-binding", Auth(dtos.ADMIN))
		{
			clusterRoleBindingWorkloadRoutes.GET("/", allClusterRoleBindings)                                          // PARAM: -
			clusterRoleBindingWorkloadRoutes.GET("/describe/:name", validateParam("name"), describeClusterRoleBinding) // PARAM: name
			clusterRoleBindingWorkloadRoutes.DELETE("/:name", validateParam("name"), deleteClusterRoleBinding)         // PARAM: name
			clusterRoleBindingWorkloadRoutes.PATCH("/", patchClusterRoleBinding)                                       // BODY: json-object
			clusterRoleBindingWorkloadRoutes.POST("/", createClusterRoleBinding)                                       // BODY: yaml-object
		}

		// volume-attachment
		volumeAttachmentWorkloadRoutes := router.Group("/volume-attachment", Auth(dtos.ADMIN))
		{
			volumeAttachmentWorkloadRoutes.GET("/", allVolumeAttachments)                                          // PARAM: -
			volumeAttachmentWorkloadRoutes.GET("/describe/:name", validateParam("name"), describeVolumeAttachment) // PARAM: name
			volumeAttachmentWorkloadRoutes.DELETE("/:name", validateParam("name"), deleteVolumeAttachment)         // PARAM: name
			volumeAttachmentWorkloadRoutes.PATCH("/", patchVolumeAttachment)                                       // BODY: json-object
			volumeAttachmentWorkloadRoutes.POST("/", createVolumeAttachment)                                       // BODY: yaml-object
		}

		// network-policy
		networkPolicyWorkloadRoutes := router.Group("/network-policy")
		{
			networkPolicyWorkloadRoutes.GET("/", Auth(dtos.USER), validateParam("namespace"), allNetworkPolicies)                                     // PARAM: namespace
			networkPolicyWorkloadRoutes.GET("/describe/:namespace/:name", Auth(dtos.USER), validateParam("namespace", "name"), describeNetworkPolicy) // PARAM: namespace, name
			networkPolicyWorkloadRoutes.DELETE("/:namespace/:name", Auth(dtos.ADMIN), validateParam("namespace", "name"), deleteNetworkPolicy)        // PARAM: namespace, name
			networkPolicyWorkloadRoutes.PATCH("/", Auth(dtos.ADMIN), patchNetworkPolicy)                                                              // BODY: json-object
			networkPolicyWorkloadRoutes.POST("/", Auth(dtos.ADMIN), createNetworkPolicy)                                                              // BODY: yaml-object
		}

		// storage-class
		storageClassWorkloadRoutes := router.Group("/storage-class")
		{
			storageClassWorkloadRoutes.GET("/", Auth(dtos.USER), allStorageClasses)                                         // PARAM: namespace
			storageClassWorkloadRoutes.GET("/describe/:name", Auth(dtos.USER), validateParam("name"), describeStorageClass) // PARAM: namespace, name
			storageClassWorkloadRoutes.DELETE("/:name", Auth(dtos.ADMIN), validateParam("name"), deleteStorageClass)        // PARAM: namespace, name
			storageClassWorkloadRoutes.PATCH("/", Auth(dtos.ADMIN), patchStorageClass)                                      // BODY: json-object
			storageClassWorkloadRoutes.POST("/", Auth(dtos.ADMIN), createStorageClass)                                      // BODY: yaml-object
		}

		// crds
		crdsWorkloadRoutes := router.Group("/crds", Auth(dtos.ADMIN))
		{
			crdsWorkloadRoutes.GET("/", allCrds)                                                            // PARAM: -
			crdsWorkloadRoutes.GET("/describe/:name", validateParam("name"), Auth(dtos.ADMIN), describeCrd) // PARAM: name
			crdsWorkloadRoutes.DELETE("/:name", validateParam("name"), Auth(dtos.ADMIN), deleteCrd)         // PARAM: name
			crdsWorkloadRoutes.PATCH("/", Auth(dtos.ADMIN), patchCrd)                                       // BODY: json-object
			crdsWorkloadRoutes.POST("/", Auth(dtos.ADMIN), createCrd)                                       // BODY: yaml-object
		}

		// endpoints
		endpointsWorkloadRoutes := router.Group("/endpoints", Auth(dtos.USER))
		{
			endpointsWorkloadRoutes.GET("/", allEndpoints)                                                                  // PARAM: namespace
			endpointsWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeEndpoint) // PARAM: namespace, name
			endpointsWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteEndpoint)         // PARAM: namespace, name
			endpointsWorkloadRoutes.PATCH("/", patchEndpoint)                                                               // BODY: json-object
			endpointsWorkloadRoutes.POST("/", createEndpoint)                                                               // BODY: yaml-object
		}

		// leases
		leasesWorkloadRoutes := router.Group("/leases", Auth(dtos.USER))
		{
			leasesWorkloadRoutes.GET("/", allLeases)                                                                  // PARAM: namespace
			leasesWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeLease) // PARAM: namespace, name
			leasesWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteLease)         // PARAM: namespace, name
			leasesWorkloadRoutes.PATCH("/", patchLease)                                                               // BODY: json-object
			leasesWorkloadRoutes.POST("/", createLease)                                                               // BODY: yaml-object
		}

		// priority-classes
		priorityClassesWorkloadRoutes := router.Group("/priority-classes", Auth(dtos.ADMIN))
		{
			priorityClassesWorkloadRoutes.GET("/", allPriorityClasses)                                         // PARAM: -
			priorityClassesWorkloadRoutes.GET("/describe/:name", validateParam("name"), describePriorityClass) // PARAM: name
			priorityClassesWorkloadRoutes.DELETE("/:name", validateParam("name"), deletePriorityClass)         // PARAM: name
			priorityClassesWorkloadRoutes.PATCH("/", patchPriorityClass)                                       // BODY: json-object
			priorityClassesWorkloadRoutes.POST("/", createPriorityClass)                                       // BODY: yaml-object
		}

		// volume-snapshots
		volumeSnapshotsWorkloadRoutes := router.Group("/volume-snapshots", Auth(dtos.USER))
		{
			volumeSnapshotsWorkloadRoutes.GET("/", allVolumeSnapshots)                                                                  // PARAM: namespace
			volumeSnapshotsWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeVolumeSnapshot) // PARAM: namespace, name
			volumeSnapshotsWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteVolumeSnapshot)         // PARAM: namespace, name
			volumeSnapshotsWorkloadRoutes.PATCH("/", patchVolumeSnapshot)                                                               // BODY: json-object
			volumeSnapshotsWorkloadRoutes.POST("/", createVolumeSnapshot)                                                               // BODY: yaml-object
		}

		// resource-quota
		resourceQuotaWorkloadRoutes := router.Group("/resource-quota", Auth(dtos.ADMIN))
		{
			resourceQuotaWorkloadRoutes.GET("/", allResourceQuotas)                                                                  // PARAM: namespace
			resourceQuotaWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeResourceQuota) // PARAM: namespace, name
			resourceQuotaWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteResourceQuota)         // PARAM: namespace, name
			resourceQuotaWorkloadRoutes.PATCH("/", patchResourceQuota)                                                               // BODY: json-object
			resourceQuotaWorkloadRoutes.POST("/", createResourceQuota)                                                               // BODY: yaml-object
		}
	}
}

// GENERAL
// @Tags General
// @Produce json
// @Success 200 {array} kubernetes.K8sNewWorkload
// @Router /workload/templates [get]
// @Security Bearer
func allWorkloadTemplates(c *gin.Context) {
	c.JSON(http.StatusOK, kubernetes.ListCreateTemplates())
}

// @Tags General
// @Produce json
// @Success 200 {array} string
// @Router /workload/available-resources [get]
// @Security Bearer
func allKubernetesResources(c *gin.Context) {
	user, err := CheckUserAuthorization(c)
	if err != nil || user == nil {
		utils.MalformedMessage(c, "User not found.")
		return
	}
	c.JSON(http.StatusOK, kubernetes.WorkloadsForAccesslevel(user.AccessLevel))
}

// ---------------------- NAMESPACES ----------------------

// NAMESPACES
// @Tags Workloads
// @Produce json
// @Success 200 {array} v1.Namespace
// @Router /workload/namespace [get]
// @Param namespace query string false "name of the namespace"
// @Security Bearer
func allNamespaces(c *gin.Context) {
	c.JSON(http.StatusOK, kubernetes.ListK8sNamespaces(""))
}

// NAMESPACES
// @Tags Workloads
// @Produce json
// @Success 200 {array} v1.Namespace
// @Router /workload/namespace/{name} [get]
// @Param name path string false  "name of the namespace"
// @Security Bearer
func describeNamespaces(c *gin.Context) {
	name := c.Param("name")
	c.JSON(http.StatusOK, kubernetes.DescribeK8sNamespace(name))
}

// NAMESPACES
// @Tags Workloads
// @Produce json
// @Success 201 {object} utils.K8sWorkloadResult
// @Router /workload/namespace [post]
// @Security Bearer
func createNamespace(c *gin.Context) {
	var data v1.Namespace
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		return
	}
	c.JSON(http.StatusCreated, kubernetes.CreateK8sNamespace(data))
}

// NAMESPACES
// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/namespace/{name} [delete]
// @Param name path string false  "name of the namespace"
// @Security Bearer
func deleteNamespace(c *gin.Context) {
	name := c.Param("name")
	err := kubernetes.DeleteK8sNamespaceBy(name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// ---------------------- PODS ----------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/pod [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
func allPods(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sPods(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/pod/describe/{namespace}/{name} [get]
// @Param namespace path string true  "namespace name"
// @Param name path string true  "pod name"
// @Security Bearer
func describePod(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sPod(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/pod/{namespace}/{name} [delete]
// @Param namespace path string true "namespace name"
// @Param name path string true "pod name"
// @Security Bearer
func deletePod(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sPodBy(namespace, name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/pod [patch]
// @Security Bearer
func patchPod(c *gin.Context) {
	var data v1.Pod
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sPod(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/pod [post]
// @Security Bearer
func createPod(c *gin.Context) {
	var data v1.Pod
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sPod(data))

}

// ---------------------- DEPLOYMENTS ----------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/deployment [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
func allDeployments(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sDeployments(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/deployment/describe/{namespace}/{name} [get]
// @Param namespace path string true  "namespace name"
// @Param name path string true  "deployment name"
// @Security Bearer
func describeDeployment(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sDeployment(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/deployment/{namespace}/{name} [delete]
// @Param namespace path string true  "namespace name"
// @Param name path string true  "deployment name"
// @Security Bearer
func deleteDeployment(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sDeploymentBy(namespace, name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/deployment [patch]
// @Security Bearer
func patchDeployment(c *gin.Context) {
	var data v1Apps.Deployment
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sDeployment(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/deployment [post]
// @Security Bearer
func createDeployment(c *gin.Context) {
	var data v1Apps.Deployment
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sDeployment(data))
}

// ---------------------- SERVICES ----------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/service [get]
// @Security Bearer
// @Param namespace query string false  "namespace name"
func allServices(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sServices(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/service/describe/{namespace}/{name} [get]
// @Param namespace path string true  "namespace name"
// @Param name path string true  "service name"
// @Security Bearer
func describeService(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sService(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/service/{namespace}/{name [delete]
// @Security Bearer
// @Param namespace path string true  "namespace name"
// @Param name path string true  "service name"
func deleteService(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sServiceBy(namespace, name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/service [patch]
// @Security Bearer
func patchService(c *gin.Context) {
	var data v1.Service
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sService(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/service [post]
// @Security Bearer
func createService(c *gin.Context) {
	var data v1.Service
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sService(data))
}

// ---------------------- INGRESSES ----------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/ingress [get]
// @Security Bearer
// @Param namespace query string false  "namespace name"
func allIngresses(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sIngresses(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/ingress/describe/{namespace}/{name} [get]
// @Param namespace path string true  "namespace name"
// @Param name path string true  "ingress name"
// @Security Bearer
func describeIngress(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sIngress(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/ingress/{namespace}/{name} [delete]
// @Security Bearer
// @Param namespace path string true  "namespace name"
// @Param name path string true  "ingress name"
func deleteIngress(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sIngressBy(namespace, name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/ingress [patch]
// @Security Bearer
func patchIngress(c *gin.Context) {
	var data v1Networking.Ingress
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sIngress(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/ingress [post]
// @Security Bearer
func createIngress(c *gin.Context) {
	var data v1Networking.Ingress
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sIngress(data))
}

// ---------------------- CONFIGMAPS ----------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/configmap [get]
// @Security Bearer
// @Param namespace query string false  "namespace name"
func allConfigmaps(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sConfigmaps(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/configmap/describe/{namespace}/{name} [get]
// @Param namespace path string true "namespace"
// @Param name path string true "configmap name"
// @Security Bearer
func describeConfigmap(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sConfigmap(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/configmap/{namespace}/{name} [delete]
// @Security Bearer
// @Param namespace path string true "namespace"
// @Param name path string true "configmap name"
func deleteConfigmap(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sConfigmapBy(namespace, name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/configmap [patch]
// @Security Bearer
func patchConfigmap(c *gin.Context) {
	var data v1.ConfigMap
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sConfigMap(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/configmap [post]
// @Security Bearer
func createConfigmap(c *gin.Context) {
	var data v1.ConfigMap
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sConfigMap(data))
}

// ---------------------- SECRETS ----------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/secret [get]
// @Security Bearer
// @Param namespace query string false  "namespace name"
func allSecrets(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sSecrets(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/secret/describe/{namespace}/{name} [get]
// @Param namespace path string true  "namespace name"
// @Param name path string true  "secret name"
// @Security Bearer
func describeSecret(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sSecret(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/secret/{namespace}/{name} [delete]
// @Param namespace path string true  "namespace name"
// @Param name path string true  "secret name"
// @Security Bearer
func deleteSecret(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sSecretBy(namespace, name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/secret [patch]
// @Security Bearer
func patchSecret(c *gin.Context) {
	var data v1.Secret
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sSecret(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/secret [post]
// @Security Bearer
func createSecret(c *gin.Context) {
	var data v1.Secret
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sSecret(data))
}

// ---------------------- NODES ----------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/node [get]
// @Security Bearer
func allNodes(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.ListK8sNodes())
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/node/describe/{name} [get]
// @Param name path string true  "node name"
// @Security Bearer
func describeNode(c *gin.Context) {
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sNode(name))
}

// ---------------------- DEAMONSETS ----------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/daemon-set [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
func allDaemonSets(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sDaemonsets(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/daemon-set/describe/{namespace}/{name} [get]
// @Param namespace path string true "namespace"
// @Param namespace path string true "name"
// @Security Bearer
func describeDaemonSet(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sDaemonSet(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/daemon-set/{namespace}/{name} [delete]
// @Param namespace path string true "namespace"
// @Param namespace path string true "name"
// @Security Bearer
func deleteDaemonSet(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sDaemonSetBy(namespace, name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/daemon-set [patch]
// @Security Bearer
func patchDaemonSet(c *gin.Context) {
	var data v1Apps.DaemonSet
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sDaemonSet(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/daemon-set [post]
// @Security Bearer
func createDaemonSet(c *gin.Context) {
	var data v1Apps.DaemonSet
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sDaemonSet(data))
}

// ---------------------- STATEFULSETS ----------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/stateful-set [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
func allStatefulSets(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllStatefulSets(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/stateful-set/describe/{namespace}/{name} [get]
// @Param namespace path string true "namespace"
// @Param name path string true "stateful-set name"
// @Security Bearer
func describeStatefulSet(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sStatefulset(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/stateful-set/{namespace}/{name} [delete]
// @Param namespace path string true "namespace"
// @Param name path string true "stateful-set name"
// @Security Bearer
func deleteStatefulSet(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sStatefulsetBy(namespace, name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/stateful-set [patch]
// @Security Bearer
func patchStatefulSet(c *gin.Context) {
	var data v1Apps.StatefulSet
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sStatefulset(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/stateful-set [post]
// @Security Bearer
func createStatefulSet(c *gin.Context) {
	var data v1Apps.StatefulSet
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sStatefulset(data))
}

// ---------------------- JOBS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/job [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
func allJobs(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllJobs(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/job/describe/{namespace}/{name} [get]
// @Param namespace path string true "namespace"
// @Param name path string true "job name"
// @Security Bearer
func describeJob(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sJob(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/job/{namespace}/{name} [delete]
// @Param namespace path string true "namespace"
// @Param name path string true "job name"
// @Security Bearer
func deleteJob(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sJobBy(namespace, name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/job [patch]
// @Security Bearer
func patchJob(c *gin.Context) {
	var data v1Job.Job
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sJob(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/job [post]
// @Security Bearer
func createJob(c *gin.Context) {
	var data v1Job.Job
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sJob(data))
}

// ---------------------- CRONJOBS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/cron-job [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
func allCronJobs(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllCronjobs(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/cron-job/describe/{namespace}/{name} [get]
// @Param namespace path string true "namespace"
// @Param name path string true "cronjob name"
// @Security Bearer
func describeCronJob(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sCronJob(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/cron-job/{namespace}/{name} [delete]
// @Param namespace path string true "namespace"
// @Param name path string true "cronjob name"
// @Security Bearer
func deleteCronJob(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sCronJobBy(namespace, name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/cron-job [patch]
// @Security Bearer
func patchCronJob(c *gin.Context) {
	var data v1Job.CronJob
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sCronJob(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/cron-job [post]
// @Security Bearer
func createCronJob(c *gin.Context) {
	var data v1Job.CronJob
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sCronJob(data))
}

// ---------------------- REPLICASETS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/replica-set [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
func allReplicasets(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sReplicasets(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/replica-set/describe/{namespace}/{name} [get]
// @Param namespace path string true "namespace"
// @Param name path string true "replica-set name"
// @Security Bearer
func describeReplicaset(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sReplicaset(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/replica-set/{namespace}/{name} [delete]
// @Param namespace path string true "namespace"
// @Param name path string true "replica-set name"
// @Security Bearer
func deleteReplicaset(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sReplicasetBy(namespace, name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/replica-set [patch]
// @Security Bearer
func patchReplicaset(c *gin.Context) {
	var data v1Apps.ReplicaSet
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sReplicaset(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/replica-set [post]
// @Security Bearer
func createReplicaset(c *gin.Context) {
	var data v1Apps.ReplicaSet
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sReplicaSet(data))
}

// ---------------------- PERSISTENT VOLUMES ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/persistent-volume [get]
// @Security Bearer
func allPersistentVolumes(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllPersistentVolumes())
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/persistent-volume/describe/{name} [get]
// @Param name path string true "persistent-volume name"
// @Security Bearer
func describePersistentVolume(c *gin.Context) {
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sPersistentVolume(name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/persistent-volume/{name} [delete]
// @Param name path string true "persistent-volume name"
// @Security Bearer
func deletePersistentVolume(c *gin.Context) {
	name := c.Param("name")
	err := kubernetes.DeleteK8sPersistentVolumeBy(name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/persistent-volume [patch]
// @Security Bearer
func patchPersistentVolume(c *gin.Context) {
	var data v1.PersistentVolume
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sPersistentVolume(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/persistent-volume [post]
// @Security Bearer
func createPersistentVolume(c *gin.Context) {
	var data v1.PersistentVolume
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sPersistentVolume(data))
}

// ---------------------- PERSISTENT VOLUME CLAIMS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/persistent-volume-claim [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
func allPersistentVolumeClaims(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sPersistentVolumeClaims(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/persistent-volume-claim/describe/{namespace}/{name} [get]
// @Param namespace path string true "namespace"
// @Param name path string true "persistent-volume-claim name"
// @Security Bearer
func describePersistentVolumeClaim(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sPersistentVolumeClaim(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/persistent-volume-claim/{namespace}/{name} [delete]
// @Param namespace path string true "namespace"
// @Param name path string true "persistent-volume-claim name"
// @Security Bearer
func deletePersistentVolumeClaim(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sPersistentVolumeClaimBy(namespace, name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/persistent-volume-claim [patch]
// @Security Bearer
func patchPersistentVolumeClaim(c *gin.Context) {
	var data v1.PersistentVolumeClaim
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sPersistentVolumeClaim(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/persistent-volume-claim [post]
// @Security Bearer
func createPersistentVolumeClaim(c *gin.Context) {
	var data v1.PersistentVolumeClaim
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sPersistentVolumeClaim(data))
}

// ---------------------- HORIZONTAL POD AUTOSCALER ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/horizontal-pod-autoscaler [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
func allHpas(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllHpas(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/horizontal-pod-autoscaler/describe/{namespace}/{name} [get]
// @Param namespace path string true "namespace"
// @Param name path string true "hpa name"
// @Security Bearer
func describeHpa(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sHpa(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/horizontal-pod-autoscaler/{namespace}/{name} [delete]
// @Param namespace path string true "namespace"
// @Param name path string true "hpa name"
// @Security Bearer
func deleteHpa(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sHpaBy(namespace, name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/horizontal-pod-autoscaler [patch]
// @Security Bearer
func patchHpa(c *gin.Context) {
	var data v2Scale.HorizontalPodAutoscaler
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sHpa(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/horizontal-pod-autoscaler [post]
// @Security Bearer
func createHpa(c *gin.Context) {
	var data v2Scale.HorizontalPodAutoscaler
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sHpa(data))
}

// ---------------------- EVENTS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/event [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
func allEvents(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllEvents(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/event/describe/{namespace}/{name} [get]
// @Param namespace path string true "namespace"
// @Param name path string true "event name"
// @Security Bearer
func describeEvent(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sEvent(namespace, name))
}

// ---------------------- CERTIFICATES ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/certificate [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
func allCertificates(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sCertificates(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/certificate/describe/{namespace}/{name} [get]
// @Param namespace path string true "namespace"
// @Param name path string true "certificate name"
// @Security Bearer
func describeCertificate(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sCertificate(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/certificate/{namespace}/{name} [delete]
// @Param namespace path string true "namespace"
// @Param name path string true "certificate name"
// @Security Bearer
func deleteCertificate(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sCertificateBy(namespace, name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/certificate [patch]
// @Security Bearer
func patchCertificate(c *gin.Context) {
	var data cmapi.Certificate
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sCertificate(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/certificate [post]
// @Security Bearer
func createCertificate(c *gin.Context) {
	var data cmapi.Certificate
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sCertificate(data))
}

// ---------------------- CERTIFICATE REQUESTS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/certificate-request [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
func allCertificateRequests(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllCertificateSigningRequests(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/certificate-request/describe/{namespace}/{name} [get]
// @Param namespace path string true "namespace"
// @Param name path string true "certificate request name"
// @Security Bearer
func describeCertificateRequest(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sCertificateSigningRequest(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/certificate-request/{namespace}/{name} [delete]
// @Param namespace path string true "namespace name"
// @Param name path string true "certificate request name"
// @Security Bearer
func deleteCertificateRequest(c *gin.Context) {

	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sCertificateSigningRequestBy(namespace, name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/certificate-request [patch]
// @Security Bearer
func patchCertificateRequest(c *gin.Context) {
	var data cmapi.CertificateRequest
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sCertificateSigningRequest(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/certificate-request [post]
// @Security Bearer
func createCertificateRequest(c *gin.Context) {
	var data cmapi.CertificateRequest
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sCertificateSigningRequest(data))
}

// ---------------------- ORDERS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/orders [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
func allOrders(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllOrders(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/orders/describe/{namespace}/{name} [get]
// @Param namespace path string true "namespace"
// @Param name path string true "order name"
// @Security Bearer
func describeOrder(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sOrder(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/orders/{namespace}/{name} [delete]
// @Param namespace path string true "namespace name"
// @Param name path string true "order name"
// @Security Bearer
func deleteOrder(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sOrderBy(namespace, name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/orders [patch]
// @Security Bearer
func patchOrder(c *gin.Context) {
	var data v1Cert.Order
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sOrder(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/orders [post]
// @Security Bearer
func createOrder(c *gin.Context) {
	var data v1Cert.Order
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sOrder(data))
}

// ---------------------- ISSUERS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/issuer [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
func allIssuers(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllIssuer(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/issuer/describe/{namespace}/{name} [get]
// @Param namespace path string true "namespace"
// @Param name path string true "issuer name"
// @Security Bearer
func describeIssuer(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sIssuer(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/issuer/{namespace}/{name} [delete]
// @Param namespace path string true "namespace name"
// @Param name path string true "issuer name"
// @Security Bearer
func deleteIssuer(c *gin.Context) {

	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sIssuerBy(namespace, name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/issuer [patch]
// @Security Bearer
func patchIssuer(c *gin.Context) {
	var data cmapi.Issuer
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sIssuer(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/issuer [post]
// @Security Bearer
func createIssuer(c *gin.Context) {
	var data cmapi.Issuer
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sIssuer(data))
}

// ---------------------- CLUSTER ISSUERS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/cluster-issuer [get]
// @Security Bearer
func allClusterIssuers(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllClusterIssuers())
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/cluster-issuer/describe/{name} [get]
// @Param name path string true "cluster-issuer name"
// @Security Bearer
func describeClusterIssuer(c *gin.Context) {
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sClusterIssuer(name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/cluster-issuer/{name} [delete]
// @Param name path string true "cluster-issuer name"
// @Security Bearer
func deleteClusterIssuer(c *gin.Context) {
	name := c.Param("name")
	err := kubernetes.DeleteK8sClusterIssuerBy(name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/cluster-issuer [patch]
// @Security Bearer
func patchClusterIssuer(c *gin.Context) {
	var data cmapi.ClusterIssuer
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sClusterIssuer(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/cluster-issuer [post]
// @Security Bearer
func createClusterIssuer(c *gin.Context) {
	var data cmapi.ClusterIssuer
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sClusterIssuer(data))
}

// ---------------------- SERVICE ACCOUNTS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/service-account [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
func allServiceAccounts(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllServiceAccounts(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/service-account/describe/{namespace}/{name} [get]
// @Param namespace path string true "namespace"
// @Param name path string true "service-account name"
// @Security Bearer
func describeServiceAccount(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sServiceAccount(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/service-account/{namespace}/{name} [delete]
// @Param namespace path string true "namespace name"
// @Param name path string true "service-account name"
// @Security Bearer
func deleteServiceAccount(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sServiceAccountBy(namespace, name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/service-account [patch]
// @Security Bearer
func patchServiceAccount(c *gin.Context) {
	var data v1.ServiceAccount
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sServiceAccount(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/service-account [post]
// @Security Bearer
func createServiceAccount(c *gin.Context) {
	var data v1.ServiceAccount
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sServiceAccount(data))
}

// ---------------------- ROLES ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/role [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
func allRoles(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllRoles(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/role/describe/{namespace}/{name} [get]
// @Param namespace path string true "namespace"
// @Param name path string true "role name"
// @Security Bearer
func describeRole(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sRole(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/role/{namespace}/{name} [delete]
// @Param namespace path string true "namespace name"
// @Param name path string true "role name"
// @Security Bearer
func deleteRole(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sRoleBy(namespace, name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/role [patch]
// @Security Bearer
func patchRole(c *gin.Context) {
	var data v1Rbac.Role
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sRole(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/role [post]
// @Security Bearer
func createRole(c *gin.Context) {
	var data v1Rbac.Role
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sRole(data))
}

// ---------------------- ROLE BINDINGS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/role-binding [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
func allRoleBindings(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllRoleBindings(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/role-binding/describe/{namespace}/{name} [get]
// @Param namespace path string true "namespace"
// @Param name path string true "role-binding name"
// @Security Bearer
func describeRoleBinding(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sRoleBinding(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/role-binding/{namespace}/{name} [delete]
// @Param namespace path string true "namespace name"
// @Param name path string true "role-binding name"
// @Security Bearer
func deleteRoleBinding(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sRoleBindingBy(namespace, name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/role-binding [patch]
// @Security Bearer
func patchRoleBinding(c *gin.Context) {
	var data v1Rbac.RoleBinding
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sRoleBinding(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/role-binding [post]
// @Security Bearer
func createRoleBinding(c *gin.Context) {
	var data v1Rbac.RoleBinding
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sRoleBinding(data))
}

// ---------------------- CLUSTER ROLES ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/cluster-role [get]
// @Security Bearer
func allClusterRoles(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllClusterRoles())
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/cluster-role/describe/{name} [get]
// @Param name path string true "cluster-role name"
// @Security Bearer
func describeClusterRole(c *gin.Context) {
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sClusterRole(name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/cluster-role/{name} [delete]
// @Param name path string true "cluster-role name"
// @Security Bearer
func deleteClusterRole(c *gin.Context) {

	name := c.Param("name")
	err := kubernetes.DeleteK8sClusterRoleBy(name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/cluster-role [patch]
// @Security Bearer
func patchClusterRole(c *gin.Context) {
	var data v1Rbac.ClusterRole
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sClusterRole(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/cluster-role [post]
// @Security Bearer
func createClusterRole(c *gin.Context) {
	var data v1Rbac.ClusterRole
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sClusterRole(data))
}

// ---------------------- CLUSTER ROLE BINDINGS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/cluster-role-binding [get]
// @Security Bearer
func allClusterRoleBindings(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllClusterRoleBindings())
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/cluster-role-binding/describe/{name} [get]
// @Param name path string true "cluster-role-binding name"
// @Security Bearer
func describeClusterRoleBinding(c *gin.Context) {
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sClusterRoleBinding(name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/cluster-role-binding/{name} [delete]
// @Param name path string true "cluster-role-binding name"
// @Security Bearer
func deleteClusterRoleBinding(c *gin.Context) {
	name := c.Param("name")
	err := kubernetes.DeleteK8sClusterRoleBindingBy(name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/cluster-role-binding [patch]
// @Security Bearer
func patchClusterRoleBinding(c *gin.Context) {
	var data v1Rbac.ClusterRoleBinding
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sClusterRoleBinding(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/cluster-role-binding [post]
// @Security Bearer
func createClusterRoleBinding(c *gin.Context) {
	var data v1Rbac.ClusterRoleBinding
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sClusterRoleBinding(data))
}

// ---------------------- VOLUME ATTACHMENTS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/volume-attachment [get]
// @Security Bearer
func allVolumeAttachments(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllVolumeAttachments())
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/volume-attachment/describe/{name} [get]
// @Param name path string true "volume-attachment name"
// @Security Bearer
func describeVolumeAttachment(c *gin.Context) {
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sVolumeAttachment(name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/volume-attachment/{name} [delete]
// @Param name path string true "volume-attachment name"
// @Security Bearer
func deleteVolumeAttachment(c *gin.Context) {
	name := c.Param("name")
	err := kubernetes.DeleteK8sVolumeAttachmentBy(name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/volume-attachment [patch]
// @Security Bearer
func patchVolumeAttachment(c *gin.Context) {
	var data v1Storage.VolumeAttachment
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sVolumeAttachment(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/volume-attachment [post]
// @Security Bearer
func createVolumeAttachment(c *gin.Context) {
	var data v1Storage.VolumeAttachment
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sVolumeAttachment(data))
}

// ---------------------- NETWORK POLICIES ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/network-policy [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
func allNetworkPolicies(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllNetworkPolicies(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/network-policy/describe/{namespace}/{name} [get]
// @Param namespace path string true "namespace"
// @Param name path string true "network-policy name"
// @Security Bearer
func describeNetworkPolicy(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sNetworkPolicy(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/network-policy/{namespace}/{name} [delete]
// @Param namespace path string true "namespace name"
// @Param name path string true "network-policy name"
// @Security Bearer
func deleteNetworkPolicy(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sNetworkPolicyBy(namespace, name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/network-policy [patch]
// @Security Bearer
func patchNetworkPolicy(c *gin.Context) {
	var data v1Networking.NetworkPolicy
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sNetworkPolicy(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/network-policy [post]
// @Security Bearer
func createNetworkPolicy(c *gin.Context) {
	var data v1Networking.NetworkPolicy
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sNetworkpolicy(data))
}

// ---------------------- STORAGECLASSES ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/storage-class [get]
// @Security Bearer
func allStorageClasses(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllStorageClasses())
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/storage-class/describe/{name} [get]
// @Param name path string true "storage-class name"
// @Security Bearer
func describeStorageClass(c *gin.Context) {
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sStorageClass(name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/storage-class/{namespace}/{name} [delete]
// @Param name path string true "storage-class name"
// @Security Bearer
func deleteStorageClass(c *gin.Context) {
	name := c.Param("name")
	err := kubernetes.DeleteK8sStorageClassBy(name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/storage-class [patch]
// @Security Bearer
func patchStorageClass(c *gin.Context) {
	var data v1Storage.StorageClass
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sStorageClass(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/storage-class [post]
// @Security Bearer
func createStorageClass(c *gin.Context) {
	var data v1Storage.StorageClass
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sStorageClass(data))
}

// ---------------------- CUSTOM RESSOURCE DEFINITIONS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/crds [get]
// @Security Bearer
func allCrds(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllCustomResourceDefinitions())
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/crds/describe/{name} [get]
// @Param name path string true "crds name"
// @Security Bearer
func describeCrd(c *gin.Context) {
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sCustomResourceDefinition(name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/crds/{name} [delete]
// @Param name path string true "crds name"
// @Security Bearer
func deleteCrd(c *gin.Context) {
	// name := c.Param("name")
	// TODO
	// err := kubernetes.DeleteK8sCustomResourceDefinition(name)
	// if err != nil {
	// 	utils.MalformedMessage(c, err.Error())
	// 	return
	// }
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/crds [patch]
// @Security Bearer
func patchCrd(c *gin.Context) {
	var data apiExt.CustomResourceDefinition
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sCustomResourceDefinition(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/crds [post]
// @Security Bearer
func createCrd(c *gin.Context) {
	var data apiExt.CustomResourceDefinition
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sCustomResourceDefinition(data))
}

// ---------------------- ENDPOINTS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/endpoints [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
func allEndpoints(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllEndpoints(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/endpoints/describe/{namespace}/{name} [get]
// @Param namespace path string true "namespace"
// @Param name path string true "endpoint name"
// @Security Bearer
func describeEndpoint(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sEndpoint(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/endpoints/{namespace}/{name} [delete]
// @Param namespace path string true "namespace name"
// @Param name path string true "endpoints request name"
// @Security Bearer
func deleteEndpoint(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sEndpointBy(namespace, name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/endpoints [patch]
// @Security Bearer
func patchEndpoint(c *gin.Context) {
	var data v1.Endpoints
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sEndpoint(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/endpoints [post]
// @Security Bearer
func createEndpoint(c *gin.Context) {
	var data v1.Endpoints
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sEndpoint(data))
}

// ---------------------- LEASES ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/leases [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
func allLeases(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllLeases(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/leases/describe/{namespace}/{name} [get]
// @Param namespace path string true "namespace"
// @Param name path string true "lease name"
// @Security Bearer
func describeLease(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sLease(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/leases/{namespace}/{name} [delete]
// @Param namespace path string true "namespace name"
// @Param name path string true "lease name"
// @Security Bearer
func deleteLease(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sLeaseBy(namespace, name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/leases [patch]
// @Security Bearer
func patchLease(c *gin.Context) {
	var data v1Coordination.Lease
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sLease(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/leases [post]
// @Security Bearer
func createLease(c *gin.Context) {
	var data v1Coordination.Lease
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sLease(data))
}

// ---------------------- PRIORITY CLASSES ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/priority-classes [get]
// @Security Bearer
func allPriorityClasses(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllPriorityClasses())
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/priority-classes/describe/{name} [get]
// @Param name path string true "priority-classes name"
// @Security Bearer
func describePriorityClass(c *gin.Context) {
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sPriorityClass(name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/priority-classes/{namespace}/{name} [delete]
// @Param name path string true "priority-class name"
// @Security Bearer
func deletePriorityClass(c *gin.Context) {
	name := c.Param("name")
	err := kubernetes.DeleteK8sPriorityClassBy(name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/priority-classes [patch]
// @Security Bearer
func patchPriorityClass(c *gin.Context) {
	var data v1Scheduling.PriorityClass
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sPriorityClass(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/priority-classes [post]
// @Security Bearer
func createPriorityClass(c *gin.Context) {
	var data v1Scheduling.PriorityClass
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sPriorityClass(data))
}

// ---------------------- VOLUME SNAPSHOTS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/volume-snapshots [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
func allVolumeSnapshots(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllVolumeSnapshots(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/volume-snapshots/describe/{namespace}/{name} [get]
// @Param namespace path string true "namespace"
// @Param name path string true "volume-snapshot name"
// @Security Bearer
func describeVolumeSnapshot(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sVolumeSnapshot(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/volume-snapshots/{namespace}/{name} [delete]
// @Param namespace path string true "namespace name"
// @Param name path string true "volume-snapshots name"
// @Security Bearer
func deleteVolumeSnapshot(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sVolumeSnapshotBy(namespace, name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/volume-snapshots [patch]
// @Security Bearer
func patchVolumeSnapshot(c *gin.Context) {
	var data v6Snap.VolumeSnapshot
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sVolumeSnapshot(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/volume-snapshots [post]
// @Security Bearer
func createVolumeSnapshot(c *gin.Context) {
	var data v6Snap.VolumeSnapshot
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sVolumeSnapshot(data))
}

// ---------------------- RESOURCE QUOTAS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/resource-quota [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
func allResourceQuotas(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllResourceQuotas(namespace))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/resource-quota/describe/{namespace}/{name} [get]
// @Param namespace path string true "namespace"
// @Param name path string true "resource-quota name"
// @Security Bearer
func describeResourceQuota(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sResourceQuota(namespace, name))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /workload/resource-quota/{namespace}/{name} [delete]
// @Param namespace path string true "namespace name"
// @Param name path string true "resource-quota name"
// @Security Bearer
func deleteResourceQuota(c *gin.Context) {

	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sResourceQuotaBy(namespace, name)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/resource-quota [patch]
// @Security Bearer
func patchResourceQuota(c *gin.Context) {
	var data v1.ResourceQuota
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sResourceQuota(data))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /workload/resource-quota [post]
// @Security Bearer
func createResourceQuota(c *gin.Context) {
	var data v1.ResourceQuota
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sResourceQuota(data))
}
