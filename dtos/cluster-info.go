package dtos

import "github.com/mogenius/punq/utils"

type ClusterInfoDto struct {
	ClusterStatus ClusterStatusDto      `json:"clusterStatus"`
	NodeStats     []NodeStat            `json:"nodeStats"`
	Country       *utils.CountryDetails `json:"country"`
}
