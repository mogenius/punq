/*
Copyright © 2022 mogenius, Benedikt Iltisberger
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/mogenius/punq/kubernetes"
	"github.com/mogenius/punq/utils"
	"github.com/mogenius/punq/version"

	"github.com/spf13/cobra"
)

var updateOperatorImageCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the container image of the punq operator.",
	Long:  `Upgrade the container image of the punq operator within your currently selected kubernetes context.`,
	Run: func(cmd *cobra.Command, args []string) {
		var vers string
		var err error
		if version.Branch == "main" {
			vers, err = utils.CurrentReleaseVersion()
		} else {
			vers, err = utils.CurrentPreReleaseVersion()
		}

		if err != nil {
			utils.PrintError(err.Error())
			return
		}

		operatorVer, err := kubernetes.GetCurrentOperatorVersion()
		if err != nil {
			utils.PrintError(err.Error())
			return
		}

		utils.PrintInfo(fmt.Sprintf("\nYour version:    %s", operatorVer))
		utils.PrintInfo(fmt.Sprintf("Current version: %s", vers))

		if vers == operatorVer && !forceUpgrade {
			utils.PrintInfo("You are already on the latest version.")
			return
		} else {
			yellow := color.New(color.FgYellow).SprintFunc()
			if !utils.ConfirmTask(fmt.Sprintf("Do you really want to upgrade to image '%s' in context '%s'?", yellow(vers), yellow(kubernetes.CurrentContextName()))) {
				os.Exit(0)
			}

			imageName := fmt.Sprintf("ghcr.io/mogenius/punq:%s", vers)
			err = kubernetes.UpdateDeploymentImage(utils.CONFIG.Kubernetes.OwnNamespace, version.Name, imageName, nil)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
			}
			fmt.Printf("✅ Deployment %s/%s updated to image %s\n", utils.CONFIG.Kubernetes.OwnNamespace, version.Name, vers)
		}
	},
}

func init() {
	updateOperatorImageCmd.PersistentFlags().BoolVarP(&forceUpgrade, "force", "f", false, "Force upgrade of deployment image.")
	rootCmd.AddCommand(updateOperatorImageCmd)
}
