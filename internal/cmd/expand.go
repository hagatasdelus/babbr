package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hagatasdelus/babbr/internal/config"
	"github.com/hagatasdelus/babbr/internal/expand"
)

var expandCmd = &cobra.Command{
	Use:   "expand",
	Short: "(Internal) Expand an abbreviation based on buffer content",
	Long:  "Internal command used by the shell integration to expand abbreviations",
	RunE: func(cmd *cobra.Command, args []string) error {
		leftBuffer, _ := cmd.Flags().GetString("lbuffer")
		rightBuffer, _ := cmd.Flags().GetString("rbuffer")

		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		expander := expand.NewExpander(cfg)
		result, err := expander.Expand(expand.ExpandRequest{
			LeftBuffer:  leftBuffer,
			RightBuffer: rightBuffer,
		})
		if err != nil {
			return fmt.Errorf("failed to expand: %w", err)
		}

		fullLine := result.NewLeftBuffer + result.NewRightBuffer
		safeFullLine := strings.ReplaceAll(fullLine, "'", "'\"'\"'")
		fmt.Printf("READLINE_LINE='%s'\n", safeFullLine)
		fmt.Printf("READLINE_POINT=%d\n", result.CursorOffset)
		if result.SetCursor {
			fmt.Printf("SET_CURSOR=1\n")
		}

		return nil
	},
}

func init() {
	expandCmd.Flags().String("lbuffer", "", "Left buffer content")
	expandCmd.Flags().String("rbuffer", "", "Right buffer content")
	rootCmd.AddCommand(expandCmd)
}
