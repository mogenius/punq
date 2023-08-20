/*
Copyright © 2022 mogenius, Benedikt Iltisberger
*/
package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"
	"github.com/mogenius/punq/kubernetes"
	"github.com/mogenius/punq/services"

	"github.com/mogenius/punq/utils"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "🚀🚀🚀",
	Long:  `test`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.InitConfigYaml(true, nil, false)
		go services.InitGin()

		yellow := color.New(color.FgYellow).SprintFunc()
		if !utils.ConfirmTask(fmt.Sprintf("Do you realy want to deploy punq to '%s' context?", yellow(kubernetes.CurrentContextName())), 1) {
			os.Exit(0)
		}

		kubernetes.Deploy(yellow(kubernetes.CurrentContextName()))

		//utils.OpenBrowser(fmt.Sprintf("http://%s:%s/punq", "0.0.0.0", utils.CONFIG.Browser.Port))

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
