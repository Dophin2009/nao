package data

import (
	"fmt"

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

// Identifier returns the ID of the MediaProducer
func (mp *MediaProducer) Identifier() int {
	return mp.ID
}

// SetIdentifier sets the ID of the MediaProducer
func (mp *MediaProducer) SetIdentifier(ID int) {
	mp.ID = ID
}

// Ver returns the verison of the MediaProducer
func (mp *MediaProducer) Ver() int {
	return mp.Version
}

// UpdateVer increments the version of the
// Character by one
func (mp *MediaProducer) UpdateVer() {
	mp.Version++
}

// Validate returns an error if the MediaProducer is
// not valid for the database
func (mp *MediaProducer) Validate(tx *bolt.Tx) (err error) {
	// Check if Media with ID specified in new MediaProducer exists
	// Get Media bucket, exit if error
	mb, err := bucket(mediaBucketName, tx)
	if err != nil {
		return err
	}
	_, err = get(mp.MediaID, mb)
	if err != nil {
		return err
	}

	// Check if Producer with ID specified in new MediaProducer exists
	// Get Producer bucket, exit if error
	pb, err := bucket(producerBucketName, tx)
	if err != nil {
		return err
	}
	_, err = get(mp.ProducerID, pb)
	if err != nil {
		return err
	}

	return nil
}

const mediaProducerBucketName = "MediaProducer"

// MediaProducerGet retrieves a single instance of MediaProducer with
// the given ID
func MediaProducerGet(ID int, db *bolt.DB) (mp MediaProducer, err error) {
	err = getByID(ID, &mp, mediaProducerBucketName, db)
	return
}

// MediaProducerGetAll retrieves all persisted MediaProducer values
func MediaProducerGetAll(db *bolt.DB) (list []MediaProducer, err error) {
	return MediaProducerGetFilter(db, func(mp *MediaProducer) bool { return true })
}

// MediaProducerGetByMedia retrieves a list of instances of MediaProducer
// with the given Media ID
func MediaProducerGetByMedia(mID int, db *bolt.DB) (list []MediaProducer, err error) {
	return MediaProducerGetFilter(db, func(mp *MediaProducer) bool {
		return mp.MediaID == mID
	})
}

// MediaProducerGetByProducer retrieves a list of instances of MediaProducer
// with the given Producer ID
func MediaProducerGetByProducer(pID int, db *bolt.DB) (list []MediaProducer, err error) {
	return MediaProducerGetFilter(db, func(mp *MediaProducer) bool {
		return mp.ProducerID == pID
	})
}

// MediaProducerGetFilter retrieves all persisted MediaProducer values
func MediaProducerGetFilter(db *bolt.DB, filter func(mp *MediaProducer) bool) (list []MediaProducer, err error) {
	ilist, err := getFilter(&MediaProducer{}, func(entity Idenitifiable) (bool, error) {
		mp, ok := entity.(*MediaProducer)
		if !ok {
			return false, fmt.Errorf("type assertion failed: entity is not a MediaProducer")
		}
		return filter(mp), nil
	}, mediaProducerBucketName, db)

	list = make([]MediaProducer, len(ilist))
	for i, mp := range ilist {
		list[i] = *mp.(*MediaProducer)
	}

	return
}

// MediaProducerCreate persists a new instance of MediaProducer
func MediaProducerCreate(mp *MediaProducer, db *bolt.DB) error {
	return create(mp, mediaProducerBucketName, db)
}

// MediaProducerUpdate updates the properties of an existing
// persisted Producer instance
func MediaProducerUpdate(mp *MediaProducer, db *bolt.DB) error {
	return update(mp, mediaProducerBucketName, db)
}
