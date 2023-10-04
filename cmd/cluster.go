/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/mogenius/punq/services"
	"github.com/mogenius/punq/structs"
	"github.com/spf13/cobra"

	"github.com/mogenius/punq/kubernetes"
)

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Cluster related tasks.",
	Long:  `Cluster related tasks.`,
}

var clusterContextCmd = &cobra.Command{
	Use:   "context",
	Short: "Print context information and exit.",
	Long:  `Print context information and exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		// init contexts
		kubernetes.ContextAddMany(services.ListContexts())

		structs.PrettyPrint(kubernetes.ContextForId(contextId))
	},
}

var clusterInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print information and exit.",
	Long:  `Print information and exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		// init contexts
		kubernetes.ContextAddMany(services.ListContexts())

		structs.PrettyPrint(kubernetes.ClusterInfo(&contextId))
	},
}

func init() {
	rootCmd.AddCommand(clusterCmd)
	clusterCmd.AddCommand(clusterContextCmd)
	clusterCmd.AddCommand(clusterInfoCmd)
}
