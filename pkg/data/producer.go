package data

import (
	"errors"
	"fmt"
	"strings"

	json "github.com/json-iterator/go"
)

// Producer represents a single studio, producer, licensor, etc.
type Producer struct {
	Titles []Title
	Types  []string
	Meta   ModelMetadata
}

// Metadata return Meta.
func (p *Producer) Metadata() *ModelMetadata {
	return &p.Meta
}

// ProducerService performs operations on Producer.
type ProducerService struct{}

// Create persists the given Producer.
func (ser *ProducerService) Create(p *Producer, tx Tx) (int, error) {
	return tx.Database().Create(p, ser, tx)
}

// Update rplaces the value of the Producer with the given ID.
func (ser *ProducerService) Update(p *Producer, tx Tx) error {
	return tx.Database().Update(p, ser, tx)
}

// Delete deletes the Producer with the given ID.
func (ser *ProducerService) Delete(id int, tx Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// GetAll retrieves all persisted values of Producer.
func (ser *ProducerService) GetAll(first *int, skip *int, tx Tx) ([]*Producer, error) {
	vlist, err := tx.Database().GetAll(first, skip, ser, tx)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to Producers: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of Producer that pass the filter.
func (ser *ProducerService) GetFilter(
	first *int, skip *int, tx Tx, keep func(p *Producer) bool,
) ([]*Producer, error) {
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
		return nil, fmt.Errorf("failed to map Models to Producers: %w", err)
	}
	return list, nil
}

// GetMultiple retrieves the persisted Producer values specified by the
// given IDs that pass the filter.
func (ser *ProducerService) GetMultiple(
	ids []int, first *int, skip *int, tx Tx, keep func(p *Producer) bool,
) ([]*Producer, error) {
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
		return nil, fmt.Errorf("failed to map Models to Producers: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted Producer with the given ID.
func (ser *ProducerService) GetByID(id int, tx Tx) (*Producer, error) {
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

// Bucket returns the name of the bucket for Producer.
func (ser *ProducerService) Bucket() string {
	return "Producer"
}

// Clean cleans the given Producer for storage.
func (ser *ProducerService) Clean(m Model, _ Tx) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	for i, t := range e.Types {
		e.Types[i] = strings.Trim(t, " ")
	}
	return nil
}

// Validate returns an error if the Producer is not valid for the database.
func (ser *ProducerService) Validate(m Model, _ Tx) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Initialize sets initial values for some properties.
func (ser *ProducerService) Initialize(_ Model, _ Tx) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing Producer
// in updates.
func (ser *ProducerService) PersistOldProperties(_ Model, _ Model, _ Tx) error {
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

// mapfromModel returns a list of Producer type asserted from the given list of
// Model.
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
