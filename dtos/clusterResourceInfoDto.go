package dtos

import "github.com/mogenius/punq/utils"

type ClusterResourceInfoDto struct {
	LoadBalancerExternalIps []string              `json:"loadBalancerExternalIps"`
	NodeStats               []NodeStat            `json:"nodeStats"`
	Country                 *utils.CountryDetails `json:"country"`
}
