package persistence

import (
  "github.com/boltdb/bolt"
  "log"
)

type BoltDb struct {
  DataPath string
}

func (db BoltDb) Find(bucketName []byte, key []byte) []byte {
  var copyBytes []byte
  myDb, err := bolt.Open(db.DataPath, 0600, nil)
  if err != nil {
    log.Fatal(err)
  }
  defer myDb.Close()
  myDb.Update(func(tx *bolt.Tx) error {
    if err != nil {
      log.Fatal(err)
    }
    bucket, err := tx.CreateBucketIfNotExists(bucketName)
    if err != nil {
      log.Fatal(err)
    }
    someBytes := bucket.Get(key)
    l := len(someBytes)
    copyBytes = make([]byte,l, 2*l)

    copy(copyBytes, someBytes)
    return nil
  })
  return copyBytes
}

func (db BoltDb) Save(v Saveable) Saveable {
  dataStore, err := bolt.Open(db.DataPath, 0600, nil)
  if err != nil {
    log.Fatal(err)
  }

  defer dataStore.Close()

  err = dataStore.Update(func(tx *bolt.Tx) error {
    bucket, err := tx.CreateBucketIfNotExists([]byte(v.TableName()))
    if err != nil {
      log.Fatal(err)
    }
    if !v.Persisted() {
      var id uint64
      id, _ = bucket.NextSequence()
      v = v.SetId(id)
    }

    err = bucket.Put(v.PrimaryKey(), v.Marshal())
    if err != nil {
      log.Fatal(err)
    }
    return nil
  })
  return v
}
