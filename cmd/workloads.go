package cmd

import (
	"github.com/mogenius/punq/kubernetes"
	"github.com/mogenius/punq/utils"
	"github.com/spf13/cobra"
)

var workloadCmd = &cobra.Command{
	Use:   "workloads",
	Short: "Manage kubernetes workloads.",
	Long:  `The workload command lets you manage all workloads on your cluster.`,
}

var listWorkloadsCmd = &cobra.Command{
	Use:   "list",
	Short: "List punq supported workloads.",
	Long:  `The workloads command lets you list all workloads managed by punq.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.InitConfigYaml(false, nil, false)

		kubernetes.ListWorkloads()
	},
}

var listTemplatesCmd = &cobra.Command{
	Use:   "list-templates",
	Short: "List punq supported templates to create workloads.",
	Long:  `The list-templates command lets you list all workloads which can be created punq using a template.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.InitConfigYaml(false, nil, false)

		kubernetes.ListTemplatesTerminal()
	},
}

func init() {
	workloadCmd.AddCommand(listWorkloadsCmd)
	workloadCmd.AddCommand(listTemplatesCmd)

	rootCmd.AddCommand(workloadCmd)
}
