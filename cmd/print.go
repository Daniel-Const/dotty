package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/Daniel-Const/dotty/core"
)

func init() {
	rootCmd.AddCommand(printCmd)
}

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print a profiles dotfiles map",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pName := args[0]
		p := core.NewProfile(pName)
		_, err := p.LoadMap()
		if err != nil {
			fmt.Printf("Error!")
			log.Fatal(err)
		}
		p.Print()
		return nil
	},
}
