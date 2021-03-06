package data

import (
	"errors"
	"fmt"

	"github.com/Dophin2009/nao/pkg/models"
	"github.com/Dophin2009/nao/pkg/db"
	json "github.com/json-iterator/go"
)

// UserPersonService performs operations on UserPerson.
type UserPersonService struct {
	UserService   *UserService
	PersonService *PersonService
	Hooks         db.PersistHooks
}

// NewUserPersonService returns a UserPersonService.
func NewUserPersonService(hooks db.PersistHooks, userService *UserService,
	personService *PersonService) *UserPersonService {
	userPersonService := &UserPersonService{
		UserService:   userService,
		PersonService: personService,
		Hooks:         hooks,
	}

	// Add hook to delete UserPerson on User deletion
	deleteUserPersonOnDeleteUser := func(um db.Model, _ db.Service, tx db.Tx) error {
		uID := um.Metadata().ID
		err := userPersonService.DeleteByUser(uID, tx)
		if err != nil {
			return fmt.Errorf("failed to delete UserPerson by User ID %d: %w",
				uID, err)
		}
		return nil
	}
	uSerHooks := userService.PersistHooks()
	uSerHooks.PreDeleteHooks =
		append(uSerHooks.PreDeleteHooks, deleteUserPersonOnDeleteUser)

	// Add hook to delete UserPerson on Person deletion
	deleteUserPersonOnDeletePerson := func(pm db.Model, _ db.Service, tx db.Tx) error {
		pID := pm.Metadata().ID
		err := userPersonService.DeleteByPerson(pID, tx)
		if err != nil {
			return fmt.Errorf("failed to delete UserPerson by Person ID %d: %w",
				pID, err)
		}
		return nil
	}
	pSerHooks := personService.PersistHooks()
	pSerHooks.PreDeleteHooks =
		append(pSerHooks.PreDeleteHooks, deleteUserPersonOnDeletePerson)

	return userPersonService
}

// Create persists the given UserPerson.
func (ser *UserPersonService) Create(up *models.UserPerson, tx db.Tx) (int, error) {
	return tx.Database().Create(up, ser, tx)
}

// Update ruplaces the value of the UserPerson with the given ID.
func (ser *UserPersonService) Update(up *models.UserPerson, tx db.Tx) error {
	return tx.Database().Update(up, ser, tx)
}

// Delete deletes the UserPerson with the given ID.
func (ser *UserPersonService) Delete(id int, tx db.Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// DeleteByUser deletes the UserPersons with the given User ID.
func (ser *UserPersonService) DeleteByUser(uID int, tx db.Tx) error {
	return tx.Database().DeleteFilter(ser, tx, func(m db.Model) bool {
		up, err := ser.AssertType(m)
		if err != nil {
			return false
		}
		return up.UserID == uID
	})
}

// DeleteByPerson deletes the UserPersons with the given Person ID.
func (ser *UserPersonService) DeleteByPerson(pID int, tx db.Tx) error {
	return tx.Database().DeleteFilter(ser, tx, func(m db.Model) bool {
		up, err := ser.AssertType(m)
		if err != nil {
			return false
		}
		return up.PersonID == pID
	})
}

// GetAll retrieves all persisted values of UserPerson.
func (ser *UserPersonService) GetAll(first *int, skip *int, tx db.Tx) ([]*models.UserPerson, error) {
	vlist, err := tx.Database().GetAll(first, skip, ser, tx)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to UserPerson: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of UserPerson that pass the
// filter.
func (ser *UserPersonService) GetFilter(
	first *int, skip *int, tx db.Tx, keep func(up *models.UserPerson) bool,
) ([]*models.UserPerson, error) {
	vlist, err := tx.Database().GetFilter(first, skip, ser, tx,
		func(m db.Model) bool {
			up, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(up)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to UserPerson: %w", err)
	}
	return list, nil
}

// GetMultiple retrieves the persisted UserPerson values specified by the
// given IDs that pass the filter.
func (ser *UserPersonService) GetMultiple(
	ids []int, tx db.Tx, keep func(up *models.UserPerson) bool,
) ([]*models.UserPerson, error) {
	vlist, err := tx.Database().GetMultiple(ids, ser, tx,
		func(m db.Model) bool {
			up, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(up)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to UserPersons: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted UserPerson with the given ID.
func (ser *UserPersonService) GetByID(id int, tx db.Tx) (*models.UserPerson, error) {
	m, err := tx.Database().GetByID(id, ser, tx)
	if err != nil {
		return nil, err
	}

	up, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return up, nil
}

// GetByUser retrieves the persisted UserPerson with the given User ID.
func (ser *UserPersonService) GetByUser(
	uID int, first *int, skip *int, tx db.Tx,
) ([]*models.UserPerson, error) {
	return ser.GetFilter(first, skip, tx, func(up *models.UserPerson) bool {
		return up.UserID == uID
	})
}

// GetByPerson retrieves the persisted UserPerson with the given Person ID.
func (ser *UserPersonService) GetByPerson(
	cID int, first *int, skip *int, tx db.Tx,
) ([]*models.UserPerson, error) {
	return ser.GetFilter(first, skip, tx, func(up *models.UserPerson) bool {
		return up.PersonID == cID
	})
}

// Bucket returns the name of the bucket for UserPerson.
func (ser *UserPersonService) Bucket() string {
	return "UserPerson"
}

// Clean cleans the given UserPerson for storage.
func (ser *UserPersonService) Clean(m db.Model, _ db.Tx) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Validate returns an error if the UserPerson is not valid for the database.
func (ser *UserPersonService) Validate(m db.Model, tx db.Tx) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	db := tx.Database()

	// Check if User with ID specified in UserPerson exists
	_, err = db.GetRawByID(e.UserID, ser.UserService, tx)
	if err != nil {
		return fmt.Errorf("failed to get User with ID %d: %w", e.UserID, err)
	}

	// Check if Person with ID specified in UserPerson exists
	_, err = db.GetRawByID(e.PersonID, ser.PersonService, tx)
	if err != nil {
		return fmt.Errorf(
			"failed to get Person with ID %d: %w", e.PersonID, err)
	}

	return nil
}

// Initialize sets initial values for some properties.
func (ser *UserPersonService) Initialize(_ db.Model, _ db.Tx) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing
// UserPerson in updates.
func (ser *UserPersonService) PersistOldProperties(_ db.Model, _ db.Model, _ db.Tx) error {
	return nil
}

// PersistHooks returns the persistence hook functions.
func (ser *UserPersonService) PersistHooks() *db.PersistHooks {
	return &ser.Hooks
}

// Marshal transforms the given UserPerson into JSON.
func (ser *UserPersonService) Marshal(m db.Model) ([]byte, error) {
	up, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	v, err := json.Marshal(up)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONMarshal, err)
	}

	return v, nil
}

// Unmarshal parses the given JSON into UserPerson.
func (ser *UserPersonService) Unmarshal(buf []byte) (db.Model, error) {
	var up models.UserPerson
	err := json.Unmarshal(buf, &up)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &up, nil
}

// AssertType exposes the given db.Model as a UserPerson.
func (ser *UserPersonService) AssertType(m db.Model) (*models.UserPerson, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	up, ok := m.(*models.UserPerson)
	if !ok {
		return nil,
			fmt.Errorf("model: %w", errors.New("not of UserPerson type"))
	}
	return up, nil
}

// mapfromModel returns a list of UserPerson type asserted from the given
// list of db.Model.
func (ser *UserPersonService) mapFromModel(vlist []db.Model) ([]*models.UserPerson, error) {
	list := make([]*models.UserPerson, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}
