package commands

import (
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
)

var EwaCmd = &cobra.Command{
  Use:   "ewa",
  Short: "ewa saves your stuffs",
  Long: `ewa is the main command, used to save all the stuff you do
during the day at the command-line`,
  PersistentPreRun: func(cmd *cobra.Command, args []string) {
    viper.SetConfigType("yaml")
    viper.SetConfigName(".ewa.config")
    viper.AddConfigPath("$HOME")
    err := viper.ReadInConfig()
    CheckErr(err, "config read error:")
  },

}


// TODO
// add config file using viper default to ~/.ewa.config
// allow for 
//   - note file (default to ~/.ewa/notes.dat)

//  add an init method that creates the default config
//  and default dirs and files

// check for config values in PersisitentPreRun 
