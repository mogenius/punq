package cmd

import (
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/services"
	"github.com/mogenius/punq/structs"
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
		dtos.ListUsers(users, showPasswords)
	},
}

var addUserCmd = &cobra.Command{
	Use:   "add",
	Short: "Add punq user.",
	Long:  `The add command lets you add a user into punq.`,
	Run: func(cmd *cobra.Command, args []string) {
		if email == "" {
			logger.Log.Fatal("-email cannot be empty.")
		}
		if displayName == "" {
			logger.Log.Fatal("-displayname cannot be empty.")
		}
		if password == "" {
			logger.Log.Fatal("-password cannot be empty.")
		}
		selectedAccess := dtos.READER // default level
		if accessLevel != "" {
			selectedAccess = dtos.AccessLevelFromString(accessLevel)
		}

		services.AddUser(dtos.PunqUserCreateInput{
			Email:       email,
			Password:    password,
			DisplayName: displayName,
			AccessLevel: selectedAccess,
		})
	},
}

var updateUserCmd = &cobra.Command{
	Use:   "update",
	Short: "Update punq user.",
	Long:  `The update command lets you update a user in punq.`,
	Run: func(cmd *cobra.Command, args []string) {
		if userId == "" {
			logger.Log.Fatal("Please selecte a userId to update a user.")
		}
		if email == "" && displayName == "" && password == "" && accessLevel == "" {
			logger.Log.Fatal("One of the following options must be used to update a user: -email -displayname -password -accesslevel")
		}

		user := dtos.PunqUserUpdateInput{
			Id: userId,
		}
		// if user == nil {
		// 	logger.Log.Fatalf("Selected userId '%s' not found.", userId)
		// }
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

		_, err := services.UpdateUser(user)
		if err != nil {
			logger.Log.Fatalf(err.Error())
		}
	},
}

var deleteUserCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete punq user.",
	Long:  `The delete command lets you delete a specific user in punq.`,
	Run: func(cmd *cobra.Command, args []string) {
		if userId == "" {
			logger.Log.Fatal("-userid cannot be empty.")
		}

		result := services.DeleteUser(userId)
		structs.PrettyPrint(result)
	},
}

var getUserCmd = &cobra.Command{
	Use:   "get",
	Short: "Get specific punq user.",
	Long:  `The get command lets you get a specific user of punq.`,
	Run: func(cmd *cobra.Command, args []string) {
		if userId == "" {
			logger.Log.Fatal("-userid cannot be empty.")
		}

		user, err := services.GetUser(userId)
		if err != nil {
			logger.Log.Fatal(err)
		}
		structs.PrettyPrint(user)
	},
}

func init() {
	userCmd.AddCommand(listUserCmd)
	listUserCmd.Flags().BoolVarP(&showPasswords, "show-passwords", "s", false, "Display current passwords")

	userCmd.AddCommand(updateUserCmd)
	updateUserCmd.Flags().StringVarP(&userId, "userid", "u", "", "UserId of the user")
	updateUserCmd.Flags().StringVarP(&email, "email", "e", "", "E-Mail address of the user")
	updateUserCmd.Flags().StringVarP(&displayName, "displayname", "j", "", "Display name of the user")
	updateUserCmd.Flags().StringVarP(&password, "password", "p", "", "Password of the user")
	updateUserCmd.Flags().StringVarP(&accessLevel, "accesslevel", "a", "", "AccessLeve of the user")

	userCmd.AddCommand(addUserCmd)
	addUserCmd.Flags().StringVarP(&email, "email", "e", "", "E-Mail address of the new user")
	addUserCmd.Flags().StringVarP(&displayName, "displayname", "j", "", "Display name of the new user")
	addUserCmd.Flags().StringVarP(&password, "password", "p", "", "Password of the new user")
	addUserCmd.Flags().StringVarP(&accessLevel, "accesslevel", "a", "", "AccessLeve of the new user")

	userCmd.AddCommand(deleteUserCmd)
	deleteUserCmd.Flags().StringVarP(&userId, "userid", "u", "", "UserId of the user")

	userCmd.AddCommand(getUserCmd)
	getUserCmd.Flags().StringVarP(&userId, "userid", "u", "", "UserId of the user")

	rootCmd.AddCommand(userCmd)
}
