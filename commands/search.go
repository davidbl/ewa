package commands

import (
  "github.com/spf13/cobra"
  "github.com/boltdb/bolt"
  "encoding/gob"
  "bytes"
  "fmt"
)

func init() {
  EwaCmd.AddCommand(searchCmd)
}

var searchCmd = &cobra.Command{
  Use: "search",
  Short: "find notes, using tags",
  Long: `find notes, using tags, eg 'ewa search foo bar baz' will return
the set of notes that were tagged with either foo or bar or baz`,
  Run:  func(cmd *cobra.Command, args []string) {
    findNotes(args)
  },
}
func findNotes(tags []string) {
  db, err := bolt.Open(DataPath(), 0600, nil)
  CheckErrFatal(err, "db file open err")

  defer db.Close()

  err = db.Update(func(tx *bolt.Tx) error {
    c := tx.Bucket(config.TagBucketName).Cursor()

    // find all the values (ids) from the given tags
    idSet := make(map[uint64]uint64)
    for _, tag := range tags {
      prefix := []byte(tag)
      for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k,v = c.Next() {
        exTagBuf := bytes.NewBuffer(v)
        tagDec := gob.NewDecoder(exTagBuf)
        var exTag Tag
        err = tagDec.Decode(&exTag)
        CheckErrFatal(err, "tag decode error:")
        for _,id := range exTag.NoteIds {
          idSet[id] = id
        }
      }
    }
    noteCursor := tx.Bucket(config.NoteBucketName).Cursor()
    for k,_ := range idSet {
      _, noteBytes := noteCursor.Seek(Itob(k))
      noteBuf := bytes.NewBuffer(noteBytes)
      noteDec := gob.NewDecoder(noteBuf)
      var exNote Note
      err = noteDec.Decode(&exNote)
      CheckErrFatal(err, "note decode error:")
      fmt.Printf("note: %s (on %s)\n", exNote.Note, exNote.Timestamp)
    }
    return nil
  })
}
