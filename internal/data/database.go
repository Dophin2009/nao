package data

import (
	"os"

	bolt "go.etcd.io/bbolt"
)

// Buckets provides an array of all the buckets in the database
func Buckets() []string {
	return []string{
		MediaBucket, ProducerBucket, GenreBucket, EpisodeBucket,
		CharacterBucket, PersonBucket, UserBucket, MediaProducerBucket,
		MediaRelationBucket, MediaGenreBucket, MediaCharacterBucket,
		UserMediaBucket, UserMediaListBucket, JWTBucket,
	}
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
				_, err = tx.CreateBucketIfNotExists([]byte(bucket))
				if err != nil {
					return err
				}
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

// Bucket returns the database bucket with the
// given name
func Bucket(name string, tx *bolt.Tx) (bucket *bolt.Bucket, err error) {
	bucket = tx.Bucket([]byte(name))
	return
}
