package dtos

import "time"

type ClusterStatusDto struct {
	ClusterName                  string `json:"clusterName"`
	Pods                         int    `json:"pods"`
	CpuInMilliCores              int    `json:"cpu"`
	CpuLimitInMilliCores         int    `json:"cpuLimit"`
	MemoryInBytes                int64  `json:"memoryInBytes"`
	MemoryLimitInBytes           int64  `json:"memoryLimitInBytes"`
	EphemeralStorageLimitInBytes int64  `json:"ephemeralStorageLimitInBytes"`
	CurrentTime                  string `json:"currentTime"`
	KubernetesVersion            string `json:"kubernetesVersion"`
	Platform                     string `json:"platform"`
}

func ClusterStatusDtoExmapleData() ClusterStatusDto {
	return ClusterStatusDto{
		ClusterName:                  "clusterName",
		Pods:                         1,
		CpuInMilliCores:              1,
		CpuLimitInMilliCores:         1,
		MemoryInBytes:                123,
		MemoryLimitInBytes:           1456,
		EphemeralStorageLimitInBytes: 166,
		CurrentTime:                  time.Now().Format(time.RFC3339),
		KubernetesVersion:            "v1.27.2",
		Platform:                     "linux/arm64",
	}
}
