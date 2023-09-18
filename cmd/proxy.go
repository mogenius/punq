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

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "🚀🚀🚀 Port-Forward to the application within your currently selected kubernetes context. 🚀🚀🚀",
	Long: `
	Port-Forward to punq within your currently selected kubernetes context.`,
	Run: func(cmd *cobra.Command, args []string) {
		yellow := color.New(color.FgYellow).SprintFunc()
		if !utils.ConfirmTask(fmt.Sprintf("Do you realy want to connect to punq in '%s' context?", yellow(kubernetes.CurrentContextName())), 1) {
			os.Exit(0)
		}

		// FORWARD BACKEND
		readyBackendCh := make(chan struct{})
		stopBackendCh := make(chan struct{}, 1)
		go kubernetes.StartPortForward(utils.CONFIG.Backend.Port, utils.CONFIG.Backend.Port, readyBackendCh, stopBackendCh, &contextId)

		// FORWARD FRONTEND
		readyFrontendCh := make(chan struct{})
		stopFrontendCh := make(chan struct{}, 1)
		go kubernetes.StartPortForward(utils.CONFIG.Frontend.Port, utils.CONFIG.Frontend.Port, readyFrontendCh, stopFrontendCh, &contextId)

		select {
		case <-readyBackendCh:
			fmt.Printf("Backend %s:%d is ready! 🚀🚀🚀\n", utils.CONFIG.Backend.Host, utils.CONFIG.Backend.Port)
			break
		case <-stopBackendCh:
			break
		}

		select {
		case <-readyFrontendCh:
			url := fmt.Sprintf("http://%s:%d", utils.CONFIG.Frontend.Host, utils.CONFIG.Frontend.Port)
			fmt.Printf("Frontend %s:%d is ready! 🚀🚀🚀\n", utils.CONFIG.Frontend.Host, utils.CONFIG.Frontend.Port)
			utils.OpenBrowser(url)
			fmt.Printf("We opened a browser with %s for you.\n", url)
			break
		case <-stopFrontendCh:
			break
		}

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)
}
