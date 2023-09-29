/*
Copyright © 2022 mogenius, Benedikt Iltisberger
*/
package cmd

import (
	"os"

	"github.com/mogenius/punq/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "All general system commands.",
}

var resetConfig = &cobra.Command{
	Use:   "reset-config",
	Short: "Remove all components from your cluster.",
	Long: `
	This cmd removes all remaining parts of the daemonset, configs, etc. from your cluster. 
	This can be used if something went wrong during automatic cleanup.`,
	Run: func(cmd *cobra.Command, args []string) {
		yellow := color.New(color.FgYellow).SprintFunc()
		if !utils.ConfirmTask(yellow("Do you realy want to reset your configuration file to default?"), 1) {
			os.Exit(0)
		}
		utils.DeleteCurrentConfig()
	},
}

func init() {
	rootCmd.AddCommand(systemCmd)
	systemCmd.AddCommand(resetConfig)
}
