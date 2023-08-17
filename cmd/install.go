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

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the application into your cluster without auto-removal.",
	Long: `
	This cmd installs the application permanently into you cluster. 
	Please run cleanup if you want to remove it again.`,
	Run: func(cmd *cobra.Command, args []string) {
		yellow := color.New(color.FgYellow).SprintFunc()
		if !utils.ConfirmTask(fmt.Sprintf("Do you realy want to install punq to '%s' context?", yellow(kubernetes.CurrentContextName())), 1) {
			os.Exit(0)
		}

		kubernetes.Deploy()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
