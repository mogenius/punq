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

		readyCh := make(chan struct{})
		stopCh := make(chan struct{}, 1)
		go kubernetes.StartPortForward(utils.CONFIG.Browser.Port, utils.CONFIG.Kubernetes.ContainerPort, readyCh, stopCh, &contextId)

		url := fmt.Sprintf("http://%s:%d/punq", utils.CONFIG.Browser.Host, utils.CONFIG.Browser.Port)

		select {
		case <-readyCh:
			utils.OpenBrowser(url)
			fmt.Printf("We opened a browser with %s for you.\n", url)
			break
		case <-stopCh:
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
