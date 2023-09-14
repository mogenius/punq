package dtos

import (
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

type PunqContext struct {
	Id      string       `json:"id" validate:"required"`
	Name    string       `json:"name" validate:"required"`
	Context string       `json:"context" validate:"required"`
	Access  []PunqAccess `json:"access" validate:"required"`
}

func (c *PunqContext) AddAccess(userId string, accessLevel AccessLevel) {
	for _, access := range c.Access {
		if access.UserId == userId {
			// UPDATE EXISTING
			access.Level = accessLevel
			return
		}
	}
	// CREATE NEW
	c.Access = append(c.Access, PunqAccess{
		UserId: userId,
		Level:  accessLevel,
	})
}

func (c *PunqContext) RemoveAccess(userId string) {
	resultingArray := []PunqAccess{}
	for _, access := range c.Access {
		if access.UserId != userId {
			resultingArray = append(resultingArray, access)
		}
	}
	c.Access = resultingArray
}

func ListContexts(contexts []PunqContext) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "ID", "Name", "Access"})
	for index, context := range contexts {
		accessStr := "*"
		accessEntries := []string{}
		for _, access := range context.Access {
			accessEntries = append(accessEntries, fmt.Sprintf("%s (%s)", access.UserId, access.Level))
		}
		if len(accessEntries) > 0 {
			accessStr = strings.Join(accessEntries, ", ")
		}
		t.AppendRow(
			table.Row{index + 1, context.Id, context.Name, accessStr},
		)
	}
	t.Render()
}
