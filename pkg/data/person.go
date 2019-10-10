package data

import (
	"encoding/json"

	bolt "go.etcd.io/bbolt"
)

// Person represents a single person
type Person struct {
	ID          int
	Names       []Info
	Information []Info
	Version     int
}

// PersonBucketName provides the database bucket name
// for the Person entity
const personBucketName = "Person"

// PersonGet retrieves a single instance of Person with
// the given ID
func PersonGet(ID int, db *bolt.DB) (p Person, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get Person bucket, exit if error
		b, err := bucket(personBucketName, tx)
		if err != nil {
			return err
		}

		// Get Person by ID, exit if error
		v, err := get(ID, b)
		if err != nil {
			return err
		}
		return json.Unmarshal(v, &p)
	})

	return
}

// PersonGetAll retrieves all persisted Person values
func PersonGetAll(db *bolt.DB) (list []Person, err error) {
	return PersonGetFilter(db, func(p *Person) bool { return true })
}

// PersonGetFilter retrieves all persisted Person values
// that pass the filter
func PersonGetFilter(db *bolt.DB, filter func(p *Person) bool) (list []Person, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get Person bucket, exit if error
		b, err := bucket(personBucketName, tx)
		if err != nil {
			return err
		}

		// Unmarshal and add Persons to slice,
		// exit if error
		return b.ForEach(func(k, v []byte) error {
			var p Person
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

// PersonCreate persists a new instance of Person
func PersonCreate(p *Person, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Get Person bucket, exit if error
		b, err := bucket(personBucketName, tx)
		if err != nil {
			return err
		}

		// Get next ID in sequence and assign to Person
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		p.ID = int(id)

		// Save Person in bucket
		buf, err := json.Marshal(p)
		if err != nil {
			return err
		}

		return b.Put(itob(p.ID), buf)
	})
}

// PersonUpdate updates the properties of an existing
// persisted Person instance
func PersonUpdate(p *Person, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Get Person bucket, exit if error
		b, err := bucket(personBucketName, tx)
		if err != nil {
			return err
		}

		// Check if Person with ID exists
		o, err := get(p.ID, b)
		if err != nil {
			return err
		}

		// Replace properties of new with immutable
		// ones of old
		old := Person{}
		err = json.Unmarshal([]byte(o), &old)
		// Update version
		p.Version = old.Version + 1

		// Save Person
		buf, err := json.Marshal(p)
		if err != nil {
			return err
		}

		return b.Put(itob(p.ID), buf)
	})
}
