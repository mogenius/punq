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

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the application into your cluster without auto-removal.",
	Long: `
	This cmd installs the application permanently into you cluster. 
	Please run cleanup if you want to remove it again.`,
	Run: func(cmd *cobra.Command, args []string) {
		yellow := color.New(color.FgYellow).SprintFunc()
		if !utils.ConfirmTask(fmt.Sprintf("Do you really want to install punq to '%s' context?", yellow(kubernetes.CurrentContextName()))) {
			os.Exit(0)
		}

		clusterName := kubernetes.CurrentContextName()

		kubernetes.Deploy(clusterName, ingressHostname)
		services.InitUserService()
		services.InitAuthService()
		services.CreateAdminUser()

		fmt.Printf("\n🚀🚀🚀 Successfully installed punq in '%s'.\n\n", clusterName)
	},
}

func init() {
	installCmd.Flags().StringVarP(&ingressHostname, "ingress", "i", "", "Ingress hostname for operator.")
	rootCmd.AddCommand(installCmd)
}
