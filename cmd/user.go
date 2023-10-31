package cmd

import (
	"fmt"
	"strings"

	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/services"
	"github.com/mogenius/punq/utils"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage punq users.",
	Long:  `The user command lets you manage all user related task like add, remove, list users.`,
}

var listUserCmd = &cobra.Command{
	Use:   "list",
	Short: "List punq users.",
	Long:  `The list command lets you list all users of punq.`,
	Run: func(cmd *cobra.Command, args []string) {
		users := services.ListUsers()
		dtos.ListUsers(users)
	},
}

var addUserCmd = &cobra.Command{
	Use:   "add",
	Short: "Add punq user.",
	Long:  `The add command lets you add a user into punq.`,
	Run: func(cmd *cobra.Command, args []string) {
		selectedAccess := dtos.READER // default level
		if accessLevel != "" {
			selectedAccess = dtos.AccessLevelFromString(accessLevel)
		} else {
			selectedAccess = dtos.READER
		}

		firstname := utils.RandomFirstName()
		middlename := utils.RandomMiddleName()
		lastname := utils.RandomLastName()

		if email == "" {
			email = fmt.Sprintf("%s-%s@punq.dev", strings.ToLower(firstname), strings.ToLower(lastname))
		}

		if password == "" {
			password = utils.NanoId()
		}

		if displayName == "" {
			displayName = fmt.Sprintf("%s %s %s", firstname, middlename, lastname)
		}

		newUser := dtos.PunqUserCreateInput{
			Email:       email,
			Password:    password,
			DisplayName: displayName,
			AccessLevel: selectedAccess,
		}

		newPunqUser, err := services.AddUser(newUser)
		if err != nil {
			utils.FatalError(err.Error())
		} else {
			utils.PrintInfo("User added succesfully ✅.")
			newPunqUser.Password = password // print password to terminal once
			newPunqUser.PrintToTerminalWithPwd()
		}
	},
}

var updateUserCmd = &cobra.Command{
	Use:   "update",
	Short: "Update punq user.",
	Long:  `The update command lets you update a user in punq.`,
	Run: func(cmd *cobra.Command, args []string) {
		RequireStringFlag(userId, "user-id")

		if email == "" && displayName == "" && password == "" && accessLevel == "" {
			utils.FatalError("One of the following options must be used to update a user: -email -displayname -password -accesslevel")
		}

		user, err := services.GetUser(userId)
		if err != nil || user == nil {
			utils.FatalError(fmt.Sprintf("Selected userId '%s' not found.", userId))
		}

		if displayName != "" {
			user.DisplayName = displayName
		}
		if password != "" {
			user.Password = password
		}
		if email != "" {
			user.Email = email
		}
		if accessLevel != "" {
			user.AccessLevel = dtos.AccessLevelFromString(accessLevel)
		}

		_, err = services.UpdateUser(*user)
		if err != nil {
			utils.FatalError(err.Error())
		} else {
			utils.PrintInfo("User updated succesfully ✅.")
		}
	},
}

var deleteUserCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete punq user.",
	Long:  `The delete command lets you delete a specific user in punq.`,
	Run: func(cmd *cobra.Command, args []string) {
		RequireStringFlag(userId, "user-id")

		err := services.DeleteUser(userId)
		if err != nil {
			utils.FatalError(err.Error())
		}
		utils.PrintInfo(fmt.Sprintf("User %s successfully deleted.", userId))
	},
}

var getUserCmd = &cobra.Command{
	Use:   "get",
	Short: "Get specific punq user.",
	Long:  `The get command lets you get a specific user of punq.`,
	Run: func(cmd *cobra.Command, args []string) {
		RequireStringFlag(userId, "user-id")

		user, err := services.GetUser(userId)
		if err != nil {
			utils.FatalError(err.Error())
		}
		user.PrintToTerminal()
	},
}

func init() {
	userCmd.AddCommand(listUserCmd)

	userCmd.AddCommand(updateUserCmd)
	updateUserCmd.Flags().StringVarP(&userId, "user-id", "u", "", "UserId of the user")
	updateUserCmd.Flags().StringVarP(&email, "email", "e", "", "E-Mail address of the user")
	updateUserCmd.Flags().StringVarP(&displayName, "displayname", "j", "", "Display name of the user")
	updateUserCmd.Flags().StringVarP(&password, "password", "p", "", "Password of the user")
	updateUserCmd.Flags().StringVarP(&accessLevel, "accesslevel", "a", "", "AccessLevel of the user: valid values are READER, USER, or ADMIN")

	userCmd.AddCommand(addUserCmd)
	addUserCmd.Flags().StringVarP(&email, "email", "e", "", "E-Mail address of the new user")
	addUserCmd.Flags().StringVarP(&displayName, "displayname", "j", "", "Display name of the new user")
	addUserCmd.Flags().StringVarP(&password, "password", "p", "", "Password of the new user")
	addUserCmd.Flags().StringVarP(&accessLevel, "accesslevel", "a", "", "AccessLevel of the new user: valid values are READER, USER, or ADMIN")

	userCmd.AddCommand(deleteUserCmd)
	deleteUserCmd.Flags().StringVarP(&userId, "userid", "u", "", "UserId of the user")

	userCmd.AddCommand(getUserCmd)
	getUserCmd.Flags().StringVarP(&userId, "userid", "u", "", "UserId of the user")

	rootCmd.AddCommand(userCmd)
}
