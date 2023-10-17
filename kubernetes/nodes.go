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
	v1metrics "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

func GetNodeStats(contextId *string) []dtos.NodeStat {
	result := []dtos.NodeStat{}
	nodes := ListNodes(contextId)
	nodeMetrics := ListNodeMetricss(contextId)

	for index, node := range nodes {

		utilizedCores := float64(0)
		utilizedMemory := int64(0)
		if len(nodeMetrics) > 0 {
			// Find the corresponding node metrics
			var nodeMetric *v1metrics.NodeMetrics
			for _, nm := range nodeMetrics {
				if nm.Name == node.Name {
					nodeMetric = &nm
					break
				}
			}

			// CPU
			cpuUsageDec := nodeMetric.Usage.Cpu().AsDec()
			cpuUsage, works := cpuUsageDec.Unscaled()
			if !works {
				logger.Log.Errorf("Failed to get CPU usage for node %s", node.Name)
			}
			if cpuUsage == 0 {
				cpuUsage = 1
			}
			utilizedCores = float64(cpuUsage) / 1000000000

			// Memory
			utilizedMemory, works = nodeMetric.Usage.Memory().AsInt64()
			if !works {
				logger.Log.Errorf("Failed to get MEMORY usage for node %s", node.Name)
			}
		}

		mem, _ := node.Status.Capacity.Memory().AsInt64()
		cpu, _ := node.Status.Capacity.Cpu().AsInt64()
		maxPods, _ := node.Status.Capacity.Pods().AsInt64()
		ephemeral, _ := node.Status.Capacity.StorageEphemeral().AsInt64()

		nodeStat := dtos.NodeStat{
			Name:                  fmt.Sprintf("Node-%d", index+1),
			MaschineId:            node.Status.NodeInfo.MachineID,
			CpuInCores:            cpu,
			CpuInCoresUtilized:    utilizedCores,
			MemoryInBytes:         mem,
			MemoryInBytesUtilized: utilizedMemory,
			EphemeralInBytes:      ephemeral,
			MaxPods:               maxPods,
			KubletVersion:         node.Status.NodeInfo.KubeletVersion,
			OsType:                node.Status.NodeInfo.OperatingSystem,
			OsImage:               node.Status.NodeInfo.OSImage,
			Architecture:          node.Status.NodeInfo.Architecture,
		}
		result = append(result, nodeStat)
		//nodeStat.PrintPretty()
	}
	return result
}

func ListK8sNodes(contextId *string) utils.K8sWorkloadResult {
	provider, err := NewKubeProvider(contextId)
	if provider == nil || err != nil {
		err := fmt.Errorf("failed to load provider")
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
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return nil, err
	}
	return provider.ClientSet.CoreV1().Nodes().Get(context.TODO(), name, metav1.GetOptions{})
}

func DeleteK8sNode(name string, contextId *string) error {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return err
	}
	return provider.ClientSet.CoreV1().Nodes().Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func DescribeK8sNode(name string, contextId *string) utils.K8sWorkloadResult {
	cmd := exec.Command("kubectl", "describe", "node", name, ContextFlag(contextId))

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorf("Failed to execute command (%s): %v", cmd.String(), err)
		logger.Log.Errorf("Error: %s", string(output))
		return WorkloadResult(nil, string(output))
	}
	return WorkloadResult(string(output), nil)
}
