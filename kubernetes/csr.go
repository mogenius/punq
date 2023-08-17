package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	cmapi "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllCertificateSigningRequests(namespaceName string) K8sWorkloadResult {
	result := []cmapi.CertificateRequest{}

	provider := NewKubeProviderCertManager()
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

func UpdateK8sCertificateSigningRequest(data cmapi.CertificateRequest) K8sWorkloadResult {
	kubeProvider := NewKubeProviderCertManager()
	certificateClient := kubeProvider.ClientSet.CertmanagerV1().CertificateRequests(data.Namespace)
	_, err := certificateClient.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sCertificateSigningRequest(data cmapi.CertificateRequest) K8sWorkloadResult {
	kubeProvider := NewKubeProviderCertManager()
	certificateClient := kubeProvider.ClientSet.CertmanagerV1().CertificateRequests(data.Namespace)
	err := certificateClient.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sCertificateSigningRequest(name string) K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "csr", name)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func NewK8sCertificateSigningRequest() K8sNewWorkload {
	return NewWorkload(
		RES_CERTIFICATE_REQUEST,
		utils.InitCertificateSigningRequestYaml(),
		"A CertificateSigningRequest is used to request a digital certificate based on a newly created or existing private key. In this example, a CSR named 'mycsr' is created with a specific certificate request (which should be your own, the one provided here is a placeholder), the standard Kubernetes signer, and for client auth usage. Please note that this example contains a placeholder for spec.request. You would need to replace this with the base64-encoded representation of your actual certificate signing request.")
}
