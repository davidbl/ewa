package commands

import (
  "github.com/spf13/cobra"
  "github.com/boltdb/bolt"
  "strings"
  "bytes"
  "encoding/gob"
  "encoding/binary"
  "time"
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

type Tag struct {
  TagText string
  NoteIds []uint64
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


func writeNote(noteBody string) {
  db, err := bolt.Open(DataPath(), 0600, nil)
  CheckErrFatal(err, "db file open err")

  defer db.Close()

  var noteBuf bytes.Buffer
  noteEnc := gob.NewEncoder(&noteBuf)

  err = db.Update(func(tx *bolt.Tx) error {
    bucket, err := tx.CreateBucketIfNotExists(noteBucketName)
    if err != nil {
      return err
    }
    id, _ := bucket.NextSequence()

    note := Note{id, noteBody, time.Now().UTC()}
    err = noteEnc.Encode(note)
    CheckErrFatal(err, "note encode error:")

    err = bucket.Put(itob(id), noteBuf.Bytes())
    if err != nil {
      return err
    }
    config.Log.Println("saving note:", note)

    // save the tags, if any
    tagBucket, err := tx.CreateBucketIfNotExists(tagBucketName)
    if err != nil {
      return err
    }

    tagStrings := strings.Split(tags,",")
    for _, tag := range tagStrings {
      exTagVal := tagBucket.Get([]byte(tag))
      if exTagVal != nil {
        exTagBuf := bytes.NewBuffer(exTagVal)
        tagDec := gob.NewDecoder(exTagBuf)
        var exTag Tag
        err = tagDec.Decode(&exTag)
        exTag.NoteIds = append(exTag.NoteIds, id)
        var tagBuf bytes.Buffer
        tagEnc := gob.NewEncoder(&tagBuf)
        err = tagEnc.Encode(exTag)
        err = tagBucket.Put([]byte(tag),tagBuf.Bytes())
        if err != nil {
          return err
        }
        config.Log.Println("updated tag", exTag)
      } else {
        var tagBuf bytes.Buffer
        tagEnc := gob.NewEncoder(&tagBuf)
        t := Tag{tag,[]uint64{id}}
        err = tagEnc.Encode(t)
        CheckErrFatal(err, "tag encode error:")
        err = tagBucket.Put([]byte(tag), tagBuf.Bytes())
        if err != nil {
          return err
        }
        config.Log.Println("saving tag", tag)
      }
    }
    return nil
  })
}
