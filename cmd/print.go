package cmd

import (
  "fmt"

  "github.com/spf13/cobra"

  "github.com/Daniel-Const/dotty/core"
)

func init() {
    rootCmd.AddCommand(printCommand)
}

var printCommand = &cobra.Command {
    Use:   "print",
    Short: "Print a profiles dotfiles map",
    Args: cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        pName := args[0]
        p := core.NewProfile(pName)
        p.Load()
        for i := range p.Dots {
            fmt.Printf("File: %s, Destination: %s\n", p.Dots[i].SrcPath, p.Dots[i].DestPath)
        }
        return nil
  },
}
