package dtos

import (
	"fmt"
	"punq/utils"
)

type NodeStat struct {
	Name             string `json:"name" validate:"required"`
	MaschineId       string `json:"maschineId" validate:"required"`
	Cpus             int64  `json:"cpus" validate:"required"`
	MemoryInBytes    int64  `json:"memoryInBytes" validate:"required"`
	EphemeralInBytes int64  `json:"ephemeralInBytes" validate:"required"`
	MaxPods          int64  `json:"maxPods" validate:"required"`
	KubletVersion    string `json:"kubletVersion" validate:"required"`
	OsType           string `json:"osType" validate:"required"`
	OsImage          string `json:"osImage" validate:"required"`
	Architecture     string `json:"architecture" validate:"required"`
}

func (o *NodeStat) PrintPretty() {
	fmt.Printf("%s: %s %s [%s/%s] - CPUs: %d, RAM: %s, Ephemeral: %s, MaxPods: %d\n",
		o.Name,
		o.KubletVersion,
		o.OsImage,
		o.OsType,
		o.Architecture,
		o.Cpus,
		utils.BytesToHumanReadable(o.MemoryInBytes),
		utils.BytesToHumanReadable(o.EphemeralInBytes),
		o.MaxPods,
	)
}
