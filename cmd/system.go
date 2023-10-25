/*
Copyright © 2022 mogenius, Benedikt Iltisberger
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mogenius/punq/kubernetes"
	"github.com/mogenius/punq/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "All general system commands.",
}

var resetConfig = &cobra.Command{
	Use:   "reset-config",
	Short: "Remove all components from your cluster.",
	Long: `
	This cmd removes all remaining parts of the daemonset, configs, etc. from your cluster. 
	This can be used if something went wrong during automatic cleanup.`,
	Run: func(cmd *cobra.Command, args []string) {
		yellow := color.New(color.FgYellow).SprintFunc()
		if !utils.ConfirmTask(yellow("Do you really want to reset your configuration file to default?")) {
			os.Exit(0)
		}
		utils.DeleteCurrentConfig()
	},
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the system for all required components and offer healing.",
	Run: func(cmd *cobra.Command, args []string) {
		// check internet access
		// check for kubectl
		// check kubernetes version
		// check for ingresscontroller
		// check for metrics server
		// check for helm
		// check for cluster provider
		// check for api versions

		contextName := kubernetes.CurrentContextName()

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Check", "Status", "Message"})
		// check internet access
		inetResult, inetErr := utils.CheckInternetAccess()
		t.AppendRow(
			table.Row{"Internet Access", utils.StatusEmoji(inetResult), StatusMessage(inetErr, "Check your internet connection.", "Internet access works.")},
		)
		t.AppendSeparator()

		// check for punq
		punqInstalledVersion, punqInstalledErr := kubernetes.IsPunqInstalled()
		if punqInstalledErr != nil {
			punqInstalledErr = fmt.Errorf("punq is not installed in context '%s'", contextName)
		}
		t.AppendRow(
			table.Row{"punq installed", utils.StatusEmoji(punqInstalledVersion != ""), StatusMessage(punqInstalledErr, "Please run 'punq install -i punq.localhost'", fmt.Sprintf("Version '%s' in '%s' found.", punqInstalledVersion, contextName))},
		)
		t.AppendSeparator()

		// check for kubectl
		kubectlResult, kubectlOutput, kubectlErr := utils.IsKubectlInstalled()
		t.AppendRow(
			table.Row{"kubectl", utils.StatusEmoji(kubectlResult), StatusMessage(kubectlErr, "Plase install kubectl (https://kubernetes.io/docs/tasks/tools/) on your system to proceed.", kubectlOutput)},
		)
		t.AppendSeparator()

		// check kubernetes version
		kubernetesVersion := kubernetes.KubernetesVersion(nil)
		kubernetesVersionResult := kubernetesVersion != nil
		t.AppendRow(
			table.Row{"Kubernetes Version", utils.StatusEmoji(kubernetesVersionResult), StatusMessage(kubectlErr, "Cannot determin version of kubernetes.", fmt.Sprintf("Version: %s\nPlatform: %s", kubernetesVersion.String(), kubernetesVersion.Platform))},
		)
		t.AppendSeparator()

		// check for ingresscontroller
		ingressType, ingressTypeErr := kubernetes.DetermineIngressControllerType(nil)
		t.AppendRow(
			table.Row{"Ingress Controller", utils.StatusEmoji(ingressTypeErr == nil), StatusMessage(ingressTypeErr, "Cannot determin ingress controller type.", ingressType.String())},
		)
		t.AppendSeparator()

		// check for metrics server
		metricsResult, metricsVersion, metricsErr := kubernetes.IsMetricsServerAvailable(nil)
		t.AppendRow(
			table.Row{"Metrics Server", utils.StatusEmoji(metricsResult), StatusMessage(metricsErr, "kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml\nNote: Running docker-desktop? Please add '- --kubelet-insecure-tls' to the args sction in the deployment of metrics-server.", metricsVersion)},
		)
		t.AppendSeparator()

		// check for helm
		helmResult, helmOutput, helmErr := utils.IsHelmInstalled()
		t.AppendRow(
			table.Row{"Helm", utils.StatusEmoji(helmResult), StatusMessage(helmErr, "Plase install helm (https://helm.sh/docs/intro/install/) on your system to proceed.", helmOutput)},
		)
		t.AppendSeparator()

		// check cluster provider
		clusterProvOutput, clusterProvErr := kubernetes.GuessClusterProvider(nil)
		t.AppendRow(
			table.Row{"Cluster Provider", utils.StatusEmoji(clusterProvErr == nil), StatusMessage(clusterProvErr, "We could not determine the provider of this cluster.", string(clusterProvOutput))},
		)
		t.AppendSeparator()

		// API Versions
		apiVerResult, apiVerErr := kubernetes.ApiVersions(nil)
		apiVersStr := ""
		for _, entry := range apiVerResult {
			apiVersStr += fmt.Sprintf("%s\n", entry)
		}
		apiVersStr = strings.TrimRight(apiVersStr, "\n\r")
		t.AppendRow(
			table.Row{"API Versions", utils.StatusEmoji(len(apiVerResult) > 0), StatusMessage(apiVerErr, "Cannot determin API versions.", apiVersStr)},
		)
		t.Render()
	},
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print information and exit.",
	Long:  `Print information and exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintSettings()
	},
}

var ingressControllerCmd = &cobra.Command{
	Use:   "ingress-controller-type",
	Short: "Print ingress-controller-type and exit.",
	Long:  `Print ingress-controller-type and exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		ingressType, err := kubernetes.DetermineIngressControllerType(nil)
		if err != nil {
			utils.PrintError(err.Error())
		}
		utils.PrintInfo(fmt.Sprintf("Ingress Controller Type: %s", ingressType.String()))
	},
}

func init() {
	rootCmd.AddCommand(systemCmd)
	systemCmd.AddCommand(resetConfig)
	systemCmd.AddCommand(infoCmd)
	systemCmd.AddCommand(ingressControllerCmd)
	systemCmd.AddCommand(checkCmd)
}

// UTILS

func StatusMessage(err error, solution string, successMsg string) string {
	if err != nil {
		return fmt.Sprintf("Error: %s\nSolution: %s", err.Error(), solution)
	}
	return successMsg
}
