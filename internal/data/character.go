package data

import (
	"errors"

	json "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
)

// Character represents a single character.
type Character struct {
	ID          int
	Names       map[string]string
	Information map[string]string
	Version     int
	Model
}

// Iden returns the ID.
func (c *Character) Iden() int {
	return c.ID
}

// CharacterBucket is the name of the database bucket
// for Character.
const CharacterBucket = "Character"

// CharacterService performs operations on Characters.
type CharacterService struct {
	DB *bolt.DB
	Service
}

// Create persists the given Character.
func (ser *CharacterService) Create(c *Character) error {
	return Create(c, ser)
}

// Update replaces the value of the Character with the
// given ID.
func (ser *CharacterService) Update(c *Character) error {
	return Update(c, ser)
}

// Delete deletes the Character with the given ID.
func (ser *CharacterService) Delete(id int) error {
	return Delete(id, ser)
}

// GetAll retrieves all persisted values of Character.
func (ser *CharacterService) GetAll(first int, prefixID *int) ([]*Character, error) {
	vlist, err := GetAll(ser, first, prefixID)
	if err != nil {
		return nil, err
	}

	return ser.mapFromModel(vlist)
}

// GetFilter retrieves all persisted values of Character that
// pass the filter.
func (ser *CharacterService) GetFilter(first int, prefixID *int, keep func(c *Character) bool) ([]*Character, error) {
	vlist, err := GetFilter(ser, first, prefixID, func(m Model) bool {
		c, err := ser.AssertType(m)
		if err != nil {
			return false
		}
		return keep(c)
	})
	if err != nil {
		return nil, err
	}

	return ser.mapFromModel(vlist)
}

// GetByID retrieves the persisted Character with the given ID.
func (ser *CharacterService) GetByID(id int) (*Character, error) {
	m, err := GetByID(id, ser)
	if err != nil {
		return nil, err
	}

	c, err := ser.AssertType(m)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Database returns the database reference.
func (ser *CharacterService) Database() *bolt.DB {
	return ser.DB
}

// Bucket returns the name of the bucket for Media.
func (ser *CharacterService) Bucket() string {
	return CharacterBucket
}

// Clean cleans the given Character for storage
func (ser *CharacterService) Clean(m Model) error {
	_, err := ser.AssertType(m)
	return err
}

// Validate returns an error if the Character is
// not valid for the database.
func (ser *CharacterService) Validate(m Model) error {
	_, err := ser.AssertType(m)
	return err
}

// Initialize sets initial values for some properties.
func (ser *CharacterService) Initialize(m Model, id int) error {
	c, err := ser.AssertType(m)
	if err != nil {
		return err
	}
	c.ID = id
	c.Version = 0
	return nil
}

// PersistOldProperties maintains certain properties
// of the existing Character in updates.
func (ser *CharacterService) PersistOldProperties(n Model, o Model) error {
	nc, err := ser.AssertType(n)
	if err != nil {
		return err
	}
	oc, err := ser.AssertType(o)
	if err != nil {
		return err
	}
	nc.Version = oc.Version + 1
	return nil
}

// Marshal transforms the given Character into JSON.
func (ser *CharacterService) Marshal(m Model) ([]byte, error) {
	c, err := ser.AssertType(m)
	if err != nil {
		return nil, err
	}

	v, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Unmarshal parses the given JSON into Character.
func (ser *CharacterService) Unmarshal(buf []byte) (Model, error) {
	var c Character
	err := json.Unmarshal(buf, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// AssertType exposes the given Model as a Character.
func (ser *CharacterService) AssertType(m Model) (*Character, error) {
	if m == nil {
		return nil, errors.New("model must not be nil")
	}

	c, ok := m.(*Character)
	if !ok {
		return nil, errors.New("model must be of Character type")
	}
	return c, nil
}

// mapFromModel returns a list of Character type asserted
// from the given list of Model.
func (ser *CharacterService) mapFromModel(vlist []Model) ([]*Character, error) {
	list := make([]*Character, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}
