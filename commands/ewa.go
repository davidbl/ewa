package commands

import (
  "github.com/spf13/cobra"
  "os"
  "os/user"
  "fmt"
  "path"
)

var (
  config Config
)

type Config struct {
  DataDir string
  DataFile string
}

var EwaCmd = &cobra.Command{
  Use:   "ewa",
  Short: "ewa saves your stuffs",
  Long: `ewa is the main command, used to save all the stuff you do
during the day at the command-line`,
  PersistentPreRun: func(cmd *cobra.Command, args []string) {
    setConfig()
  },

}

func DataPath() string {
  return path.Join(config.DataDir, config.DataFile)
}

func setConfig() {
  config.DataFile = "ewa.db"
  if os.Getenv("EWA_DATADIR") != "" {
    config.DataDir = os.Getenv("EWA_DATADIR")
  } else {
    usr, err := user.Current()
    CheckErr(err, "unable to get current user")
    fmt.Println("home dir:", usr.HomeDir)
    config.DataDir = usr.HomeDir
  }
  _, err := os.Stat(config.DataDir)
  if os.IsNotExist(err) {
    fmt.Println("Creating missing data directory", config.DataDir)
    err = os.MkdirAll(config.DataDir, 0755)
  }
  if err != nil {
    panic(fmt.Sprintf("%s", err))
  }
}
