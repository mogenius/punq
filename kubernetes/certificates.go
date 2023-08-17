package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	cmapi "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllCertificates(namespaceName string) []cmapi.Certificate {
	result := []cmapi.Certificate{}

	provider := NewKubeProviderCertManager()
	certificatesList, err := provider.ClientSet.CertmanagerV1().Certificates(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllCertificates ERROR: %s", err.Error())
		return result
	}

	for _, certificate := range certificatesList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, certificate.ObjectMeta.Namespace) {
			result = append(result, certificate)
		}
	}
	return result
}

func GetCertificate(namespaceName string, resourceName string) (*cmapi.Certificate, error) {
	provider := NewKubeProviderCertManager()
	certificate, err := provider.ClientSet.CertmanagerV1().Certificates(namespaceName).Get(context.TODO(), resourceName, metav1.GetOptions{})
	if err != nil {
		logger.Log.Errorf("GetCertificate ERROR: %s", err.Error())
		return nil, err
	}
	return certificate, nil
}

func AllK8sCertificates(namespaceName string) K8sWorkloadResult {
	result := []cmapi.Certificate{}

	provider := NewKubeProviderCertManager()
	certificatesList, err := provider.ClientSet.CertmanagerV1().Certificates(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllCertificates ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, certificate := range certificatesList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, certificate.ObjectMeta.Namespace) {
			result = append(result, certificate)
		}
	}
	return WorkloadResult(result, nil)
}

func UpdateK8sCertificate(data cmapi.Certificate) K8sWorkloadResult {
	kubeProvider := NewKubeProviderCertManager()
	certificateClient := kubeProvider.ClientSet.CertmanagerV1().Certificates(data.Namespace)
	_, err := certificateClient.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sCertificate(data cmapi.Certificate) K8sWorkloadResult {
	kubeProvider := NewKubeProviderCertManager()
	certificateClient := kubeProvider.ClientSet.CertmanagerV1().Certificates(data.Namespace)
	err := certificateClient.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sCertificate(namespace string, name string) K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "certificate", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func NewK8sCertificate() K8sNewWorkload {
	return NewWorkload(
		RES_CERTIFICATE,
		utils.InitCertificateYaml(),
		"A Certificate resource in cert-manager is used to request, manage, and store TLS certificates from certificate authorities. In this example, a Certificate named 'my-certificate' is created. It requests a TLS certificate for the domain names 'example.com' and 'www.example.com'. The certificate will be issued and managed by the ClusterIssuer named 'my-cluster-issuer'. The resulting certificate will be stored in a Secret named 'my-certificate-secret'. Please note that this is a simplified example, and the actual configuration may vary depending on the specific certificate authority and issuer being used. Always refer to the documentation for the certificate manager you are using and follow the guidelines provided by the certificate authority.")
}
