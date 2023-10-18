/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"runtime"

	"github.com/mogenius/punq/utils"
	"github.com/mogenius/punq/version"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information and exit.",
	Long:  `Print version information and exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		PrintVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&checkForUpdates, "checkUpdates", "u", false, "Check for punq updates.")
}

func PrintVersion() {
	utils.PrintLogo()
	fmt.Println("")
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("        CLI:            %s\n", yellow(version.Ver))
	fmt.Printf("        OperatorImage:  %s\n", yellow(version.OperatorImage))
	fmt.Printf("        Branch:         %s\n", yellow(version.Branch))
	fmt.Printf("        Commit:         %s\n", yellow(version.GitCommitHash))
	fmt.Printf("        Timestamp:      %s\n", yellow(version.BuildTimestamp))
	fmt.Printf("        Arch:           %s/%s\n", yellow(runtime.GOOS), yellow(runtime.GOARCH))
	fmt.Println("")

	if checkForUpdates {
		utils.IsNewReleaseAvailable()
		fmt.Println("")
	}
}
