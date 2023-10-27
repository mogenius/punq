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
		if !utils.ConfirmTask(fmt.Sprintf("Do you really want to connect to punq in '%s' context?", yellow(kubernetes.CurrentContextName()))) {
			os.Exit(0)
		}

		// FORWARD BACKEND
		readyBackendCh := make(chan struct{})
		stopBackendCh := make(chan struct{}, 1)
		backendUrl := fmt.Sprintf("http://%s:%d", utils.CONFIG.Backend.Host, utils.CONFIG.Backend.Port)
		go kubernetes.StartPortForward(utils.CONFIG.Backend.Port, utils.CONFIG.Backend.Port, readyBackendCh, stopBackendCh, &contextId)

		// FORWARD FRONTEND
		readyFrontendCh := make(chan struct{})
		stopFrontendCh := make(chan struct{}, 1)
		frontendUrl := fmt.Sprintf("http://%s:%d", utils.CONFIG.Frontend.Host, utils.CONFIG.Frontend.Port)
		go kubernetes.StartPortForward(utils.CONFIG.Frontend.Port, utils.CONFIG.Frontend.Port, readyFrontendCh, stopFrontendCh, &contextId)

		// FORWARD WEBSOCKET
		readyWebsocketCh := make(chan struct{})
		stopWebsocketCh := make(chan struct{}, 1)
		websocketUrl := fmt.Sprintf("ws://%s:%d", utils.CONFIG.Websocket.Host, utils.CONFIG.Websocket.Port)
		go kubernetes.StartPortForward(utils.CONFIG.Websocket.Port, utils.CONFIG.Websocket.Port, readyWebsocketCh, stopWebsocketCh, &contextId)

		select {
		case <-readyBackendCh:
			fmt.Printf("Backend %s is ready! 🚀🚀🚀\n", backendUrl)
			break
		case <-stopBackendCh:
			break
		}

		select {
		case <-readyFrontendCh:
			fmt.Printf("Frontend %s is ready! 🚀🚀🚀\n", frontendUrl)
			utils.OpenBrowser("http://localhost:8888")
			break
		case <-stopFrontendCh:
			break
		}

		select {
		case <-readyWebsocketCh:
			fmt.Printf("Websocket %s is ready! 🚀🚀🚀\n", websocketUrl)
			break
		case <-stopWebsocketCh:
			break
		}

		kubernetes.Proxy(backendUrl, frontendUrl, websocketUrl)

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)
}
