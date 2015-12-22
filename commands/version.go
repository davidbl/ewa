package commands

import (
  "fmt"
  "github.com/spf13/cobra"

)

var version string = "0.0.1"

func init() {
  EwaCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
  Use: "version",
  Short: "Display the version number of Ewa",
  Long: `Display the version number of Ewa that
is currently running on this system`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Printf("Ewa CLI note-saving system v%v\n", version)
  },
}

