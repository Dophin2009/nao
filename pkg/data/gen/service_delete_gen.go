package data

import (
	"github.com/cheekybits/genny/generic"
	bolt "go.etcd.io/bbolt"
)

//go:generate genny -in=service_delete_gen.go -out=service_delete.gen.go gen "EntityType=Media,Episode,Character,Genre,Producer,Person,User,MediaRelation,MediaCharacter,MediaGenre,MediaProducer,UserMedia,UserMediaList"

// EntityType is a generic placeholder for all entity types;
// it is assumed that EntityType structs have an ID and Version,
// both of which should be of int type.
type EntityType generic.Type

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
