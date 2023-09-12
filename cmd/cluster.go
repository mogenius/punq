/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/structs"
	"github.com/spf13/cobra"

	"github.com/mogenius/punq/kubernetes"
)

// versionCmd represents the version command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Print information and exit.",
	Long:  `Print information and exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		if contextId == "" {
			logger.Log.Fatal("contextId cannot be empty.")
		}
		fmt.Println("Cluster information for context: " + contextId)
		structs.PrettyPrint(kubernetes.ContextForId(contextId))
	},
}

func init() {
	rootCmd.AddCommand(clusterCmd)
	clusterCmd.Flags().StringVarP(&contextId, "context", "c", "", "Define a context")
}
