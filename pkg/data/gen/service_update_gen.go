package data

import (
	"github.com/cheekybits/genny/generic"
	json "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
)

//go:generate genny -in=service_update_gen.go -out=service_update.gen.go gen "EntityType=Media,Episode,Character,Genre,Producer,Person,User,MediaRelation,MediaCharacter,MediaGenre,MediaProducer,UserMedia,UserMediaList"

// EntityType is a generic placeholder for all entity types;
// it is assumed that EntityType structs have an ID and Version,
// both of which should be of int type.
type EntityType generic.Type

// Update replaces the persisted instance
// of EntityType with the given EntityType of
// the same ID.
func (ser *EntityTypeService) Update(e *EntityType) (err error) {
	err = ser.Clean(e)
	if err != nil {
		return err
	}

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
		ser.persistOldProperties(&o, e)

		// Save Entity
		buf, err := json.Marshal(e)
		if err != nil {
			return err
		}

		return b.Put(itob(e.ID), buf)
	})
}
