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
		vers, err := utils.CurrentReleaseVersion()
		if err != nil {
			utils.PrintError(err.Error())
			return
		}

		utils.PrintInfo(fmt.Sprintf("\nYour version:    v%s", version.Ver))
		utils.PrintInfo(fmt.Sprintf("Current version: %s", vers))

		if vers == version.Ver {
			utils.PrintInfo("You are already on the latest version.")
			return
		} else {
			yellow := color.New(color.FgYellow).SprintFunc()
			if !utils.ConfirmTask(fmt.Sprintf("Do you really want to upgrade to image '%s' in context '%s'?", yellow(vers), yellow(kubernetes.CurrentContextName())), 1) {
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
	rootCmd.AddCommand(updateOperatorImageCmd)
}
