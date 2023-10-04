/*
Copyright © 2022 mogenius, Benedikt Iltisberger
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mogenius/punq/kubernetes"
	"github.com/mogenius/punq/version"

	"github.com/mogenius/punq/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var updateOperatorImageCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the container image of the punq operator within your currently selected kubernetes context.",
	Long: `
	Upgrade container image of the punq operator within your currently selected kubernetes context..`,
	Run: func(cmd *cobra.Command, args []string) {
		if utils.CONFIG.Misc.Stage == "prod" {
			utils.FatalError("You are running in production mode. Please switch to a development mode first.")
		}

		currentVersionCmd := exec.Command("sh", "-c", "git describe --tags $(git rev-list --tags --max-count=1)")

		output, err := currentVersionCmd.CombinedOutput()
		vers := strings.TrimSpace(string(output))
		if err != nil {
			utils.PrintError(fmt.Sprintf("Failed to execute command (%s): %v", vers, err))
			utils.PrintError(string(output))
			return
		}

		yellow := color.New(color.FgYellow).SprintFunc()
		if !utils.ConfirmTask(fmt.Sprintf("Do you really want to upgrade to image '%s' in context '%s'?", yellow(vers), yellow(kubernetes.CurrentContextName())), 1) {
			os.Exit(0)
		}

		imageName := fmt.Sprintf("ghcr.io/mogenius/punq-dev:%s", vers)

		err = kubernetes.UpdateDeploymentImage(utils.CONFIG.Kubernetes.OwnNamespace, version.Name, imageName, nil)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		}
		fmt.Printf("✅ Deployment %s/%s updated to image %s\n", utils.CONFIG.Kubernetes.OwnNamespace, version.Name, vers)
	},
}

func init() {
	updateOperatorImageCmd.Hidden = true
	rootCmd.AddCommand(updateOperatorImageCmd)
}
