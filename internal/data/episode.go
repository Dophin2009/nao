package data

import (
	"errors"
	"fmt"
	"time"

	json "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
)

// TODO: User rating/comments/etc. of Episodes

// Episode represents a single episode or chapter for some media.
type Episode struct {
	Titles   []Title
	Synopses []Title
	Date     *time.Time
	Duration *int
	Filler   bool
	Recap    bool
	Meta     ModelMetadata
}

// Metadata returns Meta.
func (ep *Episode) Metadata() *ModelMetadata {
	return &ep.Meta
}

// EpisodeSet is an ordered list of episodes.
type EpisodeSet struct {
	MediaID      int
	Descriptions []Title
	Episodes     []int
	Meta         ModelMetadata
}

// Metadata returns the Meta.
func (set *EpisodeSet) Metadata() *ModelMetadata {
	return &set.Meta
}

// EpisodeBucket is the name of the database bucket for Episodes.
const EpisodeBucket = "Episode"

// EpisodeService performs operations on Episodes.
type EpisodeService struct{}

// Create persists the given Episode.
func (ser *EpisodeService) Create(ep *Episode, tx Tx) (int, error) {
	return tx.Database().Create(ep, ser, tx)
}

// Update replaces the value of the Episode with the given ID.
func (ser *EpisodeService) Update(ep *Episode, tx Tx) error {
	return tx.Database().Update(ep, ser, tx)
}

// Delete deletes the Episode with the given ID.
func (ser *EpisodeService) Delete(id int, tx Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// GetAll retrieves all persisted values of Episode.
func (ser *EpisodeService) GetAll(first *int, skip *int, tx Tx) ([]*Episode, error) {
	vlist, err := tx.Database().GetAll(first, skip, ser, tx)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to Episodes: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of Episode that pass the filter.
func (ser *EpisodeService) GetFilter(
	first *int, skip *int, tx Tx, keep func(ep *Episode) bool,
) ([]*Episode, error) {
	vlist, err := tx.Database().GetFilter(first, skip, ser, tx,
		func(m Model) bool {
			ep, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(ep)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to Episodes: %w", err)
	}
	return list, nil
}

// GetMultiple retrieves the persisted Episode values specified by the given
// IDs that pass the filter.
func (ser *EpisodeService) GetMultiple(
	ids []int, first *int, skip *int, tx Tx, keep func(ep *Episode) bool,
) ([]*Episode, error) {
	vlist, err := tx.Database().GetMultiple(ids, first, skip, ser, tx,
		func(m Model) bool {
			ep, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(ep)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to Episodes: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted Episode with the given ID.
func (ser *EpisodeService) GetByID(id int, tx Tx) (*Episode, error) {
	m, err := tx.Database().GetByID(id, ser, tx)
	if err != nil {
		return nil, err
	}

	ep, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return ep, nil
}

// Bucket returns the name of the bucket for Episode.
func (ser *EpisodeService) Bucket() string {
	return EpisodeBucket
}

// Clean cleans the given Episode for storage.
func (ser *EpisodeService) Clean(m Model, _ Tx) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Validate returns an error if the Episode is not valid for the database.
func (ser *EpisodeService) Validate(m Model, _ Tx) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Initialize sets initial values for some properties.
func (ser *EpisodeService) Initialize(_ Model, _ Tx) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing Episode in
// updates.
func (ser *EpisodeService) PersistOldProperties(_ Model, _ Model, _ Tx) error {
	return nil
}

// Marshal transforms the given Episode into JSON.
func (ser *EpisodeService) Marshal(m Model) ([]byte, error) {
	ep, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	v, err := json.Marshal(ep)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONMarshal, err)
	}

	return v, nil
}

// Unmarshal parses the given JSON into Episode.
func (ser *EpisodeService) Unmarshal(buf []byte) (Model, error) {
	var ep Episode
	err := json.Unmarshal(buf, &ep)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &ep, nil
}

// AssertType exposes the Model as an Episode.
func (ser *EpisodeService) AssertType(m Model) (*Episode, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	ep, ok := m.(*Episode)
	if !ok {
		return nil, fmt.Errorf("model: %w", errors.New("not of Episode type"))
	}
	return ep, nil
}

// mapfromModel returns a list of Episode type asserted from the given list of
// Model.
func (ser *EpisodeService) mapFromModel(vlist []Model) ([]*Episode, error) {
	list := make([]*Episode, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}

// EpisodeSetBucket is the name of the database bucket for EpisodeSets.
const EpisodeSetBucket = "EpisodeSet"

// EpisodeSetService performs operations on EpisodeSets.
type EpisodeSetService struct {
	DB *bolt.DB
	Service
}

// Create persists the given EpisodeSet.
func (ser *EpisodeSetService) Create(set *EpisodeSet, tx Tx) (int, error) {
	return tx.Database().Create(set, ser, tx)
}

// Update replaces the value of the EpisodeSet with the given ID.
func (ser *EpisodeSetService) Update(set *EpisodeSet, tx Tx) error {
	return tx.Database().Update(set, ser, tx)
}

// Delete deletes the EpisodeSet with the given ID.
func (ser *EpisodeSetService) Delete(id int, tx Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// GetAll retrieves all persisted values of EpisodeSet.
func (ser *EpisodeSetService) GetAll(first *int, skip *int, tx Tx) ([]*EpisodeSet, error) {
	vlist, err := tx.Database().GetAll(first, skip, ser, tx)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to EpisodeSets: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of EpisodeSet that pass the filter.
func (ser *EpisodeSetService) GetFilter(
	first *int, skip *int, tx Tx, keep func(*EpisodeSet) bool,
) ([]*EpisodeSet, error) {
	vlist, err := tx.Database().GetFilter(first, skip, ser, tx,
		func(m Model) bool {
			set, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(set)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to EpisodesSets: %w", err)
	}
	return list, nil
}

// GetMultiple retrieves the persisted EpisodeSet values specified by the given
// IDs that pass the filter.
func (ser *EpisodeSetService) GetMultiple(
	ids []int, first *int, skip *int, tx Tx, keep func(set *EpisodeSet) bool,
) ([]*EpisodeSet, error) {
	vlist, err := tx.Database().GetMultiple(ids, first, skip, ser, tx,
		func(m Model) bool {
			set, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(set)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to EpisodeSets: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted EpisodeSet with the given ID.
func (ser *EpisodeSetService) GetByID(id int, tx Tx) (*EpisodeSet, error) {
	m, err := tx.Database().GetByID(id, ser, tx)
	if err != nil {
		return nil, err
	}

	set, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return set, nil
}

// GetByMedia retrieves a list of instances of EpisodeSet with the given Media
// ID.
func (ser *EpisodeSetService) GetByMedia(
	mID int, first *int, skip *int, tx Tx,
) ([]*EpisodeSet, error) {
	return ser.GetFilter(first, skip, tx, func(set *EpisodeSet) bool {
		return set.MediaID == mID
	})
}

// Bucket returns the name of the bucket for EpisodeSet.
func (ser *EpisodeSetService) Bucket() string {
	return EpisodeBucket
}

// Clean cleans the given EpisodeSet for storage.
func (ser *EpisodeSetService) Clean(m Model, _ Tx) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Validate returns an error if the Episodeset is not valid for the database.
func (ser *EpisodeSetService) Validate(m Model, _ Tx) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Initialize sets initial values for some properties.
func (ser *EpisodeSetService) Initialize(_ Model, _ Tx) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing EpisodeSet
// in updates.
func (ser *EpisodeSetService) PersistOldProperties(_ Model, _ Model, _ Tx) error {
	return nil
}

// Marshal transforms the given EpisodeSet into JSON.
func (ser *EpisodeSetService) Marshal(m Model) ([]byte, error) {
	set, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	v, err := json.Marshal(set)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONMarshal, err)
	}

	return v, nil
}

// Unmarshal parses the given JSON into EpisodeSet.
func (ser *EpisodeSetService) Unmarshal(buf []byte) (Model, error) {
	var set EpisodeSet
	err := json.Unmarshal(buf, &set)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &set, nil
}

// AssertType exposes the Model as an EpisodeSet.
func (ser *EpisodeSetService) AssertType(m Model) (*EpisodeSet, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	set, ok := m.(*EpisodeSet)
	if !ok {
		return nil, fmt.Errorf("model: %w", errors.New("not of EpisodeSet type"))
	}
	return set, nil
}

// mapfromModel returns a list of EpisodeSet type asserted from the given list
// of Model.
func (ser *EpisodeSetService) mapFromModel(vlist []Model) ([]*EpisodeSet, error) {
	list := make([]*EpisodeSet, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}
