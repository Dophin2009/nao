package data

import (
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

// Media represents a single instance of a media
type Media struct {
	ID              int
	Titles          []Info
	Synopses        []Info
	Background      []Info
	StartDate       *time.Time
	EndDate         *time.Time
	SeasonPremiered *Season
	Type            string
	Source          string
	Version         int
}

// Identifier returns the ID of the Media
func (m *Media) Identifier() int {
	return m.ID
}

// SetIdentifier sets the ID of the Media
func (m *Media) SetIdentifier(ID int) {
	m.ID = ID
}

// Ver returns the verison of the Media
func (m *Media) Ver() int {
	return m.Version
}

// UpdateVer increments the version of the
// Media by one
func (m *Media) UpdateVer() {
	m.Version++
}

// Validate returns an error if the Media is
// not valid for the database
func (m *Media) Validate(tx *bolt.Tx) (err error) {
	return nil
}

// Season contains information about the quarter
// and year
type Season struct {
	Quarter Quarter
	Year    int
}

// Quarter represents the quarter of the year
// by integer
type Quarter int

const (
	// Winter is the first quarter of the year,
	// encapsulating the months January, February,
	// and March
	Winter Quarter = iota

	// Spring is the second quarter of the year,
	// encapsulating the months April, May, and June
	Spring

	// Summer is the third quarter of the year,
	// encapsulating the months July, August, and
	// September
	Summer

	// Fall is the fouth quarter of the year,
	// encapsulating the months October,
	// November, and December
	Fall
)

const mediaBucketName = "Media"

// MediaGet retrieves a single instance of Media with
// the given ID
func MediaGet(ID int, db *bolt.DB) (m Media, err error) {
	err = getByID(ID, &m, mediaBucketName, db)
	return
}

// MediaGetAll retrieves all persisted Media values
func MediaGetAll(db *bolt.DB) (list []Media, err error) {
	return MediaGetFilter(db, func(m *Media) bool { return true })
}

// MediaGetFilter retrieves all persisted Media values
// that pass the filter
func MediaGetFilter(db *bolt.DB, filter func(m *Media) bool) (list []Media, err error) {
	ilist, err := getFilter(&Media{}, func(entity Idenitifiable) (bool, error) {
		m, ok := entity.(*Media)
		if !ok {
			return false, fmt.Errorf("type assertion failed: entity is not a Media")
		}
		return filter(m), nil
	}, mediaBucketName, db)

	list = make([]Media, len(ilist))
	for i, m := range ilist {
		list[i] = *m.(*Media)
	}

	return
}

// MediaCreate persists a new instance of Media
func MediaCreate(m *Media, db *bolt.DB) error {
	return create(m, mediaBucketName, db)
}

// MediaUpdate updates the properties of an existing
// persisted Media instance
func MediaUpdate(m *Media, db *bolt.DB) error {
	return update(m, mediaBucketName, db)
}
