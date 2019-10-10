package data

import (
	"encoding/json"

	bolt "go.etcd.io/bbolt"
)

// Character represents a single  character
type Character struct {
	ID          int
	Names       []Info
	Information []Info
	Version     int
}

// CharacterBucketName provides the database bucket name
// for the Character entity
const characterBucketName = "Character"

// CharacterGet retrieves a single instance of Character with
// the given ID
func CharacterGet(ID int, db *bolt.DB) (p Character, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get Character bucket, exit if error
		b, err := bucket(characterBucketName, tx)
		if err != nil {
			return err
		}

		// Get Character by ID, exit if error
		v, err := get(ID, b)
		if err != nil {
			return err
		}
		return json.Unmarshal(v, &p)
	})

	return
}

// CharacterGetAll retrieves all persisted Character values
func CharacterGetAll(db *bolt.DB) (list []Character, err error) {
	return CharacterGetFilter(db, func(p *Character) bool { return true })
}

// CharacterGetFilter retrieves all persisted Character values
// that pass the filter
func CharacterGetFilter(db *bolt.DB, filter func(p *Character) bool) (list []Character, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get Character bucket, exit if error
		b, err := bucket(characterBucketName, tx)
		if err != nil {
			return err
		}

		// Unmarshal and add Characters to slice,
		// exit if error
		return b.ForEach(func(k, v []byte) error {
			var p Character
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

// CharacterCreate persists a new instance of Character
func CharacterCreate(p *Character, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Get Character bucket, exit if error
		b, err := bucket(characterBucketName, tx)
		if err != nil {
			return err
		}

		// Get next ID in sequence and assign to Character
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		p.ID = int(id)

		// Save Character in bucket
		buf, err := json.Marshal(p)
		if err != nil {
			return err
		}

		return b.Put(itob(p.ID), buf)
	})
}

// CharacterUpdate updates the properties of an existing
// persisted Character instance
func CharacterUpdate(p *Character, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Get Character bucket, exit if error
		b, err := bucket(characterBucketName, tx)
		if err != nil {
			return err
		}

		// Check if Character with ID exists
		o, err := get(p.ID, b)
		if err != nil {
			return err
		}

		// Replace properties of new with immutable
		// ones of old
		old := Character{}
		err = json.Unmarshal([]byte(o), &old)
		// Update version
		p.Version = old.Version + 1

		// Save Character
		buf, err := json.Marshal(p)
		if err != nil {
			return err
		}

		return b.Put(itob(p.ID), buf)
	})
}
