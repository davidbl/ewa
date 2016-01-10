package commands

import (
  "fmt"
  "ewa/envar"
  "github.com/spf13/cobra"
)

func init() {
  EwaCmd.AddCommand(environCmd)
}

var environCmd = &cobra.Command{
  Use: "environ",
  Short: "list the environment variables that control Ewa",
  Long: "list the environment variables that control Ewa",
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println( envar.Help())
  },
}

