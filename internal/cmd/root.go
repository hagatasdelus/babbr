package cmd

import (
	"github.com/hagatasdelus/babbr/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "babbr",
	Short:   "Fish shell-style abbreviations for bash",
	Long:    "A tool that provides fish shell-style abbreviation functionality for bash shell",
	Version: version.Version,
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}
