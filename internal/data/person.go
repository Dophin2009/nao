package data

import (
	"errors"
	"fmt"

	json "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
)

// TODO: User rating/favoriting/comments/etc. of Persons

// Person represents a single person
type Person struct {
	ID          int
	Names       []Title
	Information []Title
	Version     int
	Model
}

// Iden returns the ID.
func (p *Person) Iden() int {
	return p.ID
}

// PersonBucket is the name of the database bucket for
// Person.
const PersonBucket = "Person"

// PersonService performs operations on Persons.
type PersonService struct {
	DB *bolt.DB
	Service
}

// Create persists the given Person.
func (ser *PersonService) Create(p *Person) error {
	return Create(p, ser)
}

// Update rplaces the value of the Person with the
// given ID.
func (ser *PersonService) Update(p *Person) error {
	return Update(p, ser)
}

// Delete deletes the Person with the given ID.
func (ser *PersonService) Delete(id int) error {
	return Delete(id, ser)
}

// GetAll retrieves all persisted values of Person.
func (ser *PersonService) GetAll(first *int, skip *int) ([]*Person, error) {
	vlist, err := GetAll(ser, first, skip)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to Persons: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of Person that
// pass the filter.
func (ser *PersonService) GetFilter(first *int, skip *int, keep func(p *Person) bool) ([]*Person, error) {
	vlist, err := GetFilter(ser, first, skip, func(m Model) bool {
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
func (ser *PersonService) GetByID(id int) (*Person, error) {
	m, err := GetByID(id, ser)
	if err != nil {
		return nil, err
	}

	p, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return p, nil
}

// Database returns the database reference.
func (ser *PersonService) Database() *bolt.DB {
	return ser.DB
}

// Bucket returns the name of the bucket for Person.
func (ser *PersonService) Bucket() string {
	return PersonBucket
}

// Clean cleans the given Person for storage
func (ser *PersonService) Clean(m Model) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Validate returns an error if the Person is
// not valid for the database.
func (ser *PersonService) Validate(m Model) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Initialize sets initial values for some properties.
func (ser *PersonService) Initialize(m Model, id int) error {
	p, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	p.ID = id
	p.Version = 0
	return nil
}

// PersistOldProperties maintains certain properties
// of the existing Person in updates.
func (ser *PersonService) PersistOldProperties(n Model, o Model) error {
	np, err := ser.AssertType(n)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	op, err := ser.AssertType(o)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	np.Version = op.Version + 1
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

// mapfromModel returns a list of Person type
// asserted from the given list of Model.
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
