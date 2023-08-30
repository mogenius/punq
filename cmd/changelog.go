/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/mogenius/punq/utils"

	"github.com/spf13/cobra"
)

var changeLogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Print changelog information and exit.",
	Long:  `Print changelog information and exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintChangeLog()
	},
}

func init() {
	rootCmd.AddCommand(changeLogCmd)
}
