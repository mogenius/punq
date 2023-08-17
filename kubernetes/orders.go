package kubernetes

import (
	"context"
	"os/exec"

	"punq/logger"
	"punq/utils"

	v1 "github.com/cert-manager/cert-manager/pkg/apis/acme/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllOrders(namespaceName string) K8sWorkloadResult {
	result := []v1.Order{}

	provider := NewKubeProviderCertManager()
	orderList, err := provider.ClientSet.AcmeV1().Orders(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllCertificateSigningRequests ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, certificate := range orderList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, certificate.ObjectMeta.Namespace) {
			result = append(result, certificate)
		}
	}
	return WorkloadResult(result, nil)
}

func UpdateK8sOrder(data v1.Order) K8sWorkloadResult {
	kubeProvider := NewKubeProviderCertManager()
	orderClient := kubeProvider.ClientSet.AcmeV1().Orders(data.Namespace)
	_, err := orderClient.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sOrder(data v1.Order) K8sWorkloadResult {
	kubeProvider := NewKubeProviderCertManager()
	orderClient := kubeProvider.ClientSet.AcmeV1().Orders(data.Namespace)
	err := orderClient.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sOrder(namespace string, name string) K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "order", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func NewK8sOrder() K8sNewWorkload {
	return NewWorkload(
		RES_ORDER,
		utils.InitOrderYaml(),
		"An ORDER is referring to a Custom Resource Definition (CRD) or a resource from the Kubernetes extension cert-manager. Order is a resource used to represent an order with an ACME server (like Let's Encrypt) for a TLS certificate. Once an Order resource is created, cert-manager will attempt to fulfill the Order by obtaining a certificate.")
}
