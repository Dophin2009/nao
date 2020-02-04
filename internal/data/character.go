package data

import (
	"errors"
	"fmt"

	json "github.com/json-iterator/go"
)

// TODO: User ratings/favoriting/comments/etc. of Characters

// Character represents a single character.
type Character struct {
	Names       []Title
	Information []Title
	Meta        ModelMetadata
}

// Metadata returns Meta.
func (c *Character) Metadata() *ModelMetadata {
	return &c.Meta
}

// CharacterBucket is the name of the database bucket for Character.
const CharacterBucket = "Character"

// CharacterService performs operations on Characters.
type CharacterService struct{}

// Create persists the given Character.
func (ser *CharacterService) Create(c *Character, tx Tx) (int, error) {
	return tx.Database().Create(c, ser, tx)
}

// Update replaces the value of the Character with the given ID.
func (ser *CharacterService) Update(c *Character, tx Tx) error {
	return tx.Database().Update(c, ser, tx)
}

// Delete deletes the Character with the given ID.
func (ser *CharacterService) Delete(id int, tx Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// GetAll retrieves all persisted values of Character.
func (ser *CharacterService) GetAll(first *int, skip *int, tx Tx) ([]*Character, error) {
	vlist, err := tx.Database().GetAll(first, skip, ser, tx)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to Characters: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of Character that pass the filter.
func (ser *CharacterService) GetFilter(
	first *int, skip *int, tx Tx,
	keep func(c *Character) bool,
) ([]*Character, error) {
	vlist, err := tx.Database().GetFilter(first, skip, ser, tx,
		func(m Model) bool {
			c, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(c)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to Characters: %w", err)
	}
	return list, nil
}

// GetMultiple retrieves the persisted Character values specified by the given
// IDs that pass the filter.
func (ser *CharacterService) GetMultiple(
	ids []int, first *int, skip *int, tx Tx, keep func(c *Character) bool,
) ([]*Character, error) {
	vlist, err := tx.Database().GetMultiple(ids, first, skip, ser, tx,
		func(m Model) bool {
			c, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(c)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to Characters: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted Character with the given ID.
func (ser *CharacterService) GetByID(id int, tx Tx) (*Character, error) {
	m, err := tx.Database().GetByID(id, ser, tx)
	if err != nil {
		return nil, err
	}

	c, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return c, nil
}

// Bucket returns the name of the bucket for Media.
func (ser *CharacterService) Bucket() string {
	return CharacterBucket
}

// Clean cleans the given Character for storage
func (ser *CharacterService) Clean(m Model, _ Tx) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Validate returns an error if the Character is not valid for the database.
func (ser *CharacterService) Validate(m Model, _ Tx) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Initialize sets initial values for some properties.
func (ser *CharacterService) Initialize(_ Model, _ Tx) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing Character
// in updates.
func (ser *CharacterService) PersistOldProperties(_ Model, _ Model, _ Tx) error {
	return nil
}

// Marshal transforms the given Character into JSON.
func (ser *CharacterService) Marshal(m Model) ([]byte, error) {
	c, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	v, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONMarshal, err)
	}

	return v, nil
}

// Unmarshal parses the given JSON into Character.
func (ser *CharacterService) Unmarshal(buf []byte) (Model, error) {
	var c Character
	err := json.Unmarshal(buf, &c)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &c, nil
}

// AssertType exposes the given Model as a Character.
func (ser *CharacterService) AssertType(m Model) (*Character, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	c, ok := m.(*Character)
	if !ok {
		return nil, fmt.Errorf("model: %w", errors.New("not of Character type"))
	}
	return c, nil
}

// mapFromModel returns a list of Character type asserted from the given list
// of Model.
func (ser *CharacterService) mapFromModel(vlist []Model) ([]*Character, error) {
	list := make([]*Character, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}
