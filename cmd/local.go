/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/mogenius/punq/logger"
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
		logger.Log.Noticef("Initialized operator with %d contexts.", len(contexts))

		utils.OpenBrowser(fmt.Sprintf("http://%s:%d/punq", utils.CONFIG.Browser.Host, utils.CONFIG.Browser.Port))

		operator.InitGin()
	},
}

func init() {
	rootCmd.AddCommand(localCmd)
}
