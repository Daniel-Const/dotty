package cmd

import (
  "fmt"

  "github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "0.0.1",
  Long:  `Still in core development`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("Dotty v0.0.1")
  },
}
