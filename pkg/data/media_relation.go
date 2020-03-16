package data

import (
	"errors"
	"fmt"
	"strings"

	json "github.com/json-iterator/go"
	"github.com/Dophin2009/nao/pkg/db"
)

// MediaRelation represents a relationship between single instances of Media
// and Producer.
type MediaRelation struct {
	OwnerID      int
	RelatedID    int
	Relationship string
	Meta         db.ModelMetadata
}

// Metadata returns Meta.
func (mr *MediaRelation) Metadata() *db.ModelMetadata {
	return &mr.Meta
}

// MediaRelationService performs operations on MediaRelation.
type MediaRelationService struct {
	MediaService *MediaService
	Hooks        db.PersistHooks
}

// NewMediaRelationService returns a MediaRelationService.
func NewMediaRelationService(hooks db.PersistHooks, mediaService *MediaService) *MediaRelationService {
	mediaRelationService := &MediaRelationService{
		MediaService: mediaService,
		Hooks:        hooks,
	}

	deleteMediaRelationOnDeleteMedia := func(mdm db.Model, _ db.Service, tx db.Tx) error {
		mID := mdm.Metadata().ID
		err := mediaRelationService.DeleteByOwner(mID, tx)
		if err != nil {
			return fmt.Errorf("failed to delete MediaRelation by Owner ID %d: %w",
				mID, err)
		}
		err = mediaRelationService.DeleteByRelated(mID, tx)
		if err != nil {
			return fmt.Errorf("failed to delete MediaRelation by Related ID %d: %w",
				mID, err)
		}
		return nil
	}
	mdSerHooks := mediaService.PersistHooks()
	mdSerHooks.PreDeleteHooks =
		append(mdSerHooks.PreDeleteHooks, deleteMediaRelationOnDeleteMedia)

	return mediaRelationService
}

// Create persists the given MediaRelation.
func (ser *MediaRelationService) Create(mr *MediaRelation, tx db.Tx) (int, error) {
	return tx.Database().Create(mr, ser, tx)
}

// Update rmrlaces the value of the MediaRelation with the given ID.
func (ser *MediaRelationService) Update(mr *MediaRelation, tx db.Tx) error {
	return tx.Database().Update(mr, ser, tx)
}

// Delete deletes the MediaRelation with the given ID.
func (ser *MediaRelationService) Delete(id int, tx db.Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// DeleteByOwner deletes the MediaRelation with the given Owner ID.
func (ser *MediaRelationService) DeleteByOwner(mID int, tx db.Tx) error {
	return tx.Database().DeleteFilter(ser, tx, func(m db.Model) bool {
		mr, err := ser.AssertType(m)
		if err != nil {
			return false
		}
		return mr.OwnerID == mID
	})
}

// DeleteByRelated deletes the MediaRelation with the given Related ID.
func (ser *MediaRelationService) DeleteByRelated(mID int, tx db.Tx) error {
	return tx.Database().DeleteFilter(ser, tx, func(m db.Model) bool {
		mr, err := ser.AssertType(m)
		if err != nil {
			return false
		}
		return mr.RelatedID == mID
	})
}

// GetAll retrieves all persisted values of MediaRelation.
func (ser *MediaRelationService) GetAll(first *int, skip *int, tx db.Tx) ([]*MediaRelation, error) {
	vlist, err := tx.Database().GetAll(first, skip, ser, tx)
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
	first *int, skip *int, tx db.Tx, keep func(mr *MediaRelation) bool,
) ([]*MediaRelation, error) {
	vlist, err := tx.Database().GetFilter(first, skip, ser, tx, func(m db.Model) bool {
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

// GetMultiple retrieves the persisted MediaRelation values specified by the
// given IDs that pass the filter.
func (ser *MediaRelationService) GetMultiple(
	ids []int, tx db.Tx, keep func(mr *MediaRelation) bool,
) ([]*MediaRelation, error) {
	vlist, err := tx.Database().GetMultiple(ids, ser, tx,
		func(m db.Model) bool {
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
func (ser *MediaRelationService) GetByID(id int, tx db.Tx) (*MediaRelation, error) {
	m, err := tx.Database().GetByID(id, ser, tx)
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
	mID int, first *int, skip *int, tx db.Tx,
) ([]*MediaRelation, error) {
	return ser.GetFilter(first, skip, tx, func(mr *MediaRelation) bool {
		return mr.OwnerID == mID
	})
}

// GetByRelated retrieves a list of instances of MediaRelation with the given
// related Media ID.
func (ser *MediaRelationService) GetByRelated(
	mID int, first *int, skip *int, tx db.Tx,
) ([]*MediaRelation, error) {
	return ser.GetFilter(first, skip, tx, func(mr *MediaRelation) bool {
		return mr.RelatedID == mID
	})
}

// GetByRelationship retrieves a list of instances of Media Relation with the
// given relationship.
func (ser *MediaRelationService) GetByRelationship(
	relationship string, first *int, skip *int, tx db.Tx,
) ([]*MediaRelation, error) {
	return ser.GetFilter(first, skip, tx, func(mr *MediaRelation) bool {
		return mr.Relationship == relationship
	})
}

// Bucket returns the name of the bucket for MediaRelation.
func (ser *MediaRelationService) Bucket() string {
	return "MediaRelation"
}

// Clean cleans the given MediaRelation for storage.
func (ser *MediaRelationService) Clean(m db.Model, _ db.Tx) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	e.Relationship = strings.Trim(e.Relationship, " ")
	return nil
}

// Validate returns an error if the MediaRelation is not valid for the
// database.
func (ser *MediaRelationService) Validate(m db.Model, tx db.Tx) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	db := tx.Database()

	// Check if owning Media with ID specified in new MediaRelation exists
	_, err = db.GetRawByID(e.OwnerID, ser.MediaService, tx)
	if err != nil {
		return fmt.Errorf("failed to get Media with ID %d: %w", e.OwnerID, err)
	}

	// Check if related Media with ID specified in new MediaRelation exists
	_, err = db.GetRawByID(e.RelatedID, ser.MediaService, tx)
	if err != nil {
		return fmt.Errorf("failed to get Media with ID %d: %w", e.RelatedID, err)
	}

	return nil
}

// Initialize sets initial values for some properties.
func (ser *MediaRelationService) Initialize(_ db.Model, _ db.Tx) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing
// MediaRelation in updates.
func (ser *MediaRelationService) PersistOldProperties(_ db.Model, _ db.Model, _ db.Tx) error {
	return nil
}

// PersistHooks returns the persistence hook functions.
func (ser *MediaRelationService) PersistHooks() *db.PersistHooks {
	return &ser.Hooks
}

// Marshal transforms the given MediaRelation into JSON.
func (ser *MediaRelationService) Marshal(m db.Model) ([]byte, error) {
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
func (ser *MediaRelationService) Unmarshal(buf []byte) (db.Model, error) {
	var mr MediaRelation
	err := json.Unmarshal(buf, &mr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &mr, nil
}

// AssertType exposes the given Model as a MediaRelation.
func (ser *MediaRelationService) AssertType(m db.Model) (*MediaRelation, error) {
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
func (ser *MediaRelationService) mapFromModel(vlist []db.Model) ([]*MediaRelation, error) {
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
