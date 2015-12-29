package commands

import (
  "encoding/binary"
  "encoding/gob"
  "github.com/boltdb/bolt"
  "bytes"
  "time"
)

func CheckErr(err error, msg string) {
  if err != nil {
    config.Log.Printf("%s: %v\n", msg, err)
  }
}

func CheckErrFatal(err error, msg string) {
  if err != nil {
    config.Log.Fatalf("%s: %v\n", msg, err)
  }
}

func Itob(v uint64) []byte {
  b := make([]byte, 8)
  binary.BigEndian.PutUint64(b, v)
  return b
}

func NoteById(id uint64, tx *bolt.Tx) Note {
  bucket := tx.Bucket(config.NoteBucketName)
  noteBytes := bucket.Get(Itob(id))
  return NoteFromByte(noteBytes)
}

func NoteSave(noteBody string, tx *bolt.Tx) Note {
  var noteBuf bytes.Buffer
  noteEnc := gob.NewEncoder(&noteBuf)
  bucket, err := tx.CreateBucketIfNotExists(config.NoteBucketName)
  CheckErrFatal(err, "bucket error:")
  id, _ := bucket.NextSequence()

  note := Note{id, noteBody, time.Now().UTC()}
  err = noteEnc.Encode(note)
  CheckErrFatal(err, "note encode error:")

  err = bucket.Put(Itob(id), noteBuf.Bytes())
  CheckErrFatal(err, "note save error:")
  return note
}

func NoteFromByte(noteBytes []byte) Note {
  var note Note

  noteBuf := bytes.NewBuffer(noteBytes)
  noteDec := gob.NewDecoder(noteBuf)
  err := noteDec.Decode(&note)
  CheckErrFatal(err, "note decode error:")

  return note
}

func TagUpdate(tagBytes []byte, note Note, tagBucket *bolt.Bucket) Tag {
  var exTag Tag
  var tagBuf bytes.Buffer

  exTagBuf := bytes.NewBuffer(tagBytes)
  tagDec := gob.NewDecoder(exTagBuf)
  err := tagDec.Decode(&exTag)

  exTag.NoteIds = append(exTag.NoteIds, note.Id)

  tagEnc := gob.NewEncoder(&tagBuf)
  err = tagEnc.Encode(exTag)
  CheckErrFatal(err, "tag encoding error:")

  err = tagBucket.Put([]byte(exTag.TagText),tagBuf.Bytes())
  CheckErrFatal(err, "tag Update error:")

  return exTag
}

func TagCreate(tag string, note Note, tagBucket *bolt.Bucket) Tag {
  var tagBuf bytes.Buffer

  t := Tag{tag,[]uint64{note.Id}}
  tagEnc := gob.NewEncoder(&tagBuf)
  err := tagEnc.Encode(t)
  CheckErrFatal(err, "tag encode error:")

  err = tagBucket.Put([]byte(tag), tagBuf.Bytes())
  CheckErrFatal(err, "tag save error:")

  return t
}

func TagFromByte(tagBytes []byte) Tag {
  var tag Tag

  exTagBuf := bytes.NewBuffer(tagBytes)
  tagDec := gob.NewDecoder(exTagBuf)
  err := tagDec.Decode(&tag)
  CheckErrFatal(err, "tag decode error:")

  return tag
}
