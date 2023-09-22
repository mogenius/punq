package dtos

import (
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mogenius/punq/utils"
)

type PunqContext struct {
	Id          string       `json:"id" validate:"required"`
	Name        string       `json:"name" validate:"required"`
	ContextHash string       `json:"contextHash" validate:"required"`
	Context     string       `json:"context" validate:"required"`
	Access      []PunqAccess `json:"access" validate:"required"`
}

func CreateContext(id string, name string, context string, access []PunqAccess) PunqContext {
	ctx := PunqContext{}

	ctx.Name = name

	if id == "" {
		ctx.Id = utils.NanoId()
	} else {
		ctx.Id = id
	}

	ctx.Context = context
	ctx.ContextHash = utils.HashString(context)

	if len(access) > 0 {
		ctx.Access = access
	} else {
		ctx.Access = []PunqAccess{}
	}

	return ctx
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

func (c *PunqContext) PrintToTerminal() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Name", "Access", "Hash"})
	accessStr := "*"
	accessEntries := []string{}
	for _, access := range c.Access {
		accessEntries = append(accessEntries, fmt.Sprintf("%s (%d)", access.UserId, access.Level))
	}
	if len(accessEntries) > 0 {
		accessStr = strings.Join(accessEntries, ", ")
	}
	t.AppendRow(
		table.Row{c.Id, c.Name, accessStr, c.ContextHash},
	)
	t.Render()
}

func ListContextsToTerminal(contexts []PunqContext) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "ID", "Name", "Access", "Hash"})
	for index, context := range contexts {
		accessStr := "*"
		accessEntries := []string{}
		for _, access := range context.Access {
			accessEntries = append(accessEntries, fmt.Sprintf("%s (%d)", access.UserId, access.Level))
		}
		if len(accessEntries) > 0 {
			accessStr = strings.Join(accessEntries, ", ")
		}
		t.AppendRow(
			table.Row{index + 1, context.Id, context.Name, accessStr, context.ContextHash},
		)
	}
	t.Render()
}
