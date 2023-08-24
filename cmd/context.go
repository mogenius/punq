package cmd

import (
	"encoding/base64"
	"os"

	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/services"
	"github.com/mogenius/punq/structs"
	"github.com/mogenius/punq/utils"
	"github.com/spf13/cobra"
)

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Manage punq contexts.",
	Long:  `The context command lets you manage all context related tasks like add, remove, list contexts.`,
}

var listContextCmd = &cobra.Command{
	Use:   "list",
	Short: "List punq contexts.",
	Long:  `The list command lets you list all contexts managed by punq.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.InitConfigYaml(false, nil, false)

		contexts := services.ListContexts()
		dtos.ListContexts(contexts)
	},
}

var addContextCmd = &cobra.Command{
	Use:   "add",
	Short: "Add punq context.",
	Long:  `The add command lets you add a context into punq.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.InitConfigYaml(false, nil, false)

		if filePath == "" {
			logger.Log.Fatal("-filePath cannot be empty.")
		}

		// load file
		dataBytes, err := os.ReadFile(filePath)
		if err != nil {
			logger.Log.Fatalf("Error reading file '%s': %s", filePath, err.Error())
		}
		encodedData := base64.StdEncoding.EncodeToString(dataBytes)

		newContext := dtos.PunqContext{
			Id:            utils.NanoId(),
			ContextBase64: encodedData,
		}

		services.AddContext(newContext)
	},
}

var addContextAccessCmd = &cobra.Command{
	Use:   "add-access",
	Short: "Add access to punq context.",
	Long:  `The add-access command lets you add a user + access level to a context in punq.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.InitConfigYaml(false, nil, false)

		if contextId == "" {
			logger.Log.Fatal("-context-id cannot be empty.")
		}

		if accessLevel == "" {
			logger.Log.Fatal("-access-level cannot be empty.")
		}

		if userId == "" {
			logger.Log.Fatal("-user-id cannot be empty.")
		}

		ctx := services.GetContext(contextId)
		if ctx == nil {
			logger.Log.Fatalf("context '%s' not found.", contextId)
		}

		ctx.AddAccess(userId, dtos.AccessLevelFromString(accessLevel))
		services.UpdateContext(*ctx)
	},
}

var removeContextAccessCmd = &cobra.Command{
	Use:   "remove-access",
	Short: "Remove access from punq context.",
	Long:  `The remove-access command lets you remove a users access level from a context in punq.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.InitConfigYaml(false, nil, false)

		if contextId == "" {
			logger.Log.Fatal("-context-id cannot be empty.")
		}

		if userId == "" {
			logger.Log.Fatal("-user-id cannot be empty.")
		}

		ctx := services.GetContext(contextId)
		if ctx == nil {
			logger.Log.Fatalf("context '%s' not found.", contextId)
		}

		ctx.RemoveAccess(userId)
		services.UpdateContext(*ctx)
	},
}

var deleteContextCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete punq context.",
	Long:  `The delete command lets you delete a specific context in punq.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.InitConfigYaml(false, nil, false)

		if contextId == "" {
			logger.Log.Fatal("-contextid cannot be empty.")
		}

		result := services.DeleteContext(contextId)
		structs.PrettyPrint(result)
	},
}

var getContextCmd = &cobra.Command{
	Use:   "get",
	Short: "Get specific punq context.",
	Long:  `The get command lets you get a specific context from punq.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.InitConfigYaml(false, nil, false)

		if contextId == "" {
			logger.Log.Fatal("-contextid cannot be empty.")
		}

		ctx := services.GetContext(contextId)
		if ctx == nil {
			logger.Log.Errorf("No context found for '%s'.", contextId)
		} else {
			structs.PrettyPrint(ctx)
		}
	},
}

func init() {
	contextCmd.AddCommand(listContextCmd)

	contextCmd.AddCommand(addContextAccessCmd)
	addContextAccessCmd.Flags().StringVarP(&contextId, "contextid", "c", "", "ContextId of the context")
	addContextAccessCmd.Flags().StringVarP(&userId, "user-id", "u", "", "Id of the user you want to add")
	addContextAccessCmd.Flags().StringVarP(&accessLevel, "access-level", "l", string(dtos.ADMIN), "Access level of the user you want to add (READER, USER, ADMIN)")

	contextCmd.AddCommand(removeContextAccessCmd)
	removeContextAccessCmd.Flags().StringVarP(&contextId, "contextid", "c", "", "ContextId of the context")
	removeContextAccessCmd.Flags().StringVarP(&userId, "user-id", "u", "", "Id of the user you want to add")

	contextCmd.AddCommand(addContextCmd)
	addContextCmd.Flags().StringVarP(&email, "filepath", "f", "", "FilePath to the context you want to add")

	contextCmd.AddCommand(deleteContextCmd)
	deleteContextCmd.Flags().StringVarP(&contextId, "contextid", "c", "", "ContextId of the context")

	contextCmd.AddCommand(getContextCmd)
	getContextCmd.Flags().StringVarP(&contextId, "contextid", "c", "", "ContextId of the context")

	rootCmd.AddCommand(contextCmd)
}
