package dtos

// NAME       RESTRICTION
// READER     (no logs, no yaml, no secrets, no exec, no useradmin)
// USER       (no secrets, no exec, no useradmin)
// ADMIN      (none)

type AccessLevel string

const (
	READER AccessLevel = "READER"
	USER   AccessLevel = "USER"
	ADMIN  AccessLevel = "ADMIN"
	// and so on...
)

type PunqAccess struct {
	UserId string      `json:"userId" validate:"required"`
	Level  AccessLevel `json:"level" validate:"required"`
}

// func ListAccess(groups []PunqUser, showPasswords bool) {
// 	t := table.NewWriter()
// 	t.SetOutputMirror(os.Stdout)
// 	if showPasswords {
// 		t.AppendHeader(table.Row{"#", "ID", "DisplayName", "Email", "Password"})
// 		for index, user := range users {
// 			t.AppendRow(
// 				table.Row{index + 1, user.Id, user.DisplayName, user.Email, user.Password},
// 			)
// 		}
// 	} else {
// 		t.AppendHeader(table.Row{"#", "ID", "DisplayName", "Email"})
// 		for index, user := range users {
// 			t.AppendRow(
// 				table.Row{index + 1, user.Id, user.DisplayName, user.Email},
// 			)
// 		}
// 	}
// 	t.Render()
// }
