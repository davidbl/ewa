package commands

import (
  "github.com/spf13/cobra"
  "os"
  "os/user"
  "path"
  "log"
  "io"
  "ewa/persistence"
)

const (
  LogDestinationNone = 0
  LogDestinationStdOut = 1
  LogDestinationFile = 2
  LogDestinationBoth = 3
)

var (
  config Config
)

type Config struct {
  DataDir string
  DataFile string
  LogFile string
  LogDestination int
  Log *log.Logger
  NoteBucketName []byte
  TagBucketName []byte
  Store persistence.Persistor
}

var EwaCmd = &cobra.Command{
  Use:   "ewa",
  Short: "ewa saves your stuffs",
  Long: `ewa is the main command, used to save all the stuff you do
during the day at the command-line`,
  PersistentPreRun: func(cmd *cobra.Command, args []string) {
    setConfig()
  },
  PersistentPostRun: func(cmd *cobra.Command, args []string) {
    shutDown()
  },
}

func DataPath() string {
  return path.Join(config.DataDir, config.DataFile)
}

func shutDown() {
  config.Store.Close()
}

func setConfig() {
  config.TagBucketName = []byte("tags")
  config.NoteBucketName = []byte("notes")

  if os.Getenv("EWA_LOGDESTINATION") != "" {
    switch os.Getenv("EWA_LOGDESTINATION") {
    case "0", "NONE": config.LogDestination = LogDestinationNone
    case "1", "STDOUT": config.LogDestination = LogDestinationStdOut
    case "2", "FILE": config.LogDestination = LogDestinationFile
    case "3", "BOTH": config.LogDestination = LogDestinationBoth
    }
  } else {
    config.LogDestination = LogDestinationFile
  }
  // set log to stdout during initialization of config - this will be changed later
  config.Log = log.New(os.Stdout, "config: ", log.Lshortfile|log.Ldate|log.Ltime)

  config.DataFile = "ewa.db"
  if os.Getenv("EWA_DATADIR") != "" {
    config.DataDir = os.Getenv("EWA_DATADIR")
  } else {
    usr, err := user.Current()
    CheckErrFatal(err, "unable to get current user")
    config.DataDir = usr.HomeDir
  }
  _, err := os.Stat(config.DataDir)
  if os.IsNotExist(err) {
    config.Log.Println("Creating missing data directory", config.DataDir)
    err = os.MkdirAll(config.DataDir, 0755)
  }
  if err != nil {
    config.Log.Fatal(err)
  }

  db, err := persistence.Initialize(DataPath())
  if err != nil {
    config.Log.Fatal(err)
  }

  config.Store = db

  // logging
  if os.Getenv("EWA_LOGLOCATION") != "" {
    config.LogFile = os.Getenv("EWA_LOGLOCATION")
  } else {
    config.LogFile = path.Join(config.DataDir,"ewa.log")
  }
  var multi io.Writer
  switch config.LogDestination {
  case LogDestinationNone: multi = io.MultiWriter()
  case LogDestinationStdOut: multi = io.MultiWriter(os.Stdout)
  case LogDestinationFile: multi = io.MultiWriter(openLogFile())
  case LogDestinationBoth: multi = io.MultiWriter(openLogFile(), os.Stdout)
  }
  config.Log = log.New(multi, "log: ", log.Lshortfile|log.Ldate|log.Ltime)
  config.Log.Println("initializing")
}

func openLogFile() *os.File {
  file, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
  if err != nil {
      config.Log.Fatal("Failed to open log file", config.LogFile, ":", err)
    }
  return file
}
