package data

import (
	"errors"
	"fmt"

	json "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
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
	DB *bolt.DB
	Service
}

// Create persists the given UserCharacter.
func (ser *UserCharacterService) Create(uc *UserCharacter) error {
	return Create(uc, ser)
}

// Update ruclaces the value of the UserCharacter with the given ID.
func (ser *UserCharacterService) Update(uc *UserCharacter) error {
	return Update(uc, ser)
}

// Delete deletes the UserCharacter with the given ID.
func (ser *UserCharacterService) Delete(id int) error {
	return Delete(id, ser)
}

// GetAll retrieves all persisted values of UserCharacter.
func (ser *UserCharacterService) GetAll(first *int, skip *int) ([]*UserCharacter, error) {
	vlist, err := GetAll(ser, first, skip)
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
	first *int, skip *int, keep func(uc *UserCharacter) bool,
) ([]*UserCharacter, error) {
	vlist, err := GetFilter(ser, first, skip, func(m Model) bool {
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
	ids []int, first *int, skip *int, keep func(uc *UserCharacter) bool,
) ([]*UserCharacter, error) {
	vlist, err := GetMultiple(ser, ids, first, skip, func(m Model) bool {
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
func (ser *UserCharacterService) GetByID(id int) (*UserCharacter, error) {
	m, err := GetByID(id, ser)
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
	uID int, first *int, skip *int,
) ([]*UserCharacter, error) {
	return ser.GetFilter(first, skip, func(uc *UserCharacter) bool {
		return uc.UserID == uID
	})
}

// GetByCharacter retrieves the persisted UserCharacter with the given Character ID.
func (ser *UserCharacterService) GetByCharacter(
	cID int, first *int, skip *int,
) ([]*UserCharacter, error) {
	return ser.GetFilter(first, skip, func(uc *UserCharacter) bool {
		return uc.CharacterID == cID
	})
}

// Database returns the database reference.
func (ser *UserCharacterService) Database() *bolt.DB {
	return ser.DB
}

// Bucket returns the name of the bucket for UserCharacter.
func (ser *UserCharacterService) Bucket() string {
	return UserCharacterBucket
}

// Clean cleans the given UserCharacter for storage.
func (ser *UserCharacterService) Clean(m Model) error {
	_, err := ser.AssertType(m)
	return err
}

// Validate returns an error if the UserCharacter is not valid for the database.
func (ser *UserCharacterService) Validate(m Model) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	return ser.DB.View(func(tx *bolt.Tx) error {
		// Check if User with ID specified in UserCharacter exists
		// Get User bucket, exit if error
		ub, err := Bucket(UserBucket, tx)
		if err != nil {
			return fmt.Errorf("%s %q: %w", errmsgBucketOpen, UserBucket, err)
		}
		_, err = get(e.UserID, ub)
		if err != nil {
			return fmt.Errorf("failed to get User with ID %d: %w", e.UserID, err)
		}

		// Check if Character with ID specified in UserCharacter exists
		// Get Character bucket, exit if error
		cb, err := Bucket(CharacterBucket, tx)
		if err != nil {
			return fmt.Errorf("%s %q: %w", errmsgBucketOpen, CharacterBucket, err)
		}
		_, err = get(e.CharacterID, cb)
		if err != nil {
			return fmt.Errorf(
				"failed to get Character with ID %d: %w", e.CharacterID, err)
		}

		return nil
	})
}

// Initialize sets initial values for some properties.
func (ser *UserCharacterService) Initialize(m Model) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing
// UserCharacter in updates.
func (ser *UserCharacterService) PersistOldProperties(n Model, o Model) error {
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
