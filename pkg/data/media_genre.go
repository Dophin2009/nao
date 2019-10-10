package data

import (
	"encoding/json"

	bolt "go.etcd.io/bbolt"
)

// MediaGenre represents a relationship between single
// instances of Media and Genre
type MediaGenre struct {
	ID      int
	MediaID int
	GenreID int
	Version int
}

const mediaGenreBucketName = "MediaGenre"

// MediaGenreGet retrieves a single instance of MediaGenre with
// the given ID
func MediaGenreGet(ID int, db *bolt.DB) (mg MediaGenre, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get MediaGenre bucket, exit if error
		b, err := bucket(mediaGenreBucketName, tx)
		if err != nil {
			return err
		}

		// Get MediaGenre by ID, exit if error
		v, err := get(ID, b)
		if err != nil {
			return err
		}
		return json.Unmarshal(v, &mg)
	})

	return
}

// MediaGenreGetAll retrieves all persisted MediaGenre values
func MediaGenreGetAll(db *bolt.DB) (list []MediaGenre, err error) {
	return MediaGenreGetFilter(db, func(mg *MediaGenre) bool { return true })
}

// MediaGenreGetByMedia retrieves a list of instances of MediaGenre
// with the given Media ID
func MediaGenreGetByMedia(mID int, db *bolt.DB) (list []MediaGenre, err error) {
	return MediaGenreGetFilter(db, func(mg *MediaGenre) bool {
		return mg.MediaID == mID
	})
}

// MediaGenreGetByGenre retrieves a list of instances of MediaGenre
// with the given Genre ID
func MediaGenreGetByGenre(gID int, db *bolt.DB) (list []MediaGenre, err error) {
	return MediaGenreGetFilter(db, func(mg *MediaGenre) bool {
		return mg.GenreID == gID
	})
}

// MediaGenreGetFilter retrieves all persisted MediaGenre values
func MediaGenreGetFilter(db *bolt.DB, filter func(mg *MediaGenre) bool) (list []MediaGenre, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get MediaGenre bucket, exit if error
		b, err := bucket(mediaGenreBucketName, tx)
		if err != nil {
			return err
		}

		// Unmarshal and add all MediaGenres to slice,
		// exit if error
		b.ForEach(func(k, v []byte) error {
			mg := MediaGenre{}
			err = json.Unmarshal(v, &mg)
			if err != nil {
				return err
			}

			if filter(&mg) {
				list = append(list, mg)
			}
			return err
		})

		return nil
	})

	return
}

// MediaGenreCreate persists a new instance of MediaGenre
func MediaGenreCreate(mg *MediaGenre, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Get MediaGenre bucket, exit if error
		b, err := bucket(mediaGenreBucketName, tx)
		if err != nil {
			return err
		}

		// Check if MediaGenre properties are valid
		err = MediaGenreCheckRelatedIDs(mg, tx)
		if err != nil {
			return err
		}

		// Get next ID in sequence and
		// assign to MediaGenre
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		mg.ID = int(id)

		// Save MediaGenre in bucket
		buf, err := json.Marshal(mg)
		if err != nil {
			return err
		}

		return b.Put(itob(mg.ID), buf)
	})
}

// MediaGenreUpdate updates the properties of an existing
// persisted Genre instance
func MediaGenreUpdate(mg *MediaGenre, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Get MediaGenre bucket, exit if error
		b, err := bucket(mediaGenreBucketName, tx)
		if err != nil {
			return err
		}

		// Check if MediaGenre with ID exists
		o, err := get(mg.ID, b)
		if err != nil {
			return err
		}

		// Check if MediaGenre properties are valid
		err = MediaGenreCheckRelatedIDs(mg, tx)
		if err != nil {
			return err
		}

		// Replace properties of new with immutable
		// ones of old
		old := MediaGenre{}
		err = json.Unmarshal([]byte(o), &old)
		// Update version
		mg.Version = old.Version + 1

		// Save MediaGenre
		buf, err := json.Marshal(mg)
		if err != nil {
			return err
		}

		return b.Put(itob(mg.ID), buf)
	})
}

// MediaGenreCheckRelatedIDs checks if the entities specified
// by the related entity IDs exist for a MediaGenre
func MediaGenreCheckRelatedIDs(mg *MediaGenre, tx *bolt.Tx) (err error) {
	// Check if Media with ID specified in new MediaGenre exists
	// Get Media bucket, exit if error
	mb, err := bucket(mediaBucketName, tx)
	if err != nil {
		return err
	}
	_, err = get(mg.MediaID, mb)
	if err != nil {
		return err
	}

	// Check if Genre with ID specified in new MediaGenre exists
	// Get Genre bucket, exit if error
	gb, err := bucket(genreBucketName, tx)
	if err != nil {
		return err
	}
	_, err = get(mg.GenreID, gb)
	if err != nil {
		return err
	}

	return nil
}
