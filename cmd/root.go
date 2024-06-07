package cmd

import (
	"fmt"
    "os"

	"github.com/spf13/cobra"
    tea "github.com/charmbracelet/bubbletea"

    "github.com/Daniel-Const/dotty/tui"
)

var (
	rootCmd = &cobra.Command{
		Use:   "dotty",
		Short: "A dotfiles manager app",
		Long:  `Dotty is a dotfiles manager for deploying and updating multiple collections of dotfiles`,
        RunE:   func(cmd *cobra.Command, args []string) error {
            // Initialise logging
            f, err := tea.LogToFile("debug.log", "debug")
            if err != nil {
                fmt.Println(err)
            }
            defer f.Close()

            p := tea.NewProgram(
                tui.NewModel([]tui.Command{
                    tui.NewCommand("Deploy", DeployCmd.Short),
                    tui.NewCommand("Load"  , LoadCmd.Short),
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

func Execute() error {
	return rootCmd.Execute()
}

