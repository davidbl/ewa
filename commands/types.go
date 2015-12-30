package commands

import (
  "time"
  "fmt"
)

type Note struct {
  Id  uint64
  Note string
  Timestamp time.Time
}

func (n Note) String() string {
  return fmt.Sprintf("note: %s\n{id: %d, created: %s}", n.Note, n.Id, n.Timestamp)
}

type Tag struct {
  TagText string
  NoteIds []uint64
}
