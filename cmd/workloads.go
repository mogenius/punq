package cmd

import (
	"fmt"

	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/kubernetes"
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
		kubernetes.ListWorkloadsOnTerminal(dtos.ADMIN)
	},
}

var listTemplatesCmd = &cobra.Command{
	Use:   "list-templates",
	Short: "List punq supported templates to create workloads.",
	Long:  `The list-templates command lets you list all workloads which can be created punq using a template.`,
	Run: func(cmd *cobra.Command, args []string) {
		kubernetes.ListTemplatesTerminal()
	},
}

var podsCmd = &cobra.Command{
	Use:   "pods",
	Short: "pods related commands.",
	Long:  `Similar to kubectl, punq can list workloads in an orderly fashion.`,
}
var podsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pods.",
	Long:  `Similar to kubectl, punq can list workloads in an orderly fashion.`,
	Run: func(cmd *cobra.Command, args []string) {
		kubernetes.ListPodsTerminal(namespace, &contextId)
	},
}
var podsDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe pod.",
	Long:  `Similar to kubectl, punq can describe workloads in an orderly fashion.`,
	Run: func(cmd *cobra.Command, args []string) {
		RequireStringFlag(namespace, "namespace")
		RequireStringFlag(resource, "resource")
		RequireStringFlag(contextId, "context-id")

		wl := kubernetes.DescribeK8sPod(namespace, resource, &contextId)
		fmt.Println(wl.Result)
	},
}

var podDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete pod.",
	Long:  `Similar to kubectl, punq can delete workloads in an orderly fashion.`,
	Run: func(cmd *cobra.Command, args []string) {
		RequireStringFlag(namespace, "namespace")
		RequireStringFlag(resource, "resource")
		RequireStringFlag(contextId, "context-id")

		pod := kubernetes.GetPod(namespace, resource, &contextId)
		if pod != nil {
			kubernetes.DeleteK8sPod(*pod, &contextId)
		} else {
			fmt.Printf("Pod %s/%s not found.\n", namespace, resource)
		}
	},
}

func init() {
	workloadCmd.AddCommand(listWorkloadsCmd)
	workloadCmd.AddCommand(listTemplatesCmd)

	workloadCmd.AddCommand(podsCmd)
	podsCmd.AddCommand(podsListCmd)
	podsListCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Define a namespace")
	podsListCmd.Flags().StringVarP(&resource, "resource", "r", "", "Define a resource name")
	podsCmd.AddCommand(podsDescribeCmd)
	podsDescribeCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Define a namespace")
	podsDescribeCmd.Flags().StringVarP(&resource, "resource", "r", "", "Define a resource name")

	podsCmd.AddCommand(podDeleteCmd)
	podDeleteCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Define a namespace")
	podDeleteCmd.Flags().StringVarP(&resource, "resource", "r", "", "Define a resource name")

	rootCmd.AddCommand(workloadCmd)
}
