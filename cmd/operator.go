package cmd

import (
	"github.com/mogenius/punq/operator"

	"github.com/mogenius/punq/utils"

	"github.com/spf13/cobra"
)

var operatorCmd = &cobra.Command{
	Use:   "operator",
	Short: "Run the operator inside the cluster!",
	Long:  `Run the operator inside the cluster!`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.InitConfigYaml(true, nil, true)
		println("\n###############################################")
		utils.IsNewReleaseIsAvailable()
		println("###############################################\n")
		operator.InitGin()
	},
}

func init() {
	rootCmd.AddCommand(operatorCmd)
}
