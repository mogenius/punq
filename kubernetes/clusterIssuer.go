package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	cmapi "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllClusterIssuers() utils.HttpResult {
	result := []cmapi.ClusterIssuer{}

	provider := NewKubeProviderCertManager()
	issuersList, err := provider.ClientSet.CertmanagerV1().ClusterIssuers().List(context.TODO(), metav1.ListOptions{})
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

func UpdateK8sClusterIssuer(data cmapi.ClusterIssuer) utils.HttpResult {
	kubeProvider := NewKubeProviderCertManager()
	client := kubeProvider.ClientSet.CertmanagerV1().ClusterIssuers()
	_, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sClusterIssuer(data cmapi.ClusterIssuer) utils.HttpResult {
	kubeProvider := NewKubeProviderCertManager()
	client := kubeProvider.ClientSet.CertmanagerV1().ClusterIssuers()
	err := client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sClusterIssuer(name string) utils.HttpResult {
	cmd := exec.Command("kubectl", "describe", "clusterissuer", name)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sClusterIssuer(data cmapi.ClusterIssuer) utils.HttpResult {
	kubeProvider := NewKubeProviderCertManager()
	client := kubeProvider.ClientSet.CertmanagerV1().ClusterIssuers()
	_, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func NewK8sClusterIssuer() K8sNewWorkload {
	return NewWorkload(
		RES_CLUSTER_ISSUER,
		utils.InitClusterIssuerYaml(),
		"A ClusterIssuer is a custom resource definition (CRD) in cert-manager, which is a native Kubernetes certificate management controller. It represents a certificate authority (CA) that can generate signed certificates at the cluster level. In this example, a ClusterIssuer named 'my-cluster-issuer' is created. This issuer uses the Let's Encrypt ACME server for the ACME protocol. It will use the secret 'my-cluster-issuer-account-key' to store the ACME account's private key and uses the HTTP-01 challenge for domain validation with the nginx ingress class.")
}
