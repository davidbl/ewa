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

var tags string

func init() {
  EwaCmd.AddCommand(noteCmd)
  noteCmd.Flags().StringVarP(&tags, "tags", "t", "", "comma-separated list of tags (no spaces)")
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

var tagBucketName = []byte("tags")
var noteBucketName = []byte("notes")


func writeNote(note string) {
  db, err := bolt.Open(DataPath(), 0600, nil)
  CheckErr(err, "db file open err")

  defer db.Close()

  var buf bytes.Buffer

  enc := gob.NewEncoder(&buf)
  fmt.Println("tags:", tags)

  err = db.Update(func(tx *bolt.Tx) error {
    bucket, err := tx.CreateBucketIfNotExists(noteBucketName)
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

    // save the tags, if any
    tagBucket, err := tx.CreateBucketIfNotExists(tagBucketName)
    if err != nil {
      return err
    }
    tagStrings := strings.Split(tags,",")
    for _, tag := range tagStrings {
      exTagVal := tagBucket.Get([]byte(tag))
      if exTagVal != nil {
        appendand := fmt.Sprintf(",%d",id)
        newVal := make([]byte, len(exTagVal))
        copy(newVal,exTagVal)
        newVal = append(newVal, appendand...)
        err = tagBucket.Put([]byte(tag),newVal)
        if err != nil {
          return err
        }
      } else {
        idVal := fmt.Sprintf("%d",id)
        err = tagBucket.Put([]byte(tag), []byte(idVal))
        if err != nil {
          return err
        }
      }
    }
    return nil
  })
}
