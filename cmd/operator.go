package cmd

import (
	"github.com/mogenius/punq/services"

	"github.com/mogenius/punq/utils"

	"github.com/spf13/cobra"
)

var operatorCmd = &cobra.Command{
	Use:   "operator",
	Short: "Run the operator inside the cluster!",
	Long:  `Run the operator inside the cluster!`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.InitConfigYaml(true, nil, true)
		services.InitGin()
	},
}

func init() {
	rootCmd.AddCommand(operatorCmd)
}
