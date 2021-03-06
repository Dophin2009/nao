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
	DeletedAt *time.Time
	Version   int
}

// DatabaseService provides
type DatabaseService struct {
	DatabaseDriver
}

// Create persists a new instance of a Model type.
func (dbs *DatabaseService) Create(m Model, ser Service, tx Tx) (int, error) {
	// Check service
	err := CheckService(ser)
	if err != nil {
		return 0, err
	}

	// Verify validity of model
	err = ser.Validate(m, tx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errmsgModelValidation, err)
	}

	// Clean model
	err = ser.Clean(m, tx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errmsgModelCleaning, err)
	}

	// Initialize metadata
	meta := m.Metadata()
	meta.CreatedAt = time.Now()
	meta.UpdatedAt = meta.CreatedAt
	meta.Version = 0
	err = ser.Initialize(m, tx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errmsgModelInitialize, err)
	}

	// Call hooks to run before create
	hooks := ser.PersistHooks()
	if hooks != nil {
		err := hooks.PreCreateHook(m, ser, tx)
		if err != nil {
			return 0, fmt.Errorf("failed to run pre-create hooks: %w", err)
		}
	}

	// Persist
	id, err := dbs.DatabaseDriver.Create(m, ser, tx)
	if err != nil {
		return 0, err
	}

	// Call hooks to run after create
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
	// Check service
	err := CheckService(ser)
	if err != nil {
		return err
	}

	// Check if entity with ID exists
	o, err := dbs.DatabaseDriver.GetByID(m.Metadata().ID, ser, tx)
	if err != nil {
		return fmt.Errorf("failed to get by id %d: %w", m.Metadata().ID, err)
	}

	// Verify validity of model
	err = ser.Validate(m, tx)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelValidation, err)
	}

	// Prepare
	err = ser.Clean(m, tx)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelCleaning, err)
	}

	// Replace properties of updated with certain frozen
	// ones of old
	meta := m.Metadata()
	meta.UpdatedAt = time.Now()
	meta.Version = meta.Version + 1
	err = ser.PersistOldProperties(m, o, tx)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelPersistOld, err)
	}

	// Call hooks to run before update
	hooks := ser.PersistHooks()
	if hooks != nil {
		err := hooks.PreUpdateHook(m, ser, tx)
		if err != nil {
			return fmt.Errorf("failed to run pre-update hooks: %w", err)
		}
	}

	// Update in database
	err = dbs.DatabaseDriver.Update(m, ser, tx)
	if err != nil {
		return err
	}

	// Call hooks to run after update
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
	// Check service
	err := CheckService(ser)
	if err != nil {
		return err
	}

	hooks := ser.PersistHooks()

	// Get existing value
	m, err := dbs.DatabaseDriver.GetByID(id, ser, tx)
	if err != nil {
		return err
	}

	// Call hooks to run before deletion
	if hooks != nil {
		err = hooks.PreDeleteHook(m, ser, tx)
		if err != nil {
			return fmt.Errorf("failed to run pre-delete hooks: %w", err)
		}
	}

	// Delete
	err = dbs.DatabaseDriver.Delete(id, ser, tx)
	if err != nil {
		return err
	}

	// Call hooks to run after deletion
	if hooks != nil {
		err = ser.PersistHooks().PostDeleteHook(m, ser, tx)
		if err != nil {
			return fmt.Errorf("failed to run post-delete hooks: %w", err)
		}
	}

	return nil
}

// DeleteMultiple deletes the existing persisted instances of a Model
// type specified by the given IDs.
func (dbs *DatabaseService) DeleteMultiple(ids []int, first *int,
	ser Service, tx Tx, iff func(Model) bool) error {
	err := dbs.DoMultiple(ids, ser, tx, dbs.deleteWrapper(), iff)
	if err != nil {
		return err
	}

	return nil
}

// DeleteFilter deletes all the persisted instances of a Model type
// that pass the filer function.
func (dbs *DatabaseService) DeleteFilter(ser Service, tx Tx,
	iff func(Model) bool) error {
	err := dbs.DoEach(nil, nil, ser, tx, dbs.deleteWrapper(), iff)
	if err != nil {
		return err
	}

	return nil
}

func (dbs *DatabaseService) deleteWrapper() func(m Model, ser Service, tx Tx) (exit bool, err error) {
	return func(m Model, ser Service, tx Tx) (exit bool, err error) {
		err = dbs.Delete(m.Metadata().ID, ser, tx)
		if err != nil {
			return true, err
		}
		return false, nil
	}
}

// GetMultiple retrieves the persisted instances of a Model type with the given
// IDs.
//
// See GetFilter for details on `first` and `skip`.
func (dbs *DatabaseService) GetMultiple(ids []int, ser Service, tx Tx,
	keep func(m Model) bool) ([]Model, error) {
	list := []Model{}
	collect := func(m Model, _ Service, _ Tx) (exit bool, err error) {
		list = append(list, m)
		return false, nil
	}

	err := dbs.DoMultiple(ids, ser, tx, collect, keep)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetAll retrieves all persisted instances of a Model type with the given data
// layer service.
//
// See GetFilter for details on `first` and `skip`.
func (dbs *DatabaseService) GetAll(first *int, skip *int, ser Service, tx Tx) ([]Model, error) {
	return dbs.GetFilter(first, skip, ser, tx, nil)
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
func (dbs *DatabaseService) GetFilter(first *int, skip *int, ser Service, tx Tx,
	keep func(m Model) bool) ([]Model, error) {
	list := []Model{}
	collect := func(m Model, ser Service, tx Tx) (exit bool, err error) {
		// Append element to list
		list = append(list, m)
		return false, nil
	}

	err := dbs.DoEach(first, skip, ser, tx, collect, keep)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// DatabaseDriver defines generic CRUD logic for a database backend.
type DatabaseDriver interface {
	Transaction(writable bool, logic func(Tx) error) error
	Close() error

	DoMultiple(ids []int, ser Service, tx Tx,
		do func(Model, Service, Tx) (exit bool, err error), iff func(Model) bool) error
	DoEach(first *int, skip *int, ser Service, tx Tx,
		do func(Model, Service, Tx) (exit bool, err error), iff func(Model) bool) error
	FindFirst(ser Service, tx Tx, match func(Model) (exit bool, err error)) (Model, error)

	Create(m Model, ser Service, tx Tx) (int, error)
	Update(m Model, ser Service, tx Tx) error
	// Delete marks the model with the given ID as deleted.
	Delete(id int, ser Service, tx Tx) error
	GetByID(id int, ser Service, tx Tx) (Model, error)
	GetRawByID(id int, ser Service, tx Tx) ([]byte, error)
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

// CheckService returns an error if the given service or its DB are nil.
func CheckService(ser Service) error {
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
