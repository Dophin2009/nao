package data

import (
	"encoding/json"
	"time"

	bolt "go.etcd.io/bbolt"
)

// Media represents a single instance of a media
type Media struct {
	ID         int
	Titles     []Info
	Synopses   []Info
	Background []Info
	StartDate  *time.Time
	EndDate    *time.Time
	Type       string
	Source     string
}

const mediaBucketName = "Media"

// MediaGetAll retrieves all persisted Media values
func MediaGetAll(db *bolt.DB) (list []Media, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get Media bucket, exit if error
		b, err := bucket(mediaBucketName, tx)

		if err != nil {
			return err
		}

		// Unmarshal and add all Media to slice,
		// exit if error
		return b.ForEach(func(k, v []byte) error {
			m := Media{}
			err = json.Unmarshal(v, &m)
			if err != nil {
				return err
			}

			list = append(list, m)
			return err
		})
	})

	return
}

// MediaGet retrieves a single instance of Media with
// the given ID
func MediaGet(ID int, db *bolt.DB) (m Media, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get Media bucket, exit if error
		b, err := bucket(mediaBucketName, tx)
		if err != nil {
			return err
		}

		// Get Media by ID, exit if error
		v, err := get(ID, b)
		if err != nil {
			return err
		}

		return json.Unmarshal(v, &m)
	})

	return
}

// MediaCreate persists a new instance of Media
func MediaCreate(m *Media, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Get Media bucket, exit if error
		b, err := bucket(mediaBucketName, tx)
		if err != nil {
			return err
		}

		// Get next ID in sequence and
		// assign to Media
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		m.ID = int(id)

		// Save Media in bucket
		buf, err := json.Marshal(m)
		if err != nil {
			return err
		}

		return b.Put(itob(m.ID), buf)
	})
}

// MediaUpdate updates the properties of an existing
// persisted Media instance
func MediaUpdate(m *Media, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Get Media bucket, exit if error
		b, err := bucket(mediaBucketName, tx)
		if err != nil {
			return err
		}

		// Check if Media with ID exists
		_, err = get(m.ID, b)
		if err != nil {
			return err
		}

		// Replace properties of new with immutable
		// ones of old (none yet)

		// Save Media
		buf, err := json.Marshal(m)
		if err != nil {
			return err
		}

		return b.Put(itob(m.ID), buf)
	})
}
