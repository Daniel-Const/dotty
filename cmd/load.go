package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/Daniel-Const/dotty/core"
)

func init() {
    rootCmd.AddCommand(LoadCmd)
}

var LoadCmd = &cobra.Command {
    Use:    "load",
    Short:  "Load dotfiles from destination paths into a profile",
    Long:   "Copies the dotfiles / dotfile directories from a map file into a profile",
    Args:   cobra.ExactArgs(1),
    RunE:   func(cmd *cobra.Command, args []string) error {
        p, err := core.NewProfile(args[0]).LoadMap()
        if err != nil {
            log.Fatal(err)
        }
        return p.Load()
    },
}
