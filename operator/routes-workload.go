package operator

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/mogenius/punq/services"
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
		workloadRoutes.GET("/templates", Auth(dtos.READER), allWorkloadTemplates)
		workloadRoutes.GET("/available-resources", Auth(dtos.READER), allKubernetesResources)

		// namespace
		namespaceWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_NAMESPACE)), Auth(dtos.USER), RequireContextId())
		{
			namespaceWorkloadRoutes.GET("/", allNamespaces)                                                    // PARAM: -
			namespaceWorkloadRoutes.GET("/describe/:name", validateParam("name"), describeNamespaces)          // PARAM: name
			namespaceWorkloadRoutes.DELETE("/:name", Auth(dtos.ADMIN), validateParam("name"), deleteNamespace) // PARAM: name
			namespaceWorkloadRoutes.PATCH("/", patchNamespace)                                                 // BODY: json-object
			namespaceWorkloadRoutes.POST("/", createNamespace)                                                 // BODY: yaml-object
		}

		// pod
		podWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_POD)), Auth(dtos.USER), RequireContextId())
		{
			podWorkloadRoutes.GET("/", allPods)
			podWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describePod) // PARAM: namespace
			podWorkloadRoutes.GET("/logs/:namespace/:name", validateParam("namespace", "name"), logsPod)         // PARAM: namespace
			podWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deletePod)         // PARAM: namespace, name
			podWorkloadRoutes.PATCH("/", patchPod)                                                               // BODY: json-object
			podWorkloadRoutes.POST("/", createPod)                                                               // BODY: yaml-object
		}

		// deployment
		deploymentWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_DEPLOYMENT)), Auth(dtos.USER), RequireContextId())
		{
			deploymentWorkloadRoutes.GET("/", allDeployments)                                                                  // PARAM: namespace
			deploymentWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeDeployment) // PARAM: namespace, name
			deploymentWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteDeployment)         // PARAM: namespace, name
			deploymentWorkloadRoutes.PATCH("/", patchDeployment)                                                               // BODY: json-object
			deploymentWorkloadRoutes.POST("/´", createDeployment)                                                              // BODY: yaml-object
		}

		// service
		serviceWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_SERVICE)), Auth(dtos.USER), RequireContextId())
		{
			serviceWorkloadRoutes.GET("/", allServices)                                                                  // PARAM: namespace
			serviceWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeService) // PARAM: namespace, name
			serviceWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteService)         // PARAM: namespace, name
			serviceWorkloadRoutes.PATCH("/", patchService)                                                               // BODY: json-object
			serviceWorkloadRoutes.POST("/", createService)                                                               // BODY: yaml-object
		}

		// ingress
		ingressWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_INGRESS)), Auth(dtos.USER), RequireContextId())
		{
			ingressWorkloadRoutes.GET("/", allIngresses)                                                                 // PARAM: namespace
			ingressWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeIngress) // PARAM: namespace, name
			ingressWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteIngress)         // PARAM: json-object
			ingressWorkloadRoutes.PATCH("/", patchIngress)                                                               // BODY: json-object
			ingressWorkloadRoutes.POST("/", createIngress)                                                               // BODY: yaml-object
		}

		// configmap
		configmapWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_CONFIG_MAP)), Auth(dtos.USER), RequireContextId())
		{
			configmapWorkloadRoutes.GET("/", allConfigmaps)                                                                  // PARAM: namespace
			configmapWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeConfigmap) // PARAM: namespace, name
			configmapWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteConfigmap)         // PARAM: namespace, name
			configmapWorkloadRoutes.PATCH("/", patchConfigmap)                                                               // BODY: json-object
			workloadRoutes.POST("/´", createConfigmap)                                                                       // BODY: yaml-object
		}

		// secret
		secretWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_SECRET)), Auth(dtos.ADMIN), RequireContextId())
		{
			secretWorkloadRoutes.GET("/", allSecrets)                                                                  // PARAM: namespace
			secretWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeSecret) // PARAM: namespace, name
			secretWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteSecret)         // PARAM: namespace, name
			secretWorkloadRoutes.PATCH("/", patchSecret)                                                               // BODY: json-object
			secretWorkloadRoutes.POST("/", createSecret)                                                               // BODY: yaml-object
		}

		// node
		nodeWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_NODE)), Auth(dtos.USER), RequireContextId())
		{
			nodeWorkloadRoutes.GET("/", allNodes)                                          // -
			nodeWorkloadRoutes.GET("/describe/:name", validateParam("name"), describeNode) // PARAM: namespace
		}

		// daemon-set
		daemonSetWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_DAEMON_SET)), Auth(dtos.USER), RequireContextId())
		{
			daemonSetWorkloadRoutes.GET("/", allDaemonSets)                                                                  // PARAM: namespace
			daemonSetWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeDaemonSet) // PARAM: namespace, name
			daemonSetWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteDaemonSet)         // PARAM: namespace, name
			daemonSetWorkloadRoutes.PATCH("/", patchDaemonSet)                                                               // BODY: json-object
			daemonSetWorkloadRoutes.POST("/", createDaemonSet)                                                               // BODY: yaml-object

		}

		// stateful-set
		statefulSetWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_STATEFUL_SET)), Auth(dtos.USER), RequireContextId())
		{
			statefulSetWorkloadRoutes.GET("/", allStatefulSets)                                                                  // PARAM: namespace
			statefulSetWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeStatefulSet) // PARAM: namespace, name
			statefulSetWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteStatefulSet)         // PARAM: namespace, name
			statefulSetWorkloadRoutes.PATCH("/", patchStatefulSet)                                                               // BODY: json-object
			statefulSetWorkloadRoutes.POST("/", createStatefulSet)                                                               // BODY: yaml-object
		}

		// job
		jobWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_JOB)), Auth(dtos.USER), RequireContextId())
		{
			jobWorkloadRoutes.GET("/", allJobs)                                                                  // PARAM: namespace
			jobWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeJob) // PARAM: namespace, name
			jobWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteJob)         // PARAM: namespace, name
			jobWorkloadRoutes.PATCH("/", patchJob)                                                               // BODY: json-object
			jobWorkloadRoutes.POST("/", createJob)                                                               // BODY: yaml-object
		}

		// cron-job
		cronJobWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_CRON_JOB)), Auth(dtos.USER), RequireContextId())
		{
			cronJobWorkloadRoutes.GET("/", allCronJobs)                                                                  // PARAM: namespace
			cronJobWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeCronJob) // PARAM: namespace, name
			cronJobWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteCronJob)         // PARAM: namespace, name
			cronJobWorkloadRoutes.PATCH("/", patchCronJob)                                                               // BODY: json-object
			cronJobWorkloadRoutes.POST("/", createCronJob)                                                               // BODY: yaml-object
		}

		// replicaset
		replicaSetWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_REPLICA_SET)), Auth(dtos.USER), RequireContextId())
		{
			replicaSetWorkloadRoutes.GET("/", allReplicasets)                                                                  // PARAM: namespace
			replicaSetWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeReplicaset) // PARAM: namespace, name
			replicaSetWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteReplicaset)         // PARAM: namespace, name
			replicaSetWorkloadRoutes.PATCH("(", patchReplicaset)                                                               // BODY: json-object
			replicaSetWorkloadRoutes.POST("/", createReplicaset)                                                               // BODY: yaml-object
		}

		// persistent-volume
		persistentVolumeWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_PERSISTENT_VOLUME)), Auth(dtos.ADMIN), RequireContextId())
		{
			persistentVolumeWorkloadRoutes.GET("/", allPersistentVolumes)                                          // PARAM: -
			persistentVolumeWorkloadRoutes.GET("/describe/:name", validateParam("name"), describePersistentVolume) // PARAM: name
			persistentVolumeWorkloadRoutes.DELETE("/:name", validateParam("name"), deletePersistentVolume)         // PARAM: name
			persistentVolumeWorkloadRoutes.PATCH("/", patchPersistentVolume)                                       // BODY: json-object
			persistentVolumeWorkloadRoutes.POST("/", createPersistentVolume)                                       // BODY: yaml-object
		}

		// persistent-volume-claim
		persistentVolumeClaimWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_PERSISTENT_VOLUME_CLAIM)), RequireContextId())
		{
			persistentVolumeClaimWorkloadRoutes.GET("/", Auth(dtos.USER), allPersistentVolumeClaims)                                                                  // PARAM: namespace
			persistentVolumeClaimWorkloadRoutes.GET("/describe/:namespace/:name", Auth(dtos.USER), validateParam("namespace", "name"), describePersistentVolumeClaim) // PARAM: namespace, name
			persistentVolumeClaimWorkloadRoutes.DELETE("/:namespace/:name", Auth(dtos.ADMIN), validateParam("namespace", "name"), deletePersistentVolumeClaim)        // PARAM: namespace, name
			persistentVolumeClaimWorkloadRoutes.PATCH("/", Auth(dtos.ADMIN), patchPersistentVolumeClaim)                                                              // BODY: json-object
			persistentVolumeClaimWorkloadRoutes.POST("/", Auth(dtos.ADMIN), createPersistentVolumeClaim)                                                              // BODY: yaml-object
		}

		// horizontal-pod-autoscaler
		horizontalPodAutoscalerWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_HORIZONTAL_POD_AUTOSCALER)), RequireContextId())
		{
			horizontalPodAutoscalerWorkloadRoutes.GET("/", Auth(dtos.USER), allHpas)                                                                  // PARAM: namespace
			horizontalPodAutoscalerWorkloadRoutes.GET("/describe/:namespace/:name", Auth(dtos.USER), validateParam("namespace", "name"), describeHpa) // PARAM: namespace, name
			horizontalPodAutoscalerWorkloadRoutes.DELETE("/:namespace/:name", Auth(dtos.ADMIN), validateParam("namespace", "name"), deleteHpa)        // PARAM: namespace, name
			horizontalPodAutoscalerWorkloadRoutes.PATCH("/", Auth(dtos.ADMIN), patchHpa)                                                              // BODY: json-object
			horizontalPodAutoscalerWorkloadRoutes.POST("/", Auth(dtos.ADMIN), createHpa)                                                              // BODY: yaml-object
		}

		// event
		eventWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_EVENT)), Auth(dtos.USER), RequireContextId())
		{
			eventWorkloadRoutes.GET("/", allEvents)                                                                  // PARAM: namespace
			eventWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeEvent) // PARAM: namespace, name
		}

		// certificate
		certificateWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_CERTIFICATE)), Auth(dtos.USER), RequireContextId())
		{
			certificateWorkloadRoutes.GET("/", allCertificates)                                                                  // PARAM: namespace
			certificateWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeCertificate) // PARAM: namespace, name
			certificateWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteCertificate)         // PARAM: namespace, name
			certificateWorkloadRoutes.PATCH("/", patchCertificate)                                                               // BODY: json-object
			certificateWorkloadRoutes.POST("/", createCertificate)                                                               // BODY: yaml-object
		}

		// certificate-request
		certificateRequestWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_CERTIFICATE_REQUEST)), Auth(dtos.USER), RequireContextId())
		{
			certificateRequestWorkloadRoutes.GET("/", allCertificateRequests)                                                                  // PARAM: namespace
			certificateRequestWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeCertificateRequest) // PARAM: namespace, name
			certificateRequestWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteCertificateRequest)         // PARAM: namespace, name
			certificateRequestWorkloadRoutes.PATCH("/", patchCertificateRequest)                                                               // BODY: json-object
			certificateRequestWorkloadRoutes.POST("/", createCertificateRequest)                                                               // BODY: yaml-object
		}

		// orders
		ordersWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_ORDER)), Auth(dtos.USER), RequireContextId())
		{
			ordersWorkloadRoutes.GET("/", allOrders)                                                                  // PARAM: namespace
			ordersWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeOrder) // PARAM: namespace, name
			ordersWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteOrder)         // PARAM: namespace, name
			ordersWorkloadRoutes.PATCH("/", patchOrder)                                                               // BODY: json-object
			ordersWorkloadRoutes.POST("/", createOrder)                                                               // BODY: yaml-object
		}

		// issuer
		issuerWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_ISSUER)), Auth(dtos.USER), RequireContextId())
		{
			issuerWorkloadRoutes.GET("/", allIssuers)                                                                  // PARAM: namespace
			issuerWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeIssuer) // PARAM: namespace, name
			issuerWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteIssuer)         // PARAM: namespace, name
			issuerWorkloadRoutes.PATCH("/", patchIssuer)                                                               // BODY: json-object
			issuerWorkloadRoutes.POST("/", createIssuer)                                                               // BODY: yaml-object
		}

		// cluster-issuer
		clusterIssuerWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_CLUSTER_ISSUER)), Auth(dtos.ADMIN), RequireContextId())
		{
			clusterIssuerWorkloadRoutes.GET("/", allClusterIssuers)                                          // PARAM: -
			clusterIssuerWorkloadRoutes.GET("/describe/:name", validateParam("name"), describeClusterIssuer) // PARAM: name
			clusterIssuerWorkloadRoutes.DELETE("/:name", validateParam("name"), deleteClusterIssuer)         // PARAM: name
			clusterIssuerWorkloadRoutes.PATCH("/", patchClusterIssuer)                                       // BODY: json-object
			clusterIssuerWorkloadRoutes.POST("/", createClusterIssuer)                                       // BODY: yaml-object
		}

		// service-account
		serviceAccountWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_SERVICE_ACCOUNT)), Auth(dtos.ADMIN), RequireContextId())
		{
			serviceAccountWorkloadRoutes.GET("/", allServiceAccounts)                                                                  // PARAM: namespace
			serviceAccountWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeServiceAccount) // PARAM: namespace, name
			serviceAccountWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteServiceAccount)         // PARAM: namespace, name
			serviceAccountWorkloadRoutes.PATCH("/", patchServiceAccount)                                                               // BODY: json-object
			serviceAccountWorkloadRoutes.POST("/", createServiceAccount)                                                               // BODY: yaml-object
		}

		// role
		roleWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_ROLE)), RequireContextId())
		{
			roleWorkloadRoutes.GET("/", Auth(dtos.USER), allRoles)                                                                  // PARAM: namespace
			roleWorkloadRoutes.GET("/describe/:namespace/:name", Auth(dtos.USER), validateParam("namespace", "name"), describeRole) // PARAM: namespace, name
			roleWorkloadRoutes.DELETE("/:namespace/:name", Auth(dtos.ADMIN), validateParam("namespace", "name"), deleteRole)        // PARAM: namespace, name
			roleWorkloadRoutes.PATCH("/", Auth(dtos.ADMIN), patchRole)                                                              // BODY: json-object
			roleWorkloadRoutes.POST("/", Auth(dtos.ADMIN), createRole)                                                              // BODY: yaml-object
		}

		// role-binding
		roleBindingWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_ROLE_BINDING)), RequireContextId())
		{
			roleBindingWorkloadRoutes.GET("/", Auth(dtos.USER), allRoleBindings)                                                                  // PARAM: namespace
			roleBindingWorkloadRoutes.GET("/describe/:namespace/:name", Auth(dtos.USER), validateParam("namespace", "name"), describeRoleBinding) // PARAM: namespace, name
			roleBindingWorkloadRoutes.DELETE("/:namespace/:name", Auth(dtos.ADMIN), validateParam("namespace", "name"), deleteRoleBinding)        // PARAM: namespace, name
			roleBindingWorkloadRoutes.PATCH("/", Auth(dtos.ADMIN), patchRoleBinding)                                                              // BODY: json-object
			roleBindingWorkloadRoutes.POST("/", Auth(dtos.ADMIN), createRoleBinding)                                                              // BODY: yaml-object
		}

		// cluster-role
		clusterRoleWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_CLUSTER_ROLE)), Auth(dtos.ADMIN), RequireContextId())
		{
			clusterRoleWorkloadRoutes.GET("/", allClusterRoles)                                          // PARAM: -
			clusterRoleWorkloadRoutes.GET("/describe/:name", validateParam("name"), describeClusterRole) // PARAM: name
			clusterRoleWorkloadRoutes.DELETE("/:name", validateParam("name"), deleteClusterRole)         // PARAM: name
			clusterRoleWorkloadRoutes.PATCH("/cluster_role", patchClusterRole)                           // BODY: json-object
			clusterRoleWorkloadRoutes.POST("/cluster_role", createClusterRole)                           // BODY: yaml-object
		}

		// cluster-role-binding
		clusterRoleBindingWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_CLUSTER_ROLE_BINDING)), Auth(dtos.ADMIN), RequireContextId())
		{
			clusterRoleBindingWorkloadRoutes.GET("/", allClusterRoleBindings)                                          // PARAM: -
			clusterRoleBindingWorkloadRoutes.GET("/describe/:name", validateParam("name"), describeClusterRoleBinding) // PARAM: name
			clusterRoleBindingWorkloadRoutes.DELETE("/:name", validateParam("name"), deleteClusterRoleBinding)         // PARAM: name
			clusterRoleBindingWorkloadRoutes.PATCH("/", patchClusterRoleBinding)                                       // BODY: json-object
			clusterRoleBindingWorkloadRoutes.POST("/", createClusterRoleBinding)                                       // BODY: yaml-object
		}

		// volume-attachment
		volumeAttachmentWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_VOLUME_ATTACHMENT)), Auth(dtos.ADMIN), RequireContextId())
		{
			volumeAttachmentWorkloadRoutes.GET("/", allVolumeAttachments)                                          // PARAM: -
			volumeAttachmentWorkloadRoutes.GET("/describe/:name", validateParam("name"), describeVolumeAttachment) // PARAM: name
			volumeAttachmentWorkloadRoutes.DELETE("/:name", validateParam("name"), deleteVolumeAttachment)         // PARAM: name
			volumeAttachmentWorkloadRoutes.PATCH("/", patchVolumeAttachment)                                       // BODY: json-object
			volumeAttachmentWorkloadRoutes.POST("/", createVolumeAttachment)                                       // BODY: yaml-object
		}

		// network-policy
		networkPolicyWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_NETWORK_POLICY)), RequireContextId())
		{
			networkPolicyWorkloadRoutes.GET("/", Auth(dtos.USER), allNetworkPolicies)                                                                 // PARAM: namespace
			networkPolicyWorkloadRoutes.GET("/describe/:namespace/:name", Auth(dtos.USER), validateParam("namespace", "name"), describeNetworkPolicy) // PARAM: namespace, name
			networkPolicyWorkloadRoutes.DELETE("/:namespace/:name", Auth(dtos.ADMIN), validateParam("namespace", "name"), deleteNetworkPolicy)        // PARAM: namespace, name
			networkPolicyWorkloadRoutes.PATCH("/", Auth(dtos.ADMIN), patchNetworkPolicy)                                                              // BODY: json-object
			networkPolicyWorkloadRoutes.POST("/", Auth(dtos.ADMIN), createNetworkPolicy)                                                              // BODY: yaml-object
		}

		// storage-class
		storageClassWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_STORAGE_CLASS)), RequireContextId())
		{
			storageClassWorkloadRoutes.GET("/", Auth(dtos.USER), allStorageClasses)                                         // PARAM: namespace
			storageClassWorkloadRoutes.GET("/describe/:name", Auth(dtos.USER), validateParam("name"), describeStorageClass) // PARAM: namespace, name
			storageClassWorkloadRoutes.DELETE("/:name", Auth(dtos.ADMIN), validateParam("name"), deleteStorageClass)        // PARAM: namespace, name
			storageClassWorkloadRoutes.PATCH("/", Auth(dtos.ADMIN), patchStorageClass)                                      // BODY: json-object
			storageClassWorkloadRoutes.POST("/", Auth(dtos.ADMIN), createStorageClass)                                      // BODY: yaml-object
		}

		// crds
		crdsWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_CUSTOM_RESOURCE_DEFINITION)), Auth(dtos.ADMIN), RequireContextId())
		{
			crdsWorkloadRoutes.GET("/", allCrds)                                                            // PARAM: -
			crdsWorkloadRoutes.GET("/describe/:name", validateParam("name"), Auth(dtos.ADMIN), describeCrd) // PARAM: name
			crdsWorkloadRoutes.DELETE("/:name", validateParam("name"), Auth(dtos.ADMIN), deleteCrd)         // PARAM: name
			crdsWorkloadRoutes.PATCH("/", Auth(dtos.ADMIN), patchCrd)                                       // BODY: json-object
			crdsWorkloadRoutes.POST("/", Auth(dtos.ADMIN), createCrd)                                       // BODY: yaml-object
		}

		// endpoints
		endpointsWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_ENDPOINT)), Auth(dtos.USER), RequireContextId())
		{
			endpointsWorkloadRoutes.GET("/", allEndpoints)                                                                  // PARAM: namespace
			endpointsWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeEndpoint) // PARAM: namespace, name
			endpointsWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteEndpoint)         // PARAM: namespace, name
			endpointsWorkloadRoutes.PATCH("/", patchEndpoint)                                                               // BODY: json-object
			endpointsWorkloadRoutes.POST("/", createEndpoint)                                                               // BODY: yaml-object
		}

		// leases
		leasesWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_LEASE)), Auth(dtos.USER), RequireContextId())
		{
			leasesWorkloadRoutes.GET("/", allLeases)                                                                  // PARAM: namespace
			leasesWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeLease) // PARAM: namespace, name
			leasesWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteLease)         // PARAM: namespace, name
			leasesWorkloadRoutes.PATCH("/", patchLease)                                                               // BODY: json-object
			leasesWorkloadRoutes.POST("/", createLease)                                                               // BODY: yaml-object
		}

		// priority-classes
		priorityClassesWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_PRIORITY_CLASS)), Auth(dtos.ADMIN), RequireContextId())
		{
			priorityClassesWorkloadRoutes.GET("/", allPriorityClasses)                                         // PARAM: -
			priorityClassesWorkloadRoutes.GET("/describe/:name", validateParam("name"), describePriorityClass) // PARAM: name
			priorityClassesWorkloadRoutes.DELETE("/:name", validateParam("name"), deletePriorityClass)         // PARAM: name
			priorityClassesWorkloadRoutes.PATCH("/", patchPriorityClass)                                       // BODY: json-object
			priorityClassesWorkloadRoutes.POST("/", createPriorityClass)                                       // BODY: yaml-object
		}

		// volume-snapshots
		volumeSnapshotsWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_VOLUME_SNAPSHOT)), Auth(dtos.USER), RequireContextId())
		{
			volumeSnapshotsWorkloadRoutes.GET("/", allVolumeSnapshots)                                                                  // PARAM: namespace
			volumeSnapshotsWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeVolumeSnapshot) // PARAM: namespace, name
			volumeSnapshotsWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteVolumeSnapshot)         // PARAM: namespace, name
			volumeSnapshotsWorkloadRoutes.PATCH("/", patchVolumeSnapshot)                                                               // BODY: json-object
			volumeSnapshotsWorkloadRoutes.POST("/", createVolumeSnapshot)                                                               // BODY: yaml-object
		}

		// resource-quota
		resourceQuotaWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_RESOURCE_QUOTA)), Auth(dtos.ADMIN), RequireContextId())
		{
			resourceQuotaWorkloadRoutes.GET("/", allResourceQuotas)                                                                  // PARAM: namespace
			resourceQuotaWorkloadRoutes.GET("/describe/:namespace/:name", validateParam("namespace", "name"), describeResourceQuota) // PARAM: namespace, name
			resourceQuotaWorkloadRoutes.DELETE("/:namespace/:name", validateParam("namespace", "name"), deleteResourceQuota)         // PARAM: namespace, name
			resourceQuotaWorkloadRoutes.PATCH("/", patchResourceQuota)                                                               // BODY: json-object
			resourceQuotaWorkloadRoutes.POST("/", createResourceQuota)                                                               // BODY: yaml-object
		}

		// ingress-classes
		ingressClassesWorkloadRoutes := workloadRoutes.Group(fmt.Sprintf("/%s", strings.ToLower(kubernetes.RES_INGRESS_CLASS)), Auth(dtos.ADMIN), RequireContextId())
		{
			ingressClassesWorkloadRoutes.GET("/", allIngressClasses)                                         // PARAM: -
			ingressClassesWorkloadRoutes.GET("/describe/:name", validateParam("name"), describeIngressClass) // PARAM: name
			ingressClassesWorkloadRoutes.DELETE("/:name", validateParam("name"), deleteIngressClass)         // PARAM: name
			ingressClassesWorkloadRoutes.PATCH("/", patchIngressClass)                                       // BODY: json-object
			ingressClassesWorkloadRoutes.POST("/", createIngressClass)                                       // BODY: yaml-object
		}
	}
}

// @Tags General
// @Produce json
// @Success 200 {array} kubernetes.K8sNewWorkload
// @Router /backend/workload/templates/ [get]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allWorkloadTemplates(c *gin.Context) {
	c.JSON(http.StatusOK, kubernetes.ListCreateTemplates())
}

// @Tags General
// @Produce json
// @Success 200 {array} string
// @Router /backend/workload/available-resources/ [get]
// @Security Bearer
// @Param string header string true "X-Context-Id"
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
// @Router /backend/workload/namespace/ [get]
// @Param namespace query string false "name of the namespace"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allNamespaces(c *gin.Context) {
	c.JSON(http.StatusOK, kubernetes.ListK8sNamespaces("", services.GetGinContextId(c)))
}

// NAMESPACES
// @Tags Workloads
// @Produce json
// @Success 200 {array} v1.Namespace
// @Router /backend/workload/namespace/describe/{name}/ [get]
// @Param name path string true "name of the namespace"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeNamespaces(c *gin.Context) {
	name := c.Param("name")
	c.JSON(http.StatusOK, kubernetes.DescribeK8sNamespace(name, services.GetGinContextId(c)))
}

// NAMESPACES
// @Tags Workloads
// @Produce json
// @Success 201 {object} utils.K8sWorkloadResult
// @Router /backend/workload/namespace [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
// @Body {"type": "object"}
func createNamespace(c *gin.Context) {
	var data v1.Namespace
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		return
	}
	c.JSON(http.StatusCreated, kubernetes.CreateK8sNamespace(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/namespace [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchNamespace(c *gin.Context) {
	var data v1.Namespace
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sNamespace(data, services.GetGinContextId(c)))
}

// NAMESPACES
// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/namespace/{name} [delete]
// @Param name path string false  "name of the namespace"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteNamespace(c *gin.Context) {
	name := c.Param("name")
	err := kubernetes.DeleteK8sNamespaceBy(name, services.GetGinContextId(c))
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
// @Router /backend/workload/pod/ [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allPods(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sPods(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/pod/describe/{namespace}/{name}/ [get]
// @Param namespace path string true  "namespace name"
// @Param name path string true  "pod name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describePod(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sPod(namespace, name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce text/event-stream
// @Success 200 {string} string "streaming data"
// @Router /backend/workload/pod/logs/{namespace}/{name}/ [get]
// @Param namespace path string true  "namespace name"
// @Param name path string true  "pod name"
// @Param since-seconds query string false  "since-seconds"
// @Security Bearer
// @Param X-Context-Id header string true "X-Context-Id"
func logsPod(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	i, err := strconv.ParseInt(c.Query("since-seconds"), 10, 64)
	if err != nil {
		i = -1
	}

	req, err := kubernetes.StreamLog(namespace, name, i, services.GetGinContextId(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	stream, err := req.Stream(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	defer stream.Close()

	// set header
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	// get CloseNotify channel
	closeChannel := c.Writer.CloseNotify()

	go func() {
		<-closeChannel
	}()

	buf := make([]byte, 1024)
	c.Stream(func(w io.Writer) bool {
		n, err := stream.Read(buf)
		if err != nil {
			c.SSEvent("error", err.Error())
			return false
		}
		message := string(buf[:n])
		c.SSEvent("message", message)
		return true
	})
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/pod/{namespace}/{name} [delete]
// @Param namespace path string true "namespace name"
// @Param name path string true "pod name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deletePod(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sPodBy(namespace, name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/pod [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchPod(c *gin.Context) {
	var data v1.Pod
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sPod(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/pod [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createPod(c *gin.Context) {
	var data v1.Pod
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sPod(data, services.GetGinContextId(c)))

}

// ---------------------- DEPLOYMENTS ----------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/deployment/ [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allDeployments(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sDeployments(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/deployment/describe/{namespace}/{name}/ [get]
// @Param namespace path string true  "namespace name"
// @Param name path string true  "deployment name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeDeployment(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sDeployment(namespace, name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/deployment/{namespace}/{name} [delete]
// @Param namespace path string true  "namespace name"
// @Param name path string true  "deployment name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteDeployment(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sDeploymentBy(namespace, name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/deployment [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchDeployment(c *gin.Context) {
	var data v1Apps.Deployment
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sDeployment(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/deployment [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createDeployment(c *gin.Context) {
	var data v1Apps.Deployment
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sDeployment(data, services.GetGinContextId(c)))
}

// ---------------------- SERVICES ----------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/service/ [get]
// @Security Bearer
// @Param string header string true "X-Context-Id"
// @Param namespace query string false  "namespace name"
func allServices(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sServices(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/service/describe/{namespace}/{name}/ [get]
// @Param namespace path string true  "namespace name"
// @Param name path string true  "service name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeService(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sService(namespace, name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/service/{namespace}/{name [delete]
// @Security Bearer
// @Param string header string true "X-Context-Id"
// @Param namespace path string true  "namespace name"
// @Param name path string true  "service name"
func deleteService(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sServiceBy(namespace, name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/service [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchService(c *gin.Context) {
	var data v1.Service
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sService(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/service [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createService(c *gin.Context) {
	var data v1.Service
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sService(data, services.GetGinContextId(c)))
}

// ---------------------- INGRESSES ----------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/ingress/ [get]
// @Security Bearer
// @Param string header string true "X-Context-Id"
// @Param namespace query string false  "namespace name"
func allIngresses(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sIngresses(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/ingress/describe/{namespace}/{name}/ [get]
// @Param namespace path string true  "namespace name"
// @Param name path string true  "ingress name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeIngress(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sIngress(namespace, name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/ingress/{namespace}/{name} [delete]
// @Security Bearer
// @Param string header string true "X-Context-Id"
// @Param namespace path string true  "namespace name"
// @Param name path string true  "ingress name"
func deleteIngress(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sIngressBy(namespace, name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/ingress [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchIngress(c *gin.Context) {
	var data v1Networking.Ingress
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sIngress(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/ingress [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createIngress(c *gin.Context) {
	var data v1Networking.Ingress
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sIngress(data, services.GetGinContextId(c)))
}

// ---------------------- CONFIGMAPS ----------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/configmap/ [get]
// @Security Bearer
// @Param string header string true "X-Context-Id"
// @Param namespace query string false  "namespace name"
func allConfigmaps(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sConfigmaps(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/configmap/describe/{namespace}/{name}/ [get]
// @Param namespace path string true "namespace"
// @Param name path string true "configmap name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeConfigmap(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sConfigmap(namespace, name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/configmap/{namespace}/{name} [delete]
// @Security Bearer
// @Param string header string true "X-Context-Id"
// @Param namespace path string true "namespace"
// @Param name path string true "configmap name"
func deleteConfigmap(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sConfigmapBy(namespace, name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/configmap [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchConfigmap(c *gin.Context) {
	var data v1.ConfigMap
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sConfigMap(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/configmap [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createConfigmap(c *gin.Context) {
	var data v1.ConfigMap
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sConfigMap(data, services.GetGinContextId(c)))
}

// ---------------------- SECRETS ----------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/secret/ [get]
// @Security Bearer
// @Param string header string true "X-Context-Id"
// @Param namespace query string false  "namespace name"
func allSecrets(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sSecrets(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/secret/describe/{namespace}/{name}/ [get]
// @Param namespace path string true  "namespace name"
// @Param name path string true  "secret name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeSecret(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sSecret(namespace, name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/secret/{namespace}/{name} [delete]
// @Param namespace path string true  "namespace name"
// @Param name path string true  "secret name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteSecret(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sSecretBy(namespace, name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/secret [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchSecret(c *gin.Context) {
	var data v1.Secret
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sSecret(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/secret [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
// @Body {"type": "object"}
func createSecret(c *gin.Context) {
	var data v1.Secret
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sSecret(data, services.GetGinContextId(c)))
}

// ---------------------- NODES ----------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/node/ [get]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allNodes(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.ListK8sNodes(services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/node/describe/{name}/ [get]
// @Param name path string true  "node name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeNode(c *gin.Context) {
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sNode(name, services.GetGinContextId(c)))
}

// ---------------------- DEAMONSETS ----------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/daemon-set/ [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allDaemonSets(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sDaemonsets(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/daemon-set/describe/{namespace}/{name}/ [get]
// @Param namespace path string true "namespace"
// @Param namespace path string true "name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeDaemonSet(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sDaemonSet(namespace, name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/daemon-set/{namespace}/{name} [delete]
// @Param namespace path string true "namespace"
// @Param namespace path string true "name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteDaemonSet(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sDaemonSetBy(namespace, name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/daemon-set [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchDaemonSet(c *gin.Context) {
	var data v1Apps.DaemonSet
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sDaemonSet(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/daemon-set [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createDaemonSet(c *gin.Context) {
	var data v1Apps.DaemonSet
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sDaemonSet(data, services.GetGinContextId(c)))
}

// ---------------------- STATEFULSETS ----------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/stateful-set/ [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allStatefulSets(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sStatefulSets(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/stateful-set/describe/{namespace}/{name}/ [get]
// @Param namespace path string true "namespace"
// @Param name path string true "stateful-set name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeStatefulSet(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sStatefulset(namespace, name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/stateful-set/{namespace}/{name} [delete]
// @Param namespace path string true "namespace"
// @Param name path string true "stateful-set name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteStatefulSet(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sStatefulsetBy(namespace, name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/stateful-set [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchStatefulSet(c *gin.Context) {
	var data v1Apps.StatefulSet
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sStatefulset(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/stateful-set [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createStatefulSet(c *gin.Context) {
	var data v1Apps.StatefulSet
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sStatefulset(data, services.GetGinContextId(c)))
}

// ---------------------- JOBS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/job/ [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allJobs(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sJobs(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/job/describe/{namespace}/{name}/ [get]
// @Param namespace path string true "namespace"
// @Param name path string true "job name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeJob(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sJob(namespace, name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/job/{namespace}/{name} [delete]
// @Param namespace path string true "namespace"
// @Param name path string true "job name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteJob(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sJobBy(namespace, name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/job [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchJob(c *gin.Context) {
	var data v1Job.Job
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sJob(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/job [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createJob(c *gin.Context) {
	var data v1Job.Job
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sJob(data, services.GetGinContextId(c)))
}

// ---------------------- CRONJOBS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/cron-job/ [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allCronJobs(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sCronjobs(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/cron-job/describe/{namespace}/{name}/ [get]
// @Param namespace path string true "namespace"
// @Param name path string true "cronjob name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeCronJob(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sCronJob(namespace, name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/cron-job/{namespace}/{name} [delete]
// @Param namespace path string true "namespace"
// @Param name path string true "cronjob name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteCronJob(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sCronJobBy(namespace, name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/cron-job [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchCronJob(c *gin.Context) {
	var data v1Job.CronJob
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sCronJob(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/cron-job [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createCronJob(c *gin.Context) {
	var data v1Job.CronJob
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sCronJob(data, services.GetGinContextId(c)))
}

// ---------------------- REPLICASETS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/replica-set/ [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allReplicasets(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sReplicasets(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/replica-set/describe/{namespace}/{name}/ [get]
// @Param namespace path string true "namespace"
// @Param name path string true "replica-set name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeReplicaset(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sReplicaset(namespace, name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/replica-set/{namespace}/{name} [delete]
// @Param namespace path string true "namespace"
// @Param name path string true "replica-set name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteReplicaset(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sReplicasetBy(namespace, name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/replica-set [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchReplicaset(c *gin.Context) {
	var data v1Apps.ReplicaSet
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sReplicaset(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/replica-set [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createReplicaset(c *gin.Context) {
	var data v1Apps.ReplicaSet
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sReplicaSet(data, services.GetGinContextId(c)))
}

// ---------------------- PERSISTENT VOLUMES ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/persistent-volume/ [get]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allPersistentVolumes(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sPersistentVolumes(services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/persistent-volume/describe/{name}/ [get]
// @Param name path string true "persistent-volume name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describePersistentVolume(c *gin.Context) {
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sPersistentVolume(name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/persistent-volume/{name} [delete]
// @Param name path string true "persistent-volume name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deletePersistentVolume(c *gin.Context) {
	name := c.Param("name")
	err := kubernetes.DeleteK8sPersistentVolumeBy(name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/persistent-volume [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchPersistentVolume(c *gin.Context) {
	var data v1.PersistentVolume
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sPersistentVolume(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/persistent-volume [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createPersistentVolume(c *gin.Context) {
	var data v1.PersistentVolume
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sPersistentVolume(data, services.GetGinContextId(c)))
}

// ---------------------- PERSISTENT VOLUME CLAIMS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/persistent-volume-claim/ [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allPersistentVolumeClaims(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sPersistentVolumeClaims(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/persistent-volume-claim/describe/{namespace}/{name}/ [get]
// @Param namespace path string true "namespace"
// @Param name path string true "persistent-volume-claim name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describePersistentVolumeClaim(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sPersistentVolumeClaim(namespace, name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/persistent-volume-claim/{namespace}/{name} [delete]
// @Param namespace path string true "namespace"
// @Param name path string true "persistent-volume-claim name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deletePersistentVolumeClaim(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sPersistentVolumeClaimBy(namespace, name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/persistent-volume-claim [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchPersistentVolumeClaim(c *gin.Context) {
	var data v1.PersistentVolumeClaim
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sPersistentVolumeClaim(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/persistent-volume-claim [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createPersistentVolumeClaim(c *gin.Context) {
	var data v1.PersistentVolumeClaim
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sPersistentVolumeClaim(data, services.GetGinContextId(c)))
}

// ---------------------- HORIZONTAL POD AUTOSCALER ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/horizontal-pod-autoscaler/ [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allHpas(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sHpas(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/horizontal-pod-autoscaler/describe/{namespace}/{name}/ [get]
// @Param namespace path string true "namespace"
// @Param name path string true "hpa name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeHpa(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sHpa(namespace, name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/horizontal-pod-autoscaler/{namespace}/{name} [delete]
// @Param namespace path string true "namespace"
// @Param name path string true "hpa name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteHpa(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sHpaBy(namespace, name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/horizontal-pod-autoscaler [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchHpa(c *gin.Context) {
	var data v2Scale.HorizontalPodAutoscaler
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sHpa(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/horizontal-pod-autoscaler [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createHpa(c *gin.Context) {
	var data v2Scale.HorizontalPodAutoscaler
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sHpa(data, services.GetGinContextId(c)))
}

// ---------------------- EVENTS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/event/ [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allEvents(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sEvents(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/event/describe/{namespace}/{name}/ [get]
// @Param namespace path string true "namespace"
// @Param name path string true "event name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeEvent(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sEvent(namespace, name, services.GetGinContextId(c)))
}

// ---------------------- CERTIFICATES ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/certificate/ [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allCertificates(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sCertificates(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/certificate/describe/{namespace}/{name}/ [get]
// @Param namespace path string true "namespace"
// @Param name path string true "certificate name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeCertificate(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sCertificate(namespace, name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/certificate/{namespace}/{name} [delete]
// @Param namespace path string true "namespace"
// @Param name path string true "certificate name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteCertificate(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sCertificateBy(namespace, name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/certificate [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchCertificate(c *gin.Context) {
	var data cmapi.Certificate
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sCertificate(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/certificate [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createCertificate(c *gin.Context) {
	var data cmapi.Certificate
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sCertificate(data, services.GetGinContextId(c)))
}

// ---------------------- CERTIFICATE REQUESTS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/certificate-request/ [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allCertificateRequests(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sCertificateSigningRequests(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/certificate-request/describe/{namespace}/{name}/ [get]
// @Param namespace path string true "namespace"
// @Param name path string true "certificate request name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeCertificateRequest(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sCertificateSigningRequest(namespace, name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/certificate-request/{namespace}/{name} [delete]
// @Param namespace path string true "namespace name"
// @Param name path string true "certificate request name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteCertificateRequest(c *gin.Context) {

	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sCertificateSigningRequestBy(namespace, name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/certificate-request [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchCertificateRequest(c *gin.Context) {
	var data cmapi.CertificateRequest
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sCertificateSigningRequest(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/certificate-request [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createCertificateRequest(c *gin.Context) {
	var data cmapi.CertificateRequest
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sCertificateSigningRequest(data, services.GetGinContextId(c)))
}

// ---------------------- ORDERS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/orders/ [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allOrders(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sOrders(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/orders/describe/{namespace}/{name}/ [get]
// @Param namespace path string true "namespace"
// @Param name path string true "order name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeOrder(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sOrder(namespace, name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/orders/{namespace}/{name} [delete]
// @Param namespace path string true "namespace name"
// @Param name path string true "order name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteOrder(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sOrderBy(namespace, name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/orders [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchOrder(c *gin.Context) {
	var data v1Cert.Order
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sOrder(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/orders [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createOrder(c *gin.Context) {
	var data v1Cert.Order
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sOrder(data, services.GetGinContextId(c)))
}

// ---------------------- ISSUERS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/issuer/ [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allIssuers(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sIssuer(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/issuer/describe/{namespace}/{name}/ [get]
// @Param namespace path string true "namespace"
// @Param name path string true "issuer name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeIssuer(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sIssuer(namespace, name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/issuer/{namespace}/{name} [delete]
// @Param namespace path string true "namespace name"
// @Param name path string true "issuer name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteIssuer(c *gin.Context) {

	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sIssuerBy(namespace, name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/issuer [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchIssuer(c *gin.Context) {
	var data cmapi.Issuer
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sIssuer(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/issuer [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createIssuer(c *gin.Context) {
	var data cmapi.Issuer
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sIssuer(data, services.GetGinContextId(c)))
}

// ---------------------- CLUSTER ISSUERS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/cluster-issuer/ [get]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allClusterIssuers(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sClusterIssuers(services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/cluster-issuer/describe/{name}/ [get]
// @Param name path string true "cluster-issuer name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeClusterIssuer(c *gin.Context) {
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sClusterIssuer(name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/cluster-issuer/{name} [delete]
// @Param name path string true "cluster-issuer name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteClusterIssuer(c *gin.Context) {
	name := c.Param("name")
	err := kubernetes.DeleteK8sClusterIssuerBy(name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/cluster-issuer [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchClusterIssuer(c *gin.Context) {
	var data cmapi.ClusterIssuer
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sClusterIssuer(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/cluster-issuer [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createClusterIssuer(c *gin.Context) {
	var data cmapi.ClusterIssuer
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sClusterIssuer(data, services.GetGinContextId(c)))
}

// ---------------------- SERVICE ACCOUNTS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/service-account/ [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allServiceAccounts(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sServiceAccounts(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/service-account/describe/{namespace}/{name}/ [get]
// @Param namespace path string true "namespace"
// @Param name path string true "service-account name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeServiceAccount(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sServiceAccount(namespace, name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/service-account/{namespace}/{name} [delete]
// @Param namespace path string true "namespace name"
// @Param name path string true "service-account name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteServiceAccount(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sServiceAccountBy(namespace, name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/service-account [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchServiceAccount(c *gin.Context) {
	var data v1.ServiceAccount
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sServiceAccount(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/service-account [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createServiceAccount(c *gin.Context) {
	var data v1.ServiceAccount
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sServiceAccount(data, services.GetGinContextId(c)))
}

// ---------------------- ROLES ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/role/ [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allRoles(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sRoles(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/role/describe/{namespace}/{name}/ [get]
// @Param namespace path string true "namespace"
// @Param name path string true "role name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeRole(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sRole(namespace, name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/role/{namespace}/{name} [delete]
// @Param namespace path string true "namespace name"
// @Param name path string true "role name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteRole(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sRoleBy(namespace, name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/role [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchRole(c *gin.Context) {
	var data v1Rbac.Role
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sRole(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/role [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createRole(c *gin.Context) {
	var data v1Rbac.Role
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sRole(data, services.GetGinContextId(c)))
}

// ---------------------- ROLE BINDINGS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/role-binding/ [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allRoleBindings(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sRoleBindings(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/role-binding/describe/{namespace}/{name}/ [get]
// @Param namespace path string true "namespace"
// @Param name path string true "role-binding name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeRoleBinding(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sRoleBinding(namespace, name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/role-binding/{namespace}/{name} [delete]
// @Param namespace path string true "namespace name"
// @Param name path string true "role-binding name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteRoleBinding(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sRoleBindingBy(namespace, name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/role-binding [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchRoleBinding(c *gin.Context) {
	var data v1Rbac.RoleBinding
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sRoleBinding(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/role-binding [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createRoleBinding(c *gin.Context) {
	var data v1Rbac.RoleBinding
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sRoleBinding(data, services.GetGinContextId(c)))
}

// ---------------------- CLUSTER ROLES ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/cluster-role/ [get]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allClusterRoles(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sClusterRoles(services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/cluster-role/describe/{name}/ [get]
// @Param name path string true "cluster-role name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeClusterRole(c *gin.Context) {
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sClusterRole(name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/cluster-role/{name} [delete]
// @Param name path string true "cluster-role name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteClusterRole(c *gin.Context) {

	name := c.Param("name")
	err := kubernetes.DeleteK8sClusterRoleBy(name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/cluster-role [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchClusterRole(c *gin.Context) {
	var data v1Rbac.ClusterRole
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sClusterRole(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/cluster-role [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createClusterRole(c *gin.Context) {
	var data v1Rbac.ClusterRole
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sClusterRole(data, services.GetGinContextId(c)))
}

// ---------------------- CLUSTER ROLE BINDINGS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/cluster-role-binding/ [get]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allClusterRoleBindings(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sClusterRoleBindings(services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/cluster-role-binding/describe/{name}/ [get]
// @Param name path string true "cluster-role-binding name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeClusterRoleBinding(c *gin.Context) {
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sClusterRoleBinding(name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/cluster-role-binding/{name} [delete]
// @Param name path string true "cluster-role-binding name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteClusterRoleBinding(c *gin.Context) {
	name := c.Param("name")
	err := kubernetes.DeleteK8sClusterRoleBindingBy(name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/cluster-role-binding [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchClusterRoleBinding(c *gin.Context) {
	var data v1Rbac.ClusterRoleBinding
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sClusterRoleBinding(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/cluster-role-binding [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createClusterRoleBinding(c *gin.Context) {
	var data v1Rbac.ClusterRoleBinding
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sClusterRoleBinding(data, services.GetGinContextId(c)))
}

// ---------------------- VOLUME ATTACHMENTS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/volume-attachment/ [get]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allVolumeAttachments(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sVolumeAttachments(services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/volume-attachment/describe/{name}/ [get]
// @Param name path string true "volume-attachment name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeVolumeAttachment(c *gin.Context) {
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sVolumeAttachment(name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/volume-attachment/{name} [delete]
// @Param name path string true "volume-attachment name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteVolumeAttachment(c *gin.Context) {
	name := c.Param("name")
	err := kubernetes.DeleteK8sVolumeAttachmentBy(name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/volume-attachment [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchVolumeAttachment(c *gin.Context) {
	var data v1Storage.VolumeAttachment
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sVolumeAttachment(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/volume-attachment [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createVolumeAttachment(c *gin.Context) {
	var data v1Storage.VolumeAttachment
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sVolumeAttachment(data, services.GetGinContextId(c)))
}

// ---------------------- NETWORK POLICIES ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/network-policy/ [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allNetworkPolicies(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sNetworkPolicies(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/network-policy/describe/{namespace}/{name}/ [get]
// @Param namespace path string true "namespace"
// @Param name path string true "network-policy name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeNetworkPolicy(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sNetworkPolicy(namespace, name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/network-policy/{namespace}/{name} [delete]
// @Param namespace path string true "namespace name"
// @Param name path string true "network-policy name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteNetworkPolicy(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sNetworkPolicyBy(namespace, name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/network-policy [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchNetworkPolicy(c *gin.Context) {
	var data v1Networking.NetworkPolicy
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sNetworkPolicy(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/network-policy [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createNetworkPolicy(c *gin.Context) {
	var data v1Networking.NetworkPolicy
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sNetworkpolicy(data, services.GetGinContextId(c)))
}

// ---------------------- STORAGECLASSES ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/storage-class/ [get]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allStorageClasses(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sStorageClasses(services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/storage-class/describe/{name}/ [get]
// @Param name path string true "storage-class name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeStorageClass(c *gin.Context) {
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sStorageClass(name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/storage-class/{namespace}/{name} [delete]
// @Param name path string true "storage-class name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteStorageClass(c *gin.Context) {
	name := c.Param("name")
	err := kubernetes.DeleteK8sStorageClassBy(name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/storage-class [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchStorageClass(c *gin.Context) {
	var data v1Storage.StorageClass
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sStorageClass(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/storage-class [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createStorageClass(c *gin.Context) {
	var data v1Storage.StorageClass
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sStorageClass(data, services.GetGinContextId(c)))
}

// ---------------------- CUSTOM RESSOURCE DEFINITIONS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/crds/ [get]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allCrds(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllCustomResourceDefinitions())
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/crds/describe/{name}/ [get]
// @Param name path string true "crds name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeCrd(c *gin.Context) {
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sCustomResourceDefinition(name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/crds/{name} [delete]
// @Param name path string true "crds name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
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
// @Router /backend/workload/crds [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
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
// @Router /backend/workload/crds [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
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
// @Router /backend/workload/endpoints/ [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allEndpoints(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sEndpoints(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/endpoints/describe/{namespace}/{name}/ [get]
// @Param namespace path string true "namespace"
// @Param name path string true "endpoint name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeEndpoint(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sEndpoint(namespace, name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/endpoints/{namespace}/{name} [delete]
// @Param namespace path string true "namespace name"
// @Param name path string true "endpoints request name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteEndpoint(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sEndpointBy(namespace, name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/endpoints [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchEndpoint(c *gin.Context) {
	var data v1.Endpoints
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sEndpoint(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/endpoints [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createEndpoint(c *gin.Context) {
	var data v1.Endpoints
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sEndpoint(data, services.GetGinContextId(c)))
}

// ---------------------- LEASES ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/leases/ [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allLeases(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sLeases(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/leases/describe/{namespace}/{name}/ [get]
// @Param namespace path string true "namespace"
// @Param name path string true "lease name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeLease(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sLease(namespace, name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/leases/{namespace}/{name} [delete]
// @Param namespace path string true "namespace name"
// @Param name path string true "lease name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteLease(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sLeaseBy(namespace, name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/leases [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchLease(c *gin.Context) {
	var data v1Coordination.Lease
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sLease(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/leases [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createLease(c *gin.Context) {
	var data v1Coordination.Lease
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sLease(data, services.GetGinContextId(c)))
}

// ---------------------- PRIORITY CLASSES ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/priority-classes/ [get]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allPriorityClasses(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sPriorityClasses(services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/priority-classes/describe/{name}/ [get]
// @Param name path string true "priority-classes name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describePriorityClass(c *gin.Context) {
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sPriorityClass(name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/priority-classes/{namespace}/{name} [delete]
// @Param name path string true "priority-class name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deletePriorityClass(c *gin.Context) {
	name := c.Param("name")
	err := kubernetes.DeleteK8sPriorityClassBy(name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/priority-classes [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchPriorityClass(c *gin.Context) {
	var data v1Scheduling.PriorityClass
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sPriorityClass(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/priority-classes [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createPriorityClass(c *gin.Context) {
	var data v1Scheduling.PriorityClass
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sPriorityClass(data, services.GetGinContextId(c)))
}

// ---------------------- VOLUME SNAPSHOTS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/volume-snapshots/ [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allVolumeSnapshots(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sVolumeSnapshots(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/volume-snapshots/describe/{namespace}/{name}/ [get]
// @Param namespace path string true "namespace"
// @Param name path string true "volume-snapshot name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeVolumeSnapshot(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sVolumeSnapshot(namespace, name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/volume-snapshots/{namespace}/{name} [delete]
// @Param namespace path string true "namespace name"
// @Param name path string true "volume-snapshots name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteVolumeSnapshot(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sVolumeSnapshotBy(namespace, name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/volume-snapshots [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
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
// @Router /backend/workload/volume-snapshots [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createVolumeSnapshot(c *gin.Context) {
	var data v6Snap.VolumeSnapshot
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sVolumeSnapshot(data, services.GetGinContextId(c)))
}

// ---------------------- RESOURCE QUOTAS ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/resource-quota/ [get]
// @Param namespace query string false "namespace name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allResourceQuotas(c *gin.Context) {
	namespace := c.Query("namespace")
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sResourceQuotas(namespace, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/resource-quota/describe/{namespace}/{name}/ [get]
// @Param namespace path string true "namespace"
// @Param name path string true "resource-quota name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeResourceQuota(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sResourceQuota(namespace, name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/resource-quota/{namespace}/{name} [delete]
// @Param namespace path string true "namespace name"
// @Param name path string true "resource-quota name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteResourceQuota(c *gin.Context) {

	namespace := c.Param("namespace")
	name := c.Param("name")
	err := kubernetes.DeleteK8sResourceQuotaBy(namespace, name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/resource-quota [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchResourceQuota(c *gin.Context) {
	var data v1.ResourceQuota
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sResourceQuota(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/resource-quota [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createResourceQuota(c *gin.Context) {
	var data v1.ResourceQuota
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sResourceQuota(data, services.GetGinContextId(c)))
}

// ---------------------- INGRESS CLASSES ----------------------------

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/ingress-class/ [get]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func allIngressClasses(c *gin.Context) {
	utils.HttpRespondForWorkloadResult(c, kubernetes.AllK8sIngressClasses(services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/ingress-classes/describe/{name}/ [get]
// @Param name path string true "ingress-class name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func describeIngressClass(c *gin.Context) {
	name := c.Param("name")
	utils.HttpRespondForWorkloadResult(c, kubernetes.DescribeK8sIngressClass(name, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200
// @Router /backend/workload/ingress-class/{name} [delete]
// @Param name path string true "ingress-class name"
// @Security Bearer
// @Param string header string true "X-Context-Id"
func deleteIngressClass(c *gin.Context) {
	name := c.Param("name")
	err := kubernetes.DeleteK8sIngressClassBy(name, services.GetGinContextId(c))
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/ingress-class [patch]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func patchIngressClass(c *gin.Context) {
	var data v1Networking.IngressClass
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.UpdateK8sIngressClass(data, services.GetGinContextId(c)))
}

// @Tags Workloads
// @Produce json
// @Success 200 {object} utils.K8sWorkloadResult
// @Router /backend/workload/ingress-class [post]
// @Security Bearer
// @Param string header string true "X-Context-Id"
func createIngressClass(c *gin.Context) {
	var data v1Networking.IngressClass
	err := c.MustBindWith(&data, binding.YAML)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.HttpRespondForWorkloadResult(c, kubernetes.CreateK8sIngressClass(data, services.GetGinContextId(c)))
}
