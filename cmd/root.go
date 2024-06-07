package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "dotty",
		Short: "A dotfiles manager app",
		Long:  `Dotty is a dotfiles manager for deploying and updating multiple collections of dotfiles`,
        RunE:   func(cmd *cobra.Command, args []string) error {
            // TODO Implement bubbletea TUI
            fmt.Println("TODO...")
            return nil
        },
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

