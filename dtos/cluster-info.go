package dtos

type ClusterInfoDto struct {
	ClusterStatus ClusterStatusDto `json:"clusterStatus"`
	NodeStats     []NodeStat       `json:"nodeStats"`
}
