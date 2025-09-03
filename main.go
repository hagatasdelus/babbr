package main

import (
	"fmt"
	"os"

	"github.com/hagatasdelus/babbr/internal/cmd"
)

func main() {
	cmd.SetVersion(version)
	if err := cmd.GetRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
