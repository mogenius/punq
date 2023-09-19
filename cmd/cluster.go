/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/services"
	"github.com/mogenius/punq/structs"
	"github.com/spf13/cobra"

	"github.com/mogenius/punq/kubernetes"
)

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Print information and exit.",
	Long:  `Print information and exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		if contextId == "" {
			logger.Log.Fatal("contextId cannot be empty.")
		}

		// init contexts
		contexts := services.ListContexts()
		kubernetes.ContextUpdateLocalCache(contexts)

		fmt.Println("Cluster information for context: " + contextId)
		structs.PrettyPrint(kubernetes.ContextForId(contextId))
	},
}

var clusterInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print information and exit.",
	Long:  `Print information and exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		// init contexts
		contexts := services.ListContexts()
		kubernetes.ContextUpdateLocalCache(contexts)

		fmt.Println("Cluster information for context: " + contextId)
		structs.PrettyPrint(kubernetes.ClusterInfo(&contextId))
	},
}

func init() {
	rootCmd.AddCommand(clusterCmd)
	clusterCmd.Flags().StringVarP(&contextId, "context", "c", "", "Define a context")

	clusterCmd.AddCommand(clusterInfoCmd)
}
