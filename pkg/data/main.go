package data

import (
	bolt "go.etcd.io/bbolt"
)

// Buckets provides an array of all the buckets in the database
func Buckets() []string {
	return []string{mediaBucketName, producerBucketName, genreBucketName,
		episodeBucketName, characterBucketName, personBucketName,
		mediaProducerBucketName, mediaRelationBucketName, mediaGenreBucketName,
		mediaCharacterBucketName}
}

// ConnectDatabase connects to the database file at the given path
// and return a bolt.DB struct
func ConnectDatabase(dbPath string, create bool) (*bolt.DB, error) {
	// open database connection
	db, err := bolt.Open(dbPath, 0600, nil)

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
