package commands

import (
  "github.com/spf13/cobra"
  "github.com/boltdb/bolt"
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

  err = db.View(func(tx *bolt.Tx) error {
    c := tx.Bucket(config.TagBucketName).Cursor()

    // find all the values (ids) from the given tags
    idSet := make(map[uint64]uint64)
    for _, tag := range tags {
      prefix := []byte(tag)
      for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k,v = c.Next() {
        exTag := TagFromByte(v)
        for _,id := range exTag.NoteIds {
          idSet[id] = id
        }
      }
    }
    for k,_ := range idSet {
      note := NoteById(k, tx)
      fmt.Printf("note: %s (id: %d, on %s)\n", note.Note, note.Id, note.Timestamp)
    }
    return nil
  })
}
