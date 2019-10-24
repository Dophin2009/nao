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
