package data

import (
	"errors"
	"fmt"
	"strings"

	json "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
)

// MediaCharacter represents a relationship between single instances of Media
// and Character.
type MediaCharacter struct {
	MediaID       int
	CharacterID   *int
	CharacterRole *string
	PersonID      *int
	PersonRole    *string
	Meta          ModelMetadata
}

// Metadata returns Meta.
func (mc *MediaCharacter) Metadata() *ModelMetadata {
	return &mc.Meta
}

// MediaCharacterBucket is the name of the database bucket for MediaCharacter.
const MediaCharacterBucket = "MediaCharacter"

// MediaCharacterService performs operations on MediaCharacter.
type MediaCharacterService struct {
	DB *bolt.DB
	Service
}

// Create persists the given MediaCharacter.
func (ser *MediaCharacterService) Create(mc *MediaCharacter) error {
	return Create(mc, ser)
}

// Update rmclaces the value of the MediaCharacter with the given ID.
func (ser *MediaCharacterService) Update(mc *MediaCharacter) error {
	return Update(mc, ser)
}

// Delete deletes the MediaCharacter with the given ID.
func (ser *MediaCharacterService) Delete(id int) error {
	return Delete(id, ser)
}

// GetAll retrieves all persisted values of MediaCharacter.
func (ser *MediaCharacterService) GetAll(first *int, skip *int) ([]*MediaCharacter, error) {
	vlist, err := GetAll(ser, first, skip)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to MediaCharacters: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of MediaCharacter that pass the
// filter.
func (ser *MediaCharacterService) GetFilter(
	first *int, skip *int, keep func(mc *MediaCharacter) bool,
) ([]*MediaCharacter, error) {
	vlist, err := GetFilter(ser, first, skip, func(m Model) bool {
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
		return nil, fmt.Errorf("failed to map Models to MediaCharacters: %w", err)
	}
	return list, nil
}

// GetMultiple retrieves the persisted MediaCharacter values specified by the
// given IDs that pass the filter.
func (ser *MediaCharacterService) GetMultiple(
	ids []int, first *int, skip *int, keep func(mc *MediaCharacter) bool,
) ([]*MediaCharacter, error) {
	vlist, err := GetMultiple(ser, ids, first, skip, func(m Model) bool {
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
		return nil, fmt.Errorf("failed to map Models to MediaCharacters: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted MediaCharacter with the given ID.
func (ser *MediaCharacterService) GetByID(id int) (*MediaCharacter, error) {
	m, err := GetByID(id, ser)
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
	mID int, first *int, skip *int,
) ([]*MediaCharacter, error) {
	return ser.GetFilter(first, skip, func(mc *MediaCharacter) bool {
		return mc.MediaID == mID
	})
}

// GetByCharacter retrieves a list of instances of MediaCharacter with the
// given Character ID.
func (ser *MediaCharacterService) GetByCharacter(
	cID int, first *int, skip *int,
) ([]*MediaCharacter, error) {
	return ser.GetFilter(first, skip, func(mc *MediaCharacter) bool {
		return *mc.CharacterID == cID
	})
}

// GetByPerson retrieves a list of instances of MediaCharacter with the given
// Person ID.
func (ser *MediaCharacterService) GetByPerson(pID int, first *int, skip *int) ([]*MediaCharacter, error) {
	return ser.GetFilter(first, skip, func(mc *MediaCharacter) bool {
		return *mc.CharacterID == pID
	})
}

// Database returns the database reference.
func (ser *MediaCharacterService) Database() *bolt.DB {
	return ser.DB
}

// Bucket returns the name of the bucket for MediaCharacter.
func (ser *MediaCharacterService) Bucket() string {
	return MediaCharacterBucket
}

// Clean cleans the given MediaCharacter for storage.
func (ser *MediaCharacterService) Clean(m Model) error {
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
func (ser *MediaCharacterService) Validate(m Model) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	return ser.DB.View(func(tx *bolt.Tx) error {
		// Check if Media with ID specified in new MediaCharacter exists
		// Get Media bucket, exit if error
		mb, err := Bucket(MediaBucket, tx)
		if err != nil {
			return fmt.Errorf("%s %q: %w", errmsgBucketOpen, MediaBucket, err)
		}
		_, err = get(e.MediaID, mb)
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

			// Get Character bucket, exit if error
			cID := *e.CharacterID
			cb, err := Bucket(CharacterBucket, tx)
			if err != nil {
				return fmt.Errorf("%s %q: %w", errmsgBucketOpen, CharacterBucket, err)
			}
			_, err = get(cID, cb)
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
			// Get Person bucket, exit if error
			pb, err := Bucket(PersonBucket, tx)
			if err != nil {
				return err
			}
			_, err = get(*e.PersonID, pb)
			if err != nil {
				return err
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
	})
}

// Initialize sets initial values for some properties.
func (ser *MediaCharacterService) Initialize(m Model) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing
// MediaCharacter in updates.
func (ser *MediaCharacterService) PersistOldProperties(n Model, o Model) error {
	return nil
}

// Marshal transforms the given MediaCharacter into JSON.
func (ser *MediaCharacterService) Marshal(m Model) ([]byte, error) {
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
func (ser *MediaCharacterService) Unmarshal(buf []byte) (Model, error) {
	var mc MediaCharacter
	err := json.Unmarshal(buf, &mc)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &mc, nil
}

// AssertType exposes the given Model as a MediaCharacter.
func (ser *MediaCharacterService) AssertType(m Model) (*MediaCharacter, error) {
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
// list of Model.
func (ser *MediaCharacterService) mapFromModel(vlist []Model) ([]*MediaCharacter, error) {
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
