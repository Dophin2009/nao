package data

import (
	"errors"
	"fmt"

	"github.com/Dophin2009/nao/pkg/models"
	"github.com/Dophin2009/nao/pkg/db"
	json "github.com/json-iterator/go"
)

// UserMediaListService performs operations on UserMediaList.
type UserMediaListService struct {
	UserService      *UserService
	UserMediaService *UserMediaService
	Hooks            db.PersistHooks
}

// NewUserMediaListService returns a UserMediaListService.
func NewUserMediaListService(hooks db.PersistHooks, userService *UserService,
	userMediaService *UserMediaService) *UserMediaListService {
	userMediaListService := &UserMediaListService{
		UserService:      userService,
		UserMediaService: userMediaService,
		Hooks:            hooks,
	}

	// Add hook to delete UserMediaList on User deletion
	deleteUserMediaListOnDeleteUser := func(um db.Model, _ db.Service, tx db.Tx) error {
		uID := um.Metadata().ID
		err := userMediaListService.DeleteByUser(uID, tx)
		if err != nil {
			return fmt.Errorf("failed to delete UserMediaList by User ID %d: %w",
				uID, err)
		}
		return nil
	}
	uSerHooks := userService.PersistHooks()
	uSerHooks.PreDeleteHooks =
		append(uSerHooks.PreDeleteHooks, deleteUserMediaListOnDeleteUser)

	// Add hook to update UserMediaList on UserMedia deletion
	updateUserMediaListOnDeleteUserMedia := func(umm db.Model, _ db.Service, tx db.Tx) error {
		umID := umm.Metadata().ID
		err := tx.Database().DoEach(nil, nil, userMediaListService, tx,
			func(m db.Model, _ db.Service, tx db.Tx) (exit bool, err error) {
				uml, err := userMediaListService.AssertType(m)
				if err != nil {
					return true, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
				}

				// Find ID of UserMedia to be deleted in the list
				rmID := -1
				for _, id := range uml.UserMedia {
					if id == umID {
						rmID = id
						break
					}
				}
				// If UserMedia ID not found, move onto next UserMediaList
				if rmID < 0 {
					return false, nil
				}

				// Remove ID from UserMedia list
				uml.UserMedia = append(uml.UserMedia[:rmID], uml.UserMedia[rmID+1:]...)

				// Update persisted value
				err = tx.Database().Update(uml, userMediaListService, tx)
				if err != nil {
					return true, fmt.Errorf("failed to update UserMediaList: %w", err)
				}
				return false, nil
			}, nil)
		if err != nil {
			return err
		}
		return nil
	}
	umSerHooks := userMediaService.PersistHooks()
	umSerHooks.PreDeleteHooks =
		append(umSerHooks.PreDeleteHooks, updateUserMediaListOnDeleteUserMedia)

	return userMediaListService
}

// Create persists the given UserMediaList.
func (ser *UserMediaListService) Create(uml *models.UserMediaList, tx db.Tx) (int, error) {
	return tx.Database().Create(uml, ser, tx)
}

// Update rumllaces the value of the UserMediaList with the given ID.
func (ser *UserMediaListService) Update(uml *models.UserMediaList, tx db.Tx) error {
	return tx.Database().Update(uml, ser, tx)
}

// Delete deletes the UserMediaList with the given ID.
func (ser *UserMediaListService) Delete(id int, tx db.Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// DeleteByUser deletes the UserMediaLists by the given User ID.
func (ser *UserMediaListService) DeleteByUser(uID int, tx db.Tx) error {
	return tx.Database().DeleteFilter(ser, tx, func(m db.Model) bool {
		uml, err := ser.AssertType(m)
		if err != nil {
			return false
		}
		return uml.UserID == uID
	})
}

// GetAll retrieves all persisted values of UserMediaList.
func (ser *UserMediaListService) GetAll(first *int, skip *int, tx db.Tx) ([]*models.UserMediaList, error) {
	vlist, err := tx.Database().GetAll(first, skip, ser, tx)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to UserMediaLists: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of UserMediaList that pass the
// filter.
func (ser *UserMediaListService) GetFilter(
	first *int, skip *int, tx db.Tx, keep func(uml *models.UserMediaList) bool,
) ([]*models.UserMediaList, error) {
	vlist, err := tx.Database().GetFilter(first, skip, ser, tx,
		func(m db.Model) bool {
			uml, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(uml)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to UserMediaLists: %w", err)
	}
	return list, nil
}

// GetMultiple retrieves the persisted UserMediaList values specified by the
// given IDs that pass the filter.
func (ser *UserMediaListService) GetMultiple(
	ids []int, tx db.Tx, keep func(uml *models.UserMediaList) bool,
) ([]*models.UserMediaList, error) {
	vlist, err := tx.Database().GetMultiple(ids, ser, tx, func(m db.Model) bool {
		uml, err := ser.AssertType(m)
		if err != nil {
			return false
		}
		return keep(uml)
	})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to UserMediaLists: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted UserMediaList with the given ID.
func (ser *UserMediaListService) GetByID(id int, tx db.Tx) (*models.UserMediaList, error) {
	m, err := tx.Database().GetByID(id, ser, tx)
	if err != nil {
		return nil, err
	}

	uml, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return uml, nil
}

// Bucket returns the name of the bucket for UserMediaList.
func (ser *UserMediaListService) Bucket() string {
	return "UserMediaList"
}

// Clean cleans the given UserMediaList for storage.
func (ser *UserMediaListService) Clean(m db.Model, _ db.Tx) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Validate returns an error if the UserMediaList is not valid for the
// database.
func (ser *UserMediaListService) Validate(m db.Model, tx db.Tx) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	db := tx.Database()

	// Check if User with ID specified in UserMediaList exists
	_, err = db.GetRawByID(e.UserID, ser.UserService, tx)
	if err != nil {
		return fmt.Errorf("failed to get User with ID %d: %w", e.UserID, err)
	}

	// Check if UserMedia with IDs specified in UserMediaList exist
	for _, umID := range e.UserMedia {
		_, err = db.GetRawByID(umID, ser.UserMediaService, tx)
		if err != nil {
			return fmt.Errorf("failed to get UserMedia with ID %d: %w", umID, err)
		}
	}

	return nil
}

// Initialize sets initial values for some properties.
func (ser *UserMediaListService) Initialize(_ db.Model, _ db.Tx) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing
// UserMediaList in updates.
func (ser *UserMediaListService) PersistOldProperties(_ db.Model, _ db.Model, _ db.Tx) error {
	return nil
}

// PersistHooks returns the persistence hook functions.
func (ser *UserMediaListService) PersistHooks() *db.PersistHooks {
	return &ser.Hooks
}

// Marshal transforms the given UserMediaList into JSON.
func (ser *UserMediaListService) Marshal(m db.Model) ([]byte, error) {
	uml, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	v, err := json.Marshal(uml)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONMarshal, err)
	}

	return v, nil
}

// Unmarshal parses the given JSON into UserMediaList.
func (ser *UserMediaListService) Unmarshal(buf []byte) (db.Model, error) {
	var uml models.UserMediaList
	err := json.Unmarshal(buf, &uml)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &uml, nil
}

// AssertType exposes the given db.Model as a UserMediaList.
func (ser *UserMediaListService) AssertType(m db.Model) (*models.UserMediaList, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	uml, ok := m.(*models.UserMediaList)
	if !ok {
		return nil, fmt.Errorf("model: %w", errors.New("not of UserMediaList type"))
	}
	return uml, nil
}

// mapfromModel returns a list of UserMediaList type asserted from the given
// list of db.Model.
func (ser *UserMediaListService) mapFromModel(vlist []db.Model) ([]*models.UserMediaList, error) {
	list := make([]*models.UserMediaList, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}
