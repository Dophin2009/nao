package data

import (
	"errors"
	"fmt"

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
func (ser *MediaGenreService) GetAll(first int, prefixID *int) ([]*MediaGenre, error) {
	vlist, err := GetAll(ser, first, prefixID)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to MediaGenres: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of MediaGenre that
// pass the filter.
func (ser *MediaGenreService) GetFilter(first int, prefixID *int, keep func(mg *MediaGenre) bool) ([]*MediaGenre, error) {
	vlist, err := GetFilter(ser, first, prefixID, func(m Model) bool {
		mg, err := ser.AssertType(m)
		if err != nil {
			return false
		}
		return keep(mg)
	})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to MediaGenres: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted MediaGenre with the given ID.
func (ser *MediaGenreService) GetByID(id int) (*MediaGenre, error) {
	m, err := GetByID(id, ser)
	if err != nil {
		return nil, err
	}

	mg, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return mg, nil
}

// GetByMedia retrieves a list of instances of MediaGenre
// with the given Media ID.
func (ser *MediaGenreService) GetByMedia(mID int, first int, prefixID *int) ([]*MediaGenre, error) {
	return ser.GetFilter(first, prefixID, func(mg *MediaGenre) bool {
		return mg.MediaID == mID
	})
}

// GetByGenre retrieves a list of instances of MediaGenre
// with the given Genre ID.
func (ser *MediaGenreService) GetByGenre(gID int, first int, prefixID *int) ([]*MediaGenre, error) {
	return ser.GetFilter(first, prefixID, func(mg *MediaGenre) bool {
		return mg.GenreID == gID
	})
}

// Clean cleans the given MediaGenre for storage.
func (ser *MediaGenreService) Clean(m Model) error {
	_, err := ser.AssertType(m)
	return err
}

// Validate returns an error if the MediaGenre is
// not valid for the database.
func (ser *MediaGenreService) Validate(m Model) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return err
	}

	return ser.DB.View(func(tx *bolt.Tx) error {
		// Check if Media with ID specified in new MediaGenre exists
		// Get Media bucket, exit if error
		mb, err := Bucket(MediaBucket, tx)
		if err != nil {
			return fmt.Errorf("%s %q: %w", errmsgBucketOpen, MediaBucket, err)
		}
		_, err = get(e.MediaID, mb)
		if err != nil {
			return fmt.Errorf("failed to get Media with ID %d: %w", e.MediaID, err)
		}

		// Check if Genre with ID specified in new MediaGenre exists
		// Get Genre bucket, exit if error
		gb, err := Bucket(GenreBucket, tx)
		if err != nil {
			return fmt.Errorf("%s %q: %w", errmsgBucketOpen, GenreBucket, err)
		}
		_, err = get(e.GenreID, gb)
		if err != nil {
			return fmt.Errorf("failed to get Genre with ID %d: %w", e.GenreID, err)
		}

		return nil
	})
}

// Initialize sets initial values for some properties.
func (ser *MediaGenreService) Initialize(m Model, id int) error {
	mg, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	mg.ID = id
	mg.Version = 0
	return nil
}

// PersistOldProperties maintains certain properties
// of the existing MediaGenre in updates.
func (ser *MediaGenreService) PersistOldProperties(n Model, o Model) error {
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

// Marshal transforms the given MediaGenre into JSON.
func (ser *MediaGenreService) Marshal(m Model) ([]byte, error) {
	mg, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	v, err := json.Marshal(mg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONMarshal, err)
	}

	return v, nil
}

// Unmarshal parses the given JSON into MediaGenre.
func (ser *MediaGenreService) Unmarshal(buf []byte) (Model, error) {
	var mg MediaGenre
	err := json.Unmarshal(buf, &mg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &mg, nil
}

// AssertType exposes the given Model as a MediaGenre.
func (ser *MediaGenreService) AssertType(m Model) (*MediaGenre, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	mg, ok := m.(*MediaGenre)
	if !ok {
		return nil, fmt.Errorf("model: %w", errors.New("not of MediaGenre type"))
	}
	return mg, nil
}

// mapfromModel returns a list of MediaGenre type
// asserted from the given list of Model.
func (ser *MediaGenreService) mapFromModel(vlist []Model) ([]*MediaGenre, error) {
	list := make([]*MediaGenre, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}