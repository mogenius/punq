package kubernetes

import (
	"fmt"
	"os/exec"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"
	apiExt "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

// TODO

func AllCustomResourceDefinitions() utils.K8sWorkloadResult {
	// TODO

	// 	result := []apiExt.CustomResourceDefinition{}

	// 	provider,err := NewKubeProvider()
	// 	certificatesList, err := provider.ClientSet.ApiextensionsV1()
	// 	if err != nil {
	// 		logger.Log.Errorf("AllCertificateSigningRequests ERROR: %s", err.Error())
	// 		return WorkloadResult(nil, err)
	// 	}

	// 	for _, certificate := range certificatesList.Items {
	// 		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, certificate.ObjectMeta.Namespace) {
	// 			result = append(result, certificate)
	// 		}
	// 	}
	return WorkloadResult(nil, nil)
}

func UpdateK8sCustomResourceDefinition(data apiExt.CustomResourceDefinition) utils.K8sWorkloadResult {
	// TODO

	// providererr := NewKubeProvider()
	// certificateClient := provider.ClientSet.Ex.CertificateRequests(data.Namespace)
	// _, err := certificateClient.Update(context.TODO(), &data, metav1.UpdateOptions{})
	// if err != nil {
	// 	return WorkloadResult(nil, err)
	// }
	return WorkloadResult(nil, nil)
}

func DeleteK8sCustomResourceDefinition(data apiExt.CustomResourceDefinition) utils.K8sWorkloadResult {
	// TODO

	// providererr := NewKubeProvider()
	// certificateClient := provider.ClientSet.
	// err := certificateClient.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	// if err != nil {
	// 	return WorkloadResult(nil, err)
	// }
	return WorkloadResult(nil, nil)
}

func DescribeK8sCustomResourceDefinition(name string, contextId *string) utils.K8sWorkloadResult {
	cmd := exec.Command("kubectl", fmt.Sprintf("describe crds %s%s", name, ContextFlag(contextId)))

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sCustomResourceDefinition(data apiExt.CustomResourceDefinition) utils.K8sWorkloadResult {
	// TODO

	// providererr := NewKubeProvider()
	// client := provider.ClientSet.CoreV1().ConfigMaps(data.Namespace)
	// _, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	// if err != nil {
	// 	return WorkloadResult(nil, err)
	// }
	return WorkloadResult(nil, nil)
}

func NewK8sCustomResourceDefinition() K8sNewWorkload {
	return NewWorkload(
		RES_CUSTOM_RESOURCE_DEFINITION,
		utils.InitCustomResourceDefinitionYaml(),
		"A CustomResourceDefinition (CRD) in Kubernetes is a mechanism for defining and using custom resources. It extends the Kubernetes API to create new types of resources that can be managed just like the built-in resources. This CRD defines a new resource type CronTab in the group stable.example.com. The version is v1 which is being served and used for storage. The CronTab resource type has a specification that requires three fields: cronSpec, replicas, and image. The cronSpec field is a string that specifies the cron syntax for schedule, replicas is an integer that specifies the number of replicas, and image is the string that specifies the container image to run.")
}
