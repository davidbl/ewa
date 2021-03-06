package commands

import (
  "github.com/spf13/cobra"
  "os"
  "os/user"
  "path"
  "log"
  "io"
  "ewa/persistence"
  "ewa/envar"
)

const (
  LogDestinationNone = "NONE"
  LogDestinationStdOut = "STDOUT"
  LogDestinationFile = "FILE"
  LogDestinationBoth = "FILE_AND_STDOUT"
)

var (
  config Config
)

type Config struct {
  DataDir string
  DataFile string
  LogFile string
  LogDestination string
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

func setEnv() {
  envar.Add("EWA_DATADIR", "tmp", "Directory where data files are stored", setDataDir)
  envar.Add("EWA_DATAFILE", "ewa.db", "Name of data file", nil)
  envar.Add("EWA_TAGBUCKETNAME", "tags", "Bucket Name for tags", nil)
  envar.Add("EWA_NOTEBUCKETNAME", "notes", "Bucket Name for notes", nil)
  envar.Add("EWA_LOGDESTINATION", LogDestinationFile, "Destination of log output", pickLogDestination)
  envar.Add("EWA_LOGLOCATION", "./ewa.log", "Name of log file", nil)


  config.DataDir = envar.Get("EWA_DATADIR")
  config.DataFile = envar.Get("EWA_DATAFILE")
  config.TagBucketName = []byte(envar.Get("EWA_TAGBUCKETNAME"))
  config.NoteBucketName = []byte(envar.Get("EWA_NOTEBUCKETNAME"))
  config.LogDestination = envar.Get("EWA_LOGDESTINATION")
  config.LogFile = envar.Get("EWA_LOGLOCATION")
}

func setConfig() {
  // set log to stdout during initialization of config - this will be changed later
  config.Log = log.New(os.Stdout, "config: ", log.Lshortfile|log.Ldate|log.Ltime)

  setEnv()

  db, err := persistence.Initialize(DataPath())
  if err != nil {
    config.Log.Fatal(err)
  }

  config.Store = db

  // logging
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

func pickLogDestination(v string, defaultV string) string {
  switch v {
  default: return defaultV
  case "0", LogDestinationNone: return LogDestinationNone
  case "1", LogDestinationStdOut: return LogDestinationStdOut
  case "2", LogDestinationFile: return LogDestinationFile
  case "3", "BOTH", LogDestinationBoth: return LogDestinationBoth
  }
}

func setDataDir(v string, defaultV string) string {
  var val string
  if v != "" {
    val = v
  } else {
    usr, err := user.Current()
    if err != nil {
      config.Log.Println("unable to get current user home dir,  using default")
      val = defaultV
    } else {
      val = usr.HomeDir
    }
  }
 // create the dir, if needed
 _, err := os.Stat(val)
 if os.IsNotExist(err) {
    config.Log.Println("Creating missing data directory", val)
    err = os.MkdirAll(val, 0755)
 }
 CheckErrFatal(err, "unable to create DataDir")
 return val
}
