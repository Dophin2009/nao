package data

import (
	"errors"
	"fmt"
	"strings"

	json "github.com/json-iterator/go"
	"gitlab.com/Dophin2009/nao/pkg/db"
)

// MediaCharacter represents a relationship between single instances of Media
// and Character.
type MediaCharacter struct {
	MediaID       int
	CharacterID   *int
	CharacterRole *string
	PersonID      *int
	PersonRole    *string
	Meta          db.ModelMetadata
}

// Metadata returns Meta.
func (mc *MediaCharacter) Metadata() *db.ModelMetadata {
	return &mc.Meta
}

// MediaCharacterService performs operations on MediaCharacter.
type MediaCharacterService struct {
	MediaService     *MediaService
	CharacterService *CharacterService
	PersonService    *PersonService
}

// Create persists the given MediaCharacter.
func (ser *MediaCharacterService) Create(mc *MediaCharacter, tx db.Tx) (int, error) {
	return tx.Database().Create(mc, ser, tx)
}

// Update rmclaces the value of the MediaCharacter with the given ID.
func (ser *MediaCharacterService) Update(mc *MediaCharacter, tx db.Tx) error {
	return tx.Database().Update(mc, ser, tx)
}

// Delete deletes the MediaCharacter with the given ID.
func (ser *MediaCharacterService) Delete(id int, tx db.Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// GetAll retrieves all persisted values of MediaCharacter.
func (ser *MediaCharacterService) GetAll(first *int, skip *int, tx db.Tx) ([]*MediaCharacter, error) {
	vlist, err := tx.Database().GetAll(first, skip, ser, tx)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to MediaCharacters: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of MediaCharacter that pass the
// filter.
func (ser *MediaCharacterService) GetFilter(
	first *int, skip *int, tx db.Tx, keep func(mc *MediaCharacter) bool,
) ([]*MediaCharacter, error) {
	vlist, err := tx.Database().GetFilter(first, skip, ser, tx,
		func(m db.Model) bool {
			mc, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(mc)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to MediaCharacters: %w", err)
	}
	return list, nil
}

// GetMultiple retrieves the persisted MediaCharacter values specified by the
// given IDs that pass the filter.
func (ser *MediaCharacterService) GetMultiple(
	ids []int, first *int, skip *int, tx db.Tx, keep func(mc *MediaCharacter) bool,
) ([]*MediaCharacter, error) {
	vlist, err := tx.Database().GetMultiple(ids, first, skip, ser, tx,
		func(m db.Model) bool {
			mc, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(mc)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to MediaCharacters: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted MediaCharacter with the given ID.
func (ser *MediaCharacterService) GetByID(id int, tx db.Tx) (*MediaCharacter, error) {
	m, err := tx.Database().GetByID(id, ser, tx)
	if err != nil {
		return nil, err
	}

	mc, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return mc, nil
}

// GetByMedia retrieves a list of instances of MediaCharacter with the given
// Media ID.
func (ser *MediaCharacterService) GetByMedia(
	mID int, first *int, skip *int, tx db.Tx,
) ([]*MediaCharacter, error) {
	return ser.GetFilter(first, skip, tx, func(mc *MediaCharacter) bool {
		return mc.MediaID == mID
	})
}

// GetByCharacter retrieves a list of instances of MediaCharacter with the
// given Character ID.
func (ser *MediaCharacterService) GetByCharacter(
	cID int, first *int, skip *int, tx db.Tx,
) ([]*MediaCharacter, error) {
	return ser.GetFilter(first, skip, tx, func(mc *MediaCharacter) bool {
		return *mc.CharacterID == cID
	})
}

// GetByPerson retrieves a list of instances of MediaCharacter with the given
// Person ID.
func (ser *MediaCharacterService) GetByPerson(
	pID int, first *int, skip *int, tx db.Tx,
) ([]*MediaCharacter, error) {
	return ser.GetFilter(first, skip, tx, func(mc *MediaCharacter) bool {
		return *mc.CharacterID == pID
	})
}

// Bucket returns the name of the bucket for MediaCharacter.
func (ser *MediaCharacterService) Bucket() string {
	return "MediaCharacter"
}

// Clean cleans the given MediaCharacter for storage.
func (ser *MediaCharacterService) Clean(m db.Model, _ db.Tx) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	if e.CharacterID != nil {
		*e.CharacterRole = strings.Trim(*e.CharacterRole, " ")
	}
	if e.PersonRole != nil {
		*e.PersonRole = strings.Trim(*e.PersonRole, " ")
	}
	return nil
}

// Validate returns an error if the MediaCharacter is not valid for the
// database.
func (ser *MediaCharacterService) Validate(m db.Model, tx db.Tx) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	db := tx.Database()

	// Check if Media with ID specified in MediaCharacter exists
	_, err = db.GetRawByID(e.MediaID, ser.MediaService, tx)
	if err != nil {
		return fmt.Errorf("failed to get Media with ID %d: %w", e.MediaID, err)
	}

	// Invalid if both Character and Person are not specified
	if e.CharacterID == nil && e.PersonID == nil {
		nsterr := fmt.Errorf("character ID and person ID: %w", errNil)
		return fmt.Errorf(
			"either character ID or person ID must be specified: %w", nsterr)
	}

	// Check if Character with ID specified in new MediaCharacter exists
	// CharacterID might be not specified
	if e.CharacterID != nil {
		// CharacterRole must be present if CharacterID is specified
		if e.CharacterRole == nil {
			nsterr := fmt.Errorf("character role: %w", errNil)
			return fmt.Errorf(
				"character role must not be nil if character ID is specified: %w",
				nsterr,
			)
		}

		cID := *e.CharacterID
		_, err = db.GetRawByID(cID, ser.CharacterService, tx)
		if err != nil {
			return fmt.Errorf("failed to get Character with ID %d: %w", cID, err)
		}
	} else {
		// CharacterRole must not be specified if CharacterID is not
		if e.CharacterRole != nil {
			nsterr := fmt.Errorf("character ID: %w", errNil)
			return fmt.Errorf(
				"character role must be nil if character ID is not specified: %w",
				nsterr,
			)
		}
	}

	// Check if Person with ID specified in new MediaCharacter exists
	// PersonID may be not specified
	if e.PersonID != nil {
		// PersonRole must be present if PersonID is specified
		if e.PersonRole == nil {
			nsterr := fmt.Errorf("person role: %w", errNil)
			return fmt.Errorf(
				"person role must not be nil if person ID is specified: %w", nsterr)
		}

		pID := *e.PersonID
		_, err = db.GetRawByID(pID, ser.PersonService, tx)
		if err != nil {
			return fmt.Errorf("failed to get Person with ID %d: %w", pID, err)
		}
	} else {
		// PersonRole must not be specified if PersonID is not
		if e.PersonRole != nil {
			nsterr := fmt.Errorf("person ID: %w", errNil)
			return fmt.Errorf(
				"person role must be nil if person ID is not specified: %w", nsterr)
		}
	}

	return nil
}

// Initialize sets initial values for some properties.
func (ser *MediaCharacterService) Initialize(_ db.Model, _ db.Tx) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing
// MediaCharacter in updates.
func (ser *MediaCharacterService) PersistOldProperties(_ db.Model, _ db.Model, _ db.Tx) error {
	return nil
}

// Marshal transforms the given MediaCharacter into JSON.
func (ser *MediaCharacterService) Marshal(m db.Model) ([]byte, error) {
	mc, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	v, err := json.Marshal(mc)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONMarshal, err)
	}

	return v, nil
}

// Unmarshal parses the given JSON into MediaCharacter.
func (ser *MediaCharacterService) Unmarshal(buf []byte) (db.Model, error) {
	var mc MediaCharacter
	err := json.Unmarshal(buf, &mc)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &mc, nil
}

// AssertType exposes the given db.Model as a MediaCharacter.
func (ser *MediaCharacterService) AssertType(m db.Model) (*MediaCharacter, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	mc, ok := m.(*MediaCharacter)
	if !ok {
		return nil, fmt.Errorf("model: %w", errors.New("not of MediaCharacter type"))
	}
	return mc, nil
}

// mapfromModel returns a list of MediaCharacter type asserted from the given
// list of db.Model.
func (ser *MediaCharacterService) mapFromModel(vlist []db.Model) ([]*MediaCharacter, error) {
	list := make([]*MediaCharacter, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}
