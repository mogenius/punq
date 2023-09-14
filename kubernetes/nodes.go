package kubernetes

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/logger"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetNodeStats(contextId *string) []dtos.NodeStat {
	result := []dtos.NodeStat{}
	nodes := ListNodes(contextId)

	for index, node := range nodes {
		mem, _ := node.Status.Capacity.Memory().AsInt64()
		cpu, _ := node.Status.Capacity.Cpu().AsInt64()
		maxPods, _ := node.Status.Capacity.Pods().AsInt64()
		ephemeral, _ := node.Status.Capacity.StorageEphemeral().AsInt64()

		nodeStat := dtos.NodeStat{
			Name:             fmt.Sprintf("Node-%d", index+1),
			MaschineId:       node.Status.NodeInfo.MachineID,
			Cpus:             cpu,
			MemoryInBytes:    mem,
			EphemeralInBytes: ephemeral,
			MaxPods:          maxPods,
			KubletVersion:    node.Status.NodeInfo.KubeletVersion,
			OsType:           node.Status.NodeInfo.OperatingSystem,
			OsImage:          node.Status.NodeInfo.OSImage,
			Architecture:     node.Status.NodeInfo.Architecture,
		}
		result = append(result, nodeStat)
		nodeStat.PrintPretty()
	}
	return result
}

func ListK8sNodes(contextId *string) utils.K8sWorkloadResult {
	var provider *KubeProvider = NewKubeProvider(contextId)
	if provider == nil {
		err := fmt.Errorf("Failed to load kubeprovider.")
		logger.Log.Errorf(err.Error())
		return WorkloadResult(nil, err)
	}

	nodeMetricsList, err := provider.ClientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Errorf("ListNodeMetrics ERROR: %s", err.Error())
		return WorkloadResult(nil, err)
	}
	return WorkloadResult(nodeMetricsList.Items, nil)
}

func GetK8sNode(name string, contextId *string) (*v1.Node, error) {
	var provider *KubeProvider = NewKubeProvider(contextId)
	return provider.ClientSet.CoreV1().Nodes().Get(context.TODO(), name, metav1.GetOptions{})
}

func DeleteK8sNode(name string, contextId *string) error {
	var provider *KubeProvider = NewKubeProvider(contextId)
	return provider.ClientSet.CoreV1().Nodes().Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sNode(name string, contextId *string) utils.K8sWorkloadResult {
	cmd := exec.Command("kubectl", ContextFlag(contextId), "describe", "node", name)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}
