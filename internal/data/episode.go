package data

import (
	"errors"
	"time"

	json "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
)

// Episode represents a single episode or chapter
// for some media.
type Episode struct {
	ID       int
	MediaID  int
	Titles   []Info
	Date     *time.Time
	Synopses []Info
	Duration *uint
	Filler   bool
	Recap    bool
	Version  int
	Model
}

// Iden returns the ID.
func (ep *Episode) Iden() int {
	return ep.ID
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
func (ser *EpisodeService) GetAll() ([]*Episode, error) {
	vlist, err := GetAll(ser)
	if err != nil {
		return nil, err
	}

	return ser.mapFromModel(vlist)
}

// GetFilter retrieves all persisted values of Episode that
// pass the filter.
func (ser *EpisodeService) GetFilter(keep func(ep *Episode) bool) ([]*Episode, error) {
	vlist, err := GetFilter(ser, func(m Model) bool {
		ep, err := ser.AssertType(m)
		if err != nil {
			return false
		}
		return keep(ep)
	})
	if err != nil {
		return nil, err
	}

	return ser.mapFromModel(vlist)
}

// GetByID retrieves the persisted Episode with the given ID.
func (ser *EpisodeService) GetByID(id int) (*Episode, error) {
	m, err := GetByID(id, ser)
	if err != nil {
		return nil, err
	}

	ep, err := ser.AssertType(m)
	if err != nil {
		return nil, err
	}
	return ep, nil
}

// GetByMedia retrieves a list of instances of Episode
// with the given Media ID.
func (ser *EpisodeService) GetByMedia(mID int) ([]*Episode, error) {
	return ser.GetFilter(func(ep *Episode) bool {
		return ep.MediaID == mID
	})
}

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
	e, err := ser.AssertType(m)
	if err != nil {
		return err
	}

	if err = infoListClean(e.Titles); err != nil {
		return err
	}
	if err = infoListClean(e.Synopses); err != nil {
		return err
	}
	return nil
}

// Validate returns an error if the Episode is
// not valid for the database.
func (ser *EpisodeService) Validate(m Model) error {
	_, err := ser.AssertType(m)
	return err
}

// Initialize sets initial values for some properties.
func (ser *EpisodeService) Initialize(m Model, id int) error {
	ep, err := ser.AssertType(m)
	if err != nil {
		return err
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
		return err
	}
	oep, err := ser.AssertType(o)
	if err != nil {
		return err
	}
	nep.Version = oep.Version + 1
	return nil
}

// Marshal transforms the given Episode into JSON.
func (ser *EpisodeService) Marshal(m Model) ([]byte, error) {
	ep, err := ser.AssertType(m)
	if err != nil {
		return nil, err
	}

	v, err := json.Marshal(ep)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Unmarshal parses the given JSON into Episode.
func (ser *EpisodeService) Unmarshal(buf []byte) (Model, error) {
	var ep Episode
	err := json.Unmarshal(buf, &ep)
	if err != nil {
		return nil, err
	}
	return &ep, nil
}

// AssertType exposes the Model as an Episode.
func (ser *EpisodeService) AssertType(m Model) (*Episode, error) {
	if m == nil {
		return nil, errors.New("model must not be nil")
	}

	ep, ok := m.(*Episode)
	if !ok {
		return nil, errors.New("model must be of Episode type")
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
			return nil, err
		}
	}
	return list, nil
}
