package commands

import (
  "github.com/spf13/cobra"
  "github.com/boltdb/bolt"
  "strings"
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

var noteCmd = &cobra.Command{
  Use: "note",
  Short: "save a note",
  Long: "save a note",
  Run:  func(cmd *cobra.Command, args []string) {
    writeNote(strings.Join(args, " "))
  },
}

func writeNote(noteBody string) {
  db, err := bolt.Open(DataPath(), 0600, nil)
  CheckErrFatal(err, "db file open err")

  defer db.Close()

  err = db.Update(func(tx *bolt.Tx) error {
    note := NoteSave(noteBody, tx)
    config.Log.Println("saving note:", note)

    // save the tags, if any
    tagBucket, err := tx.CreateBucketIfNotExists(config.TagBucketName)
    if err != nil {
      return err
    }

    tagStrings := strings.Split(tags,",")
    for _, tag := range tagStrings {
      exTagVal := tagBucket.Get([]byte(tag))
      if exTagVal != nil {
        t := TagUpdate(exTagVal, note, tagBucket)
        config.Log.Println("updated tag", t)
      } else {
        t := TagCreate(tag, note, tagBucket)
        config.Log.Println("created tag", t)
      }
    }
    return nil
  })
}
