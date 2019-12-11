package data

import (
	"github.com/cheekybits/genny/generic"
	json "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
)

//go:generate genny -in=$GOFILE -out=gen_$GOFILE gen "EntityType=Media,Episode,Character,Genre,Producer,Person,MediaRelation,MediaCharacter,MediaGenre,MediaProducer,UserMedia,UserMediaList"

// EntityType is a generic placeholder for all entity types;
// it is assumed that EntityType structs have an ID and Version,
// both of which should be of int type.
type EntityType generic.Type

// Create persists a new instance of EntityType to the database.
func (ser *EntityTypeService) Create(e *EntityType) (err error) {
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
