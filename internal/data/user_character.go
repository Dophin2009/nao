package data

import (
	"errors"
	"fmt"

	"github.com/Dophin2009/nao/pkg/models"
	"github.com/Dophin2009/nao/pkg/db"
	json "github.com/json-iterator/go"
)

// UserCharacterService performs operations on UserCharacter.
type UserCharacterService struct {
	UserService      *UserService
	CharacterService *CharacterService
	Hooks            db.PersistHooks
}

// NewUserCharacterService returns a UserCharacterService.
func NewUserCharacterService(hooks db.PersistHooks, userService *UserService,
	characterService *CharacterService) *UserCharacterService {
	userCharacterService := &UserCharacterService{
		UserService:      userService,
		CharacterService: characterService,
		Hooks:            hooks,
	}

	// Add hook to delete UserCharacter on User deletion
	deleteUserCharacterOnDeleteUser := func(um db.Model, _ db.Service, tx db.Tx) error {
		uID := um.Metadata().ID
		err := userCharacterService.DeleteByUser(uID, tx)
		if err != nil {
			return fmt.Errorf("failed to delete UserCharacter by User ID %d: %w",
				uID, err)
		}
		return nil
	}
	uSerHooks := userService.PersistHooks()
	uSerHooks.PreDeleteHooks =
		append(uSerHooks.PreDeleteHooks, deleteUserCharacterOnDeleteUser)

	// Add hook to delete UserCharacter on Character deletion
	deleteUserCharacterOnDeleteCharacter := func(cm db.Model, _ db.Service, tx db.Tx) error {
		cID := cm.Metadata().ID
		err := userCharacterService.DeleteByCharacter(cID, tx)
		if err != nil {
			return fmt.Errorf("failed to delete UserCharacter by Character ID %d: %w",
				cID, err)
		}
		return nil
	}
	cSerHooks := characterService.PersistHooks()
	cSerHooks.PreDeleteHooks =
		append(cSerHooks.PreDeleteHooks, deleteUserCharacterOnDeleteCharacter)

	return userCharacterService
}

// Create persists the given UserCharacter.
func (ser *UserCharacterService) Create(uc *models.UserCharacter, tx db.Tx) (int, error) {
	return tx.Database().Create(uc, ser, tx)
}

// Update ruclaces the value of the UserCharacter with the given ID.
func (ser *UserCharacterService) Update(uc *models.UserCharacter, tx db.Tx) error {
	return tx.Database().Update(uc, ser, tx)
}

// Delete deletes the UserCharacter with the given ID.
func (ser *UserCharacterService) Delete(id int, tx db.Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// DeleteByUser deletes the UserCharacters with the given User ID.
func (ser *UserCharacterService) DeleteByUser(uID int, tx db.Tx) error {
	return tx.Database().DeleteFilter(ser, tx, func(ucm db.Model) bool {
		uc, err := ser.AssertType(ucm)
		if err != nil {
			return false
		}
		return uc.UserID == uID
	})
}

// DeleteByCharacter deletes the UserCharacters with the given Character ID.
func (ser *UserCharacterService) DeleteByCharacter(cID int, tx db.Tx) error {
	return tx.Database().DeleteFilter(ser, tx, func(ucm db.Model) bool {
		uc, err := ser.AssertType(ucm)
		if err != nil {
			return false
		}
		return uc.CharacterID == cID
	})
}

// GetAll retrieves all persisted values of UserCharacter.
func (ser *UserCharacterService) GetAll(first *int, skip *int, tx db.Tx) ([]*models.UserCharacter, error) {
	vlist, err := tx.Database().GetAll(first, skip, ser, tx)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to UserCharacter: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of UserCharacter that pass the
// filter.
func (ser *UserCharacterService) GetFilter(
	first *int, skip *int, tx db.Tx, keep func(uc *models.UserCharacter) bool,
) ([]*models.UserCharacter, error) {
	vlist, err := tx.Database().GetFilter(first, skip, ser, tx,
		func(m db.Model) bool {
			uc, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(uc)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to UserCharacter: %w", err)
	}
	return list, nil
}

// GetMultiple retrieves the persisted UserCharacter values specified by the
// given IDs that pass the filter.
func (ser *UserCharacterService) GetMultiple(
	ids []int, tx db.Tx, keep func(uc *models.UserCharacter) bool,
) ([]*models.UserCharacter, error) {
	vlist, err := tx.Database().GetMultiple(ids, ser, tx,
		func(m db.Model) bool {
			uc, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(uc)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to UserCharacters: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted UserCharacter with the given ID.
func (ser *UserCharacterService) GetByID(id int, tx db.Tx) (*models.UserCharacter, error) {
	m, err := tx.Database().GetByID(id, ser, tx)
	if err != nil {
		return nil, err
	}

	uc, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return uc, nil
}

// GetByUser retrieves the persisted UserCharacter with the given User ID.
func (ser *UserCharacterService) GetByUser(
	uID int, first *int, skip *int, tx db.Tx,
) ([]*models.UserCharacter, error) {
	return ser.GetFilter(first, skip, tx, func(uc *models.UserCharacter) bool {
		return uc.UserID == uID
	})
}

// GetByCharacter retrieves the persisted UserCharacter with the given Character ID.
func (ser *UserCharacterService) GetByCharacter(
	cID int, first *int, skip *int, tx db.Tx,
) ([]*models.UserCharacter, error) {
	return ser.GetFilter(first, skip, tx, func(uc *models.UserCharacter) bool {
		return uc.CharacterID == cID
	})
}

// Bucket returns the name of the bucket for UserCharacter.
func (ser *UserCharacterService) Bucket() string {
	return "UserCharacter"
}

// Clean cleans the given UserCharacter for storage.
func (ser *UserCharacterService) Clean(m db.Model, _ db.Tx) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Validate returns an error if the UserCharacter is not valid for the database.
func (ser *UserCharacterService) Validate(m db.Model, tx db.Tx) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	db := tx.Database()

	// Check if User with ID specified in UserCharacter exists
	// Get User bucket, exit if error
	_, err = db.GetRawByID(e.UserID, ser.UserService, tx)
	if err != nil {
		return fmt.Errorf("failed to get User with ID %d: %w", e.UserID, err)
	}

	// Check if Character with ID specified in UserCharacter exists
	// Get Character bucket, exit if error
	_, err = db.GetRawByID(e.CharacterID, ser.CharacterService, tx)
	if err != nil {
		return fmt.Errorf(
			"failed to get Character with ID %d: %w", e.CharacterID, err)
	}

	return nil
}

// Initialize sets initial values for some properties.
func (ser *UserCharacterService) Initialize(_ db.Model, _ db.Tx) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing
// UserCharacter in updates.
func (ser *UserCharacterService) PersistOldProperties(_ db.Model, _ db.Model, _ db.Tx) error {
	return nil
}

// PersistHooks returns the persistence hook functions.
func (ser *UserCharacterService) PersistHooks() *db.PersistHooks {
	return &ser.Hooks
}

// Marshal transforms the given UserCharacter into JSON.
func (ser *UserCharacterService) Marshal(m db.Model) ([]byte, error) {
	uc, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	v, err := json.Marshal(uc)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONMarshal, err)
	}

	return v, nil
}

// Unmarshal parses the given JSON into UserCharacter.
func (ser *UserCharacterService) Unmarshal(buf []byte) (db.Model, error) {
	var uc models.UserCharacter
	err := json.Unmarshal(buf, &uc)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &uc, nil
}

// AssertType exposes the given db.Model as a UserCharacter.
func (ser *UserCharacterService) AssertType(m db.Model) (*models.UserCharacter, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	uc, ok := m.(*models.UserCharacter)
	if !ok {
		return nil,
			fmt.Errorf("model: %w", errors.New("not of UserCharacter type"))
	}
	return uc, nil
}

// mapfromModel returns a list of UserCharacter type asserted from the given
// list of db.Model.
func (ser *UserCharacterService) mapFromModel(vlist []db.Model) ([]*models.UserCharacter, error) {
	list := make([]*models.UserCharacter, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}
