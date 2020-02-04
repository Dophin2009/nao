package data

import (
	"errors"
	"fmt"

	json "github.com/json-iterator/go"
)

// MediaGenre represents a relationship between single instances of Media and
// Genre.
type MediaGenre struct {
	MediaID int
	GenreID int
	Meta    ModelMetadata
}

// Metadata returns Meta.
func (mg *MediaGenre) Metadata() *ModelMetadata {
	return &mg.Meta
}

// MediaGenreBucket is the name of the database bucket for MediaGenre.
const MediaGenreBucket = "MediaGenre"

// MediaGenreService performs operations on MediaGenre.
type MediaGenreService struct {
	MediaService *MediaService
	GenreService *GenreService
}

// Create persists the given MediaGenre.
func (ser *MediaGenreService) Create(mg *MediaGenre, tx Tx) (int, error) {
	return tx.Database().Create(mg, ser, tx)
}

// Update rmglaces the value of the MediaGenre with the given ID.
func (ser *MediaGenreService) Update(mg *MediaGenre, tx Tx) error {
	return tx.Database().Update(mg, ser, tx)
}

// Delete deletes the MediaGenre with the given ID.
func (ser *MediaGenreService) Delete(id int, tx Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// GetAll retrieves all persisted values of MediaGenre.
func (ser *MediaGenreService) GetAll(first *int, skip *int, tx Tx) ([]*MediaGenre, error) {
	vlist, err := tx.Database().GetAll(first, skip, ser, tx)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to MediaGenres: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of MediaGenre that pass the filter.
func (ser *MediaGenreService) GetFilter(
	first *int, skip *int, tx Tx, keep func(mg *MediaGenre) bool,
) ([]*MediaGenre, error) {
	vlist, err := tx.Database().GetFilter(first, skip, ser, tx,
		func(m Model) bool {
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
func (ser *MediaGenreService) GetByID(id int, tx Tx) (*MediaGenre, error) {
	m, err := tx.Database().GetByID(id, ser, tx)
	if err != nil {
		return nil, err
	}

	mg, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return mg, nil
}

// GetMultiple retrieves the persisted MediaGenre values specified by the given
// IDs that pass the filter.
func (ser *MediaGenreService) GetMultiple(
	ids []int, first *int, skip *int, tx Tx, keep func(mg *MediaGenre) bool,
) ([]*MediaGenre, error) {
	vlist, err := tx.Database().GetMultiple(ids, first, skip, ser, tx,
		func(m Model) bool {
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

// GetByMedia retrieves a list of instances of MediaGenre with the given Media
// ID.
func (ser *MediaGenreService) GetByMedia(
	mID int, first *int, skip *int, tx Tx,
) ([]*MediaGenre, error) {
	return ser.GetFilter(first, skip, tx, func(mg *MediaGenre) bool {
		return mg.MediaID == mID
	})
}

// GetByGenre retrieves a list of instances of MediaGenre with the given Genre
// ID.
func (ser *MediaGenreService) GetByGenre(
	gID int, first *int, skip *int, tx Tx,
) ([]*MediaGenre, error) {
	return ser.GetFilter(first, skip, tx, func(mg *MediaGenre) bool {
		return mg.GenreID == gID
	})
}

// Bucket returns the name of the bucket for MediaGenre.
func (ser *MediaGenreService) Bucket() string {
	return MediaGenreBucket
}

// Clean cleans the given MediaGenre for storage.
func (ser *MediaGenreService) Clean(_ Model, _ Tx) error {
	return nil
}

// Validate returns an error if the MediaGenre is not valid for the database.
func (ser *MediaGenreService) Validate(m Model, tx Tx) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	db := tx.Database()

	// Check if Media with ID specified in new MediaGenre exists
	_, err = db.GetRawByID(e.MediaID, ser.MediaService, tx)
	if err != nil {
		return fmt.Errorf("failed to get Media with ID %d: %w", e.MediaID, err)
	}

	// Check if Genre with ID specified in new MediaGenre exists
	_, err = db.GetRawByID(e.GenreID, ser.GenreService, tx)
	if err != nil {
		return fmt.Errorf("failed to get Genre with ID %d: %w", e.GenreID, err)
	}

	return nil
}

// Initialize sets initial values for some properties.
func (ser *MediaGenreService) Initialize(_ Model, _ Tx) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing MediaGenre
// in updates.
func (ser *MediaGenreService) PersistOldProperties(_ Model, _ Model, _ Tx) error {
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

// mapfromModel returns a list of MediaGenre type asserted from the given list
// of Model.
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
