package data

import (
	"encoding/json"

	"github.com/cheekybits/genny/generic"
)

//go:generate genny -in=$GOFILE -out=gen_$GOFILE gen "EntityType=Media,Episode,Character,Genre,Producer,Person,User,MediaRelation,MediaCharacter,MediaGenre,MediaProducer,UserMedia,UserMediaList"

// EntityType is a generic placeholder for all entity types;
// it is assumed that EntityType structs have an ID and Version,
// both of which should be of int type.
type EntityType generic.Type

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
