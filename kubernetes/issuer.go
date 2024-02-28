package kubernetes

import (
	"context"
	"fmt"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	cmapi "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllIssuer(namespaceName string, contextId *string) utils.K8sWorkloadResult {
	result := []cmapi.Issuer{}

	provider, err := NewKubeProviderCertManager(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	issuersList, err := provider.ClientSet.CertmanagerV1().Issuers(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllIssuer ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, issuer := range issuersList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, issuer.ObjectMeta.Namespace) {
			result = append(result, issuer)
		}
	}
	return WorkloadResult(result, nil)
}

func GetIssuer(namespaceName string, name string, contextId *string) (*cmapi.Issuer, error) {
	provider, err := NewKubeProviderCertManager(contextId)
	if err != nil {
		return nil, err
	}
	return provider.ClientSet.CertmanagerV1().Issuers(namespaceName).Get(context.TODO(), name, metav1.GetOptions{})
}

func UpdateK8sIssuer(data cmapi.Issuer, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProviderCertManager(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.CertmanagerV1().Issuers(data.Namespace)
	res, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func DeleteK8sIssuer(data cmapi.Issuer, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProviderCertManager(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.CertmanagerV1().Issuers(data.Namespace)
	err = client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sIssuerBy(namespace string, name string, contextId *string) error {
	provider, err := NewKubeProviderCertManager(contextId)
	if err != nil {
		return err
	}
	client := provider.ClientSet.CertmanagerV1().Issuers(namespace)
	return client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sIssuer(namespace string, name string, contextId *string) utils.K8sWorkloadResult {
	cmd := utils.RunOnLocalShell(fmt.Sprintf("kubectl describe issuer %s -n %s%s", name, namespace, ContextFlag(contextId)))

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sIssuer(data cmapi.Issuer, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProviderCertManager(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.CertmanagerV1().Issuers(data.Namespace)
	res, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func NewK8sIssuer() K8sNewWorkload {
	return NewWorkload(
		RES_ISSUER,
		utils.InitIssuerYaml(),
		"An Issuer is a custom resource definition (CRD) in cert-manager, which is a native Kubernetes certificate management controller. It represents a certificate authority (CA) that can generate signed certificates. In this example, an Issuer named 'example-issuer' is created. This issuer uses the Let's Encrypt staging server for the ACME protocol. It will use the secret 'example-issuer-account-key' to store the ACME account's private key and uses HTTP-01 challenge for domain validation.")
}
