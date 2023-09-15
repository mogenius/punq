package dtos

import "time"

type ClusterStatusDto struct {
	ClusterName           string `json:"clusterName"`
	Pods                  int    `json:"pods"`
	CpuInMilliCores       int    `json:"cpu"`
	CpuLimitInMilliCores  int    `json:"cpuLimit"`
	Memory                string `json:"memory"`
	MemoryLimit           string `json:"memoryLimit"`
	EphemeralStorageLimit string `json:"ephemeralStorageLimit"`
	CurrentTime           string `json:"currentTime"`
	KubernetesVersion     string `json:"kubernetesVersion"`
	Platform              string `json:"platform"`
}

func ClusterStatusDtoExmapleData() ClusterStatusDto {
	return ClusterStatusDto{
		ClusterName:           "clusterName",
		Pods:                  1,
		CpuInMilliCores:       1,
		CpuLimitInMilliCores:  1,
		Memory:                "123",
		MemoryLimit:           "1456",
		EphemeralStorageLimit: "166",
		CurrentTime:           time.Now().Format(time.RFC3339),
		KubernetesVersion:     "v1.27.2",
		Platform:              "linux/arm64",
	}
}
