package cmd

import (
	"os"

	cc "github.com/ivanpirog/coloredcobra"
	"github.com/mogenius/punq/utils"
	"github.com/spf13/cobra"
)

var useClusterConfig bool
var debug bool
var customConfig string
var namespace string
var resource string
var checkForUpdates bool
var email string
var password string
var displayName string
var userId string
var showPasswords bool
var ingressHostname string
var filePath string
var contextId string
var accessLevel string

var rootCmd = &cobra.Command{
	Use:   "punq",
	Short: "Collect traffic data using pcap from a machine.",
	Long:  `Use punq to manage the workloads of your kubernetes clusters relatively neat. 🚀`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		utils.InitConfigYaml(debug, customConfig, useClusterConfig)
	},
}

func Execute() {

	cc.Init(&cc.Config{
		RootCmd:  rootCmd,
		Headings: cc.HiCyan + cc.Bold + cc.Underline,
		Commands: cc.HiYellow + cc.Bold,
		Example:  cc.Italic,
		ExecName: cc.Bold,
		Flags:    cc.Bold,
	})

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&useClusterConfig, "use-cluster-config", "x", false, "Load different default config to run in cluster")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug information")
	rootCmd.PersistentFlags().StringVarP(&customConfig, "config", "y", "", "Use config from custom location")
}
