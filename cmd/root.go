package cmd

import (
	"os"

	cc "github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"
)

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
	Use:   "github.com/mogenius/punq",
	Short: "Collect traffic data using pcap from a machine.",
	Long:  `Use punq to manage the workloads of your kubernetes clusters relatively neat. 🚀`,
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
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
