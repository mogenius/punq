/*
Copyright © 2022 mogenius, Benedikt Iltisberger
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/mogenius/punq/services"

	"github.com/mogenius/punq/kubernetes"
	"github.com/mogenius/punq/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove all punq components from your cluster.",
	Long: `
	This cmd removes all remaining parts of the daemonset, configs, etc. from your cluster. 
	This can be used if something went wrong during automatic setup/cleanup.`,
	Run: func(cmd *cobra.Command, args []string) {
		yellow := color.New(color.FgYellow).SprintFunc()
		if !utils.ConfirmTask(fmt.Sprintf("Do you really want to remove punq from '%s' context?", yellow(kubernetes.CurrentContextName()))) {
			os.Exit(0)
		}

		clusterName := kubernetes.CurrentContextName()

		kubernetes.Remove(yellow(clusterName))
		services.RemoveKeyPair()

		fmt.Printf("\n🚀🚀🚀 Successfully uninstalled punq from '%s'.\n\n", clusterName)
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
