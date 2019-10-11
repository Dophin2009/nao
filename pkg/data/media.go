package data

import (
	"encoding/json"
	"time"

	bolt "go.etcd.io/bbolt"
)

// Media represents a single instance of a media
type Media struct {
	ID              int
	Titles          []Info
	Synopses        []Info
	Background      []Info
	StartDate       *time.Time
	EndDate         *time.Time
	SeasonPremiered *Season
	Type            string
	Source          string
	Version         int
}

// Season contains information about the quarter
// and year
type Season struct {
	Quarter Quarter
	Year    int
}

// Quarter represents the quarter of the year
// by integer
type Quarter int

const (
	// Winter is the first quarter of the year,
	// encapsulating the months January, February,
	// and March
	Winter Quarter = iota

	// Spring is the second quarter of the year,
	// encapsulating the months April, May, and June
	Spring

	// Summer is the third quarter of the year,
	// encapsulating the months July, August, and
	// September
	Summer

	// Fall is the fouth quarter of the year,
	// encapsulating the months October,
	// November, and December
	Fall
)

const mediaBucketName = "Media"

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

// MediaGetAll retrieves all persisted Media values
func MediaGetAll(db *bolt.DB) (list []Media, err error) {
	return MediaGetFilter(db, func(m *Media) bool { return true })
}

// MediaGetFilter retrieves all persisted Media values
// that pass the filter
func MediaGetFilter(db *bolt.DB, filter func(m *Media) bool) (list []Media, err error) {
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

			if filter(&m) {
				list = append(list, m)
			}
			return err
		})
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
		o, err := get(m.ID, b)
		if err != nil {
			return err
		}

		// Replace properties of new with immutable
		// ones of old
		old := Media{}
		err = json.Unmarshal([]byte(o), &old)
		// Update version
		m.Version = old.Version + 1

		// Save Media
		buf, err := json.Marshal(m)
		if err != nil {
			return err
		}

		return b.Put(itob(m.ID), buf)
	})
}
