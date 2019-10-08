package data

import (
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

// Media represents a single instance of a media
type Media struct {
	ID       int
	Titles   []Info
	Synopsis string
}

// MediaBucketName provides the database bucket name
// for the Media entity
func MediaBucketName() string {
	return "Media"
}

func bucket() []byte {
	return []byte(MediaBucketName())
}

// MediaGetAll retrieves all persisted Media values
func MediaGetAll(db *bolt.DB) (list []Media, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket())
		b.ForEach(func(k, v []byte) error {
			m := Media{}
			err = json.Unmarshal(v, &m)
			list = append(list, m)
			return err
		})

		return nil
	})

	return
}

// MediaGet retrieves a single instance of Media with
// the given ID
func MediaGet(ID int, db *bolt.DB) (media Media, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket())
		v := b.Get(Itob(ID))
		return json.Unmarshal(v, &media)
	})

	return
}

// MediaCreate persists a new instance of Media
func MediaCreate(media *Media, db *bolt.DB) error {
	if media.ID != 0 {
		return fmt.Errorf("media id must be default value")
	}

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket())

		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		media.ID = int(id)

		buf, err := json.Marshal(media)
		if err != nil {
			return err
		}

		return b.Put(Itob(media.ID), buf)
	})
}
