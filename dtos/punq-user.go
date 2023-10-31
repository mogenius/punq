package dtos

import (
	"errors"
	"os"

	"golang.org/x/crypto/bcrypt"

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

type PunqUserCreateInput struct {
	Email       string      `json:"email" validate:"required"`
	Password    string      `json:"password" validate:"required"`
	DisplayName string      `json:"displayName" validate:"required"`
	AccessLevel AccessLevel `json:"accessLevel" validate:"required"`
}

func ListUsers(users []PunqUser) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "ID", "DisplayName", "Email", "AccessLevel", "Created"})
	for index, user := range users {
		t.AppendRow(
			table.Row{index + 1, user.Id, user.DisplayName, user.Email, user.AccessLevel.String(), utils.JsonStringToHumanDuration(user.Created)},
		)
	}
	t.Render()
}

func (user *PunqUser) PrintToTerminal() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "DisplayName", "Email", "AccessLevel"})
	t.AppendRow(
		table.Row{user.Id, user.DisplayName, user.Email, user.AccessLevel.String()},
	)
	t.Render()
}

func (user *PunqUser) PrintToTerminalWithPwd() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "DisplayName", "Email", "Password", "AccessLevel"})
	t.AppendRow(
		table.Row{user.Id, user.DisplayName, user.Email, user.Password, user.AccessLevel.String()},
	)
	t.Render()
}

func (user *PunqUser) PasswordCheck(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, err
	}
	return true, nil
}
