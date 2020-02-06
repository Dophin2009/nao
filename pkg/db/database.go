package db

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"
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

// DatabaseService provides
type DatabaseService struct {
	DatabaseDriver
}

// Create persists a new instance of a Model type.
func (dbs *DatabaseService) Create(m Model, ser Service, tx Tx) (int, error) {
	hooks := ser.PersistHooks()
	if hooks != nil {
		err := hooks.PreCreateHook(m, ser, tx)
		if err != nil {
			return 0, fmt.Errorf("failed to run pre-create hooks: %w", err)
		}
	}

	id, err := dbs.DatabaseDriver.Create(m, ser, tx)
	if err != nil {
		return 0, err
	}

	if hooks != nil {
		err = hooks.PostCreateHook(m, ser, tx)
		if err != nil {
			return 0, fmt.Errorf("failed to run post-create hooks: %w", err)
		}
	}

	return id, nil
}

// Update modifies an existing instance of a Model type.
func (dbs *DatabaseService) Update(m Model, ser Service, tx Tx) error {
	hooks := ser.PersistHooks()
	if hooks != nil {
		err := hooks.PreUpdateHook(m, ser, tx)
		if err != nil {
			return fmt.Errorf("failed to run pre-update hooks: %w", err)
		}
	}

	err := dbs.DatabaseDriver.Update(m, ser, tx)
	if err != nil {
		return err
	}

	if hooks != nil {
		err = hooks.PostUpdateHook(m, ser, tx)
		if err != nil {
			return fmt.Errorf("failed to run post-update hooks: %w", err)
		}
	}

	return nil
}

// Delete deletes an existing persisted instance of a Model type.
func (dbs *DatabaseService) Delete(id int, ser Service, tx Tx) error {
	hooks := ser.PersistHooks()
	m, err := dbs.DatabaseDriver.GetByID(id, ser, tx)
	if err != nil {
		return err
	}

	if hooks != nil {
		err = hooks.PreDeleteHook(m, ser, tx)
		if err != nil {
			return fmt.Errorf("failed to run pre-delete hooks: %w", err)
		}
	}

	err = dbs.DatabaseDriver.Delete(id, ser, tx)
	if err != nil {
		return err
	}

	if hooks != nil {
		err = ser.PersistHooks().PostDeleteHook(m, ser, tx)
		if err != nil {
			return fmt.Errorf("failed to run post-delete hooks: %w", err)
		}
	}

	return nil
}

// DatabaseDriver defines generic CRUD logic for a database backend.
type DatabaseDriver interface {
	Transaction(writable bool, logic func(Tx) error) error
	Close() error

	Create(m Model, ser Service, tx Tx) (int, error)
	Update(m Model, ser Service, tx Tx) error
	Delete(id int, ser Service, tx Tx) error
	GetByID(id int, ser Service, tx Tx) (Model, error)
	GetRawByID(id int, ser Service, tx Tx) ([]byte, error)
	GetMultiple(ids []int, first *int, skip *int, ser Service, tx Tx,
		keep func(Model) bool) ([]Model, error)
	GetAll(first *int, skip *int, ser Service, tx Tx) ([]Model, error)
	GetFilter(first *int, skip *int, ser Service, tx Tx,
		keep func(Model) bool) ([]Model, error)
	FindFirst(ser Service, tx Tx, match func(Model) (bool, error)) (Model, error)
}

// Tx defines a wrapper for database transactions objects.
type Tx interface {
	Database() *DatabaseService
	Unwrap() interface{}
}

var (
	// errNil is an error returned when some pointer is nil.
	errNil = errors.New("is nil")
	// errNotFound is an error returned when the requested object is not found.
	errNotFound = errors.New("not found")
	// errAlreadyExists is an error returned when a unique value already exists.
	errAlreadyExists = errors.New("already exists")
	// errInvalid is an error returned when some value is invalid.
	errInvalid = errors.New("invalid")
	// errUnwritableTx is an error returned when an update attempt was made with
	// a transaction object that does now allow updates.
	errUnwritableTx = errors.New("read-only transaction")
)

const (
	errmsgModelCleaning   = "failed to clean model"
	errmsgModelValidation = "failed to validate model"
	errmsgModelInitialize = "failed to initialize model values"
	errmsgModelPersistOld = "failed to persist old model values"
	errmsgModelMarshal    = "failed to marshal model"
	errmsgModelUnmarshal  = "failed to unmarshal model"
	errmsgBucketOpen      = "failed to open bucket"
	errmsgBucketNextSeq   = "failed to generate next sequence ID"
	errmsgBucketPut       = "failed to put value in bucket"
	errmsgBucketDelete    = "failed to delete value in bucket"
)

// checkService returns an error if the given service or its DB are nil.
func checkService(ser Service) error {
	if ser == nil {
		return fmt.Errorf("service: %w", errNil)
	}
	return nil
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
