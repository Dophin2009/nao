package data

import (
	"errors"
	"fmt"

	"github.com/Dophin2009/nao/pkg/data/models"
	"github.com/Dophin2009/nao/pkg/db"
	json "github.com/json-iterator/go"
)

// CharacterService performs operations on Characters.
type CharacterService struct {
	Hooks db.PersistHooks
}

// NewCharacterService returns a CharacterService.
func NewCharacterService(hooks db.PersistHooks) *CharacterService {
	return &CharacterService{
		Hooks: hooks,
	}
}

// Create persists the given Character.
func (ser *CharacterService) Create(c *models.Character, tx db.Tx) (int, error) {
	return tx.Database().Create(c, ser, tx)
}

// Update replaces the value of the Character with the given ID.
func (ser *CharacterService) Update(c *models.Character, tx db.Tx) error {
	return tx.Database().Update(c, ser, tx)
}

// Delete deletes the Character with the given ID.
func (ser *CharacterService) Delete(id int, tx db.Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// GetAll retrieves all persisted values of Character.
func (ser *CharacterService) GetAll(first *int, skip *int, tx db.Tx) ([]*models.Character, error) {
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
	first *int, skip *int, tx db.Tx,
	keep func(c *models.Character) bool,
) ([]*models.Character, error) {
	vlist, err := tx.Database().GetFilter(first, skip, ser, tx,
		func(m db.Model) bool {
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
	ids []int, tx db.Tx, keep func(c *models.Character) bool,
) ([]*models.Character, error) {
	vlist, err := tx.Database().GetMultiple(ids, ser, tx,
		func(m db.Model) bool {
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
func (ser *CharacterService) GetByID(id int, tx db.Tx) (*models.Character, error) {
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
	return "Character"
}

// Clean cleans the given Character for storage
func (ser *CharacterService) Clean(m db.Model, _ db.Tx) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Validate returns an error if the Character is not valid for the database.
func (ser *CharacterService) Validate(m db.Model, _ db.Tx) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Initialize sets initial values for some properties.
func (ser *CharacterService) Initialize(_ db.Model, _ db.Tx) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing Character
// in updates.
func (ser *CharacterService) PersistOldProperties(_ db.Model, _ db.Model, _ db.Tx) error {
	return nil
}

// Marshal transforms the given Character into JSON.
func (ser *CharacterService) Marshal(m db.Model) ([]byte, error) {
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

// PersistHooks returns the persistence hook functions.
func (ser *CharacterService) PersistHooks() *db.PersistHooks {
	return &ser.Hooks
}

// Unmarshal parses the given JSON into Character.
func (ser *CharacterService) Unmarshal(buf []byte) (db.Model, error) {
	var c models.Character
	err := json.Unmarshal(buf, &c)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &c, nil
}

// AssertType exposes the given Model as a Character.
func (ser *CharacterService) AssertType(m db.Model) (*models.Character, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	c, ok := m.(*models.Character)
	if !ok {
		return nil, fmt.Errorf("model: %w", errors.New("not of Character type"))
	}
	return c, nil
}

// mapFromModel returns a list of Character type asserted from the given list
// of Model.
func (ser *CharacterService) mapFromModel(vlist []db.Model) ([]*models.Character, error) {
	list := make([]*models.Character, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}
