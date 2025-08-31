package cmd

import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
)

//go:embed shell/init.bash
var bashInit string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate shell integration code",
	Long:  "Generate the shell integration code to be evaluated in bash",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print(bashInit)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
