package cmd

import (
    "github.com/spf13/cobra"
    "dotty/core"
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
        p := core.NewProfile(args[0]).LoadMap()
        return p.Load()
    },
}
