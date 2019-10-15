package data

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

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

// Idenitifiable encompasses all the entities
// persisted in the database
type Idenitifiable interface {
	Identifier() int
	SetIdentifier(int)

	Ver() int
	UpdateVer()

	Validate(tx *bolt.Tx) error
}

// getByID is a generic function that queries the given bucket
// in the given database for an entity of the given ID
func getByID(ID int, entity Idenitifiable, bucketName string, db *bolt.DB) (err error) {
	return db.View(func(tx *bolt.Tx) error {
		// Get bucket, exit if error
		b, err := bucket(bucketName, tx)
		if err != nil {
			return err
		}

		// Get entity by ID, exit if error
		v, err := get(ID, b)
		if err != nil {
			return err
		}

		return json.Unmarshal(v, entity)
	})
}

func getFilter(t Idenitifiable, filter func(e Idenitifiable) (bool, error), bucketName string, db *bolt.DB) (list []Idenitifiable, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get bucket, exit if error
		b, err := bucket(bucketName, tx)
		if err != nil {
			return err
		}

		// Unmarshal and add all entities who
		// pass filter to slice, exit if error
		return b.ForEach(func(k, v []byte) error {
			m := t
			err = json.Unmarshal(v, &m)
			if err != nil {
				return err
			}

			pass, err := filter(m)
			if err != nil {
				return err
			}
			if pass {
				list = append(list, m)
			}
			return nil
		})
	})
	return
}

// create is a generic function that persists an
// entity into the given bucket in the given
// databases
func create(entity Idenitifiable, bucketName string, db *bolt.DB) (err error) {
	return db.Update(func(tx *bolt.Tx) error {
		// Get bucket, exit if error
		b, err := bucket(bucketName, tx)
		if err != nil {
			return err
		}

		// Verify validity of struct
		err = entity.Validate(tx)
		if err != nil {
			return err
		}

		// Get next ID in sequence and
		// assign to entity
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		entity.SetIdentifier(int(id))
		entity.UpdateVer()

		// Save entity in bucket
		buf, err := json.Marshal(entity)
		if err != nil {
			return err
		}

		return b.Put(itob(entity.Identifier()), buf)
	})
}

// update is a generic function replaces the value
// of the given entity id in the given bucket
// in the given database
func update(entity Idenitifiable, bucketName string, db *bolt.DB) (err error) {
	return db.Update(func(tx *bolt.Tx) error {
		// Get bucket, exit if error
		b, err := bucket(bucketName, tx)
		if err != nil {
			return err
		}

		// Check if entity with ID exists
		o, err := get(entity.Identifier(), b)
		if err != nil {
			return err
		}

		// Verify validity of struct
		err = entity.Validate(tx)
		if err != nil {
			return err
		}

		// Replace properties of new with internal
		// ones of old
		old := entity
		err = json.Unmarshal([]byte(o), &old)
		entity.UpdateVer()

		// Save Media
		buf, err := json.Marshal(entity)
		if err != nil {
			return err
		}

		return b.Put(itob(entity.Identifier()), buf)
	})
}

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
		return nil, fmt.Errorf("entity with id %d not found", ID)
	}
	return v, nil
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
