package data

import (
	"errors"
	"fmt"

	json "github.com/json-iterator/go"
)

// TODO: User rating/favoriting/comments/etc. of Persons

// Person represents a single person
type Person struct {
	Names       []Title
	Information []Title
	Meta        ModelMetadata
}

// Metadata returns Meta.
func (p *Person) Metadata() *ModelMetadata {
	return &p.Meta
}

// PersonBucket is the name of the database bucket for Person.
const PersonBucket = "Person"

// PersonService performs operations on Persons.
type PersonService struct{}

// Create persists the given Person.
func (ser *PersonService) Create(p *Person, tx Tx) (int, error) {
	return tx.Database().Create(p, ser, tx)
}

// Update rplaces the value of the Person with the given ID.
func (ser *PersonService) Update(p *Person, tx Tx) error {
	return tx.Database().Update(p, ser, tx)
}

// Delete deletes the Person with the given ID.
func (ser *PersonService) Delete(id int, tx Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// GetAll retrieves all persisted values of Person.
func (ser *PersonService) GetAll(first *int, skip *int, tx Tx) ([]*Person, error) {
	vlist, err := tx.Database().GetAll(first, skip, ser, tx)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to Persons: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of Person that pass the filter.
func (ser *PersonService) GetFilter(
	first *int, skip *int, tx Tx, keep func(p *Person) bool,
) ([]*Person, error) {
	vlist, err := tx.Database().GetFilter(first, skip, ser, tx,
		func(m Model) bool {
			p, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(p)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to Persons: %w", err)
	}
	return list, nil
}

// GetMultiple retrieves the persisted Person values specified by the given IDs
// that pass the filter.
func (ser *PersonService) GetMultiple(
	ids []int, first *int, skip *int, tx Tx, keep func(p *Person) bool,
) ([]*Person, error) {
	vlist, err := tx.Database().GetMultiple(ids, first, skip, ser, tx,
		func(m Model) bool {
			p, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(p)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to Persons: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted Person with the given ID.
func (ser *PersonService) GetByID(id int, tx Tx) (*Person, error) {
	m, err := tx.Database().GetByID(id, ser, tx)
	if err != nil {
		return nil, err
	}

	p, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return p, nil
}

// Bucket returns the name of the bucket for Person.
func (ser *PersonService) Bucket() string {
	return PersonBucket
}

// Clean cleans the given Person for storage.
func (ser *PersonService) Clean(m Model, _ Tx) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Validate returns an error if the Person is not valid for the database.
func (ser *PersonService) Validate(m Model, _ Tx) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Initialize sets initial values for some properties.
func (ser *PersonService) Initialize(_ Model, _ Tx) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing Person in
// updates.
func (ser *PersonService) PersistOldProperties(_ Model, _ Model, _ Tx) error {
	return nil
}

// Marshal transforms the given Person into JSON.
func (ser *PersonService) Marshal(m Model) ([]byte, error) {
	p, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	v, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONMarshal, err)
	}

	return v, nil
}

// Unmarshal parses the given JSON into Person.
func (ser *PersonService) Unmarshal(buf []byte) (Model, error) {
	var p Person
	err := json.Unmarshal(buf, &p)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &p, nil
}

// AssertType exposes the given Model as a Person.
func (ser *PersonService) AssertType(m Model) (*Person, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	p, ok := m.(*Person)
	if !ok {
		return nil, fmt.Errorf("model: %w", errors.New("not of Person type"))
	}
	return p, nil
}

// mapfromModel returns a list of Person type asserted from the given list of
// Model.
func (ser *PersonService) mapFromModel(vlist []Model) ([]*Person, error) {
	list := make([]*Person, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}
