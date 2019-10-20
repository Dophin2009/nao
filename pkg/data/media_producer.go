package data

import (
	"strings"

	bolt "go.etcd.io/bbolt"
)

// MediaProducer represents a relationship between single
// instances of Media and Producer
type MediaProducer struct {
	ID         int
	MediaID    int
	ProducerID int
	Role       string
	Version    int
}

// Clean cleans the given MediaProducer for storage
func (ser *MediaProducerService) Clean(e *MediaProducer) (err error) {
	e.Role = strings.Trim(e.Role, " ")
	return nil
}

// Validate returns an error if the MediaProducer is
// not valid for the database.
func (ser *MediaProducerService) Validate(e *MediaProducer) (err error) {
	return ser.DB.View(func(tx *bolt.Tx) error {
		// Check if Media with ID specified in new MediaProducer exists
		// Get Media bucket, exit if error
		mb, err := Bucket(MediaBucketName, tx)
		if err != nil {
			return err
		}
		_, err = get(e.MediaID, mb)
		if err != nil {
			return err
		}

		// Check if Producer with ID specified in new MediaProducer exists
		// Get Producer bucket, exit if error
		pb, err := Bucket(ProducerBucketName, tx)
		if err != nil {
			return err
		}
		_, err = get(e.ProducerID, pb)
		if err != nil {
			return err
		}

		return nil
	})
}

// GetByMedia retrieves a list of instances of MediaProducer
// with the given Media ID.
func (ser *MediaProducerService) GetByMedia(mID int, db *bolt.DB) (list []MediaProducer, err error) {
	return ser.GetFilter(func(mp *MediaProducer) bool {
		return mp.MediaID == mID
	})
}

// GetByProducer retrieves a list of instances of MediaProducer
// with the given Producer ID.
func (ser *MediaProducerService) GetByProducer(pID int, db *bolt.DB) (list []MediaProducer, err error) {
	return ser.GetFilter(func(mp *MediaProducer) bool {
		return mp.ProducerID == pID
	})
}
