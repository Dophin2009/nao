package data

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Dophin2009/nao/pkg/models"
	"github.com/Dophin2009/nao/pkg/db"
	json "github.com/json-iterator/go"
)

// MediaProducerService performs operations on MediaProducer.
type MediaProducerService struct {
	MediaService    *MediaService
	ProducerService *ProducerService
	Hooks           db.PersistHooks
}

// NewMediaProducer retursn a MediaProducer.
func NewMediaProducer(hooks db.PersistHooks, mediaService *MediaService,
	producerService *ProducerService) *MediaProducerService {
	// Initialize MediaProducerService
	mediaProducerService := &MediaProducerService{
		MediaService:    mediaService,
		ProducerService: producerService,
		Hooks:           hooks,
	}

	// Add hook to delete MediaProducer on Media deletion
	deleteMediaProducerOnDeleteMedia := func(mdm db.Model, _ db.Service, tx db.Tx) error {
		mID := mdm.Metadata().ID
		err := mediaProducerService.DeleteByMedia(mID, tx)
		if err != nil {
			return fmt.Errorf("failed to delete MediaProducer by Media ID %d: %w",
				mID, err)
		}
		return nil
	}
	mdSerHooks := mediaService.PersistHooks()
	mdSerHooks.PreDeleteHooks =
		append(mdSerHooks.PreDeleteHooks, deleteMediaProducerOnDeleteMedia)

	// Add hook to delete MediaProducer on Producer deletion
	deleteMediaProducerOnDeleteProducer := func(pd db.Model, _ db.Service, tx db.Tx) error {
		pID := pd.Metadata().ID
		err := mediaProducerService.DeleteByProducer(pID, tx)
		if err != nil {
			return fmt.Errorf("failed to delete MediaProducer by Producer ID %d: %w",
				pID, err)
		}
		return nil
	}
	pSerHooks := producerService.PersistHooks()
	pSerHooks.PreDeleteHooks =
		append(pSerHooks.PreDeleteHooks, deleteMediaProducerOnDeleteProducer)

	return mediaProducerService
}

// Create persists the given MediaProducer.
func (ser *MediaProducerService) Create(mp *models.MediaProducer, tx db.Tx) (int, error) {
	return tx.Database().Create(mp, ser, tx)
}

// Update rmplaces the value of the MediaProducer with the
// given ID.
func (ser *MediaProducerService) Update(mp *models.MediaProducer, tx db.Tx) error {
	return tx.Database().Update(mp, ser, tx)
}

// Delete deletes the MediaProducer with the given ID.
func (ser *MediaProducerService) Delete(id int, tx db.Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// DeleteByMedia deletes the MediaProducers with the given Media ID.
func (ser *MediaProducerService) DeleteByMedia(mID int, tx db.Tx) error {
	return tx.Database().DeleteFilter(ser, tx, func(m db.Model) bool {
		mp, err := ser.AssertType(m)
		if err != nil {
			return false
		}

		return mp.MediaID == mID
	})
}

// DeleteByProducer deletes the MediaProducers with the given Producer ID.
func (ser *MediaProducerService) DeleteByProducer(pID int, tx db.Tx) error {
	return tx.Database().DeleteFilter(ser, tx, func(m db.Model) bool {
		mp, err := ser.AssertType(m)
		if err != nil {
			return false
		}

		return mp.ProducerID == pID
	})
}

// GetAll retrieves all persisted values of MediaProducer.
func (ser *MediaProducerService) GetAll(
	first *int, skip *int, tx db.Tx) ([]*models.MediaProducer, error) {
	vlist, err := tx.Database().GetAll(first, skip, ser, tx)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to MediaProducer: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of MediaProducer that pass the
// filter.
func (ser *MediaProducerService) GetFilter(
	first *int, skip *int, tx db.Tx, keep func(mp *models.MediaProducer) bool,
) ([]*models.MediaProducer, error) {
	vlist, err := tx.Database().GetFilter(first, skip, ser, tx,
		func(m db.Model) bool {
			mp, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(mp)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to MediaProducer: %w", err)
	}
	return list, nil
}

// GetMultiple retrieves the persisted MediaProducer values specified by the
// given IDs that pass the filter.
func (ser *MediaProducerService) GetMultiple(
	ids []int, tx db.Tx, keep func(mp *models.MediaProducer) bool,
) ([]*models.MediaProducer, error) {
	vlist, err := tx.Database().GetMultiple(ids, ser, tx,
		func(m db.Model) bool {
			mp, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(mp)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to MediaProducers: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted MediaProducer with the given ID.
func (ser *MediaProducerService) GetByID(id int, tx db.Tx) (*models.MediaProducer, error) {
	m, err := tx.Database().GetByID(id, ser, tx)
	if err != nil {
		return nil, err
	}

	mp, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return mp, nil
}

// GetByMedia retrieves a list of instances of MediaProducer with the given
// Media ID.
func (ser *MediaProducerService) GetByMedia(
	mID int, first *int, skip *int, tx db.Tx,
) ([]*models.MediaProducer, error) {
	return ser.GetFilter(first, skip, tx, func(mp *models.MediaProducer) bool {
		return mp.MediaID == mID
	})
}

// GetByProducer retrieves a list of instances of MediaProducer with the given
// Producer ID.
func (ser *MediaProducerService) GetByProducer(
	pID int, first *int, skip *int, tx db.Tx,
) ([]*models.MediaProducer, error) {
	return ser.GetFilter(first, skip, tx, func(mp *models.MediaProducer) bool {
		return mp.ProducerID == pID
	})
}

// Bucket returns the name of the bucket for MediaProducer.
func (ser *MediaProducerService) Bucket() string {
	return "MediaProducer"
}

// Clean cleans the given MediaProducer for storage.
func (ser *MediaProducerService) Clean(m db.Model, _ db.Tx) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	e.Role = strings.Trim(e.Role, " ")
	return nil
}

// Validate returns an error if the MediaProducer is not valid for the
// database.
func (ser *MediaProducerService) Validate(m db.Model, tx db.Tx) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	db := tx.Database()

	// Check if Media with ID specified in new MediaProducer exists
	_, err = db.GetRawByID(e.MediaID, ser, tx)
	if err != nil {
		return fmt.Errorf("failed to get Media with ID %d: %w", e.MediaID, err)
	}

	// Check if Producer with ID specified in new MediaProducer exists
	_, err = db.GetRawByID(e.ProducerID, ser, tx)
	if err != nil {
		return fmt.Errorf("failed to get Producer with ID %d: %w", e.ProducerID, err)
	}

	return nil
}

// Initialize sets initial values for some properties.
func (ser *MediaProducerService) Initialize(_ db.Model, _ db.Tx) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing
// MediaProducer in updates.
func (ser *MediaProducerService) PersistOldProperties(_ db.Model, _ db.Model, _ db.Tx) error {
	return nil
}

// PersistHooks returns the persistence hook functions.
func (ser *MediaProducerService) PersistHooks() *db.PersistHooks {
	return &ser.Hooks
}

// Marshal transforms the given MediaProducer into JSON.
func (ser *MediaProducerService) Marshal(m db.Model) ([]byte, error) {
	mp, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	v, err := json.Marshal(mp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONMarshal, err)
	}

	return v, nil
}

// Unmarshal parses the given JSON into MediaProducer.
func (ser *MediaProducerService) Unmarshal(buf []byte) (db.Model, error) {
	var mp models.MediaProducer
	err := json.Unmarshal(buf, &mp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &mp, nil
}

// AssertType exposes the given db.Model as a MediaProducer.
func (ser *MediaProducerService) AssertType(m db.Model) (*models.MediaProducer, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	mp, ok := m.(*models.MediaProducer)
	if !ok {
		return nil, fmt.Errorf("model: %w", errors.New("not of MediaProducer type"))
	}
	return mp, nil
}

// mapfromModel returns a list of MediaProducer type asserted from the given
// list of db.Model.
func (ser *MediaProducerService) mapFromModel(vlist []db.Model) ([]*models.MediaProducer, error) {
	list := make([]*models.MediaProducer, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}
