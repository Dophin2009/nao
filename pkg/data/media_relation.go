package data

import (
	bolt "go.etcd.io/bbolt"
)

// MediaRelation represents a relationship between single
// instances of Media and Producer.
type MediaRelation struct {
	ID           int
	OwnerID      int
	RelatedID    int
	Relationship string
	Version      int
}

// Validate returns an error if the MediaRelation is
// not valid for the database.
func (ser *MediaRelationService) Validate(e *MediaRelation) (err error) {
	return ser.DB.View(func(tx *bolt.Tx) error {
		// Get Media bucket, exit if error
		mb, err := Bucket(MediaBucketName, tx)
		if err != nil {
			return err
		}

		// Check if owning Media with ID specified in new MediaRelation exists
		_, err = get(e.OwnerID, mb)
		if err != nil {
			return err
		}

		// Check if related Media with ID specified in new MediaRelation exists
		_, err = get(e.RelatedID, mb)
		if err != nil {
			return err
		}

		return nil
	})
}

// GetByOwner retrieves a list of instances of MediaRelation
// with the given owning Media ID.
func (ser *MediaRelationService) GetByOwner(mID int, db *bolt.DB) (list []MediaRelation, err error) {
	return ser.GetFilter(func(mr *MediaRelation) bool {
		return mr.OwnerID == mID
	})
}

// GetByRelated retrieves a list of instances of MediaRelation
// with the given related Media ID.
func (ser *MediaRelationService) GetByRelated(mID int, db *bolt.DB) (list []MediaRelation, err error) {
	return ser.GetFilter(func(mr *MediaRelation) bool {
		return mr.RelatedID == mID
	})
}

// GetByRelationship retrieves a list of instances of Media Relation
// with the given relationship.
func (ser *MediaRelationService) GetByRelationship(relationship string, db *bolt.DB) (list []MediaRelation, err error) {
	return ser.GetFilter(func(mr *MediaRelation) bool {
		return mr.Relationship == relationship
	})
}
