package data

import (
	"errors"

	json "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
)

// MediaGenre represents a relationship between single
// instances of Media and Genre.
type MediaGenre struct {
	ID      int
	MediaID int
	GenreID int
	Version int
	Model
}

// Iden returns the ID.
func (mc *MediaGenre) Iden() int {
	return mc.ID
}

// MediaGenreBucket is the name of the database bucket for
// MediaGenre.
const MediaGenreBucket = "MediaGenre"

// MediaGenreService performs operations on MediaGenre.
type MediaGenreService struct {
	DB *bolt.DB
	Service
}

// Create persists the given MediaGenre.
func (ser *MediaGenreService) Create(mg *MediaGenre) error {
	return Create(mg, ser)
}

// Update rmglaces the value of the MediaGenre with the
// given ID.
func (ser *MediaGenreService) Update(mg *MediaGenre) error {
	return Update(mg, ser)
}

// Delete deletes the MediaGenre with the given ID.
func (ser *MediaGenreService) Delete(id int) error {
	return Delete(id, ser)
}

// GetAll retrieves all persisted values of MediaGenre.
func (ser *MediaGenreService) GetAll() ([]*MediaGenre, error) {
	vlist, err := GetAll(ser)
	if err != nil {
		return nil, err
	}

	return ser.mapFromModel(vlist)
}

// GetFilter retrieves all persisted values of MediaGenre that
// pass the filter.
func (ser *MediaGenreService) GetFilter(keep func(mg *MediaGenre) bool) ([]*MediaGenre, error) {
	vlist, err := GetFilter(ser, func(m Model) bool {
		mg, err := ser.assertType(m)
		if err != nil {
			return false
		}
		return keep(mg)
	})
	if err != nil {
		return nil, err
	}

	return ser.mapFromModel(vlist)
}

// GetByID retrieves the persisted MediaGenre with the given ID.
func (ser *MediaGenreService) GetByID(id int) (*MediaGenre, error) {
	m, err := GetByID(id, ser)
	if err != nil {
		return nil, err
	}

	mg, err := ser.assertType(m)
	if err != nil {
		return nil, err
	}
	return mg, nil
}

// GetByMedia retrieves a list of instances of MediaGenre
// with the given Media ID.
func (ser *MediaGenreService) GetByMedia(mID int) ([]*MediaGenre, error) {
	return ser.GetFilter(func(mg *MediaGenre) bool {
		return mg.MediaID == mID
	})
}

// GetByGenre retrieves a list of instances of MediaGenre
// with the given Genre ID.
func (ser *MediaGenreService) GetByGenre(gID int) ([]*MediaGenre, error) {
	return ser.GetFilter(func(mg *MediaGenre) bool {
		return mg.GenreID == gID
	})
}

// Clean cleans the given MediaGenre for storage.
func (ser *MediaGenreService) Clean(m Model) error {
	_, err := ser.assertType(m)
	return err
}

// Validate returns an error if the MediaGenre is
// not valid for the database.
func (ser *MediaGenreService) Validate(m Model) error {
	e, err := ser.assertType(m)
	if err != nil {
		return err
	}

	return ser.DB.View(func(tx *bolt.Tx) error {
		// Check if Media with ID specified in new MediaGenre exists
		// Get Media bucket, exit if error
		mb, err := Bucket(MediaBucket, tx)
		if err != nil {
			return err
		}
		_, err = get(e.MediaID, mb)
		if err != nil {
			return err
		}

		// Check if Genre with ID specified in new MediaGenre exists
		// Get Genre bucket, exit if error
		gb, err := Bucket(GenreBucket, tx)
		if err != nil {
			return err
		}
		_, err = get(e.GenreID, gb)
		if err != nil {
			return err
		}

		return nil
	})
}

// Initialize sets initial values for some properties.
func (ser *MediaGenreService) Initialize(m Model, id int) error {
	mg, err := ser.assertType(m)
	if err != nil {
		return err
	}
	mg.ID = id
	mg.Version = 0
	return nil
}

// PersistOldProperties maintains certain properties
// of the existing MediaGenre in updates.
func (ser *MediaGenreService) PersistOldProperties(n Model, o Model) error {
	nm, err := ser.assertType(n)
	if err != nil {
		return err
	}
	om, err := ser.assertType(o)
	if err != nil {
		return err
	}
	nm.Version = om.Version + 1
	return nil
}

// Marshal transforms the given MediaGenre into JSON.
func (ser *MediaGenreService) Marshal(m Model) ([]byte, error) {
	mg, err := ser.assertType(m)
	if err != nil {
		return nil, err
	}

	v, err := json.Marshal(mg)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Unmarshal parses the given JSON into MediaGenre.
func (ser *MediaGenreService) Unmarshal(buf []byte) (Model, error) {
	var mg MediaGenre
	err := json.Unmarshal(buf, &mg)
	if err != nil {
		return nil, err
	}
	return &mg, nil
}

func (ser *MediaGenreService) assertType(m Model) (*MediaGenre, error) {
	if m == nil {
		return nil, errors.New("model must not be nil")
	}

	mg, ok := m.(*MediaGenre)
	if !ok {
		return nil, errors.New("model must be of MediaGenre type")
	}
	return mg, nil
}

// mapfromModel returns a list of MediaGenre type
// asserted from the given list of Model.
func (ser *MediaGenreService) mapFromModel(vlist []Model) ([]*MediaGenre, error) {
	list := make([]*MediaGenre, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.assertType(v)
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}
