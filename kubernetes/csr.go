package kubernetes

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	cmapi "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllCertificateSigningRequests(namespaceName string, contextId *string) utils.K8sWorkloadResult {
	result := []cmapi.CertificateRequest{}

	provider, err := NewKubeProviderCertManager(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	certificatesList, err := provider.ClientSet.CertmanagerV1().CertificateRequests(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllCertificateSigningRequests ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, certificate := range certificatesList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, certificate.ObjectMeta.Namespace) {
			result = append(result, certificate)
		}
	}
	return WorkloadResult(result, nil)
}

func GetCertificateSigningRequest(namespaceName string, name string, contextId *string) (*cmapi.CertificateRequest, error) {
	provider, err := NewKubeProviderCertManager(contextId)
	if err != nil {
		return nil, err
	}
	return provider.ClientSet.CertmanagerV1().CertificateRequests(namespaceName).Get(context.TODO(), name, metav1.GetOptions{})
}

func UpdateK8sCertificateSigningRequest(data cmapi.CertificateRequest, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProviderCertManager(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.CertmanagerV1().CertificateRequests(data.Namespace)
	res, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func DeleteK8sCertificateSigningRequest(data cmapi.CertificateRequest, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProviderCertManager(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.CertmanagerV1().CertificateRequests(data.Namespace)
	err = client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sCertificateSigningRequestBy(namespace string, name string, contextId *string) error {
	provider, err := NewKubeProviderCertManager(contextId)
	if err != nil {
		return err
	}
	client := provider.ClientSet.CertmanagerV1().CertificateRequests(namespace)
	return client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sCertificateSigningRequest(namespace string, name string, contextId *string) utils.K8sWorkloadResult {
	cmd := exec.Command("kubectl", fmt.Sprintf("describe -n %s csr%s", name, ContextFlag(contextId)))

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sCertificateSigningRequest(data cmapi.CertificateRequest, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProviderCertManager(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.CertmanagerV1().CertificateRequests(data.Namespace)
	res, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func NewK8sCertificateSigningRequest() K8sNewWorkload {
	return NewWorkload(
		RES_CERTIFICATE_REQUEST,
		utils.InitCertificateSigningRequestYaml(),
		"A CertificateSigningRequest is used to request a digital certificate based on a newly created or existing private key. In this example, a CSR named 'mycsr' is created with a specific certificate request (which should be your own, the one provided here is a placeholder), the standard Kubernetes signer, and for client auth usage. Please note that this example contains a placeholder for spec.request. You would need to replace this with the base64-encoded representation of your actual certificate signing request.")
}
