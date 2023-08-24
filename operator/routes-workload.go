package operator

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/kubernetes"
	v1Apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	v1Networking "k8s.io/api/networking/v1"
)

func InitWorkloadRoutes(router *gin.Engine) {
	router.GET("/workload/templates", Auth(dtos.USER), allWorkloadTemplates)
	router.GET("/workload/available-resources", Auth(dtos.READER), allKubernetesResources)

	router.GET("/workload/namespace/all", Auth(dtos.USER), allNamespaces)
	router.DELETE("/workload/namespace", Auth(dtos.USER), deleteNamespace)

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

	router.GET("/workload/secret", Auth(dtos.USER), allSecrets)              // QUERY: namespace
	router.GET("/workload/secret/describe", Auth(dtos.USER), describeSecret) // QUERY: namespace, name
	router.DELETE("/workload/secret", Auth(dtos.USER), deleteSecret)         // BODY: json-object
	router.PATCH("/workload/secret", Auth(dtos.USER), patchSecret)           // BODY: json-object
	router.POST("/workload/secret", Auth(dtos.USER), createSecret)           // BODY: yaml-object

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

	router.GET("/workload/job/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/cron_job/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/replica_set/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/persistent_volume/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/persistent_volume_claim/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/horizontal_pod_autoscaler/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/event/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/certificate/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/certificaterequest/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/orders/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/issuer/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/clusterissuer/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/service_account/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/role/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/role_binding/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/cluster_role/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/cluster_role_binding/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/volume_attachment/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/network_policy/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/storageclass/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/crds/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/endpoints/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/leases/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/priorityclasses/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/volumesnapshots/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/resourcequota/all", Auth(dtos.USER), NOT_IMPLEMENED)
}

// GENERAL
func NOT_IMPLEMENED(c *gin.Context) {
	c.JSON(http.StatusTeapot, gin.H{"error": "NOT IMPLEMENTED YET! TODO"})
}

func allWorkloadTemplates(c *gin.Context) {
	c.JSON(http.StatusOK, kubernetes.ListCreateTemplates())
}

func allKubernetesResources(c *gin.Context) {
	c.JSON(http.StatusOK, kubernetes.ALL_RESOURCES)
}

// NAMESPACES
func allNamespaces(c *gin.Context) {
	c.JSON(http.StatusOK, kubernetes.ListAllNamespace())
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
