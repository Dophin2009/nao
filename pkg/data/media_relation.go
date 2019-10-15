package data

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

// MediaRelation represents a relationship between single
// instances of Media and Producer
type MediaRelation struct {
	ID           int
	OwnerID      int
	RelatedID    int
	Relationship string
	Version      int
}

// Identifier returns the ID of the MediaRelation
func (mr *MediaRelation) Identifier() int {
	return mr.ID
}

// SetIdentifier sets the ID of the MediaRelation
func (mr *MediaRelation) SetIdentifier(ID int) {
	mr.ID = ID
}

// Ver returns the verison of the MediaRelation
func (mr *MediaRelation) Ver() int {
	return mr.Version
}

// UpdateVer increments the version of the
// Character by one
func (mr *MediaRelation) UpdateVer() {
	mr.Version++
}

// Validate returns an error if the MediaRelation is
// not valid for the database
func (mr *MediaRelation) Validate(tx *bolt.Tx) (err error) {
	// Get Media bucket, exit if error
	mb, err := bucket(mediaBucketName, tx)
	if err != nil {
		return err
	}

	// Check if owning Media with ID specified in new MediaRelation exists
	_, err = get(mr.OwnerID, mb)
	if err != nil {
		return err
	}

	// Check if related Media with ID specified in new MediaRelation exists
	_, err = get(mr.RelatedID, mb)
	if err != nil {
		return err
	}

	return nil
}

const mediaRelationBucketName = "MediaRelation"

// MediaRelationGet retrieves a single instance of MediaRelation with
// the given ID
func MediaRelationGet(ID int, db *bolt.DB) (mr MediaRelation, err error) {
	err = getByID(ID, &mr, mediaRelationBucketName, db)
	return
}

// MediaRelationGetAll retrieves all persisted MediaRelation values
func MediaRelationGetAll(db *bolt.DB) (list []MediaRelation, err error) {
	return MediaRelationGetFilter(db, func(mr *MediaRelation) bool { return true })
}

// MediaRelationGetByOwner retrieves a list of instances of MediaRelation
// with the given owning Media ID
func MediaRelationGetByOwner(mID int, db *bolt.DB) (list []MediaRelation, err error) {
	return MediaRelationGetFilter(db, func(mr *MediaRelation) bool {
		return mr.OwnerID == mID
	})
}

// MediaRelationGetByRelated retrieves a list of instances of MediaRelation
// with the given related Media ID
func MediaRelationGetByRelated(mID int, db *bolt.DB) (list []MediaRelation, err error) {
	return MediaRelationGetFilter(db, func(mr *MediaRelation) bool {
		return mr.RelatedID == mID
	})
}

// MediaRelationGetByRelationship retrieves a list of instances of Media Relation
// with the given relationship
func MediaRelationGetByRelationship(relationship string, db *bolt.DB) (list []MediaRelation, err error) {
	return MediaRelationGetFilter(db, func(mr *MediaRelation) bool {
		return mr.Relationship == relationship
	})
}

// MediaRelationGetFilter retrieves all persisted MediaRelation values
func MediaRelationGetFilter(db *bolt.DB, filter func(mr *MediaRelation) bool) (list []MediaRelation, err error) {
	ilist, err := getFilter(&MediaRelation{}, func(entity Idenitifiable) (bool, error) {
		mr, ok := entity.(*MediaRelation)
		if !ok {
			return false, fmt.Errorf("type assertion failed: entity is not a MediaRelation")
		}
		return filter(mr), nil
	}, mediaRelationBucketName, db)

	list = make([]MediaRelation, len(ilist))
	for i, mr := range ilist {
		list[i] = *mr.(*MediaRelation)
	}

	return
}

// MediaRelationCreate persists a new instance of MediaRelation
func MediaRelationCreate(mr *MediaRelation, db *bolt.DB) error {
	return create(mr, mediaRelationBucketName, db)
}

// MediaRelationUpdate updates the properties of an existing
// persisted Producer instance
func MediaRelationUpdate(mr *MediaRelation, db *bolt.DB) error {
	return update(mr, mediaRelationBucketName, db)
}
