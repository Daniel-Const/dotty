package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Daniel-Const/dotty/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/Daniel-Const/dotty/core"
)

func init() {
	var path string = ""
	home, err := os.UserHomeDir()
	if err == nil {
		path = filepath.Join(home, ".config/dotty/dotty.conf")
	}
	rootCmd.PersistentFlags().String("config", path, "Path to Dotty config file")
}

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

			// Set config path argument
			var configPath string = ""
			if configFlag := cmd.Flags().Lookup("config"); configFlag != nil {
				configPath = configFlag.Value.String()
			} else {
				log.Println("Failed to load config argument")
				return nil
			}

			config, err := core.LoadConfig(configPath)
			if config == nil {
				log.Println(err)
			}

			p := tea.NewProgram(
				tui.NewModel([]tui.Command{
					{Name: "Deploy", Desc: DeployCmd.Short},
					{Name: "Load", Desc: LoadCmd.Short},
				}, config),
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
