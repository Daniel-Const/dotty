package cmd

import (
    "github.com/spf13/cobra"
    "dotty/core"
)

func init() {
    rootCmd.AddCommand(DeployCommand)
}


var DeployCommand = &cobra.Command {
    Use:    "deploy",
    Short:  "Deploy dot files in a profile",
    Long:   "Copy all of the dotfiles to the locations specified in the map file",
    Args:   cobra.ExactArgs(1),
    RunE:   func(cmd *cobra.Command, args []string) error {
        p := core.NewProfile(args[0]).Load()
        return p.Deploy()
    },
}