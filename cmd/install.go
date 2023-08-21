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

var ingressHostname string

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the application into your cluster without auto-removal.",
	Long: `
	This cmd installs the application permanently into you cluster. 
	Please run cleanup if you want to remove it again.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.InitConfigYaml(true, nil, false)
		yellow := color.New(color.FgYellow).SprintFunc()
		if !utils.ConfirmTask(fmt.Sprintf("Do you realy want to install punq to '%s' context?", yellow(kubernetes.CurrentContextName())), 1) {
			os.Exit(0)
		}

		kubernetes.Deploy(yellow(kubernetes.CurrentContextName()), ingressHostname)
	},
}

func init() {
	installCmd.Flags().StringVarP(&ingressHostname, "ingress", "i", "", "Ingress hostname for operator.")
	rootCmd.AddCommand(installCmd)
}
