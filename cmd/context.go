package cmd

import (
	"encoding/base64"
	"os"

	"github.com/google/uuid"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/services"
	"github.com/mogenius/punq/structs"
	"github.com/mogenius/punq/utils"
	"github.com/spf13/cobra"
)

var filePath string
var contextId string

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
		for _, ctx := range contexts {
			structs.PrettyPrint(ctx)
		}
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
			Id:            uuid.New().String(),
			ContextBase64: encodedData,
		}

		services.AddContext(newContext)
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

	contextCmd.AddCommand(addContextCmd)
	addContextCmd.Flags().StringVarP(&email, "filepath", "f", "", "FilePath to the context you want to add")

	contextCmd.AddCommand(deleteContextCmd)
	deleteContextCmd.Flags().StringVarP(&contextId, "contextid", "c", "", "ContextId of the context")

	contextCmd.AddCommand(getContextCmd)
	getContextCmd.Flags().StringVarP(&contextId, "contextid", "c", "", "ContextId of the context")

	rootCmd.AddCommand(contextCmd)
}
