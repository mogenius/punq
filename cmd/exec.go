/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/mogenius/punq/kubernetes"

	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Exec a command in a container.",
	Long:  `Exec a command in a container.`,
	Run: func(cmd *cobra.Command, args []string) {
		kubernetes.ExecTest()
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
