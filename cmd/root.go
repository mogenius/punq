package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	cc "github.com/ivanpirog/coloredcobra"
	mokubernetes "github.com/mogenius/punq/kubernetes"
	"github.com/mogenius/punq/utils"
	"github.com/spf13/cobra"
)

var resetConfig bool
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
		if resetConfig {
			utils.DeleteCurrentConfig()
		}
		utils.InitConfigYaml(debug, customConfig, stage)
		mokubernetes.Init(utils.CONFIG.Kubernetes.RunInCluster)
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
	rootCmd.PersistentFlags().StringVarP(&stage, "stage", "s", "", "Use different stage environment")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug information")
	rootCmd.PersistentFlags().BoolVarP(&resetConfig, "reset-config", "k", false, "Delete the current config and replace it with the default one")
	rootCmd.PersistentFlags().StringVarP(&customConfig, "config", "y", "", "Use config from custom location")
}

func FatalError(message string) {
	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf(red("Error: %s\n"), message)
	os.Exit(0)
}

func PrintError(message string) {
	red := color.New(color.FgRed).SprintFunc()
	fmt.Println(red(message))
}

func PrintInfo(message string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Println(yellow(message))
}
