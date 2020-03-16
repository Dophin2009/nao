package data

import (
	"errors"
	"fmt"

	"github.com/Dophin2009/nao/pkg/data/models"
	"github.com/Dophin2009/nao/pkg/db"
	json "github.com/json-iterator/go"
)

// UserEpisodeService performs operations on UserEpisode.
type UserEpisodeService struct {
	UserService    *UserService
	EpisodeService *EpisodeService
	Hooks          db.PersistHooks
}

// NewUserEpisodeService returns a UserEpisodeService.
func NewUserEpisodeService(hooks db.PersistHooks, userService *UserService,
	episodeService *EpisodeService) *UserEpisodeService {
	// Initiate UserEpisodeService
	userEpisodeService := &UserEpisodeService{
		UserService:    userService,
		EpisodeService: episodeService,
		Hooks:          hooks,
	}

	// Add hook to delete UserEpisode on User deletion
	deleteUserEpisodeOnDeleteUser := func(um db.Model, _ db.Service, tx db.Tx) error {
		uID := um.Metadata().ID
		err := userEpisodeService.DeleteByUser(uID, tx)
		if err != nil {
			return fmt.Errorf("failed to delete UserEpisode by User ID %d: %w",
				uID, err)
		}
		return nil
	}
	uSerHooks := userService.PersistHooks()
	uSerHooks.PreDeleteHooks =
		append(uSerHooks.PreDeleteHooks, deleteUserEpisodeOnDeleteUser)

	// Add hook to delete UserEpisode on Episode deletion
	deleteUserEpisodeOnDeleteEpisode := func(epm db.Model, _ db.Service, tx db.Tx) error {
		epID := epm.Metadata().ID
		err := userEpisodeService.DeleteByEpisode(epID, tx)
		if err != nil {
			return fmt.Errorf("failed to delete UserEpisode by Episode ID %d: %w",
				epID, err)
		}
		return nil
	}
	epSerHooks := episodeService.PersistHooks()
	epSerHooks.PreDeleteHooks =
		append(epSerHooks.PreDeleteHooks, deleteUserEpisodeOnDeleteEpisode)

	return userEpisodeService
}

// Create persists the given UserEpisode.
func (ser *UserEpisodeService) Create(uep *models.UserEpisode, tx db.Tx) (int, error) {
	return tx.Database().Create(uep, ser, tx)
}

// Update rueplaces the value of the UserEpisode with the given ID.
func (ser *UserEpisodeService) Update(uep *models.UserEpisode, tx db.Tx) error {
	return tx.Database().Update(uep, ser, tx)
}

// Delete deletes the UserEpisode with the given ID.
func (ser *UserEpisodeService) Delete(id int, tx db.Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// DeleteByUser deletes the UserEpisodes with the given User ID.
func (ser *UserEpisodeService) DeleteByUser(uID int, tx db.Tx) error {
	return tx.Database().DeleteFilter(ser, tx, func(m db.Model) bool {
		uep, err := ser.AssertType(m)
		if err != nil {
			return false
		}
		return uep.UserID == uID
	})
}

// DeleteByEpisode deletes the UserEpisodes with the given Episode ID.
func (ser *UserEpisodeService) DeleteByEpisode(epID int, tx db.Tx) error {
	return tx.Database().DeleteFilter(ser, tx, func(m db.Model) bool {
		uep, err := ser.AssertType(m)
		if err != nil {
			return false
		}
		return uep.EpisodeID == epID
	})
}

// GetAll retrieves all persisted values of UserEpisode.
func (ser *UserEpisodeService) GetAll(first *int, skip *int, tx db.Tx) ([]*models.UserEpisode, error) {
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
	first *int, skip *int, tx db.Tx, keep func(uep *models.UserEpisode) bool,
) ([]*models.UserEpisode, error) {
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
	ids []int, tx db.Tx, keep func(uep *models.UserEpisode) bool,
) ([]*models.UserEpisode, error) {
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
func (ser *UserEpisodeService) GetByID(id int, tx db.Tx) (*models.UserEpisode, error) {
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
) ([]*models.UserEpisode, error) {
	return ser.GetFilter(first, skip, tx, func(uep *models.UserEpisode) bool {
		return uep.UserID == uID
	})
}

// GetByEpisode retrieves the persisted UserEpisode with the given Episode ID.
func (ser *UserEpisodeService) GetByEpisode(
	epID int, first *int, skip *int, tx db.Tx,
) ([]*models.UserEpisode, error) {
	return ser.GetFilter(first, skip, tx, func(uep *models.UserEpisode) bool {
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
	var uep models.UserEpisode
	err := json.Unmarshal(buf, &uep)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &uep, nil
}

// AssertType exposes the given db.Model as a UserEpisode.
func (ser *UserEpisodeService) AssertType(m db.Model) (*models.UserEpisode, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	uep, ok := m.(*models.UserEpisode)
	if !ok {
		return nil,
			fmt.Errorf("model: %w", errors.New("not of UserEpisode type"))
	}
	return uep, nil
}

// mapfromModel returns a list of UserEpisode type asserted from the given
// list of db.Model.
func (ser *UserEpisodeService) mapFromModel(vlist []db.Model) ([]*models.UserEpisode, error) {
	list := make([]*models.UserEpisode, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}
