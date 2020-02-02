package data

import (
	"errors"
	"fmt"
	"time"

	json "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
)

// UserMedia represents a relationship between a User and a Media, containing
// information about the User's opinion on the Media.
type UserMedia struct {
	ID               int
	UserID           int
	MediaID          int
	Status           *WatchStatus
	Priority         *int
	Score            *int
	Recommended      *int
	WatchedInstances []WatchedInstance
	Comments         map[string]string
	UserMediaListIDs []int
	Version          int
}

// Iden returns the ID.
func (um *UserMedia) Iden() int {
	return um.ID
}

// WatchedInstance contains information about a single watch of some Media.
type WatchedInstance struct {
	Episodes  int
	Ongoing   bool
	StartDate *time.Time
	EndDate   *time.Time
	Comments  map[string]string
}

// WatchStatus is an enum that represents the status of a Media's consumption
// by a User.
type WatchStatus int

const (
	// Completed means that the User has consumed the Media in its entirety at
	// least once.
	Completed WatchStatus = iota

	// Planning means that the User is planning to consume the Media sometime in
	// the future.
	Planning

	// Dropped means that the User has never consumed the Media in its entirety
	// and abandoned it in the middle somewhere.
	Dropped

	// Hold means the User has begun consuming the Media, but has placed it on
	// hold.
	Hold
)

// UnmarshalJSON defines custom JSON deserialization for WatchStatus.
func (ws *WatchStatus) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}

	value, ok := map[string]WatchStatus{
		"Completed": Completed,
		"Planning":  Planning,
		"Dropped":   Dropped,
		"Hold":      Hold,
	}[s]
	if !ok {
		return fmt.Errorf("watch status value %q: %w", s, errInvalid)
	}
	*ws = value
	return nil
}

// MarshalJSON defines custom JSON serialization for WatchStatus.
func (ws *WatchStatus) MarshalJSON() (v []byte, err error) {
	value, ok := map[WatchStatus]string{
		Completed: "Completed",
		Planning:  "Planning",
		Dropped:   "Dropped",
		Hold:      "Hold",
	}[*ws]
	if !ok {
		return nil, fmt.Errorf("watch status value %d: %w", *ws, err)
	}

	v, err = json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONMarshal, err)
	}

	return v, nil
}

// UserMediaBucket is the name of the database bucket for UserMedia.
const UserMediaBucket = "UserMedia"

// UserMediaService performs operations on UserMedia.
type UserMediaService struct {
	DB *bolt.DB
	Service
}

// Create persists the given UserMedia.
func (ser *UserMediaService) Create(um *UserMedia) error {
	return Create(um, ser)
}

// Update rumlaces the value of the UserMedia with the given ID.
func (ser *UserMediaService) Update(um *UserMedia) error {
	return Update(um, ser)
}

// Delete deletes the UserMedia with the given ID.
func (ser *UserMediaService) Delete(id int) error {
	return Delete(id, ser)
}

// GetAll retrieves all persisted values of UserMedia.
func (ser *UserMediaService) GetAll(first *int, skip *int) ([]*UserMedia, error) {
	vlist, err := GetAll(ser, first, skip)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to UserMedia: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of UserMedia that pass the filter.
func (ser *UserMediaService) GetFilter(
	first *int, skip *int, keep func(um *UserMedia) bool,
) ([]*UserMedia, error) {
	vlist, err := GetFilter(ser, first, skip, func(m Model) bool {
		um, err := ser.AssertType(m)
		if err != nil {
			return false
		}
		return keep(um)
	})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to UserMedia: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted UserMedia with the given ID.
func (ser *UserMediaService) GetByID(id int) (*UserMedia, error) {
	m, err := GetByID(id, ser)
	if err != nil {
		return nil, err
	}

	um, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return um, nil
}

// GetByUser retrieves the persisted UserMedia with the given User ID.
func (ser *UserMediaService) GetByUser(
	uID int, first *int, skip *int,
) ([]*UserMedia, error) {
	return ser.GetFilter(first, skip, func(um *UserMedia) bool {
		return um.UserID == uID
	})
}

// GetByMedia retrieves the persisted UserMedia with the given Media ID.
func (ser *UserMediaService) GetByMedia(
	mID int, first *int, skip *int,
) ([]*UserMedia, error) {
	return ser.GetFilter(first, skip, func(um *UserMedia) bool {
		return um.MediaID == mID
	})
}

// Database returns the database reference.
func (ser *UserMediaService) Database() *bolt.DB {
	return ser.DB
}

// Bucket returns the name of the bucket for UserMedia.
func (ser *UserMediaService) Bucket() string {
	return UserMediaBucket
}

// Clean cleans the given UserMedia for storage.
func (ser *UserMediaService) Clean(m Model) error {
	_, err := ser.AssertType(m)
	return err
}

// Validate returns an error if the UserMedia is not valid for the database.
func (ser *UserMediaService) Validate(m Model) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	return ser.DB.View(func(tx *bolt.Tx) error {
		// Check if User with ID specified in UserMedia exists
		// Get User bucket, exit if error
		ub, err := Bucket(UserBucket, tx)
		if err != nil {
			return fmt.Errorf("%s %q: %w", errmsgBucketOpen, UserBucket, err)
		}
		_, err = get(e.UserID, ub)
		if err != nil {
			return fmt.Errorf("failed to get User with ID %d: %w", e.UserID, err)
		}

		// Check if Media with ID specified in MediaCharacter exists
		// Get Media bucket, exit if error
		mb, err := Bucket(MediaBucket, tx)
		if err != nil {
			return fmt.Errorf("%s %q: %w", errmsgBucketOpen, MediaBucket, err)
		}
		_, err = get(e.MediaID, mb)
		if err != nil {
			return fmt.Errorf("failed to get Media with ID %d: %w", e.MediaID, err)
		}

		// Check if UserMediaLists with IDs specified in UserMedia exists
		// Get User bucket, exit if error
		umlb, err := Bucket(UserMediaListBucket, tx)
		if err != nil {
			return fmt.Errorf("%s %q: %w", errmsgBucketOpen, UserMediaListBucket, err)
		}
		for _, listID := range e.UserMediaListIDs {
			_, err = get(listID, umlb)
			if err != nil {
				return fmt.Errorf("failed to get UserMediaList with ID %d: %w", listID, err)
			}
		}

		return nil
	})
}

// Initialize sets initial values for some properties.
func (ser *UserMediaService) Initialize(m Model, id int) error {
	md, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	md.ID = id
	md.Version = 0
	return nil
}

// PersistOldProperties maintains certain properties of the existing UserMedia
// in updates.
func (ser *UserMediaService) PersistOldProperties(n Model, o Model) error {
	nm, err := ser.AssertType(n)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	om, err := ser.AssertType(o)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	nm.Version = om.Version + 1
	return nil
}

// Marshal transforms the given UserMedia into JSON.
func (ser *UserMediaService) Marshal(m Model) ([]byte, error) {
	um, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	v, err := json.Marshal(um)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONMarshal, err)
	}

	return v, nil
}

// Unmarshal parses the given JSON into UserMedia.
func (ser *UserMediaService) Unmarshal(buf []byte) (Model, error) {
	var um UserMedia
	err := json.Unmarshal(buf, &um)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &um, nil
}

// AssertType exposes the given Model as a UserMedia.
func (ser *UserMediaService) AssertType(m Model) (*UserMedia, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	um, ok := m.(*UserMedia)
	if !ok {
		return nil, fmt.Errorf("model: %w", errors.New("not of UserMedia type"))
	}
	return um, nil
}

// mapfromModel returns a list of UserMedia type asserted from the given list
// of Model.
func (ser *UserMediaService) mapFromModel(vlist []Model) ([]*UserMedia, error) {
	list := make([]*UserMedia, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}
