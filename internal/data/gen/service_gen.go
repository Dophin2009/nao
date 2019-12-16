package data

import (
	"github.com/cheekybits/genny/generic"

	bolt "go.etcd.io/bbolt"
)

//go:generate genny -in=service_gen.go -out=service.gen.go gen "EntityType=Media,Episode,Character,Genre,Producer,Person,User,MediaRelation,MediaCharacter,MediaGenre,MediaProducer,UserMedia,UserMediaList"

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
