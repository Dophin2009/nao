package data

import (
	"encoding/json"

	"github.com/cheekybits/genny/generic"
	bolt "go.etcd.io/bbolt"
)

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "EntityType=Media,Episode,Character,Genre,Producer,Person,User,MediaRelation,MediaCharacter,MediaGenre,MediaProducer,UserMedia"

// EntityType is a generic placeholder for all entity types;
// it is assumed that EntityType structs have an ID and Version,
// both of which should be of int type.
type EntityType generic.Type

// EntityTypeService is a struct that performs
// CRUD operations on the persistence layer.
type EntityTypeService struct {
	DB *bolt.DB
}

// EntityTypeBucketName represents the database bucket
// name for the bucket that stores EntityType instances
const EntityTypeBucketName = "EntityType"

// GetByID retrieves a single isntance of EntityType
// with the given ID.
func (ser *EntityTypeService) GetByID(e *EntityType) (err error) {
	v, err := GetByID(e.ID, EntityTypeBucketName, ser.DB)
	if err != nil {
		return err
	}

	err = json.Unmarshal(v, e)
	if err != nil {
		return err
	}
	return
}

// GetAll retrieves all persisted instances of
// EntityType.
func (ser *EntityTypeService) GetAll() (list []EntityType, err error) {
	return ser.GetFilter(func(e *EntityType) bool { return true })
}

// GetFilter retrieves all persisted instances of
// EntityType that pass the filter.
func (ser *EntityTypeService) GetFilter(keep func(e *EntityType) bool) (list []EntityType, err error) {
	vlist, err := GetAll(EntityTypeBucketName, ser.DB)
	for _, v := range vlist {
		var e EntityType
		err = json.Unmarshal(v, &e)
		if err != nil {
			return nil, err
		}

		if keep(&e) {
			list = append(list, e)
		}
	}
	return
}

// Create persists a new instance of EntityType to the database.
func (ser *EntityTypeService) Create(e *EntityType) (err error) {
	// Verify validity of struct
	err = ser.Validate(e)
	if err != nil {
		return err
	}

	return ser.DB.Update(func(tx *bolt.Tx) error {
		// Get bucket, exit if error
		b, err := Bucket(EntityTypeBucketName, tx)
		if err != nil {
			return err
		}

		// Get next ID in sequence and
		// assign to entity
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		e.ID = int(id)
		e.Version = 0

		// Save entity in bucket
		buf, err := json.Marshal(e)
		if err != nil {
			return err
		}

		return b.Put(itob(int(id)), buf)
	})
}

// Update replaces the persisted instance
// of EntityType with the given EntityType of
// the same ID.
func (ser *EntityTypeService) Update(e *EntityType) (err error) {
	// Verify validity of struct
	err = ser.Validate(e)
	if err != nil {
		return err
	}

	return ser.DB.Update(func(tx *bolt.Tx) error {
		// Check if entity with ID exists
		v, err := GetByID(e.ID, EntityTypeBucketName, ser.DB)
		if err != nil {
			return err
		}
		var o EntityType
		err = json.Unmarshal(v, &o)
		if err != nil {
			return err
		}

		// Get bucket, exit if error
		b, err := Bucket(EntityTypeBucketName, tx)
		if err != nil {
			return err
		}

		// Replace properties of new with certain
		// ones of old
		e.Version = o.Version + 1

		// Save Entity
		buf, err := json.Marshal(e)
		if err != nil {
			return err
		}

		return b.Put(itob(e.ID), buf)
	})
}

// Delete removes the persisted instance
// of EntityType with the given ID
func (ser *EntityTypeService) Delete(ID int) (e EntityType, err error) {
	err = ser.DB.Update(func(tx *bolt.Tx) error {
		// Get bucket, exit if error
		b, err := Bucket(EntityTypeBucketName, tx)
		if err != nil {
			return err
		}

		// Store existing to return
		e.ID = ID
		err = ser.GetByID(&e)
		if err != nil {
			return err
		}

		// Delete
		return b.Delete(itob(ID))
	})
	return
}
