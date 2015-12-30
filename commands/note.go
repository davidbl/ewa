package commands

import (
  "github.com/spf13/cobra"
  "strings"
)

var tagFlag string

func init() {
  EwaCmd.AddCommand(noteCmd)
  noteCmd.Flags().StringVarP(&tagFlag, "tags", "t", "", "comma-separated list of tags (no spaces)")
}

var noteCmd = &cobra.Command{
  Use: "note",
  Short: "save a note",
  Long: "save a note",
  Run:  func(cmd *cobra.Command, args []string) {
    note := writeNote(strings.Join(args, " "))
    writeTags(tagFlag, note)
  },
}

func writeTags(tagString string, note Note) {
  tagStrings := strings.Split(tagString, ",")

  var t Tag
  for _, tag := range tagStrings {
    tagBytes := config.Store.Find(config.TagBucketName, []byte(tag))
    if len(tagBytes) > 0 {
      t = TagFromByte(tagBytes)
      t.NoteIds = append(t.NoteIds, note.Id)
    } else {
      t = BuildTag(tag,note.Id)
    }
    saved := config.Store.Save(t)
    config.Log.Println("saving tag:",saved)
  }
}
func writeNote(noteBody string) Note {
  note := BuildNote(noteBody)
  newNote := config.Store.Save(note).(Note)
  config.Log.Println("saving note:", newNote)
  return newNote
}
