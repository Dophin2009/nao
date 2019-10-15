package data

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

// Producer represents a single studio, producer,
// licensor, etc.
type Producer struct {
	ID      int
	Titles  []Info
	Type    string
	Version int
}

// Identifier returns the ID of the Producer
func (p *Producer) Identifier() int {
	return p.ID
}

// SetIdentifier sets the ID of the Producer
func (p *Producer) SetIdentifier(ID int) {
	p.ID = ID
}

// Ver returns the verison of the Producer
func (p *Producer) Ver() int {
	return p.Version
}

// UpdateVer increments the version of the
// Producer by one
func (p *Producer) UpdateVer() {
	p.Version++
}

// Validate returns an error if the Producer is
// not valid for the database
func (p *Producer) Validate(tx *bolt.Tx) (err error) {
	return nil
}

// ProducerBucketName provides the database bucket name
// for the Producer entity
const producerBucketName = "Producer"

// ProducerGet retrieves a single instance of Producer with
// the given ID
func ProducerGet(ID int, db *bolt.DB) (p Producer, err error) {
	err = getByID(ID, &p, producerBucketName, db)
	return
}

// ProducerGetAll retrieves all persisted Producer values
func ProducerGetAll(db *bolt.DB) (list []Producer, err error) {
	return ProducerGetFilter(db, func(p *Producer) bool { return true })
}

// ProducerGetFilter retrieves all persisted Producer values
// that pass the filter
func ProducerGetFilter(db *bolt.DB, filter func(p *Producer) bool) (list []Producer, err error) {
	ilist, err := getFilter(&Producer{}, func(entity Idenitifiable) (bool, error) {
		p, ok := entity.(*Producer)
		if !ok {
			return false, fmt.Errorf("type assertion failed: entity is not a Producer")
		}
		return filter(p), nil
	}, producerBucketName, db)

	list = make([]Producer, len(ilist))
	for i, p := range ilist {
		list[i] = *p.(*Producer)
	}

	return
}

// ProducerCreate persists a new instance of Producer
func ProducerCreate(p *Producer, db *bolt.DB) error {
	return create(p, producerBucketName, db)
}

// ProducerUpdate updates the properties of an existing
// persisted Producer instance
func ProducerUpdate(p *Producer, db *bolt.DB) error {
	return update(p, producerBucketName, db)
}
