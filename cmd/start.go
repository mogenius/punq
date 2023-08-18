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

	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
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

		kubernetes.Deploy()
		utils.OpenBrowser(fmt.Sprintf("http://%s:%s/punq", os.Getenv("API_HOST"), os.Getenv("API_PORT")))

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		logger.Log.Warning("CLEANUP Kubernetes resources ...")
		kubernetes.Remove()
		logger.Log.Info("CLEANUP finished successfully.")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports start flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
