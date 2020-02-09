package data

import (
	"errors"
	"fmt"

	json "github.com/json-iterator/go"
	"gitlab.com/Dophin2009/nao/pkg/db"
)

// UserEpisode represents a relationship between a User and an Episode,
// containing information about the User's opinion on the Episode.
type UserEpisode struct {
	UserID    int
	EpisodeID int
	Score     *int
	Comments  []Title
	Meta      db.ModelMetadata
}

// Metadata returns Meta.
func (uep *UserEpisode) Metadata() *db.ModelMetadata {
	return &uep.Meta
}

// UserEpisodeService performs operations on UserEpisode.
type UserEpisodeService struct {
	UserService    *UserService
	EpisodeService *EpisodeService
	Hooks          db.PersistHooks
}

// Create persists the given UserEpisode.
func (ser *UserEpisodeService) Create(uep *UserEpisode, tx db.Tx) (int, error) {
	return tx.Database().Create(uep, ser, tx)
}

// Update rueplaces the value of the UserEpisode with the given ID.
func (ser *UserEpisodeService) Update(uep *UserEpisode, tx db.Tx) error {
	return tx.Database().Update(uep, ser, tx)
}

// Delete deletes the UserEpisode with the given ID.
func (ser *UserEpisodeService) Delete(id int, tx db.Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// GetAll retrieves all persisted values of UserEpisode.
func (ser *UserEpisodeService) GetAll(first *int, skip *int, tx db.Tx) ([]*UserEpisode, error) {
	vlist, err := tx.Database().GetAll(first, skip, ser, tx)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to UserEpisode: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of UserEpisode that pass the
// filter.
func (ser *UserEpisodeService) GetFilter(
	first *int, skip *int, tx db.Tx, keep func(uep *UserEpisode) bool,
) ([]*UserEpisode, error) {
	vlist, err := tx.Database().GetFilter(first, skip, ser, tx,
		func(m db.Model) bool {
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
		return nil, fmt.Errorf("failed to map db.Models to UserEpisode: %w", err)
	}
	return list, nil
}

// GetMultiple retrieves the persisted UserEpisode values specified by the
// given IDs that pass the filter.
func (ser *UserEpisodeService) GetMultiple(
	ids []int, tx db.Tx, keep func(uep *UserEpisode) bool,
) ([]*UserEpisode, error) {
	vlist, err := tx.Database().GetMultiple(ids, ser, tx,
		func(m db.Model) bool {
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
		return nil, fmt.Errorf("failed to map db.Models to UserEpisodes: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted UserEpisode with the given ID.
func (ser *UserEpisodeService) GetByID(id int, tx db.Tx) (*UserEpisode, error) {
	m, err := tx.Database().GetByID(id, ser, tx)
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
	uID int, first *int, skip *int, tx db.Tx,
) ([]*UserEpisode, error) {
	return ser.GetFilter(first, skip, tx, func(uep *UserEpisode) bool {
		return uep.UserID == uID
	})
}

// GetByEpisode retrieves the persisted UserEpisode with the given Episode ID.
func (ser *UserEpisodeService) GetByEpisode(
	epID int, first *int, skip *int, tx db.Tx,
) ([]*UserEpisode, error) {
	return ser.GetFilter(first, skip, tx, func(uep *UserEpisode) bool {
		return uep.EpisodeID == epID
	})
}

// Bucket returns the name of the bucket for UserEpisode.
func (ser *UserEpisodeService) Bucket() string {
	return "UserEpisode"
}

// Clean cleans the given UserEpisode for storage.
func (ser *UserEpisodeService) Clean(m db.Model, _ db.Tx) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Validate returns an error if the UserEpisode is not valid for the database.
func (ser *UserEpisodeService) Validate(m db.Model, tx db.Tx) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	db := tx.Database()

	// Check if User with ID specified in UserEpisode exists
	_, err = db.GetRawByID(e.UserID, ser.UserService, tx)
	if err != nil {
		return fmt.Errorf("failed to get User with ID %d: %w", e.UserID, err)
	}

	// Check if Episode with ID specified in UserEpisode exists
	_, err = db.GetRawByID(e.EpisodeID, ser.EpisodeService, tx)
	if err != nil {
		return fmt.Errorf(
			"failed to get Episode with ID %d: %w", e.EpisodeID, err)
	}

	return nil
}

// Initialize sets initial values for some properties.
func (ser *UserEpisodeService) Initialize(_ db.Model, _ db.Tx) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing
// UserEpisode in updates.
func (ser *UserEpisodeService) PersistOldProperties(_ db.Model, _ db.Model, _ db.Tx) error {
	return nil
}

// PersistHooks returns the persistence hook functions.
func (ser *UserEpisodeService) PersistHooks() *db.PersistHooks {
	return &ser.Hooks
}

// Marshal transforms the given UserEpisode into JSON.
func (ser *UserEpisodeService) Marshal(m db.Model) ([]byte, error) {
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
func (ser *UserEpisodeService) Unmarshal(buf []byte) (db.Model, error) {
	var uep UserEpisode
	err := json.Unmarshal(buf, &uep)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &uep, nil
}

// AssertType exposes the given db.Model as a UserEpisode.
func (ser *UserEpisodeService) AssertType(m db.Model) (*UserEpisode, error) {
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
// list of db.Model.
func (ser *UserEpisodeService) mapFromModel(vlist []db.Model) ([]*UserEpisode, error) {
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
