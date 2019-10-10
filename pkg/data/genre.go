package data

import (
	"encoding/json"

	bolt "go.etcd.io/bbolt"
)

// Genre represents a single instance of a genre
type Genre struct {
	ID           int
	Names        []Info
	Descriptions []Info
}

const genreBucketName = "Genre"

// GenreGetAll retrieves all persisted Genre values
func GenreGetAll(db *bolt.DB) (list []Genre, err error) {
	return GenreGetFilter(db, func(g *Genre) bool { return true })
}

// GenreGetFilter retrieves all persisted Genre values
func GenreGetFilter(db *bolt.DB, filter func(g *Genre) bool) (list []Genre, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get Genre bucket, exit if error
		b, err := bucket(genreBucketName, tx)

		if err != nil {
			return err
		}

		// Unmarshal and add all Genre to slice,
		// exit if error
		return b.ForEach(func(k, v []byte) error {
			g := Genre{}
			err = json.Unmarshal(v, &g)
			if err != nil {
				return err
			}

			if filter(&g) {
				list = append(list, g)
			}
			return err
		})
	})

	return
}

// GenreGet retrieves a single instance of Genre with
// the given ID
func GenreGet(ID int, db *bolt.DB) (g Genre, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get Genre bucket, exit if error
		b, err := bucket(genreBucketName, tx)
		if err != nil {
			return err
		}

		// Get Genre by ID, exit if error
		v, err := get(ID, b)
		if err != nil {
			return err
		}

		return json.Unmarshal(v, &g)
	})

	return
}

// GenreCreate persists a new instance of Genre
func GenreCreate(g *Genre, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Get Genre bucket, exit if error
		b, err := bucket(genreBucketName, tx)
		if err != nil {
			return err
		}

		// Get next ID in sequence and
		// assign to Genre
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		g.ID = int(id)

		// Save Genre in bucket
		buf, err := json.Marshal(g)
		if err != nil {
			return err
		}

		return b.Put(itob(g.ID), buf)
	})
}

// GenreUpdate updates the properties of an existing
// persisted Genre instance
func GenreUpdate(g *Genre, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Get Genre bucket, exit if error
		b, err := bucket(genreBucketName, tx)
		if err != nil {
			return err
		}

		// Check if Genre with ID exists
		_, err = get(g.ID, b)
		if err != nil {
			return err
		}

		// Replace properties of new with immutable
		// ones of old (none yet)

		// Save Genre
		buf, err := json.Marshal(g)
		if err != nil {
			return err
		}

		return b.Put(itob(g.ID), buf)
	})
}
