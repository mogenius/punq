/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/mogenius/punq/utils"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print information and exit.",
	Long:  `Print information and exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintSettings()
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
