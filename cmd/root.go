package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "dotty",
		Short: "A dotfiles manager app",
		Long: `Dotty is a dotfiles manager for deploying and updating
        multiple collections of dotfiles`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

