package cmd

import (
	"fmt"
	"os"

	"github.com/mogenius/punq/dtos"
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
		contexts := services.ListContexts()
		dtos.ListContextsToTerminal(contexts)
	},
}

var addContextCmd = &cobra.Command{
	Use:   "add",
	Short: "Add punq context.",
	Long:  `The add command lets you add a context into punq.`,
	Run: func(cmd *cobra.Command, args []string) {
		if filePath == "" {
			FatalError("-f cannot be empty. Please select a context yaml to load.")
		}

		// load file
		dataBytes, err := os.ReadFile(filePath)
		if err != nil {
			FatalError(fmt.Sprintf("Error reading file '%s': %s", filePath, err.Error()))
		}

		// newCtx := dtos.CreateContext("", contextId, string(dataBytes), []dtos.PunqAccess{})
		// structs.PrettyPrint(newCtx)

		contexts, err := services.ParseConfigToPunqContexts(dataBytes)
		if err != nil {
			FatalError(err.Error())
		}
		dtos.ListContextsToTerminal(contexts)

		index := utils.SelectIndexInteractive("Select context to add", len(contexts))
		// one
		if index > 0 {
			selectedContext := contexts[index-1]
			//selectedContext.PrintToTerminal()
			_, err := services.AddContext(selectedContext)
			if err != nil {
				FatalError(err.Error())
			} else {
				fmt.Printf("Context '%s' added ✅.\n", selectedContext.Name)
			}
		}
		// all
		if index == -2 {
			dtos.ListContextsToTerminal(contexts)
			for _, ctx := range contexts {
				services.AddContext(ctx)
				if err != nil {
					FatalError(err.Error())
				} else {
					fmt.Printf("Context '%s' added ✅.\n", ctx.Name)
				}
			}
		}
	},
}

var addContextAccessCmd = &cobra.Command{
	Use:   "add-access",
	Short: "Add access to punq context.",
	Long:  `The add-access command lets you add a user + access level to a context in punq.`,
	Run: func(cmd *cobra.Command, args []string) {
		if contextId == "" {
			FatalError("-context-id cannot be empty.")
		}

		if accessLevel == "" {
			FatalError("-access-level cannot be empty.")
		}

		if userId == "" {
			FatalError("-user-id cannot be empty.")
		}

		ctx, _ := services.GetContext(contextId)
		if ctx == nil {
			FatalError(fmt.Sprintf("context '%s' not found.", contextId))
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
		if contextId == "" {
			FatalError("-context-id cannot be empty.")
		}

		if userId == "" {
			FatalError("-user-id cannot be empty.")
		}

		ctx, _ := services.GetContext(contextId)
		if ctx == nil {
			FatalError(fmt.Sprintf("context '%s' not found.", contextId))
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
		if contextId == "" {
			PrintError("-contextid cannot be empty.")
		}

		result, err := services.DeleteContext(contextId)
		if err != nil {
			PrintError(err.Error())
		}
		if result != nil {
			structs.PrettyPrint(result)
		}
	},
}

var getContextCmd = &cobra.Command{
	Use:   "get",
	Short: "Get specific punq context.",
	Long:  `The get command lets you get a specific context from punq.`,
	Run: func(cmd *cobra.Command, args []string) {
		if contextId == "" {
			FatalError("-contextid cannot be empty.")
		}

		ctx, _ := services.GetContext(contextId)
		if ctx == nil {
			fmt.Printf("No context found for '%s'.\n", contextId)
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
	addContextCmd.Flags().StringVarP(&contextId, "contextid", "c", "", "ContextId of the context")
	addContextCmd.Flags().StringVarP(&filePath, "filepath", "f", "", "FilePath to the context you want to add")

	contextCmd.AddCommand(deleteContextCmd)
	deleteContextCmd.Flags().StringVarP(&contextId, "contextid", "c", "", "ContextId of the context")

	contextCmd.AddCommand(getContextCmd)
	getContextCmd.Flags().StringVarP(&contextId, "contextid", "c", "", "ContextId of the context")

	rootCmd.AddCommand(contextCmd)
}
