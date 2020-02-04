package data

import (
	"encoding/binary"
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

// Buckets provides an array of all the buckets in the database
func Buckets() []string {
	return []string{
		MediaBucket, ProducerBucket, GenreBucket, EpisodeBucket, EpisodeSetBucket,
		CharacterBucket, PersonBucket, UserBucket, MediaProducerBucket,
		MediaRelationBucket, MediaGenreBucket, MediaCharacterBucket,
		UserCharacterBucket, UserEpisodeBucket, UserMediaBucket,
		UserMediaListBucket, UserPersonBucket, JWTBucket,
	}
}

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
