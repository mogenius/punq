package structs

type StatsData struct {
	Data []Stats `json:"data"`
}

type Stats struct {
	Cluster               string `json:"cluster"`
	Namespace             string `json:"namespace"`
	PodName               string `json:"podName"`
	Cpu                   int64  `json:"cpu"`
	CpuLimit              int64  `json:"cpuLimit"`
	Memory                int64  `json:"memory"`
	MemoryLimit           int64  `json:"memoryLimit"`
	EphemeralStorageLimit int64  `json:"ephemeralStorageLimit"`
	StartTime             string `json:"startTime"`
}
