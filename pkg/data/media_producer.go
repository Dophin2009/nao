package data

import (
	"encoding/json"

	bolt "go.etcd.io/bbolt"
)

// MediaProducer represents a relationship between single
// instances of Media and Producer
type MediaProducer struct {
	ID         int
	MediaID    int
	ProducerID int
	Role       string
}

const mediaProducerBucketName = "MediaProducer"

// MediaProducerGetAll retrieves all persisted MediaProducer values
func MediaProducerGetAll(db *bolt.DB) (list []MediaProducer, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get MediaProducer bucket, exit if error
		b, err := bucket(mediaProducerBucketName, tx)
		if err != nil {
			return err
		}

		// Unmarshal and add all MediaProducers to slice,
		// exit if error
		b.ForEach(func(k, v []byte) error {
			mp := MediaProducer{}
			err = json.Unmarshal(v, &mp)
			if err != nil {
				return err
			}

			list = append(list, mp)
			return err
		})

		return nil
	})

	return
}

// MediaProducerGet retrieves a single instance of MediaProducer with
// the given ID
func MediaProducerGet(ID int, db *bolt.DB) (mp MediaProducer, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get MediaProducer bucket, exit if error
		b, err := bucket(mediaProducerBucketName, tx)
		if err != nil {
			return err
		}

		// Get MediaProducer by ID, exit if error
		v, err := get(ID, b)
		if err != nil {
			return err
		}
		return json.Unmarshal(v, &mp)
	})

	return
}

// MediaProducerGetByMedia retrieves a list of instances of MediaProducer
// with the given Media ID
func MediaProducerGetByMedia(mID int, db *bolt.DB) (list []MediaProducer, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get MediaProducer bucket, exit if error
		b, err := bucket(mediaProducerBucketName, tx)
		if err != nil {
			return err
		}

		// Get MediaProducer by Media ID
		return b.ForEach(func(k, v []byte) (err error) {
			mp := MediaProducer{}
			err = json.Unmarshal(v, &mp)
			if err != nil {
				return err
			}

			if mp.MediaID == mID {
				list = append(list, mp)
			}
			return nil
		})
	})

	return
}

// MediaProducerGetByProducer retrieves a list of instances of MediaProducer
// with the given Producer ID
func MediaProducerGetByProducer(pID int, db *bolt.DB) (list []MediaProducer, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get MediaProducer bucket, exit if error
		b, err := bucket(mediaProducerBucketName, tx)
		if err != nil {
			return err
		}

		// Get MediaProducer by Producer ID
		return b.ForEach(func(k, v []byte) (err error) {
			mp := MediaProducer{}
			err = json.Unmarshal(v, &mp)
			if err != nil {
				return err
			}

			if mp.MediaID == pID {
				list = append(list, mp)
			}
			return nil
		})
	})

	return
}

// MediaProducerCreate persists a new instance of MediaProducer
func MediaProducerCreate(mp *MediaProducer, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Get MediaProducer bucket, exit if error
		b, err := bucket(mediaProducerBucketName, tx)
		if err != nil {
			return err
		}

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

		// Get next ID in sequence and
		// assign to MediaProducer
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		mp.ID = int(id)

		// Save MediaProducer in bucket
		buf, err := json.Marshal(mp)
		if err != nil {
			return err
		}

		return b.Put(itob(mp.ID), buf)
	})
}

// MediaProducerUpdate updates the properties of an existing
// persisted Producer instance
func MediaProducerUpdate(mp *MediaProducer, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Get MediaProducer bucket, exit if error
		b, err := bucket(mediaProducerBucketName, tx)
		if err != nil {
			return err
		}

		// Check if MediaProducer with ID exists
		_, err = get(mp.ID, b)
		if err != nil {
			return err
		}

		// Replace properties of new with immutable
		// ones of old (none yet)

		// Save MediaProducer
		buf, err := json.Marshal(mp)
		if err != nil {
			return err
		}

		return b.Put(itob(mp.ID), buf)
	})
}
