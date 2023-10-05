/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/mogenius/punq/operator"
	"github.com/mogenius/punq/services"
	"github.com/mogenius/punq/utils"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Run punq from your local machine.",
	Long:  `Run punq from your local machine in your current-context in kubernetes.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintLogo()
		utils.PrintSettings()

		contexts := services.ListContexts()
		utils.PrintInfo(fmt.Sprintf("Initialized operator with %d contexts.", len(contexts)))

		go operator.InitFrontend()
		utils.OpenBrowser(fmt.Sprintf("http://%s:%d", utils.CONFIG.Frontend.Host, utils.CONFIG.Frontend.Port))
		operator.InitBackend()
	},
}

func init() {
	localCmd.Hidden = true
	rootCmd.AddCommand(localCmd)
}
