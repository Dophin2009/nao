package data

import (
	bolt "go.etcd.io/bbolt"
)

// MediaGenre represents a relationship between single
// instances of Media and Genre.
type MediaGenre struct {
	ID      int
	MediaID int
	GenreID int
	Version int
}

// Clean cleans the given MediaGenre for storage
func (ser *MediaGenreService) Clean(e *MediaGenre) (err error) {
	return nil
}

// Validate returns an error if the MediaGenre is
// not valid for the database.
func (ser *MediaGenreService) Validate(e *MediaGenre) (err error) {
	return ser.DB.View(func(tx *bolt.Tx) error {
		// Check if Media with ID specified in new MediaGenre exists
		// Get Media bucket, exit if error
		mb, err := Bucket(MediaBucketName, tx)
		if err != nil {
			return err
		}
		_, err = get(e.MediaID, mb)
		if err != nil {
			return err
		}

		// Check if Genre with ID specified in new MediaGenre exists
		// Get Genre bucket, exit if error
		gb, err := Bucket(GenreBucketName, tx)
		if err != nil {
			return err
		}
		_, err = get(e.GenreID, gb)
		if err != nil {
			return err
		}

		return nil
	})
}

// persistOldProperties maintains certain properties
// of the existing entity in updates
func (ser *MediaGenreService) persistOldProperties(old *MediaGenre, new *MediaGenre) (err error) {
	new.Version = old.Version + 1
	return nil
}

// GetByMedia retrieves a list of instances of MediaGenre
// with the given Media ID.
func (ser *MediaGenreService) GetByMedia(mID int, db *bolt.DB) (list []MediaGenre, err error) {
	return ser.GetFilter(func(mg *MediaGenre) bool {
		return mg.MediaID == mID
	})
}

// GetByGenre retrieves a list of instances of MediaGenre
// with the given Genre ID.
func (ser *MediaGenreService) GetByGenre(gID int, db *bolt.DB) (list []MediaGenre, err error) {
	return ser.GetFilter(func(mg *MediaGenre) bool {
		return mg.GenreID == gID
	})
}
