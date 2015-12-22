package commands

import (
  "github.com/spf13/cobra"
  "os"
  "fmt"
  "encoding/gob"
  "time"


)

func init() {
  EwaCmd.AddCommand(searchCmd)
  searchCmd.Flags().Int64VarP(&searchOffset, "offset", "o", 0, "offset into the file")
}

var searchOffset int64


var searchCmd = &cobra.Command{
  Use: "search",
  Short: "find a note, given an offset",
  Long: "find a note, given an offset",
  Run:  func(cmd *cobra.Command, args []string) {
    filename := "notes.dat"

    f, err := os.OpenFile(filename, os.O_RDONLY, 0666)
    CheckErr(err, "OpenFile error:")

    defer f.Close()

    f.Seek(searchOffset,0)
    var n Note
    dec := gob.NewDecoder(f)
    err = dec.Decode(&n)
    CheckErr(err, "Decode error:")

    fmt.Printf("%q: on %s, %d", n.Note, n.Timestamp.Format(time.RFC3339), searchOffset)
  },
}
