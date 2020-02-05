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

// Service provides various functions to operate on Models. All implementations
// should use type assertions to guarantee prevention of runtime errors.
type Service interface {
	Bucket() string

	Clean(m Model, tx Tx) error
	Validate(m Model, tx Tx) error
	Initialize(m Model, tx Tx) error
	PersistOldProperties(n Model, o Model, tx Tx) error

	Marshal(m Model) ([]byte, error)
	Unmarshal(buf []byte) (Model, error)
}

// Database defines generic CRUD operations for opaque Model objects for a
// database.
type Database interface {
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
	Database() Database
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
