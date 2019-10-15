package data

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

// Genre represents a single instance of a genre
type Genre struct {
	ID           int
	Names        []Info
	Descriptions []Info
	Version      int
}

// Identifier returns the ID of the Genre
func (g *Genre) Identifier() int {
	return g.ID
}

// SetIdentifier sets the ID of the Genre
func (g *Genre) SetIdentifier(ID int) {
	g.ID = ID
}

// Ver returns the verison of the Genre
func (g *Genre) Ver() int {
	return g.Version
}

// UpdateVer increments the version of the
// Character by one
func (g *Genre) UpdateVer() {
	g.Version++
}

// Validate returns an error if the Episode is
// not valid for the database
func (g *Genre) Validate(tx *bolt.Tx) (err error) {
	return nil
}

const genreBucketName = "Genre"

// GenreGet retrieves a single instance of Genre with
// the given ID
func GenreGet(ID int, db *bolt.DB) (g Genre, err error) {
	err = getByID(ID, &g, genreBucketName, db)
	return
}

// GenreGetAll retrieves all persisted Genre values
func GenreGetAll(db *bolt.DB) (list []Genre, err error) {
	return GenreGetFilter(db, func(g *Genre) bool { return true })
}

// GenreGetFilter retrieves all persisted Genre values
func GenreGetFilter(db *bolt.DB, filter func(g *Genre) bool) (list []Genre, err error) {
	ilist, err := getFilter(&Genre{}, func(entity Idenitifiable) (bool, error) {
		g, ok := entity.(*Genre)
		if !ok {
			return false, fmt.Errorf("type assertion failed: entity is not a Genre")
		}
		return filter(g), nil
	}, genreBucketName, db)

	list = make([]Genre, len(ilist))
	for i, g := range ilist {
		list[i] = *g.(*Genre)
	}

	return
}

// GenreCreate persists a new instance of Genre
func GenreCreate(g *Genre, db *bolt.DB) error {
	return create(g, genreBucketName, db)
}

// GenreUpdate updates the properties of an existing
// persisted Genre instance
func GenreUpdate(g *Genre, db *bolt.DB) error {
	return update(g, genreBucketName, db)
}
