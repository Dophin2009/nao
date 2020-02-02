package data

import (
	"errors"
	"fmt"
	"strings"

	json "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
)

// MediaRelation represents a relationship between single instances of Media
// and Producer.
type MediaRelation struct {
	ID           int
	OwnerID      int
	RelatedID    int
	Relationship string
	Version      int
}

// Iden returns the ID.
func (mr *MediaRelation) Iden() int {
	return mr.ID
}

// MediaRelationBucket is the name of the database bucket for MediaRelation.
const MediaRelationBucket = "MediaRelation"

// MediaRelationService performs operations on MediaRelation.
type MediaRelationService struct {
	DB *bolt.DB
	Service
}

// Create persists the given MediaRelation.
func (ser *MediaRelationService) Create(mr *MediaRelation) error {
	return Create(mr, ser)
}

// Update rmrlaces the value of the MediaRelation with the given ID.
func (ser *MediaRelationService) Update(mr *MediaRelation) error {
	return Update(mr, ser)
}

// Delete deletes the MediaRelation with the given ID.
func (ser *MediaRelationService) Delete(id int) error {
	return Delete(id, ser)
}

// GetAll retrieves all persisted values of MediaRelation.
func (ser *MediaRelationService) GetAll(first *int, skip *int) ([]*MediaRelation, error) {
	vlist, err := GetAll(ser, first, skip)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to MediaRelations: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of MediaRelation that pass the
// filter.
func (ser *MediaRelationService) GetFilter(
	first *int, skip *int, keep func(mr *MediaRelation) bool,
) ([]*MediaRelation, error) {
	vlist, err := GetFilter(ser, first, skip, func(m Model) bool {
		mr, err := ser.AssertType(m)
		if err != nil {
			return false
		}
		return keep(mr)
	})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to MediaRelations: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted MediaRelation with the given ID.
func (ser *MediaRelationService) GetByID(id int) (*MediaRelation, error) {
	m, err := GetByID(id, ser)
	if err != nil {
		return nil, err
	}

	mr, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return mr, nil
}

// GetByOwner retrieves a list of instances of MediaRelation with the given
// owning Media ID.
func (ser *MediaRelationService) GetByOwner(
	mID int, first *int, skip *int,
) ([]*MediaRelation, error) {
	return ser.GetFilter(first, skip, func(mr *MediaRelation) bool {
		return mr.OwnerID == mID
	})
}

// GetByRelated retrieves a list of instances of MediaRelation with the given
// related Media ID.
func (ser *MediaRelationService) GetByRelated(
	mID int, first *int, skip *int,
) ([]*MediaRelation, error) {
	return ser.GetFilter(first, skip, func(mr *MediaRelation) bool {
		return mr.RelatedID == mID
	})
}

// GetByRelationship retrieves a list of instances of Media Relation with the
// given relationship.
func (ser *MediaRelationService) GetByRelationship(
	relationship string, first *int, skip *int,
) ([]*MediaRelation, error) {
	return ser.GetFilter(first, skip, func(mr *MediaRelation) bool {
		return mr.Relationship == relationship
	})
}

// Database returns the database reference.
func (ser *MediaRelationService) Database() *bolt.DB {
	return ser.DB
}

// Bucket returns the name of the bucket for MediaRelation.
func (ser *MediaRelationService) Bucket() string {
	return MediaRelationBucket
}

// Clean cleans the given MediaRelation for storage.
func (ser *MediaRelationService) Clean(m Model) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	e.Relationship = strings.Trim(e.Relationship, " ")
	return nil
}

// Validate returns an error if the MediaRelation is not valid for the
// database.
func (ser *MediaRelationService) Validate(m Model) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	return ser.DB.View(func(tx *bolt.Tx) error {
		// Get Media bucket, exit if error
		mb, err := Bucket(MediaBucket, tx)
		if err != nil {
			return fmt.Errorf("%s %q: %w", errmsgBucketOpen, MediaBucket, err)
		}

		// Check if owning Media with ID specified in new MediaRelation exists
		_, err = get(e.OwnerID, mb)
		if err != nil {
			return fmt.Errorf("failed to get Media with ID %d: %w", e.OwnerID, err)
		}

		// Check if related Media with ID specified in new MediaRelation exists
		_, err = get(e.RelatedID, mb)
		if err != nil {
			return fmt.Errorf("failed to get Media with ID %d: %w", e.RelatedID, err)
		}

		return nil
	})
}

// Initialize sets initial values for some properties.
func (ser *MediaRelationService) Initialize(m Model, id int) error {
	mr, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	mr.ID = id
	mr.Version = 0
	return nil
}

// PersistOldProperties maintains certain properties of the existing
// MediaRelation in updates.
func (ser *MediaRelationService) PersistOldProperties(n Model, o Model) error {
	nm, err := ser.AssertType(n)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	om, err := ser.AssertType(o)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	nm.Version = om.Version + 1
	return nil
}

// Marshal transforms the given MediaRelation into JSON.
func (ser *MediaRelationService) Marshal(m Model) ([]byte, error) {
	mr, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	v, err := json.Marshal(mr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONMarshal, err)
	}

	return v, nil
}

// Unmarshal parses the given JSON into MediaRelation.
func (ser *MediaRelationService) Unmarshal(buf []byte) (Model, error) {
	var mr MediaRelation
	err := json.Unmarshal(buf, &mr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &mr, nil
}

// AssertType exposes the given Model as a MediaRelation.
func (ser *MediaRelationService) AssertType(m Model) (*MediaRelation, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	mr, ok := m.(*MediaRelation)
	if !ok {
		return nil, fmt.Errorf("model: %w", errors.New("not of MediaRelation type"))
	}
	return mr, nil
}

// mapfromModel returns a list of MediaRelation type asserted from the given
// list of Model.
func (ser *MediaRelationService) mapFromModel(vlist []Model) ([]*MediaRelation, error) {
	list := make([]*MediaRelation, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}
