package data

import (
	"errors"
	"fmt"
	"strings"

	json "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
)

// MediaCharacter represents a relationship between single
// instances of Media and Character.
type MediaCharacter struct {
	ID            int
	MediaID       int
	CharacterID   *int
	CharacterRole *string
	PersonID      *int
	PersonRole    *string
	Version       int
	Model
}

// Iden returns the ID.
func (mc *MediaCharacter) Iden() int {
	return mc.ID
}

// MediaCharacterBucket is the name of the database bucket for
// MediaCharacter.
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

// Update rmclaces the value of the MediaCharacter with the
// given ID.
func (ser *MediaCharacterService) Update(mc *MediaCharacter) error {
	return Update(mc, ser)
}

// Delete deletes the MediaCharacter with the given ID.
func (ser *MediaCharacterService) Delete(id int) error {
	return Delete(id, ser)
}

// GetAll retrieves all persisted values of MediaCharacter.
func (ser *MediaCharacterService) GetAll() ([]*MediaCharacter, error) {
	vlist, err := GetAll(ser)
	if err != nil {
		return nil, err
	}

	return ser.mapFromModel(vlist)
}

// GetFilter retrieves all persisted values of MediaCharacter that
// pass the filter.
func (ser *MediaCharacterService) GetFilter(keep func(mc *MediaCharacter) bool) ([]*MediaCharacter, error) {
	vlist, err := GetFilter(ser, func(m Model) bool {
		mc, err := ser.AssertType(m)
		if err != nil {
			return false
		}
		return keep(mc)
	})
	if err != nil {
		return nil, err
	}

	return ser.mapFromModel(vlist)
}

// GetByID retrieves the persisted MediaCharacter with the given ID.
func (ser *MediaCharacterService) GetByID(id int) (*MediaCharacter, error) {
	m, err := GetByID(id, ser)
	if err != nil {
		return nil, err
	}

	mc, err := ser.AssertType(m)
	if err != nil {
		return nil, err
	}
	return mc, nil
}

// GetByMedia retrieves a list of instances of
// MediaCharacter with the given Media ID.
func (ser *MediaCharacterService) GetByMedia(mID int) ([]*MediaCharacter, error) {
	return ser.GetFilter(func(mc *MediaCharacter) bool {
		return mc.MediaID == mID
	})
}

// GetByCharacter retrieves a list of instances of
// MediaCharacter with the given Character ID.
func (ser *MediaCharacterService) GetByCharacter(cID int) ([]*MediaCharacter, error) {
	return ser.GetFilter(func(mc *MediaCharacter) bool {
		return *mc.CharacterID == cID
	})
}

// GetByPerson retrieves a list of instances of
// MediaCharacter with the given Person ID.
func (ser *MediaCharacterService) GetByPerson(pID int) ([]*MediaCharacter, error) {
	return ser.GetFilter(func(mc *MediaCharacter) bool {
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
		return err
	}

	if e.CharacterID != nil {
		*e.CharacterRole = strings.Trim(*e.CharacterRole, " ")
	}
	if e.PersonRole != nil {
		*e.PersonRole = strings.Trim(*e.PersonRole, " ")
	}
	return nil
}

// Validate returns an error if the MediaCharacter is
// not valid for the database.
func (ser *MediaCharacterService) Validate(m Model) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return err
	}

	return ser.DB.View(func(tx *bolt.Tx) error {
		// Check if Media with ID specified in new MediaCharacter exists
		// Get Media bucket, exit if error
		mb, err := Bucket(MediaBucket, tx)
		if err != nil {
			return err
		}
		_, err = get(e.MediaID, mb)
		if err != nil {
			return err
		}

		if e.CharacterID == nil && e.PersonID == nil {
			return fmt.Errorf("either character id or person id must be specified")
		}

		// Check if Character with ID specified in new MediaCharacter exists
		// CharacterID may be not specified
		if e.CharacterID != nil {
			if e.CharacterRole == nil {
				return errors.New("character role must not be nil if character id is specified")
			}
			// Get Character bucket, exit if error
			cb, err := Bucket(CharacterBucket, tx)
			if err != nil {
				return err
			}
			_, err = get(*e.CharacterID, cb)
			if err != nil {
				return err
			}
		} else {
			if e.CharacterRole != nil {
				return fmt.Errorf("character role must be nil if character id is not specified")
			}
		}

		// Check if Person with ID specified in new MediaCharacter exists
		// PersonID may be not specified
		if e.PersonID != nil {
			if e.PersonRole == nil {
				return errors.New("person role must not be nil if person id is specified")
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
			if e.PersonRole != nil {
				return fmt.Errorf("person role must be nil if person id is not specified")
			}
		}

		return nil
	})
}

// Initialize sets initial values for some properties.
func (ser *MediaCharacterService) Initialize(m Model, id int) error {
	mc, err := ser.AssertType(m)
	if err != nil {
		return err
	}
	mc.ID = id
	mc.Version = 0
	return nil
}

// PersistOldProperties maintains certain properties
// of the existing MediaCharacter in updates.
func (ser *MediaCharacterService) PersistOldProperties(n Model, o Model) error {
	nm, err := ser.AssertType(n)
	if err != nil {
		return err
	}
	om, err := ser.AssertType(o)
	if err != nil {
		return err
	}
	nm.Version = om.Version + 1
	return nil
}

// Marshal transforms the given MediaCharacter into JSON.
func (ser *MediaCharacterService) Marshal(m Model) ([]byte, error) {
	mc, err := ser.AssertType(m)
	if err != nil {
		return nil, err
	}

	v, err := json.Marshal(mc)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Unmarshal parses the given JSON into MediaCharacter.
func (ser *MediaCharacterService) Unmarshal(buf []byte) (Model, error) {
	var mc MediaCharacter
	err := json.Unmarshal(buf, &mc)
	if err != nil {
		return nil, err
	}
	return &mc, nil
}

// AssertType exposes the given Model as a MediaCharacter.
func (ser *MediaCharacterService) AssertType(m Model) (*MediaCharacter, error) {
	if m == nil {
		return nil, errors.New("model must not be nil")
	}

	mc, ok := m.(*MediaCharacter)
	if !ok {
		return nil, errors.New("model must be of MediaCharacter type")
	}
	return mc, nil
}

// mapfromModel returns a list of MediaCharacter type
// asserted from the given list of Model.
func (ser *MediaCharacterService) mapFromModel(vlist []Model) ([]*MediaCharacter, error) {
	list := make([]*MediaCharacter, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}
