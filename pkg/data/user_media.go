package data

import (
	"encoding/json"
	"errors"
	"time"

	bolt "go.etcd.io/bbolt"
)

// UserMedia represents a relationship between a User
// and a Media, containing information about the User's
// opinion on the Media.
type UserMedia struct {
	ID               int
	UserID           int
	MediaID          int
	Status           *WatchStatus
	Priority         *int
	Score            *int
	Recommended      *int
	WatchedInstances []WatchedInstance
	Comments         []Info
	Version          int
}

// WatchedInstance contains information about a single
// watch of some Media.
type WatchedInstance struct {
	Episodes  int
	Ongoing   bool
	StartDate *time.Time
	EndDate   *time.Time
	Comments  []Info
}

// WatchStatus is an enum that represents the
// status of a Media's consumption by a User.
type WatchStatus int

const (
	// Completed means that the User has consumed
	// the Media in its entirety at least once.
	Completed WatchStatus = iota

	// Planning means that the User is planning
	// to consume the Media sometime in the future.
	Planning

	// Dropped means that the User has never
	// consumed the Media in its entirety and
	// abandoned it in the middle somewhere.
	Dropped

	// Hold means the User has begun consuming
	// the Media, but has placed it on hold.
	Hold
)

// UnmarshalJSON defines custom JSON deserialization for
// WatchStatus.
func (ws *WatchStatus) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	value, ok := map[string]WatchStatus{"Completed": Completed,
		"Planning": Planning,
		"Dropped":  Dropped,
		"Hold":     Hold,
	}[s]
	if !ok {
		return errors.New("invalid watch status value '" + s + "'")
	}
	*ws = value
	return nil
}

// MarshalJSON defines custom JSON serialization for
// WatchStatus.
func (ws *WatchStatus) MarshalJSON() (v []byte, err error) {
	value, ok := map[WatchStatus]string{Completed: "Completed",
		Planning: "Planning",
		Dropped:  "Dropped",
		Hold:     "Hold",
	}[*ws]
	if !ok {
		return nil, errors.New("Invalid watch status value")
	}
	return json.Marshal(value)
}

// Validate returns an error if the UserMedia is
// not valid for the database.
func (ser *UserMediaService) Validate(e *UserMedia) (err error) {
	return ser.DB.View(func(tx *bolt.Tx) error {
		// Check if User with ID specified in UserMedia exists
		// Get User bucket, exit if error
		ub, err := Bucket(UserBucketName, tx)
		if err != nil {
			return err
		}
		_, err = get(e.UserID, ub)
		if err != nil {
			return err
		}

		// Check if Media with ID specified in MediaCharacter exists
		// Get Media bucket, exit if error
		mb, err := Bucket(MediaBucketName, tx)
		if err != nil {
			return err
		}
		_, err = get(e.MediaID, mb)
		if err != nil {
			return err
		}

		return nil
	})
}
