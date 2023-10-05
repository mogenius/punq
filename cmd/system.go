/*
Copyright © 2022 mogenius, Benedikt Iltisberger
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/mogenius/punq/kubernetes"
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
		if !utils.ConfirmTask(yellow("Do you really want to reset your configuration file to default?"), 1) {
			os.Exit(0)
		}
		utils.DeleteCurrentConfig()
	},
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print information and exit.",
	Long:  `Print information and exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintSettings()
	},
}

var ingressControllerCmd = &cobra.Command{
	Use:   "ingress-controller-type",
	Short: "Print ingress-controller-type and exit.",
	Long:  `Print ingress-controller-type and exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		ingressType, err := kubernetes.DetermineIngressControllerType(nil)
		if err != nil {
			utils.PrintError(err.Error())
		}
		utils.PrintInfo(fmt.Sprintf("Ingress Controller Type: %s", ingressType.String()))
	},
}

func init() {
	rootCmd.AddCommand(systemCmd)
	systemCmd.AddCommand(resetConfig)
	systemCmd.AddCommand(infoCmd)
	systemCmd.AddCommand(ingressControllerCmd)
}
