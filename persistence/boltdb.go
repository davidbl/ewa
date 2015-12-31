package persistence

import (
  "github.com/boltdb/bolt"
  "log"
)

type BoltDb struct {
  DataPath string
  Db  *bolt.DB
}

func Initialize(n string) (*BoltDb, error) {
  db, err := bolt.Open(n, 0600, nil)
  if err != nil {
    return nil, err
  }
  return &BoltDb{DataPath: n, Db: db }, err
}

func (db BoltDb) Close() error {
  db.Db.Close()
  return nil
}

func (db BoltDb) Find(bucketName []byte, key []byte) []byte {
  var copyBytes []byte

  db.Db.Update(func(tx *bolt.Tx) error {
    bucket, err := tx.CreateBucketIfNotExists(bucketName)
    logFatal(err)

    someBytes := bucket.Get(key)
    l := len(someBytes)
    copyBytes = make([]byte,l, 2*l)

    copy(copyBytes, someBytes)
    return nil
  })
  return copyBytes
}

func (db BoltDb) Save(v Saveable) Saveable {
  db.Db.Update(func(tx *bolt.Tx) error {
    bucket, err := tx.CreateBucketIfNotExists([]byte(v.TableName()))
    logFatal(err)

    if !v.Persisted() {
      var id uint64
      id, _ = bucket.NextSequence()
      v = v.SetId(id)
    }

    err = bucket.Put(v.PrimaryKey(), v.Marshal())
    logFatal(err)

    return nil
  })
  return v
}

func logFatal(err error) {
    if err != nil {
      log.Fatal(err)
    }
}
