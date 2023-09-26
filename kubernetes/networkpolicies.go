package kubernetes

import (
	"context"
	"os/exec"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllNetworkPolicies(namespaceName string, contextId *string) utils.K8sWorkloadResult {
	result := []v1.NetworkPolicy{}

	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	netPolist, err := provider.ClientSet.NetworkingV1().NetworkPolicies(namespaceName).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		logger.Log.Errorf("AllNetworkPolicies ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}

	for _, netpol := range netPolist.Items {
		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, netpol.ObjectMeta.Namespace) {
			result = append(result, netpol)
		}
	}
	return WorkloadResult(result, nil)
}

func GetNetworkPolicy(namespaceName string, name string, contextId *string) (*v1.NetworkPolicy, error) {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return nil, err
	}
	return provider.ClientSet.NetworkingV1().NetworkPolicies(namespaceName).Get(context.TODO(), name, metav1.GetOptions{})
}

func UpdateK8sNetworkPolicy(data v1.NetworkPolicy, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.NetworkingV1().NetworkPolicies(data.Namespace)
	res, err := client.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func DeleteK8sNetworkPolicy(data v1.NetworkPolicy, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.NetworkingV1().NetworkPolicies(data.Namespace)
	err = client.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sNetworkPolicyBy(namespace string, name string, contextId *string) error {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return err
	}
	client := provider.ClientSet.NetworkingV1().NetworkPolicies(namespace)
	return client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sNetworkPolicy(namespace string, name string, contextId *string) utils.K8sWorkloadResult {
	cmd := exec.Command("kubectl", ContextFlag(contextId), "describe", "netpol", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func CreateK8sNetworkpolicy(data v1.NetworkPolicy, contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return WorkloadResult(nil, err)
	}
	client := provider.ClientSet.NetworkingV1().NetworkPolicies(data.Namespace)
	res, err := client.Create(context.TODO(), &data, metav1.CreateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(res, nil)
}

func NewK8sNetPol() K8sNewWorkload {
	return NewWorkload(
		RES_NETWORK_POLICY,
		utils.InitNetPolYaml(),
		"A NetworkPolicy is a specification of how selections of pods are allowed to communicate with each other and other network endpoints. n this example, a NetworkPolicy named 'my-network-policy' is created. It applies to all Pods with the label role=db in the default namespace, and it sets both inbound and outbound rules.")
}
