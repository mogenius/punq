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

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information and exit.",
	Long:  `Print version information and exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintLogo()
		yellow := color.New(color.FgYellow).SprintFunc()
		fmt.Printf("        CLI:       %s\n", yellow(version.Ver))
		fmt.Printf("        Operator:  %s\n", yellow(version.Operator))
		fmt.Printf("        Branch:    %s\n", yellow(version.Branch))
		fmt.Printf("        Commit:    %s\n", yellow(version.GitCommitHash))
		fmt.Printf("        Timestamp: %s\n", yellow(version.BuildTimestamp))
		fmt.Printf("        Arch:      %s/%s\n", yellow(runtime.GOOS), yellow(runtime.GOARCH))
		fmt.Println("")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
