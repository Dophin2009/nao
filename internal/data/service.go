package data

import (
	"encoding/binary"
	"errors"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

// TODO: Implement sorting

// Model encompasses all data models.
type Model interface {
	Iden() int
}

// Service provides various functions to operate on Models.
// All implementations should use type assertions to guarantee
// prevention of runtime errors.
type Service interface {
	Database() *bolt.DB
	Bucket() string

	Clean(m Model) error
	Validate(m Model) error

	Initialize(m Model, id int) error
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

		err = ser.Initialize(m, int(id))
		if err != nil {
			return fmt.Errorf("%s: %w", errmsgModelInitialize, err)
		}

		// Save model in bucket
		buf, err := ser.Marshal(m)
		if err != nil {
			return fmt.Errorf("%s: %w", errmsgModelMarshal, err)
		}

		err = b.Put(itob(m.Iden()), buf)
		if err != nil {
			return fmt.Errorf("%s %q: %w", errmsgBucketPut, ser.Bucket(), err)
		}

		return nil
	})
}

// Update replaces the value of the model with the given
// ID.
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
		// Check if entity with ID exists
		v, err := GetRawByID(m.Iden(), ser.Bucket(), ser.Database())
		if err != nil {
			return fmt.Errorf("failed to get by id %d: %w", m.Iden(), err)
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
		err = ser.PersistOldProperties(m, o)
		if err != nil {
			return fmt.Errorf("%s: %w", errmsgModelPersistOld, err)
		}

		// Save model
		buf, err := ser.Marshal(m)
		if err != nil {
			return fmt.Errorf("%s: %w", errmsgModelMarshal, err)
		}

		err = b.Put(itob(m.Iden()), buf)
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

// GetByID retrieves the persisted Model with the given ID.
// The given service and its DB should not be nil.
func GetByID(id int, ser Service) (Model, error) {
	if ser == nil {
		return nil, fmt.Errorf("service: %w", errNil)
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

// GetRawByID is a generic function that queries the given bucket
// in the given database for an entity of the given ID.
// The given DB pointer should not be nil.
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

// GetAll retrieves all persisted instances of a Model type
// with the given data layer service.
// Collection begins on the element after the given prefix
// and continues for `first` elements.
// If prefixID is given as nil, collection begins with the
// first persisted element.
// The given service and its DB should not be nil.
func GetAll(ser Service, first int, prefixID *int) ([]Model, error) {
	return GetFilter(ser, first, prefixID, func(m Model) bool { return true })
}

// GetFilter retrieves all persisted instances of a Model
// type that pass the filter.
// Collection begins on the element after the given prefix
// and continues for `first` elements that pass the filter.
// If prefixID is given as nil, collection begins with the
// first persisted element.
// The given service and its DB should not be nil.
// The filter function should also not be nil.
func GetFilter(ser Service, first int, prefixID *int, keep func(m Model) bool) ([]Model, error) {
	// Check service
	if err := checkService(ser); err != nil {
		return nil, err
	}

	// Filter function should not be nil.
	if keep == nil {
		return nil, fmt.Errorf("filter function: %w", errNil)
	}

	// Return empty slice if number of elements to get is 0
	if first == 0 {
		return []Model{}, nil
	}

	// List to return
	var list []Model

	// Begin database transaction
	err := ser.Database().View(func(tx *bolt.Tx) error {
		// Get bucket, exit if error
		b, err := Bucket(ser.Bucket(), tx)
		if err != nil {
			return fmt.Errorf("%s %q: %w", errmsgBucketOpen, ser.Bucket(), err)
		}

		// Get cursor for bucket
		c := b.Cursor()

		// Initialize key and value to first element or
		// right after prefix
		var k, v []byte
		if prefixID == nil {
			// Begin on first element if prefix is not provided
			k, v = c.First()
			// If first key not found, database is empty;
			// return empty list
			if k == nil {
				list = []Model{}
				return nil
			}
		} else {
			// Begin on element right after prefix if prefixID is
			// provided
			k, v = c.Seek(itob(*prefixID))
			// If key not found, return error
			if k == nil {
				return fmt.Errorf("prefix key: %w", errNotFound)
			}
			k, v = c.Next()
		}

		if first < 0 {
			// If negative `first`, get all elements starting
			// from prefix
			for ; k != nil; k, v = c.Next() {
				// Unmarshal value and check filter
				m, err := ser.Unmarshal(v)
				if err != nil {
					return fmt.Errorf("%s: %w", errmsgModelUnmarshal, err)
				}

				if keep(m) {
					list = append(list, m)
				}
			}
		} else {
			// Positive `first` means to get that many elements;
			// construct a slice of that size
			list = make([]Model, first)

			// Unmarshal and add all entities to slice
			i := 0
			for ; k != nil && i < first; k, v = c.Next() {
				// Unmarshal value and check filter
				m, err := ser.Unmarshal(v)
				if err != nil {
					return fmt.Errorf("%s: %w", errmsgModelUnmarshal, err)
				}

				if keep(m) {
					list[i] = m
					i++
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}

// get returns the raw value stored at the given key in the
// given bucket.
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

// checkService returns an error if the given service or its
// DB are nil.
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
