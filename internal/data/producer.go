package data

import (
	"errors"
	"fmt"
	"strings"

	json "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
)

// Producer represents a single studio, producer,
// licensor, etc.
type Producer struct {
	ID      int
	Titles  []Title
	Types   []string
	Version int
	Model
}

// Iden returns the ID.
func (p *Producer) Iden() int {
	return p.ID
}

// ProducerBucket is the name of the database bucket for
// Producer.
const ProducerBucket = "Producer"

// ProducerService performs operations on Producer.
type ProducerService struct {
	DB *bolt.DB
	Service
}

// Create persists the given Producer.
func (ser *ProducerService) Create(p *Producer) error {
	return Create(p, ser)
}

// Update rplaces the value of the Producer with the
// given ID.
func (ser *ProducerService) Update(p *Producer) error {
	return Update(p, ser)
}

// Delete deletes the Producer with the given ID.
func (ser *ProducerService) Delete(id int) error {
	return Delete(id, ser)
}

// GetAll retrieves all persisted values of Producer.
func (ser *ProducerService) GetAll(first *int, skip *int) ([]*Producer, error) {
	vlist, err := GetAll(ser, first, skip)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to Producers: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of Producer that
// pass the filter.
func (ser *ProducerService) GetFilter(first *int, skip *int, keep func(p *Producer) bool) ([]*Producer, error) {
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
		return nil, fmt.Errorf("failed to map Models to Producers: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted Producer with the given ID.
func (ser *ProducerService) GetByID(id int) (*Producer, error) {
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
func (ser *ProducerService) Database() *bolt.DB {
	return ser.DB
}

// Bucket returns the name of the bucket for Producer.
func (ser *ProducerService) Bucket() string {
	return ProducerBucket
}

// Clean cleans the given Producer for storage.
func (ser *ProducerService) Clean(m Model) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	for i, t := range e.Types {
		e.Types[i] = strings.Trim(t, " ")
	}
	return nil
}

// Validate returns an error if the Producer is
// not valid for the database.
func (ser *ProducerService) Validate(m Model) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Initialize sets initial values for some properties.
func (ser *ProducerService) Initialize(m Model, id int) error {
	p, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	p.ID = id
	p.Version = 0
	return nil
}

// PersistOldProperties maintains certain properties
// of the existing Producer in updates
func (ser *ProducerService) PersistOldProperties(n Model, o Model) error {
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

// Marshal transforms the given Producer into JSON.
func (ser *ProducerService) Marshal(m Model) ([]byte, error) {
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

// Unmarshal parses the given JSON into Producer.
func (ser *ProducerService) Unmarshal(buf []byte) (Model, error) {
	var p Producer
	err := json.Unmarshal(buf, &p)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &p, nil
}

// AssertType exposes the given Model as a Producer.
func (ser *ProducerService) AssertType(m Model) (*Producer, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	p, ok := m.(*Producer)
	if !ok {
		return nil, fmt.Errorf("model: %w", errors.New("not of Producer type"))
	}
	return p, nil
}

// mapfromModel returns a list of Producer type
// asserted from the given list of Model.
func (ser *ProducerService) mapFromModel(vlist []Model) ([]*Producer, error) {
	list := make([]*Producer, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}
