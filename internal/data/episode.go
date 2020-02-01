package data

import (
	"errors"
	"fmt"
	"time"

	json "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
)

// TODO: User rating/comments/etc. of Episodes

// Episode represents a single episode or chapter
// for some media.
type Episode struct {
	ID       int
	Titles   []Title
	Synopses []Title
	Date     *time.Time
	Duration *int
	Filler   bool
	Recap    bool
	Version  int
	Model
}

// Iden returns the ID.
func (ep *Episode) Iden() int {
	return ep.ID
}

// EpisodeSet is an ordered list of episodes.
type EpisodeSet struct {
	ID           int
	MediaID      int
	Descriptions []Title
	Episodes     []int
	Version      int
	Model
}

// Iden returns the ID.
func (set *EpisodeSet) Iden() int {
	return set.ID
}

// EpisodeBucket is the name of the database bucket for
// Episodes.
const EpisodeBucket = "Episode"

// EpisodeService performs operations on Episodes.
type EpisodeService struct {
	DB *bolt.DB
	Service
}

// Create persists the given Episode.
func (ser *EpisodeService) Create(ep *Episode) error {
	return Create(ep, ser)
}

// Update replaces the value of the Episode with the
// given ID.
func (ser *EpisodeService) Update(ep *Episode) error {
	return Update(ep, ser)
}

// Delete deletes the Episode with the given ID.
func (ser *EpisodeService) Delete(id int) error {
	return Delete(id, ser)
}

// GetAll retrieves all persisted values of Episode.
func (ser *EpisodeService) GetAll(first int, prefixID *int) ([]*Episode, error) {
	vlist, err := GetAll(ser, first, prefixID)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to Episodes: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of Episode that
// pass the filter.
func (ser *EpisodeService) GetFilter(first int, prefixID *int, keep func(ep *Episode) bool) ([]*Episode, error) {
	vlist, err := GetFilter(ser, first, prefixID, func(m Model) bool {
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
func (ser *EpisodeService) GetByID(id int) (*Episode, error) {
	m, err := GetByID(id, ser)
	if err != nil {
		return nil, err
	}

	ep, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return ep, nil
}

// GetByMedia retrieves a list of instances of Episode
// with the given Media ID.
// func (ser *EpisodeService) GetByMedia(mID int, first int, prefixID *int) ([]*Episode, error) {
// return ser.GetFilter(first, prefixID, func(ep *Episode) bool {
// return ep.MediaID == mID
// })
// }

// Database returns the database reference.
func (ser *EpisodeService) Database() *bolt.DB {
	return ser.DB
}

// Bucket returns the name of the bucket for Episode.
func (ser *EpisodeService) Bucket() string {
	return EpisodeBucket
}

// Clean cleans the given Episode for storage
func (ser *EpisodeService) Clean(m Model) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Validate returns an error if the Episode is
// not valid for the database.
func (ser *EpisodeService) Validate(m Model) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Initialize sets initial values for some properties.
func (ser *EpisodeService) Initialize(m Model, id int) error {
	ep, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	ep.ID = id
	ep.Version = 0
	return nil
}

// PersistOldProperties maintains certain properties
// of the existing Episode in updates.
func (ser *EpisodeService) PersistOldProperties(n Model, o Model) error {
	nep, err := ser.AssertType(n)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	oep, err := ser.AssertType(o)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	nep.Version = oep.Version + 1
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

// mapfromModel returns a list of Episode type
// asserted from the given list of Model.
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

// EpisodeSetBucket is the name of the database bucket for
// EpisodeSets.
const EpisodeSetBucket = "EpisodeSet"

// EpisodeSetService performs operations on EpisodeSets.
type EpisodeSetService struct {
	DB *bolt.DB
	Service
}

// Create persists the given EpisodeSet.
func (ser *EpisodeSetService) Create(set *EpisodeSet) error {
	return Create(set, ser)
}

// Update replaces the value of the EpisodeSet with the
// given ID.
func (ser *EpisodeSetService) Update(set *EpisodeSet) error {
	return Update(set, ser)
}

// Delete deletes the EpisodeSet with the given ID.
func (ser *EpisodeSetService) Delete(id int) error {
	return Delete(id, ser)
}

// GetAll retrieves all persisted values of EpisodeSet.
func (ser *EpisodeSetService) GetAll(first int, prefixID *int) ([]*EpisodeSet, error) {
	vlist, err := GetAll(ser, first, prefixID)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to EpisodeSets: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of EpisodeSet that
// pass the filter.
func (ser *EpisodeSetService) GetFilter(first int, prefixID *int, keep func(*EpisodeSet) bool) ([]*EpisodeSet, error) {
	vlist, err := GetFilter(ser, first, prefixID, func(m Model) bool {
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

// GetByID retrieves the persisted EpisodeSet with the given
// ID.
func (ser *EpisodeSetService) GetByID(id int) (*EpisodeSet, error) {
	m, err := GetByID(id, ser)
	if err != nil {
		return nil, err
	}

	set, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return set, nil
}

// GetByMedia retrieves a list of instances of EpisodeSet
// with the given Media ID.
func (ser *EpisodeSetService) GetByMedia(mID int, first int, prefixID *int) ([]*EpisodeSet, error) {
	return ser.GetFilter(first, prefixID, func(set *EpisodeSet) bool {
		return set.MediaID == mID
	})
}

// Database returns the database reference.
func (ser *EpisodeSetService) Database() *bolt.DB {
	return ser.DB
}

// Bucket returns the name of the bucket for EpisodeSet.
func (ser *EpisodeSetService) Bucket() string {
	return EpisodeBucket
}

// Clean cleans the given EpisodeSet for storage.
func (ser *EpisodeSetService) Clean(m Model) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Validate returns an error if the Episodeset is
// not valid for the database.
func (ser *EpisodeSetService) Validate(m Model) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Initialize sets initial values for some properties.
func (ser *EpisodeSetService) Initialize(m Model, id int) error {
	set, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	set.ID = id
	set.Version = 0
	return nil
}

// PersistOldProperties maintains certain properties
// of the existing EpisodeSet in updates.
func (ser *EpisodeSetService) PersistOldProperties(n Model, o Model) error {
	nset, err := ser.AssertType(n)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	oset, err := ser.AssertType(o)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	nset.Version = oset.Version + 1
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

// mapfromModel returns a list of EpisodeSet type
// asserted from the given list of Model.
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