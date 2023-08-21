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

	"github.com/mogenius/punq/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "🚀🚀🚀 Run the application within your currently selected kubernetes context. 🚀🚀🚀",
	Long: `
	Run the application within your currently selected kubernetes context.
	App will cleanup after being terminated with CTRL+C automatically.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.InitConfigYaml(true, nil, false)
		yellow := color.New(color.FgYellow).SprintFunc()
		if !utils.ConfirmTask(fmt.Sprintf("Do you realy want to deploy punq to '%s' context?", yellow(kubernetes.CurrentContextName())), 1) {
			os.Exit(0)
		}

		kubernetes.Deploy(yellow(kubernetes.CurrentContextName()), ingressHostname)
		utils.OpenBrowser(fmt.Sprintf("http://%s:%s/punq", os.Getenv("API_HOST"), os.Getenv("API_PORT")))

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		kubernetes.Remove(yellow(kubernetes.CurrentContextName()))
	},
}

func init() {
	startCmd.Flags().StringVarP(&ingressHostname, "ingress", "i", "", "Ingress hostname for operator.")
	rootCmd.AddCommand(startCmd)
}
