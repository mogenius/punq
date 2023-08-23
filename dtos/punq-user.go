package dtos

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mogenius/punq/utils"
)

type PunqUser struct {
	Id          string      `json:"id" validate:"required"`
	Email       string      `json:"email" validate:"required"`
	Password    string      `json:"password" validate:"required"`
	DisplayName string      `json:"displayName" validate:"required"`
	AccessLevel AccessLevel `json:"accessLevel" validate:"required"`
	Created     string      `json:"createdAt" validate:"required"`
}

func ListUsers(users []PunqUser, showPasswords bool) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	if showPasswords {
		t.AppendHeader(table.Row{"#", "ID", "DisplayName", "Email", "AccessLevel", "Created", "Password"})
		for index, user := range users {
			t.AppendRow(
				table.Row{index + 1, user.Id, user.DisplayName, user.Email, user.AccessLevel, utils.JsonStringToHumanDuration(user.Created), user.Password},
			)
		}
	} else {
		t.AppendHeader(table.Row{"#", "ID", "DisplayName", "Email", "AccessLevel", "Created"})
		for index, user := range users {
			t.AppendRow(
				table.Row{index + 1, user.Id, user.DisplayName, user.Email, user.AccessLevel, utils.JsonStringToHumanDuration(user.Created)},
			)
		}
	}
	t.Render()
}
