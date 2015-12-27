package commands

import (
  "github.com/spf13/cobra"
  "github.com/boltdb/bolt"
  "strings"
  "bytes"
  "encoding/gob"
  "encoding/binary"
  "time"
  "fmt"
)

func init() {
  EwaCmd.AddCommand(noteCmd)
}

type Note struct {
  Id  uint64
  Note string
  Timestamp time.Time
}

func itob(v uint64) []byte {
  b := make([]byte, 8)
  binary.BigEndian.PutUint64(b, v)
  return b
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
  db, err := bolt.Open(DataPath(), 0600, nil)
  CheckErr(err, "db file open err")

  defer db.Close()

  var buf bytes.Buffer

  enc := gob.NewEncoder(&buf)

  err = db.Update(func(tx *bolt.Tx) error {
    bucket, err := tx.CreateBucketIfNotExists([]byte("notes"))
    if err != nil {
      return err
    }
    id, _ := bucket.NextSequence()
    note := Note{id, note, time.Now().UTC()}
    err = enc.Encode(note)
    CheckErr(err, "encode error:")
    err = bucket.Put(itob(id), buf.Bytes())
    if err != nil {
      return err
    }
    fmt.Println(note)
    return nil
  })
}
