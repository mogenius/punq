package kubernetes

import (
	"context"
	"os/exec"
	"punq/logger"
	"punq/utils"

	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AllNetworkPolicies(namespaceName string) K8sWorkloadResult {
	result := []v1.NetworkPolicy{}

	provider := NewKubeProvider()
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

func UpdateK8sNetworkPolicy(data v1.NetworkPolicy) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	netpolClient := kubeProvider.ClientSet.NetworkingV1().NetworkPolicies(data.Namespace)
	_, err := netpolClient.Update(context.TODO(), &data, metav1.UpdateOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DeleteK8sNetworkPolicy(data v1.NetworkPolicy) K8sWorkloadResult {
	kubeProvider := NewKubeProvider()
	netpolClient := kubeProvider.ClientSet.NetworkingV1().NetworkPolicies(data.Namespace)
	err := netpolClient.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nil, nil)
}

func DescribeK8sNetworkPolicy(namespace string, name string) K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "netpol", name, "-n", namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}

func NewK8sNetPol() K8sNewWorkload {
	return NewWorkload(
		RES_NETWORK_POLICY,
		utils.InitNetPolYaml(),
		"A NetworkPolicy is a specification of how selections of pods are allowed to communicate with each other and other network endpoints. n this example, a NetworkPolicy named 'my-network-policy' is created. It applies to all Pods with the label role=db in the default namespace, and it sets both inbound and outbound rules.")
}
