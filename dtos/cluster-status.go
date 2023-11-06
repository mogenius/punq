package dtos

import (
	"time"

	"github.com/mogenius/punq/utils"
)

type ClusterStatusDto struct {
	ClusterName                  string                `json:"clusterName"`
	Pods                         int                   `json:"pods"`
	PodCpuUsageInMilliCores      int                   `json:"podCpuUsageInMilliCores"`
	PodCpuLimitInMilliCores      int                   `json:"podCpuLimitInMilliCores"`
	PodMemoryUsageInBytes        int64                 `json:"podMemoryUsageInBytes"`
	PodMemoryLimitInBytes        int64                 `json:"podMemoryLimitInBytes"`
	EphemeralStorageLimitInBytes int64                 `json:"ephemeralStorageLimitInBytes"`
	CurrentTime                  string                `json:"currentTime"`
	KubernetesVersion            string                `json:"kubernetesVersion"`
	Platform                     string                `json:"platform"`
	Country                      *utils.CountryDetails `json:"country"`
}

func ClusterStatusDtoExmapleData() ClusterStatusDto {
	return ClusterStatusDto{
		ClusterName:                  "clusterName",
		Pods:                         1,
		PodCpuUsageInMilliCores:      1,
		PodCpuLimitInMilliCores:      1,
		PodMemoryUsageInBytes:        123,
		PodMemoryLimitInBytes:        1234,
		EphemeralStorageLimitInBytes: 166,
		CurrentTime:                  time.Now().Format(time.RFC3339),
		KubernetesVersion:            "v1.27.2",
		Platform:                     "linux/arm64",
		Country:                      nil,
	}
}
