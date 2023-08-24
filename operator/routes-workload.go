package operator

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/kubernetes"
	v1 "k8s.io/api/core/v1"
)

func InitWorkloadRoutes(router *gin.Engine) {
	router.GET("/workload/templates", Auth(dtos.USER), allWorkloadTemplates)
	router.GET("/workload/available-resources", Auth(dtos.READER), allKubernetesResources)

	router.GET("/workload/namespace/all", Auth(dtos.USER), allNamespaces)
	router.DELETE("/workload/namespace", Auth(dtos.USER), deleteNamespace)

	router.GET("/workload/pod/all", Auth(dtos.USER), allPods)
	router.GET("/workload/pod/describe", Auth(dtos.USER), describePod)
	router.DELETE("/workload/pod", Auth(dtos.USER), deletePod)
	router.PATCH("/workload/pod", Auth(dtos.USER), patchPod)
	router.POST("/workload/pod", Auth(dtos.USER), createPod)

	router.GET("/workload/deployment/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/service/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/ingress/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/configmap/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/secret/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/node/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/daemon_set/all", Auth(dtos.USER), NOT_IMPLEMENED)
	router.GET("/workload/stateful_set/all", Auth(dtos.USER), NOT_IMPLEMENED)
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
	c.JSON(http.StatusOK, kubernetes.AllK8sPods(namespace))
}
func describePod(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	c.JSON(http.StatusOK, kubernetes.DescribeK8sPod(namespace, name))
}
func deletePod(c *gin.Context) {
	var data v1.Pod
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, kubernetes.DeleteK8sPod(data))
}
func patchPod(c *gin.Context) {
	var data v1.Pod
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, kubernetes.UpdateK8sPod(data))
}
func createPod(c *gin.Context) {
	var data v1.Pod
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, kubernetes.CreateK8sPod(data))
}
