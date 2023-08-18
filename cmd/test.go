/*
Copyright © 2022 mogenius, Benedikt Iltisberger
*/
package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mogenius/punq/kubernetes"
	"github.com/mogenius/punq/services"

	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "🚀🚀🚀",
	Long:  `test`,
	Run: func(cmd *cobra.Command, args []string) {
		// yellow := color.New(color.FgYellow).SprintFunc()
		// if !utils.ConfirmTask(fmt.Sprintf("Do you realy want to deploy punq to '%s' context?", yellow(kubernetes.CurrentContextName())), 1) {
		// 	os.Exit(0)
		// }

		utils.InitConfigYaml(true, nil, false)

		go services.InitGin()

		utils.OpenBrowser(fmt.Sprintf("http://%s:%s/punq", "0.0.0.0", utils.CONFIG.Browser.Port))

		//kubernetes.Deploy()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		logger.Log.Warning("CLEANUP Kubernetes resources ...")
		kubernetes.Remove()
		logger.Log.Info("CLEANUP finished successfully.")
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
