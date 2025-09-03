package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "babbr",
	Short: "Fish shell-style abbreviations for bash",
	Long:  "A tool that provides fish shell-style abbreviation functionality for bash shell",
}

func GetRootCmd() *cobra.Command {
	rootCmd.SetVersionTemplate("babbr {{.Version}}\n")
	return rootCmd
}

func SetVersion(version string) {
	rootCmd.Version = version
}
