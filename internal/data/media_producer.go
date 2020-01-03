package data

import (
	"errors"
	"strings"

	json "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
)

// MediaProducer represents a relationship between single
// instances of Media and Producer
type MediaProducer struct {
	ID         int
	MediaID    int
	ProducerID int
	Role       string
	Version    int
	Model
}

// Iden returns the ID.
func (mc *MediaProducer) Iden() int {
	return mc.ID
}

// MediaProducerBucket is the name of the database bucket for
// MediaProducer.
const MediaProducerBucket = "MediaProducer"

// MediaProducerService performs operations on MediaProducer.
type MediaProducerService struct {
	DB *bolt.DB
	Service
}

// Create persists the given MediaProducer.
func (ser *MediaProducerService) Create(mp *MediaProducer) error {
	return Create(mp, ser)
}

// Update rmplaces the value of the MediaProducer with the
// given ID.
func (ser *MediaProducerService) Update(mp *MediaProducer) error {
	return Update(mp, ser)
}

// Delete deletes the MediaProducer with the given ID.
func (ser *MediaProducerService) Delete(id int) error {
	return Delete(id, ser)
}

// GetAll retrieves all persisted values of MediaProducer.
func (ser *MediaProducerService) GetAll() ([]*MediaProducer, error) {
	vlist, err := GetAll(ser)
	if err != nil {
		return nil, err
	}

	return ser.mapFromModel(vlist)
}

// GetFilter retrieves all persisted values of MediaProducer that
// pass the filter.
func (ser *MediaProducerService) GetFilter(keep func(mp *MediaProducer) bool) ([]*MediaProducer, error) {
	vlist, err := GetFilter(ser, func(m Model) bool {
		mp, err := ser.AssertType(m)
		if err != nil {
			return false
		}
		return keep(mp)
	})
	if err != nil {
		return nil, err
	}

	return ser.mapFromModel(vlist)
}

// GetByID retrieves the persisted MediaProducer with the given ID.
func (ser *MediaProducerService) GetByID(id int) (*MediaProducer, error) {
	m, err := GetByID(id, ser)
	if err != nil {
		return nil, err
	}

	mp, err := ser.AssertType(m)
	if err != nil {
		return nil, err
	}
	return mp, nil
}

// GetByMedia retrieves a list of instances of MediaProducer
// with the given Media ID.
func (ser *MediaProducerService) GetByMedia(mID int) ([]*MediaProducer, error) {
	return ser.GetFilter(func(mp *MediaProducer) bool {
		return mp.MediaID == mID
	})
}

// GetByProducer retrieves a list of instances of MediaProducer
// with the given Producer ID.
func (ser *MediaProducerService) GetByProducer(pID int) ([]*MediaProducer, error) {
	return ser.GetFilter(func(mp *MediaProducer) bool {
		return mp.ProducerID == pID
	})
}

// Database returns the database reference.
func (ser *MediaProducerService) Database() *bolt.DB {
	return ser.DB
}

// Bucket returns the name of the bucket for MediaProducer.
func (ser *MediaProducerService) Bucket() string {
	return MediaProducerBucket
}

// Clean cleans the given MediaProducer for storage.
func (ser *MediaProducerService) Clean(m Model) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return err
	}
	e.Role = strings.Trim(e.Role, " ")
	return nil
}

// Validate returns an error if the MediaProducer is
// not valid for the database.
func (ser *MediaProducerService) Validate(m Model) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return err
	}

	return ser.DB.View(func(tx *bolt.Tx) error {
		// Check if Media with ID specified in new MediaProducer exists
		// Get Media bucket, exit if error
		mb, err := Bucket(MediaBucket, tx)
		if err != nil {
			return err
		}
		_, err = get(e.MediaID, mb)
		if err != nil {
			return err
		}

		// Check if Producer with ID specified in new MediaProducer exists
		// Get Producer bucket, exit if error
		pb, err := Bucket(ProducerBucket, tx)
		if err != nil {
			return err
		}
		_, err = get(e.ProducerID, pb)
		if err != nil {
			return err
		}

		return nil
	})
}

// Initialize sets initial values for some properties.
func (ser *MediaProducerService) Initialize(m Model, id int) error {
	mp, err := ser.AssertType(m)
	if err != nil {
		return err
	}
	mp.ID = id
	mp.Version = 0
	return nil
}

// PersistOldProperties maintains certain properties
// of the existing MediaProducer in updates.
func (ser *MediaProducerService) PersistOldProperties(n Model, o Model) error {
	nm, err := ser.AssertType(n)
	if err != nil {
		return err
	}
	om, err := ser.AssertType(o)
	if err != nil {
		return err
	}
	nm.Version = om.Version + 1
	return nil
}

// Marshal transforms the given MediaProducer into JSON.
func (ser *MediaProducerService) Marshal(m Model) ([]byte, error) {
	mp, err := ser.AssertType(m)
	if err != nil {
		return nil, err
	}

	v, err := json.Marshal(mp)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Unmarshal parses the given JSON into MediaProducer.
func (ser *MediaProducerService) Unmarshal(buf []byte) (Model, error) {
	var mp MediaProducer
	err := json.Unmarshal(buf, &mp)
	if err != nil {
		return nil, err
	}
	return &mp, nil
}

// AssertType exposes the given Model as a MediaProducer.
func (ser *MediaProducerService) AssertType(m Model) (*MediaProducer, error) {
	if m == nil {
		return nil, errors.New("model must not be nil")
	}

	mp, ok := m.(*MediaProducer)
	if !ok {
		return nil, errors.New("model must be of MediaProducer type")
	}
	return mp, nil
}

// mapfromModel returns a list of MediaProducer type
// asserted from the given list of Model.
func (ser *MediaProducerService) mapFromModel(vlist []Model) ([]*MediaProducer, error) {
	list := make([]*MediaProducer, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}
