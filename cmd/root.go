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
var forceUpgrade bool
var resources []string

var cmdsWithoutContext = []string{
	"punq",
	"punq system reset-config",
	"punq changelog",
	"punq system ingress-controller-type",
	"punq install",
	"punq clean",
	"punq version",
	"punq system check",
}

var rootCmd = &cobra.Command{
	Use:   "punq",
	Short: "A slim open-source workload manager for Kubernetes with team collaboration, WebApp, and CLI. 🚀",
	Long:  `A slim open-source workload manager for Kubernetes with team collaboration, WebApp, and CLI. 🚀`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		versionFlag, err := cmd.Flags().GetBool("version")
		if err != nil {
			os.Exit(0)
		}
		if versionFlag {
			PrintVersion()
			os.Exit(0)
		}

		if cmd.CommandPath() != "punq system reset-config" {
			utils.InitConfigYaml(debug, customConfig, stage)
		}

		if !utils.ContainsEqual(cmdsWithoutContext, cmd.CommandPath()) {
			mokubernetes.InitKubernetes(utils.CONFIG.Kubernetes.RunInCluster)
			ctxs := mokubernetes.ListAllContexts()
			mokubernetes.ContextAddMany(ctxs)
			utils.PrintInfo((fmt.Sprintf("Current context: '%s'", contextId)))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintLogo()
		fmt.Println("")
		utils.PrintWelcomeMessage()
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
	rootCmd.PersistentFlags().BoolVarP(&cliVersion, "version", "v", false, "Print version info")
	rootCmd.PersistentFlags().StringVarP(&stage, "stage", "s", "", "Use different stage environment")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug information")
	rootCmd.PersistentFlags().StringVarP(&customConfig, "config", "y", "", "Use config from custom location")
	rootCmd.PersistentFlags().StringVarP(&contextId, "context-id", "c", "own-context", "Define a context-id")
}

func RequireStringFlag(str string, name string) {
	if str == "" {
		utils.FatalError(fmt.Sprintf("--%s flag is required for this command.", name))
	}
}
