package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/Daniel-Const/dotty/tui"
)

var (
	rootCmd = &cobra.Command{
		Use:   "dotty",
		Short: "A dotfiles manager app",
		Long:  `Dotty is a dotfiles manager for deploying and updating multiple collections of dotfiles`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Initialise logging
			f, err := tea.LogToFile("debug.log", "debug")
			if err != nil {
				fmt.Println(err)
			}
			defer f.Close()

			p := tea.NewProgram(
				tui.NewModel([]tui.Command{
					{Name: "Deploy", Desc: DeployCmd.Short},
					{Name: "Load", Desc: LoadCmd.Short},
				}),
				tea.WithAltScreen(),
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
