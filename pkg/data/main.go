package data

import (
	"encoding/binary"
	"fmt"
	"os"

	bolt "go.etcd.io/bbolt"
)

// Buckets provides an array of all the buckets in the database
func Buckets() []string {
	return []string{MediaBucketName, ProducerBucketName, GenreBucketName,
		EpisodeBucketName, CharacterBucketName, PersonBucketName,
		UserBucketName, MediaProducerBucketName, MediaRelationBucketName,
		MediaGenreBucketName, MediaCharacterBucketName, UserMediaBucketName}
}

// ConnectDatabase connects to the database file at the given path
// and return a bolt.DB struct
func ConnectDatabase(dbPath string, mode os.FileMode, create bool) (*bolt.DB, error) {
	// open database connection
	db, err := bolt.Open(dbPath, mode, nil)

	// if specified to create buckets, cycle through all strings in
	// Buckets() and create buckets
	if create {
		err = db.Update(func(tx *bolt.Tx) error {
			for _, bucket := range Buckets() {
				tx.CreateBucket([]byte(bucket))
			}
			return nil
		})
	}

	return db, err
}

// ClearDatabase removes all buckets in the given database
func ClearDatabase(db *bolt.DB) (err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		for _, bucket := range Buckets() {
			err = tx.DeleteBucket([]byte(bucket))
			if err != nil {
				return err
			}
		}
		return nil
	})

	return
}

// Entity encompasses all the entities
// persisted in the database
type Entity interface {
}

// Service defines methods of persistence
// service structs
type Service interface {
	GetByID(e Entity) (err error)
	GetAll() (v [][]byte, err error)
	Create(e Entity) (err error)
	Update(e Entity) (err error)

	Validate(e Entity) (err error)
}

// GetByID is a generic function that queries the given bucket
// in the given database for an entity of the given ID
func GetByID(ID int, bucketName string, db *bolt.DB) (v []byte, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get bucket, exit if error
		b, err := Bucket(bucketName, tx)
		if err != nil {
			return err
		}

		// Get entity by ID, exit if error
		v, err = get(ID, b)
		return err
	})
	return
}

// GetAll returns a list of []byte of all the
// values in the given bucket
func GetAll(bucketName string, db *bolt.DB) (list [][]byte, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get bucket, exit if error
		b, err := Bucket(bucketName, tx)
		if err != nil {
			return err
		}

		// Unmarshal and add all entities who
		// pass filter to slice, exit if error
		return b.ForEach(func(k, v []byte) error {
			list = append(list, v)
			return nil
		})
	})
	return
}

// Bucket returns the database bucket with the
// given name
func Bucket(name string, tx *bolt.Tx) (bucket *bolt.Bucket, err error) {
	bucket = tx.Bucket([]byte(name))
	return
}

func get(ID int, bucket *bolt.Bucket) (v []byte, err error) {
	if bucket == nil {
		return nil, fmt.Errorf("bucket must not be nil")
	}

	v = bucket.Get(itob(ID))
	if v == nil {
		return nil, fmt.Errorf("entity with id %d not found", ID)
	}
	return v, nil
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
