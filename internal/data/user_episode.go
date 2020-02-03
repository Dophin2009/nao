package data

import (
	"errors"
	"fmt"

	json "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
)

// UserEpisode represents a relationship between a User and an Episode,
// containing information about the User's opinion on the Episode.
type UserEpisode struct {
	UserID    int
	EpisodeID int
	Score     *int
	Comments  []Title
	Meta      ModelMetadata
}

// Metadata returns Meta.
func (uep *UserEpisode) Metadata() *ModelMetadata {
	return &uep.Meta
}

// UserEpisodeBucket is the name of the database bucket for UserEpisode.
const UserEpisodeBucket = "UserEpisode"

// UserEpisodeService performs operations on UserEpisode.
type UserEpisodeService struct {
	DB *bolt.DB
}

// Create persists the given UserEpisode.
func (ser *UserEpisodeService) Create(uep *UserEpisode) error {
	return Create(uep, ser)
}

// Update rueplaces the value of the UserEpisode with the given ID.
func (ser *UserEpisodeService) Update(uep *UserEpisode) error {
	return Update(uep, ser)
}

// Delete deletes the UserEpisode with the given ID.
func (ser *UserEpisodeService) Delete(id int) error {
	return Delete(id, ser)
}

// GetAll retrieves all persisted values of UserEpisode.
func (ser *UserEpisodeService) GetAll(first *int, skip *int) ([]*UserEpisode, error) {
	vlist, err := GetAll(ser, first, skip)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to UserEpisode: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of UserEpisode that pass the
// filter.
func (ser *UserEpisodeService) GetFilter(
	first *int, skip *int, keep func(uep *UserEpisode) bool,
) ([]*UserEpisode, error) {
	vlist, err := GetFilter(ser, first, skip, func(m Model) bool {
		uep, err := ser.AssertType(m)
		if err != nil {
			return false
		}
		return keep(uep)
	})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to UserEpisode: %w", err)
	}
	return list, nil
}

// GetMultiple retrieves the persisted UserEpisode values specified by the
// given IDs that pass the filter.
func (ser *UserEpisodeService) GetMultiple(
	ids []int, first *int, skip *int, keep func(uep *UserEpisode) bool,
) ([]*UserEpisode, error) {
	vlist, err := GetMultiple(ser, ids, first, skip, func(m Model) bool {
		uep, err := ser.AssertType(m)
		if err != nil {
			return false
		}
		return keep(uep)
	})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to UserEpisodes: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted UserEpisode with the given ID.
func (ser *UserEpisodeService) GetByID(id int) (*UserEpisode, error) {
	m, err := GetByID(id, ser)
	if err != nil {
		return nil, err
	}

	uep, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return uep, nil
}

// GetByUser retrieves the persisted UserEpisode with the given User ID.
func (ser *UserEpisodeService) GetByUser(
	uID int, first *int, skip *int,
) ([]*UserEpisode, error) {
	return ser.GetFilter(first, skip, func(uep *UserEpisode) bool {
		return uep.UserID == uID
	})
}

// GetByEpisode retrieves the persisted UserEpisode with the given Episode ID.
func (ser *UserEpisodeService) GetByEpisode(
	epID int, first *int, skip *int,
) ([]*UserEpisode, error) {
	return ser.GetFilter(first, skip, func(uep *UserEpisode) bool {
		return uep.EpisodeID == epID
	})
}

// Database returns the database reference.
func (ser *UserEpisodeService) Database() *bolt.DB {
	return ser.DB
}

// Bucket returns the name of the bucket for UserEpisode.
func (ser *UserEpisodeService) Bucket() string {
	return UserEpisodeBucket
}

// Clean cleans the given UserEpisode for storage.
func (ser *UserEpisodeService) Clean(m Model) error {
	_, err := ser.AssertType(m)
	return err
}

// Validate returns an error if the UserEpisode is not valid for the database.
func (ser *UserEpisodeService) Validate(m Model) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	return ser.DB.View(func(tx *bolt.Tx) error {
		// Check if User with ID specified in UserEpisode exists
		// Get User bucket, exit if error
		ub, err := Bucket(UserBucket, tx)
		if err != nil {
			return fmt.Errorf("%s %q: %w", errmsgBucketOpen, UserBucket, err)
		}
		_, err = get(e.UserID, ub)
		if err != nil {
			return fmt.Errorf("failed to get User with ID %d: %w", e.UserID, err)
		}

		// Check if Episode with ID specified in UserEpisode exists
		// Get Episode bucket, exit if error
		epb, err := Bucket(EpisodeBucket, tx)
		if err != nil {
			return fmt.Errorf("%s %q: %w", errmsgBucketOpen, EpisodeBucket, err)
		}
		_, err = get(e.EpisodeID, epb)
		if err != nil {
			return fmt.Errorf(
				"failed to get Episode with ID %d: %w", e.EpisodeID, err)
		}

		return nil
	})
}

// Initialize sets initial values for some properties.
func (ser *UserEpisodeService) Initialize(m Model) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing
// UserEpisode in updates.
func (ser *UserEpisodeService) PersistOldProperties(n Model, o Model) error {
	return nil
}

// Marshal transforms the given UserEpisode into JSON.
func (ser *UserEpisodeService) Marshal(m Model) ([]byte, error) {
	uep, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	v, err := json.Marshal(uep)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONMarshal, err)
	}

	return v, nil
}

// Unmarshal parses the given JSON into UserEpisode.
func (ser *UserEpisodeService) Unmarshal(buf []byte) (Model, error) {
	var uep UserEpisode
	err := json.Unmarshal(buf, &uep)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &uep, nil
}

// AssertType exposes the given Model as a UserEpisode.
func (ser *UserEpisodeService) AssertType(m Model) (*UserEpisode, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	uep, ok := m.(*UserEpisode)
	if !ok {
		return nil,
			fmt.Errorf("model: %w", errors.New("not of UserEpisode type"))
	}
	return uep, nil
}

// mapfromModel returns a list of UserEpisode type asserted from the given
// list of Model.
func (ser *UserEpisodeService) mapFromModel(vlist []Model) ([]*UserEpisode, error) {
	list := make([]*UserEpisode, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}
