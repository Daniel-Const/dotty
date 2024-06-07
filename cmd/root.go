package cmd

import (
	"fmt"
    "os"

    "dotty/tui"

	"github.com/spf13/cobra"
    tea "github.com/charmbracelet/bubbletea"
)

var (
	rootCmd = &cobra.Command{
		Use:   "dotty",
		Short: "A dotfiles manager app",
		Long:  `Dotty is a dotfiles manager for deploying and updating multiple collections of dotfiles`,
        RunE:   func(cmd *cobra.Command, args []string) error {
            p := tea.NewProgram(
                tui.New([]tui.Command{
                    tui.Command{Name: "Deploy", Desc: DeployCmd.Short},
                    tui.Command{Name: "Load",   Desc: LoadCmd.Short},
                }),
            )
            if _, err := p.Run(); err != nil {
                fmt.Printf("Failed to run tea program :( %v", err)
                os.Exit(1)
            } 
            return nil
        },
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

