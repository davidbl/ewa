package commands

import (
  "github.com/spf13/cobra"
  "strings"
  "os"
  "bytes"
  "fmt"
  "encoding/gob"
  "time"
)

func init() {
  EwaCmd.AddCommand(noteCmd)
}

type Note struct {
  Note string
  Timestamp time.Time
}


var noteCmd = &cobra.Command{
  Use: "note",
  Short: "save a note",
  Long: "save a note",
  Run:  func(cmd *cobra.Command, args []string) {
    writeNote(strings.Join(args, " "))
  },
}

func writeNote(note string) {
  filename := "notes.dat"

  var buf bytes.Buffer

  enc := gob.NewEncoder(&buf)
  f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
  CheckErr(err, "OpenFile  error:")

  defer f.Close()

  // find current offset to end of file
  curr,err := f.Seek(0,2)
  CheckErr(err, "fileSeek  error:")

  err = enc.Encode(Note{note, time.Now().UTC()})
  CheckErr(err, "encode error:")

  _, err = f.WriteAt(buf.Bytes(),curr)
  CheckErr(err, "file.WriteAt error:")

  fmt.Printf("%q (written at %d)", note, curr)
}
