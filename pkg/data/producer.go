package data

import (
	"encoding/json"

	bolt "go.etcd.io/bbolt"
)

// Producer represents a single studio, producer,
// licensor, etc.
type Producer struct {
	ID      int
	Titles  []Info
	Type    string
	Version int
}

// ProducerBucketName provides the database bucket name
// for the Producer entity
const producerBucketName = "Producer"

// ProducerGet retrieves a single instance of Producer with
// the given ID
func ProducerGet(ID int, db *bolt.DB) (p Producer, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get Producer bucket, exit if error
		b, err := bucket(producerBucketName, tx)
		if err != nil {
			return err
		}

		// Get Producer by ID, exit if error
		v, err := get(ID, b)
		if err != nil {
			return err
		}
		return json.Unmarshal(v, &p)
	})

	return
}

// ProducerGetAll retrieves all persisted Producer values
func ProducerGetAll(db *bolt.DB) (list []Producer, err error) {
	return ProducerGetFilter(db, func(p *Producer) bool { return true })
}

// ProducerGetFilter retrieves all persisted Producer values
// that pass the filter
func ProducerGetFilter(db *bolt.DB, filter func(p *Producer) bool) (list []Producer, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get Producer bucket, exit if error
		b, err := bucket(producerBucketName, tx)
		if err != nil {
			return err
		}

		// Unmarshal and add Producers to slice,
		// exit if error
		return b.ForEach(func(k, v []byte) error {
			var p Producer
			err = json.Unmarshal(v, &p)
			if err != nil {
				return err
			}

			if filter(&p) {
				list = append(list, p)
			}
			return err
		})
	})

	return
}

// ProducerCreate persists a new instance of Producer
func ProducerCreate(p *Producer, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Get Producer bucket, exit if error
		b, err := bucket(producerBucketName, tx)
		if err != nil {
			return err
		}

		// Get next ID in sequence and assign to Producer
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		p.ID = int(id)

		// Save Producer in bucket
		buf, err := json.Marshal(p)
		if err != nil {
			return err
		}

		return b.Put(itob(p.ID), buf)
	})
}

// ProducerUpdate updates the properties of an existing
// persisted Producer instance
func ProducerUpdate(p *Producer, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Get Producer bucket, exit if error
		b, err := bucket(producerBucketName, tx)
		if err != nil {
			return err
		}

		// Check if Producer with ID exists
		o, err := get(p.ID, b)
		if err != nil {
			return err
		}

		// Replace properties of new with immutable
		// ones of old
		old := Producer{}
		err = json.Unmarshal([]byte(o), &old)
		// Update version
		p.Version = old.Version + 1

		// Save Producer
		buf, err := json.Marshal(p)
		if err != nil {
			return err
		}

		return b.Put(itob(p.ID), buf)
	})
}
