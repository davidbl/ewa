package commands

import (
  "github.com/spf13/cobra"
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
    // find all the values (ids) from the given tags
    idSet := make(map[uint64]uint64)
    for _, tag := range tags {
      tagBytes := config.Store.Find(config.TagBucketName, []byte(tag))
      if len(tagBytes) > 0 {
        exTag := TagFromByte(tagBytes)
        for _,id := range exTag.NoteIds {
          idSet[id] = id
        }
      }
    }
    for k,_ := range idSet {
      noteBytes := config.Store.Find(config.NoteBucketName, Itob(k))
      if len(noteBytes) > 0 {
        note := NoteFromByte(noteBytes)
        fmt.Printf("%v\n", note)
      }
    }
}
