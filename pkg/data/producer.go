package data

import (
	"errors"
	"fmt"
	"strings"

	json "github.com/json-iterator/go"
	"gitlab.com/Dophin2009/nao/pkg/db"
)

// Producer represents a single studio, producer, licensor, etc.
type Producer struct {
	Titles []Title
	Types  []string
	Meta   db.ModelMetadata
}

// Metadata return Meta.
func (p *Producer) Metadata() *db.ModelMetadata {
	return &p.Meta
}

// ProducerService performs operations on Producer.
type ProducerService struct {
	Hooks db.PersistHooks
}

// Create persists the given Producer.
func (ser *ProducerService) Create(p *Producer, tx db.Tx) (int, error) {
	return tx.Database().Create(p, ser, tx)
}

// Update rplaces the value of the Producer with the given ID.
func (ser *ProducerService) Update(p *Producer, tx db.Tx) error {
	return tx.Database().Update(p, ser, tx)
}

// Delete deletes the Producer with the given ID.
func (ser *ProducerService) Delete(id int, tx db.Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// GetAll retrieves all persisted values of Producer.
func (ser *ProducerService) GetAll(first *int, skip *int, tx db.Tx) ([]*Producer, error) {
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
	first *int, skip *int, tx db.Tx, keep func(p *Producer) bool,
) ([]*Producer, error) {
	vlist, err := tx.Database().GetFilter(first, skip, ser, tx,
		func(m db.Model) bool {
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
	ids []int, first *int, tx db.Tx, keep func(p *Producer) bool,
) ([]*Producer, error) {
	vlist, err := tx.Database().GetMultiple(ids, first, ser, tx,
		func(m db.Model) bool {
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
func (ser *ProducerService) GetByID(id int, tx db.Tx) (*Producer, error) {
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
func (ser *ProducerService) Clean(m db.Model, _ db.Tx) error {
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
func (ser *ProducerService) Validate(m db.Model, _ db.Tx) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Initialize sets initial values for some properties.
func (ser *ProducerService) Initialize(_ db.Model, _ db.Tx) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing Producer
// in updates.
func (ser *ProducerService) PersistOldProperties(_ db.Model, _ db.Model, _ db.Tx) error {
	return nil
}

// PersistHooks returns the persistence hook functions.
func (ser *ProducerService) PersistHooks() *db.PersistHooks {
	return &ser.Hooks
}

// Marshal transforms the given Producer into JSON.
func (ser *ProducerService) Marshal(m db.Model) ([]byte, error) {
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
func (ser *ProducerService) Unmarshal(buf []byte) (db.Model, error) {
	var p Producer
	err := json.Unmarshal(buf, &p)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &p, nil
}

// AssertType exposes the given Model as a Producer.
func (ser *ProducerService) AssertType(m db.Model) (*Producer, error) {
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
func (ser *ProducerService) mapFromModel(vlist []db.Model) ([]*Producer, error) {
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
