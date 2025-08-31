package cmd

import (
	"fmt"

	"github.com/hagatasdelus/babbr/internal/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured abbreviations",
	Long:  "Display a list of all configured abbreviations and their expansions",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if len(cfg.Abbreviations) == 0 {
			fmt.Println("No abbreviations configured.")
			return nil
		}

		for _, abbr := range cfg.Abbreviations {
			if abbr.Abbr != "" {
				fmt.Printf("%-10s -> %s\n", abbr.Abbr, abbr.Snippet)
			} else if abbr.Options != nil && abbr.Options.Regex != "" {
				fmt.Printf("%-10s -> %s (regex: %s)\n", "[regex]", abbr.Snippet, abbr.Options.Regex)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
