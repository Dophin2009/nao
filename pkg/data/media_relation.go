package data

import (
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

// MediaRelation represents a relationship between single
// instances of Media and Producer
type MediaRelation struct {
	ID           int
	OwnerID      int
	RelatedID    int
	Relationship string
}

const mediaRelationBucketName = "MediaRelation"

// MediaRelationGetAll retrieves all persisted MediaRelation values
func MediaRelationGetAll(db *bolt.DB) (list []MediaRelation, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get MediaRelation bucket, exit if error
		b, err := bucket(mediaRelationBucketName, tx)
		if err != nil {
			return err
		}

		// Unmarshal and add all MediaRelations to slice,
		// exit if error
		b.ForEach(func(k, v []byte) error {
			mr := MediaRelation{}
			err = json.Unmarshal(v, &mr)
			if err != nil {
				return err
			}

			list = append(list, mr)
			return err
		})

		return nil
	})

	return
}

// MediaRelationGet retrieves a single instance of MediaRelation with
// the given ID
func MediaRelationGet(ID int, db *bolt.DB) (mr MediaRelation, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get MediaRelation bucket, exit if error
		b, err := bucket(mediaRelationBucketName, tx)
		if err != nil {
			return err
		}

		// Get MediaRelation by ID, exit if error
		v, err := get(ID, b)
		if err != nil {
			return err
		}
		return json.Unmarshal(v, &mr)
	})

	return
}

// MediaRelationGetByOwner retrieves a list of instances of MediaRelation
// with the given owning Media ID
func MediaRelationGetByOwner(mID int, db *bolt.DB) (list []MediaRelation, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get MediaRelation bucket, exit if error
		b, err := bucket(mediaRelationBucketName, tx)
		if err != nil {
			return err
		}

		// Get MediaRelation by owning Media ID
		return b.ForEach(func(k, v []byte) (err error) {
			mr := MediaRelation{}
			err = json.Unmarshal(v, &mr)
			if err != nil {
				return err
			}

			if mr.OwnerID == mID {
				list = append(list, mr)
			}
			return nil
		})
	})

	return
}

// MediaRelationGetByRelated retrieves a list of instances of MediaRelation
// with the given related Media ID
func MediaRelationGetByRelated(mID int, db *bolt.DB) (list []MediaRelation, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get MediaRelation bucket, exit if error
		b, err := bucket(mediaRelationBucketName, tx)
		if err != nil {
			return err
		}

		// Get MediaRelation by related Media ID
		return b.ForEach(func(k, v []byte) (err error) {
			mr := MediaRelation{}
			err = json.Unmarshal(v, &mr)
			if err != nil {
				return err
			}

			if mr.RelatedID == mID {
				list = append(list, mr)
			}
			return nil
		})
	})

	return
}

// MediaRelationCreate persists a new instance of MediaRelation
func MediaRelationCreate(mr *MediaRelation, db *bolt.DB) error {
	// ID must be 0
	if mr.ID != 0 {
		return fmt.Errorf("mediaRelation id must be default value")
	}

	return db.Update(func(tx *bolt.Tx) error {
		// Get MediaRelation bucket, exit if error
		b, err := bucket(mediaRelationBucketName, tx)
		if err != nil {
			return err
		}

		// Check if Media with ID specified in new MediaRelation exists
		// Get Media bucket, exit if error
		mb, err := bucket(mediaBucketName, tx)
		if err != nil {
			return err
		}
		_, err = get(mr.OwnerID, mb)
		if err != nil {
			return err
		}

		// Check if Producer with ID specified in new MediaRelation exists
		// Get Producer bucket, exit if error
		pb, err := bucket(producerBucketName, tx)
		if err != nil {
			return err
		}
		_, err = get(mr.RelatedID, pb)
		if err != nil {
			return err
		}

		// Get next ID in sequence and
		// assign to MediaRelation
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		mr.ID = int(id)

		// Save MediaRelation in bucket
		buf, err := json.Marshal(mr)
		if err != nil {
			return err
		}

		return b.Put(itob(mr.ID), buf)
	})
}
