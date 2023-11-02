/*
Copyright © 2022 mogenius, Benedikt Iltisberger
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/mogenius/punq/kubernetes"
	"github.com/mogenius/punq/utils"
	"github.com/mogenius/punq/version"

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
		if !utils.ConfirmTask(yellow("Do you really want to reset your configuration file to default?")) {
			os.Exit(0)
		}
		utils.DeleteCurrentConfig()
	},
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the system for all required components and offer healing.",
	Run: func(cmd *cobra.Command, args []string) {
		contextName := kubernetes.CurrentContextName()

		// check for punq
		punqInstalledVersion, punqInstalledErr := kubernetes.IsDeploymentInstalled(utils.CONFIG.Kubernetes.OwnNamespace, version.Name)
		if punqInstalledErr != nil {
			utils.FatalError(fmt.Sprintf("%s is not installed in context '%s'.\nPlease switch context or run 'punq install -i punq.localhost'", version.Name, contextName))
		}
		utils.PrintInfo(fmt.Sprintf("Found version '%s' of %s in '%s'.", punqInstalledVersion, version.Name, contextName))

		kubernetes.SystemCheck()
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
	systemCmd.AddCommand(checkCmd)
}
