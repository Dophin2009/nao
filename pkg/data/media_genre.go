package data

import (
	"fmt"

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

// Identifier returns the ID of the MediaGenre
func (mg *MediaGenre) Identifier() int {
	return mg.ID
}

// SetIdentifier sets the ID of the MediaGenre
func (mg *MediaGenre) SetIdentifier(ID int) {
	mg.ID = ID
}

// Ver returns the verison of the MediaGenre
func (mg *MediaGenre) Ver() int {
	return mg.Version
}

// UpdateVer increments the version of the
// Character by one
func (mg *MediaGenre) UpdateVer() {
	mg.Version++
}

// Validate returns an error if the MediaGenre is
// not valid for the database
// by the related entity IDs exist for a MediaGenre
func (mg *MediaGenre) Validate(tx *bolt.Tx) (err error) {
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

const mediaGenreBucketName = "MediaGenre"

// MediaGenreGet retrieves a single instance of MediaGenre with
// the given ID
func MediaGenreGet(ID int, db *bolt.DB) (mg MediaGenre, err error) {
	err = getByID(ID, &mg, mediaGenreBucketName, db)
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
	ilist, err := getFilter(&MediaGenre{}, func(entity Idenitifiable) (bool, error) {
		mg, ok := entity.(*MediaGenre)
		if !ok {
			return false, fmt.Errorf("type assertion failed: entity is not a MediaGenre")
		}
		return filter(mg), nil
	}, mediaGenreBucketName, db)

	list = make([]MediaGenre, len(ilist))
	for i, mg := range ilist {
		list[i] = *mg.(*MediaGenre)
	}

	return
}

// MediaGenreCreate persists a new instance of MediaGenre
func MediaGenreCreate(mg *MediaGenre, db *bolt.DB) error {
	return create(mg, mediaGenreBucketName, db)
}

// MediaGenreUpdate updates the properties of an existing
// persisted Genre instance
func MediaGenreUpdate(mg *MediaGenre, db *bolt.DB) error {
	return update(mg, mediaGenreBucketName, db)
}
