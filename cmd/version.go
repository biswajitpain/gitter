package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// This variable will be set at build time
var version = "dev"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of gitter",
	Long:  `All software has versions. This is gitter's.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gitter version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

