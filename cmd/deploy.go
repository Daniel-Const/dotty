package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/Daniel-Const/dotty/core"
)

func init() {
	rootCmd.AddCommand(DeployCmd)
}

var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy dot files in a profile",
	Long:  "Copy all of the dotfiles to the locations specified in the map file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := core.NewProfile(args[0]).LoadMap()
		if err != nil {
			log.Fatal(err)
		}
		return p.Deploy()
	},
}
