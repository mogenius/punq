package cmd

import (
	"fmt"
	"os"

	cc "github.com/ivanpirog/coloredcobra"
	mokubernetes "github.com/mogenius/punq/kubernetes"
	"github.com/mogenius/punq/utils"
	"github.com/spf13/cobra"
)

var cliVersion bool
var stage string
var debug bool
var customConfig string
var namespace string
var resource string
var checkForUpdates bool
var email string
var password string
var displayName string
var userId string
var ingressHostname string
var filePath string
var contextId string
var accessLevel string

var rootCmd = &cobra.Command{
	Use:   "punq",
	Short: "Collect traffic data using pcap from a machine.",
	Long:  `Use punq to manage the workloads of your kubernetes clusters relatively neat. 🚀`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.CommandPath() != "punq system reset-config" {
			utils.InitConfigYaml(debug, customConfig, stage)
			mokubernetes.Init(utils.CONFIG.Kubernetes.RunInCluster)

			if contextId != "" {
				ctxs := mokubernetes.ListAllContexts()
				mokubernetes.ContextAddMany(ctxs)
			}
			utils.PrintInfo((fmt.Sprintf("Selected context '%s'.", contextId)))
		}
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

	if cliVersion {
		PrintVersion()
		os.Exit(0)
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&cliVersion, "version", "v", false, "Print version info")
	rootCmd.PersistentFlags().StringVarP(&stage, "stage", "s", "", "Use different stage environment")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug information")
	rootCmd.PersistentFlags().StringVarP(&customConfig, "config", "y", "", "Use config from custom location")
	rootCmd.PersistentFlags().StringVarP(&contextId, "context", "c", "own-context", "Define a contextId")
}
