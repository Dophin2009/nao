package data

import (
	"errors"
	"fmt"

	json "github.com/json-iterator/go"
)

// UserCharacter represents a relationship between a User and a Character,
// containing information about the User's opinion on the Character.
type UserCharacter struct {
	UserID      int
	CharacterID int
	Score       *int
	Comments    []Title
	Meta        ModelMetadata
}

// Metadata returns Meta.
func (uc *UserCharacter) Metadata() *ModelMetadata {
	return &uc.Meta
}

// UserCharacterBucket is the name of the database bucket for UserCharacter.
const UserCharacterBucket = "UserCharacter"

// UserCharacterService performs operations on UserCharacter.
type UserCharacterService struct {
	UserService      *UserService
	CharacterService *CharacterService
}

// Create persists the given UserCharacter.
func (ser *UserCharacterService) Create(uc *UserCharacter, tx Tx) (int, error) {
	return tx.Database().Create(uc, ser, tx)
}

// Update ruclaces the value of the UserCharacter with the given ID.
func (ser *UserCharacterService) Update(uc *UserCharacter, tx Tx) error {
	return tx.Database().Update(uc, ser, tx)
}

// Delete deletes the UserCharacter with the given ID.
func (ser *UserCharacterService) Delete(id int, tx Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// GetAll retrieves all persisted values of UserCharacter.
func (ser *UserCharacterService) GetAll(first *int, skip *int, tx Tx) ([]*UserCharacter, error) {
	vlist, err := tx.Database().GetAll(first, skip, ser, tx)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to UserCharacter: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of UserCharacter that pass the
// filter.
func (ser *UserCharacterService) GetFilter(
	first *int, skip *int, tx Tx, keep func(uc *UserCharacter) bool,
) ([]*UserCharacter, error) {
	vlist, err := tx.Database().GetFilter(first, skip, ser, tx,
		func(m Model) bool {
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
		return nil, fmt.Errorf("failed to map Models to UserCharacter: %w", err)
	}
	return list, nil
}

// GetMultiple retrieves the persisted UserCharacter values specified by the
// given IDs that pass the filter.
func (ser *UserCharacterService) GetMultiple(
	ids []int, first *int, skip *int, tx Tx, keep func(uc *UserCharacter) bool,
) ([]*UserCharacter, error) {
	vlist, err := tx.Database().GetMultiple(ids, first, skip, ser, tx,
		func(m Model) bool {
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
		return nil, fmt.Errorf("failed to map Models to UserCharacters: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted UserCharacter with the given ID.
func (ser *UserCharacterService) GetByID(id int, tx Tx) (*UserCharacter, error) {
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
	uID int, first *int, skip *int, tx Tx,
) ([]*UserCharacter, error) {
	return ser.GetFilter(first, skip, tx, func(uc *UserCharacter) bool {
		return uc.UserID == uID
	})
}

// GetByCharacter retrieves the persisted UserCharacter with the given Character ID.
func (ser *UserCharacterService) GetByCharacter(
	cID int, first *int, skip *int, tx Tx,
) ([]*UserCharacter, error) {
	return ser.GetFilter(first, skip, tx, func(uc *UserCharacter) bool {
		return uc.CharacterID == cID
	})
}

// Bucket returns the name of the bucket for UserCharacter.
func (ser *UserCharacterService) Bucket() string {
	return UserCharacterBucket
}

// Clean cleans the given UserCharacter for storage.
func (ser *UserCharacterService) Clean(m Model, _ Tx) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelCleaning, err)
	}
	return nil
}

// Validate returns an error if the UserCharacter is not valid for the database.
func (ser *UserCharacterService) Validate(m Model, tx Tx) error {
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
func (ser *UserCharacterService) Initialize(_ Model, _ Tx) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing
// UserCharacter in updates.
func (ser *UserCharacterService) PersistOldProperties(_ Model, _ Model, _ Tx) error {
	return nil
}

// Marshal transforms the given UserCharacter into JSON.
func (ser *UserCharacterService) Marshal(m Model) ([]byte, error) {
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
func (ser *UserCharacterService) Unmarshal(buf []byte) (Model, error) {
	var uc UserCharacter
	err := json.Unmarshal(buf, &uc)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &uc, nil
}

// AssertType exposes the given Model as a UserCharacter.
func (ser *UserCharacterService) AssertType(m Model) (*UserCharacter, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	uc, ok := m.(*UserCharacter)
	if !ok {
		return nil,
			fmt.Errorf("model: %w", errors.New("not of UserCharacter type"))
	}
	return uc, nil
}

// mapfromModel returns a list of UserCharacter type asserted from the given
// list of Model.
func (ser *UserCharacterService) mapFromModel(vlist []Model) ([]*UserCharacter, error) {
	list := make([]*UserCharacter, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}
