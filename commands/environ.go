package commands

import (
  "fmt"
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
    fmt.Println( `EWA_LOGDESTINATION - controls where logs are written
  allowable values are NONE, STDOUT, FILE, BOTH. Defaults to FILE
EWA_DATADIR - sets the directory Ewa uses to store data. Defaults to home dir of current user
EWA_LOGLOCATION - sets the file Ewa uses to store logs. Defaults to 'ewa.log'
  (log file will always be stored in EWA_DATADIR)`)
  },
}

