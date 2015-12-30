package persistence

type Saveable interface {
  TableName() string
  SetId(uint64) Saveable
  GetId() uint64
  Persisted() bool
  PrimaryKey() []byte
  Marshal() []byte
}

type Persistor interface {
  Save(v Saveable) Saveable
  Find(bucketName []byte, key []byte) []byte
}

