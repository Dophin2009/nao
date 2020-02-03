package data

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

// TODO: Implement sorting

// Model encompasses all data models.
type Model interface {
	Metadata() *ModelMetadata
}

// ModelMetadata contains information about
type ModelMetadata struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int
}

// Service provides various functions to operate on Models. All implementations
// should use type assertions to guarantee prevention of runtime errors.
type Service interface {
	Database() *bolt.DB
	Bucket() string

	Clean(m Model) error
	Validate(m Model) error

	Initialize(m Model) error
	PersistOldProperties(n Model, o Model) error

	Marshal(m Model) ([]byte, error)
	Unmarshal(buf []byte) (Model, error)
}

var (
	// errNil is an error returned when some pointer is nil.
	errNil = errors.New("is nil")
	// errNotFound is an error returned when the requested
	// object is not found.
	errNotFound = errors.New("not found")
	// errAlreadyExists is an error returned when a unique
	// value already exists.
	errAlreadyExists = errors.New("already exists")
	// errInvalid is an error returned when some value is
	// invalid.
	errInvalid = errors.New("invalid")
)

const (
	errmsgModelCleaning   = "failed to clean model"
	errmsgModelValidation = "failed to validate model"
	errmsgModelInitialize = "failed to initialize model values"
	errmsgModelPersistOld = "failed to persist old model values"
	errmsgModelMarshal    = "failed to marshal model"
	errmsgModelUnmarshal  = "failed to unmarshal model"
	errmsgModelAssertType = "failed to assert type of model"
	errmsgBucketOpen      = "failed to open bucket"
	errmsgBucketNextSeq   = "failed to generate next sequence ID"
	errmsgBucketPut       = "failed to put value in bucket"
	errmsgBucketDelete    = "failed to delete value in bucket"

	errmsgJSONMarshal   = "failed to marshal to JSON"
	errmsgJSONUnmarshal = "failed to unmarshal from JSON"
)

// Create persists the given Model.
func Create(m Model, ser Service) error {
	// Check service
	if err := checkService(ser); err != nil {
		return err
	}

	err := ser.Clean(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelCleaning, err)
	}

	// Verify validity of model
	err = ser.Validate(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelValidation, err)
	}

	return ser.Database().Update(func(tx *bolt.Tx) error {
		// Get bucket, exit if error
		b, err := Bucket(ser.Bucket(), tx)
		if err != nil {
			return fmt.Errorf("%s %q: %w", errmsgBucketOpen, ser.Bucket(), err)
		}

		// Get next ID in sequence and assign to
		// model
		id, err := b.NextSequence()
		if err != nil {
			return fmt.Errorf("%s: %w", errmsgBucketNextSeq, err)
		}

		// Initialize metadata
		meta := m.Metadata()
		meta.ID = int(id)
		meta.CreatedAt = time.Now()
		meta.UpdatedAt = time.Now()
		meta.Version = 0
		err = ser.Initialize(m)
		if err != nil {
			return fmt.Errorf("%s: %w", errmsgModelInitialize, err)
		}

		// Save model in bucket
		buf, err := ser.Marshal(m)
		if err != nil {
			return fmt.Errorf("%s: %w", errmsgModelMarshal, err)
		}

		err = b.Put(itob(meta.ID), buf)
		if err != nil {
			return fmt.Errorf("%s %q: %w", errmsgBucketPut, ser.Bucket(), err)
		}

		return nil
	})
}

// Update replaces the value of the model with the given ID.
func Update(m Model, ser Service) error {
	// Check service
	if err := checkService(ser); err != nil {
		return err
	}

	err := ser.Clean(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelCleaning, err)
	}

	// Verify validity of model
	err = ser.Validate(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelValidation, err)
	}

	return ser.Database().Update(func(tx *bolt.Tx) error {
		meta := m.Metadata()

		// Check if entity with ID exists
		v, err := GetRawByID(meta.ID, ser.Bucket(), ser.Database())
		if err != nil {
			return fmt.Errorf("failed to get by id %d: %w", meta.ID, err)
		}

		// Unmarshall old
		o, err := ser.Unmarshal(v)
		if err != nil {
			return fmt.Errorf("%s: %w", errmsgModelUnmarshal, err)
		}

		// Get bucket, exit if error
		b, err := Bucket(ser.Bucket(), tx)
		if err != nil {
			return fmt.Errorf("%s %q: %w", errmsgBucketOpen, ser.Bucket(), err)
		}

		// Replace properties of updated with certain frozen
		// ones of old
		meta.UpdatedAt = time.Now()
		meta.Version = meta.Version + 1
		err = ser.PersistOldProperties(m, o)
		if err != nil {
			return fmt.Errorf("%s: %w", errmsgModelPersistOld, err)
		}

		// Save model
		buf, err := ser.Marshal(m)
		if err != nil {
			return fmt.Errorf("%s: %w", errmsgModelMarshal, err)
		}

		err = b.Put(itob(meta.ID), buf)
		if err != nil {
			return fmt.Errorf("%s %q: %w", errmsgBucketPut, ser.Bucket(), err)
		}

		return nil
	})
}

// Delete deletes the model with the given ID.
func Delete(id int, ser Service) error {
	// Check service
	if err := checkService(ser); err != nil {
		return err
	}

	return ser.Database().Update(func(tx *bolt.Tx) error {
		// Get bucket, exit if error
		b, err := Bucket(ser.Bucket(), tx)
		if err != nil {
			return fmt.Errorf("%s %q: %w", errmsgBucketOpen, ser.Bucket(), err)
		}

		// Store existing model to return
		err = b.Delete(itob(id))
		if err != nil {
			return fmt.Errorf("failed to delete by id %d: %w", id, err)
		}

		return nil
	})
}

// GetByID retrieves the persisted Model with the given ID. The given service
// and its DB should not be nil.
func GetByID(id int, ser Service) (Model, error) {
	// Check service
	if err := checkService(ser); err != nil {
		return nil, err
	}

	v, err := GetRawByID(id, ser.Bucket(), ser.Database())
	if err != nil {
		return nil, fmt.Errorf("failed to get by id %d: %w", id, err)
	}

	m, err := ser.Unmarshal(v)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelMarshal, err)
	}

	return m, nil
}

// GetRawByID is a generic function that queries the given bucket in the given
// database for an entity of the given ID. The given DB pointer should not be
// nil.
func GetRawByID(id int, bucketName string, db *bolt.DB) ([]byte, error) {
	if db == nil {
		return nil, fmt.Errorf("DB: %w", errNil)
	}

	// Raw value to return
	var v []byte

	// Begin database transaction
	err := db.View(func(tx *bolt.Tx) error {
		// Get bucket, exit if error
		b, err := Bucket(bucketName, tx)
		if err != nil {
			return fmt.Errorf("%s %q: %w", errmsgBucketOpen, bucketName, err)
		}

		// Get entity by ID, exit if error
		v = b.Get(itob(id))
		if v == nil {
			return fmt.Errorf("model with id %d: %w", id, errNotFound)
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	return v, nil
}

// GetMultiple retrieves the persisted instances of a Model type with the given
// IDs.
//
// See GetFilter for details on `first` and `skip`.
func GetMultiple(
	ser Service, ids []int, first *int, skip *int,
	keep func(m Model) bool,
) ([]Model, error) {
	// Check service
	if err := checkService(ser); err != nil {
		return nil, err
	}

	if keep == nil {
		keep = func(m Model) bool {
			return true
		}
	}

	start, end := calculatePaginationBounds(first, skip)

	// List to return
	list := []Model{}

	i := start
	for _, id := range ids {
		if i >= end {
			break
		}

		m, err := GetByID(id, ser)
		if err != nil {
			return nil, fmt.Errorf("failed to get Model by id %d: %w", id, err)
		}

		if !keep(m) {
			continue
		}

		list = append(list, m)
		i++
	}

	return list, nil
}

// GetAll retrieves all persisted instances of a Model type with the given data
// layer service.
//
// See GetFilter for details on `first` and `skip`.
func GetAll(ser Service, first *int, skip *int) ([]Model, error) {
	return GetFilter(ser, first, skip, func(m Model) bool { return true })
}

// GetFilter retrieves all persisted instances of a Model type that pass the
// filter.
//
// Collection begins on the first valid element after skipping the `skip` valid
// elements and continues for `first` valid elements that pass the filter. If
// `skip` is given as nil, collection begins with the first valid element. If
// `first` is given as nil, collection continues until the last persisted
// element is queried. The given service and its DB should not be nil. A nil
// filter function passes all.
func GetFilter(
	ser Service, first *int, skip *int,
	keep func(m Model) bool,
) ([]Model, error) {
	// Check service
	if err := checkService(ser); err != nil {
		return nil, err
	}

	if keep == nil {
		keep = func(m Model) bool {
			return true
		}
	}

	start, end := calculatePaginationBounds(first, skip)

	// List to return
	list := []Model{}

	// Begin database transaction
	err := ser.Database().View(func(tx *bolt.Tx) error {
		// Get bucket, exit if error
		b, err := Bucket(ser.Bucket(), tx)
		if err != nil {
			return fmt.Errorf("%s %q: %w", errmsgBucketOpen, ser.Bucket(), err)
		}

		// Get cursor for bucket
		c := b.Cursor()

		// Find last key; will stop iteration when reached
		var lk []byte
		d := b.Cursor()
		lk, _ = d.Last()

		// Move cursor to starting element
		var k, v []byte
		c.First()
		for i := 0; i < start; k, v = c.Next() {
			if k == nil {
				continue
			}
			i++
		}

		// Iterate until end is reached
		lastReached := false
		for i := start; i < end; k, v = c.Next() {
			if lastReached {
				break
			}

			// If key reached last, exit
			if bytes.Equal(k, lk) {
				lastReached = true
			}

			// If no key found, continue to next
			if k == nil {
				continue
			}

			// Unmarshal element
			m, err := ser.Unmarshal(v)
			if err != nil {
				return fmt.Errorf("%s: %w", errmsgModelUnmarshal, err)
			}

			// If element does not pass filter, continue to next
			if !keep(m) {
				continue
			}

			list = append(list, m)
			i++
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}

// iterateKeys iterates through the keys of the given database bucket and
// passes the value at each key to some function.
//
// The loop exits only if do returns a true exit flag. If an error is returned
// but exit is false, the error will be ignored. If an error is returned and
// exit is true, the error will be returned.
func iterateKeys(b *bolt.Bucket, do func(v []byte) (exit bool, err error)) error {
	var lk []byte
	d := b.Cursor()
	lk, _ = d.Last()

	c := b.Cursor()
	lastReached := false
	for k, v := c.First(); !lastReached; k, v = c.Next() {
		if bytes.Equal(k, lk) {
			lastReached = true
		}

		if k == nil {
			continue
		}

		exit, err := do(v)
		if !exit {
			return fmt.Errorf("failed to perform operations on value: %w", err)
		}
	}

	return nil
}

// get returns the raw value stored at the given key in the given bucket.
func get(id int, bucket *bolt.Bucket) ([]byte, error) {
	if bucket == nil {
		return nil, fmt.Errorf("bucket: %w", errNil)
	}

	v := bucket.Get(itob(id))
	if v == nil {
		return nil, fmt.Errorf("model with id %d: %w", id, errNotFound)
	}
	return v, nil
}

func calculatePaginationBounds(first *int, skip *int) (int, int) {
	// The number of elements to skip
	var start int
	if skip == nil || *skip <= 0 {
		start = 0
	} else {
		start = *skip
	}

	// When iterator reaches this number, stop
	var end int
	if first == nil || *first < 0 {
		// Return all elements if `first` is nil
		end = -1
	} else if *first == 0 {
		end = start
	} else {
		end = start + *first
	}

	return start, end
}

// checkService returns an error if the given service or its DB are nil.
func checkService(ser Service) error {
	if ser == nil {
		return fmt.Errorf("service: %w", errNil)
	}
	if ser.Database() == nil {
		return fmt.Errorf("DB: %w", errNil)
	}
	return nil
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
