package data

import (
	"encoding/binary"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

func bucket(name string, tx *bolt.Tx) (bucket *bolt.Bucket, err error) {
	bucket = tx.Bucket([]byte(name))
	return
}

func get(ID int, bucket *bolt.Bucket) (v []byte, err error) {
	if bucket == nil {
		return nil, fmt.Errorf("bucket must not be nil")
	}

	v = bucket.Get(itob(ID))
	if v == nil {
		return nil, fmt.Errorf("entity with id not found")
	}
	return v, nil
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
