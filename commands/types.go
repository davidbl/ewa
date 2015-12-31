package commands

import (
  "bytes"
  "time"
  "encoding/gob"
  "fmt"
  "ewa/persistence"
)

type Note struct {
  Id  uint64
  Note string
  Timestamp time.Time
}

func BuildNote(body string) Note {
  return Note{Note: body, Timestamp: time.Now().UTC()}
}

// Saveable interface for Note
func (n Note) TableName() string {
  return "notes"
}
func (n Note) SetId(id uint64) persistence.Saveable {
  n.Id = id
  return n
}
func (n Note) GetId() uint64 {
  return n.Id
}
func (n Note) Persisted() bool {
  return n.Id > 0
}
func (n Note) PrimaryKey() []byte {
  return Itob(n.Id)
}
func (n Note) Marshal() []byte {
  return toByte(n)
}

func (n Note) String() string {
  return fmt.Sprintf("note: %s\n{id: %d, created: %s}", n.Note, n.Id, n.Timestamp)
}

func NoteDemarshal(noteBytes []byte) Note {
  return noteFromByte(noteBytes)
}

func noteFromByte(noteBytes []byte) Note {
  var note Note
  err := decoder(noteBytes).Decode(&note)
  CheckErrFatal(err, "note decode error:")

  return note
}

type Tag struct {
  Id  uint64
  TagText string
  NoteIds []uint64
}

func BuildTag(tagText string, id uint64) Tag {
  return Tag{TagText: tagText, NoteIds: []uint64{id}}
}
// Saveable interface for Tag
func (t Tag) TableName() string {
  return "tags"
}
func (t Tag) SetId(id uint64) persistence.Saveable {
  t.Id = id
  return t
}
func (t Tag) GetId() uint64 {
  return t.Id
}
func (t Tag) Persisted() bool {
  return t.Id > 0
}
func (t Tag) PrimaryKey() []byte {
  return []byte(t.TagText)
}
func (t Tag) Marshal() []byte {
  return toByte(t)
}

func TagDemarshal(tagBytes []byte) Tag {
  return tagFromByte(tagBytes)
}

func tagFromByte(tagBytes []byte) Tag {
  var tag Tag
  err := decoder(tagBytes).Decode(&tag)
  CheckErrFatal(err, "tag decode error:")

  return tag
}

func decoder(b []byte) *gob.Decoder {
  buf := bytes.NewBuffer(b)
  return gob.NewDecoder(buf)
}

func toByte(s persistence.Saveable) []byte {
  var buf bytes.Buffer
  encoder := gob.NewEncoder(&buf)
  err := encoder.Encode(s)
  CheckErrFatal(err, "encode error")
  return buf.Bytes()
}
