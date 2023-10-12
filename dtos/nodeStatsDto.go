package dtos

import (
	"fmt"

	"github.com/mogenius/punq/utils"
)

type NodeStat struct {
	Name                  string  `json:"name" validate:"required"`
	MaschineId            string  `json:"maschineId" validate:"required"`
	CpuInCores            int64   `json:"cpuInCores" validate:"required"`
	CpuInCoresUtilized    float64 `json:"cpuInCoresUtilized" validate:"required"`
	MemoryInBytes         int64   `json:"memoryInBytes" validate:"required"`
	MemoryInBytesUtilized int64   `json:"memoryInBytesUtilized" validate:"required"`
	EphemeralInBytes      int64   `json:"ephemeralInBytes" validate:"required"`
	MaxPods               int64   `json:"maxPods" validate:"required"`
	KubletVersion         string  `json:"kubletVersion" validate:"required"`
	OsType                string  `json:"osType" validate:"required"`
	OsImage               string  `json:"osImage" validate:"required"`
	Architecture          string  `json:"architecture" validate:"required"`
}

func (o *NodeStat) PrintPretty() {
	fmt.Printf("%s: %s %s [%s/%s] - CPUs: %d, RAM: %s, Ephemeral: %s, MaxPods: %d\n",
		o.Name,
		o.KubletVersion,
		o.OsImage,
		o.OsType,
		o.Architecture,
		o.CpuInCores,
		utils.BytesToHumanReadable(o.MemoryInBytes),
		utils.BytesToHumanReadable(o.EphemeralInBytes),
		o.MaxPods,
	)
}
