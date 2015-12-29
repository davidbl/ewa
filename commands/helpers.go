package commands

import (
  "encoding/binary"
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
