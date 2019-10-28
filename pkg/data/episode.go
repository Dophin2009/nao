package data

import (
	"time"

	bolt "go.etcd.io/bbolt"
)

// Episode represents a single episode or chapter
// for some media.
type Episode struct {
	ID       int
	MediaID  int
	Titles   []Info
	Date     *time.Time
	Synopses []Info
	Duration *uint
	Filler   bool
	Recap    bool
	Version  int
}

// Clean cleans the given Episode for storage
func (ser *EpisodeService) Clean(e *Episode) (err error) {
	if err = infoListClean(e.Titles); err != nil {
		return err
	}
	if err = infoListClean(e.Synopses); err != nil {
		return err
	}
	return nil
}

// Validate returns an error if the Episode is
// not valid for the database.
func (ser *EpisodeService) Validate(e *Episode) (err error) {
	return nil
}

// persistOldProperties maintains certain properties
// of the existing entity in updates
func (ser *EpisodeService) persistOldProperties(old *Episode, new *Episode) (err error) {
	new.Version = old.Version + 1
	return nil
}

// GetByMedia retrieves a list of instances of Episode
// with the given Media ID.
func (ser *EpisodeService) GetByMedia(mID int, db *bolt.DB) (list []Episode, err error) {
	return ser.GetFilter(func(ep *Episode) bool {
		return ep.MediaID == mID
	})
}
