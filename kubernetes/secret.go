package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllSecrets(namespaceName string, contextId *string) []v1.Secret {
	result := []v1.Secret{}

	provider := NewKubeProvider(contextId)
	secretList, err := provider.ClientSet.CoreV1().Secrets(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllSecrets ERROR: %s", err.Error())
		return result
	}

	for _, secret := range secretList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, secret.ObjectMeta.Namespace) {
			result = append(result, secret)
		}
	}
	return result
}

func SecretFor(namespace string, name string, contextId *string) *v1.Secret {
	kubeProvider := NewKubeProvider(contextId)
	secretClient := kubeProvider.ClientSet.CoreV1().Secrets(namespace)
	secret, err := secretClient.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logger.Log.Errorf("SecretFor ERROR: %s", err.Error())
		return nil
	}
	return secret
}
func GetSecret(namespace string, name string, contextId *string) (*v1.Secret, error) {
	kubeProvider := NewKubeProvider(contextId)
	secretClient := kubeProvider.ClientSet.CoreV1().Secrets(namespace)
	return secretClient.Get(context.TODO(), name, metav1.GetOptions{})
}

func AllK8sSecrets(namespaceName string, contextId *string) utils.K8sWorkloadResult {
	result := []v1.Secret{}

	provider := NewKubeProvider(contextId)
	secretList, err := provider.ClientSet.CoreV1().Secrets(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllSecrets ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, secret := range secretList.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, secret.ObjectMeta.Namespace) {
			result = append(result, secret)
		}
	}
	return WorkloadResult(result, nil)
}

func UpdateK8sSecret(data v1.Secret, contextId *string) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider(contextId)
	secretClient := kubeProvider.ClientSet.CoreV1().Secrets(data.Namespace)
	res, err := secretClient.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func DeleteK8sSecret(data v1.Secret, contextId *string) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider(contextId)
	secretClient := kubeProvider.ClientSet.CoreV1().Secrets(data.Namespace)
	err := secretClient.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sSecretBy(namespace string, name string, contextId *string) error {
	kubeProvider := NewKubeProvider(contextId)
	secretClient := kubeProvider.ClientSet.CoreV1().Secrets(namespace)
	return secretClient.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sSecret(namespace string, name string, contextId *string) utils.K8sWorkloadResult {
	cmd := exec.Command("kubectl", ContextFlag(contextId), "describe", "secret", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sSecret(data v1.Secret, contextId *string) utils.K8sWorkloadResult {
	kubeProvider := NewKubeProvider(contextId)
	client := kubeProvider.ClientSet.CoreV1().Secrets(data.Namespace)
	res, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func NewK8sSecret() K8sNewWorkload {
	return NewWorkload(
		RES_SECRET,
		utils.InitSecretYaml(),
		"A Secret is an object that contains a small amount of sensitive data such as a password, a token, or a key. In this example, a secret named 'my-secret' is created with two pieces of data: username and password. The values are arbitrary and must be base64 encoded. Please note, the Secret data is not encrypted, it's just base64 encoded. So it's not secure to store highly sensitive information. You should consider additional layer of protection such as using Kubernetes RBAC to restrict access to Secrets, and/or use solutions like sealed secrets, HashiCorp Vault, or other Kubernetes native solutions like the secrets-store-csi-driver project.")
}
