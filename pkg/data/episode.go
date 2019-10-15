package data

import (
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

// Episode represents a single episode or chapter
// for some media
type Episode struct {
	ID       int
	MediaID  int
	Titles   []Info
	Date     *time.Time
	Synopses []Info
	Duration uint
	Filler   bool
	Recap    bool
	Version  int
}

// Identifier returns the ID of the Episode
func (ep *Episode) Identifier() int {
	return ep.ID
}

// SetIdentifier sets the ID of the Episode
func (ep *Episode) SetIdentifier(ID int) {
	ep.ID = ID
}

// Ver returns the verison of the Episode
func (ep *Episode) Ver() int {
	return ep.Version
}

// UpdateVer increments the version of the
// Character by one
func (ep *Episode) UpdateVer() {
	ep.Version++
}

// Validate returns an error if the Episode is
// not valid for the database
func (ep *Episode) Validate(tx *bolt.Tx) (err error) {
	return nil
}

const episodeBucketName = "Episode"

// EpisodeGet retrieves a single instance of Episode with
// the given ID
func EpisodeGet(ID int, db *bolt.DB) (ep Episode, err error) {
	err = getByID(ID, &ep, episodeBucketName, db)
	return
}

// EpisodeGetAll retrieves all persisted Episode values
func EpisodeGetAll(db *bolt.DB) (list []Episode, err error) {
	return EpisodeGetFilter(db, func(ep *Episode) bool { return true })
}

// EpisodeGetFilter retrieves all persisted Episode values
// that pass the filter
func EpisodeGetFilter(db *bolt.DB, filter func(ep *Episode) bool) (list []Episode, err error) {
	ilist, err := getFilter(&Episode{}, func(entity Idenitifiable) (bool, error) {
		ep, ok := entity.(*Episode)
		if !ok {
			return false, fmt.Errorf("type assertion failed: entity is not a Episode")
		}
		return filter(ep), nil
	}, episodeBucketName, db)

	list = make([]Episode, len(ilist))
	for i, m := range ilist {
		list[i] = *m.(*Episode)
	}

	return
}

// EpisodeGetByMedia retrieves a list of instances of Episode
// with the given Media ID
func EpisodeGetByMedia(mID int, db *bolt.DB) (list []Episode, err error) {
	return EpisodeGetFilter(db, func(ep *Episode) bool {
		return ep.MediaID == mID
	})
}

// EpisodeCreate persists a new instance of Episode
func EpisodeCreate(ep *Episode, db *bolt.DB) error {
	return create(ep, episodeBucketName, db)
}

// EpisodeUpdate updates the properties of an existing
// persisted Producer instance
func EpisodeUpdate(ep *Episode, db *bolt.DB) error {
	return update(ep, episodeBucketName, db)
}
