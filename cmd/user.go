package cmd

import (
	"github.com/google/uuid"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/services"
	"github.com/mogenius/punq/structs"
	"github.com/mogenius/punq/utils"
	"github.com/spf13/cobra"
)

var email string
var password string
var displayName string
var userId string

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage punq users.",
	Long:  `The user command lets you manage all user related task like add, remove, list users.`,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List punq users.",
	Long:  `The user command lets you list all users of punq.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.InitConfigYaml(false, nil, false)

		users := services.ListUsers()
		for _, user := range users {
			structs.PrettyPrint(user)
		}
	},
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add punq user.",
	Long:  `The add command lets you add a user into punq.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.InitConfigYaml(false, nil, false)

		if email == "" {
			logger.Log.Fatal("-email cannot be empty.")
		}
		if displayName == "" {
			logger.Log.Fatal("-displayname cannot be empty.")
		}
		if password == "" {
			logger.Log.Fatal("-password cannot be empty.")
		}

		newUser := dtos.PunqUser{
			Id:          uuid.New().String(),
			Email:       email,
			Password:    password,
			DisplayName: displayName,
		}

		services.AddUser(newUser)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete punq user.",
	Long:  `The delete command lets you delete a specific user in punq.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.InitConfigYaml(false, nil, false)

		if userId == "" {
			logger.Log.Fatal("-userid cannot be empty.")
		}

		result := services.DeleteUser(userId)
		structs.PrettyPrint(result)
	},
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get specific punq user.",
	Long:  `The get command lets you get a specific user of punq.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.InitConfigYaml(false, nil, false)

		if userId == "" {
			logger.Log.Fatal("-userid cannot be empty.")
		}

		user := services.GetUser(userId)
		structs.PrettyPrint(user)
	},
}

func init() {
	userCmd.AddCommand(listCmd)

	userCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&email, "email", "e", "", "E-Mail address of the new user")
	addCmd.Flags().StringVarP(&displayName, "displayname", "d", "", "Display name of the new user")
	addCmd.Flags().StringVarP(&password, "password", "p", "", "Password of the new user")

	userCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVarP(&userId, "userid", "u", "", "UserId of the user")

	userCmd.AddCommand(getCmd)
	getCmd.Flags().StringVarP(&userId, "userid", "u", "", "UserId of the user")

	rootCmd.AddCommand(userCmd)
}
