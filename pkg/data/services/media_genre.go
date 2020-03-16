package services

import (
	"errors"
	"fmt"

	"github.com/Dophin2009/nao/pkg/data/models"
	"github.com/Dophin2009/nao/pkg/db"
	json "github.com/json-iterator/go"
)

// MediaGenreService performs operations on MediaGenre.
type MediaGenreService struct {
	MediaService *MediaService
	GenreService *GenreService
	Hooks        db.PersistHooks
}

// NewMediaGenreService returns a MediaGenre.
func NewMediaGenreService(hooks db.PersistHooks, mediaService *MediaService,
	genreService *GenreService) *MediaGenreService {
	// Initialize MediaGenreService
	mediaGenreService := &MediaGenreService{
		MediaService: mediaService,
		GenreService: genreService,
		Hooks:        hooks,
	}

	// Add hook to delete MediaGenre on Media deletion
	deleteMediaGenreOnDeleteMedia := func(mdm db.Model, _ db.Service, tx db.Tx) error {
		mID := mdm.Metadata().ID
		err := mediaGenreService.DeleteByMedia(mID, tx)
		if err != nil {
			return fmt.Errorf("failed to delete MediaGenre by Media ID %d: %w", mID, err)
		}
		return nil
	}
	mdSerHooks := mediaService.PersistHooks()
	mdSerHooks.PreDeleteHooks =
		append(mdSerHooks.PreDeleteHooks, deleteMediaGenreOnDeleteMedia)

	// Add hook to delete MediaGenre on Genre deletion
	deleteMediaGenreOnDeleteGenre := func(gm db.Model, _ db.Service, tx db.Tx) error {
		gID := gm.Metadata().ID
		err := mediaGenreService.DeleteByGenre(gID, tx)
		if err != nil {
			return fmt.Errorf("failed to delete MediaGenre by Genre ID %d: %w", gID, err)
		}
		return nil
	}
	gSerHooks := genreService.PersistHooks()
	gSerHooks.PreDeleteHooks =
		append(gSerHooks.PreDeleteHooks, deleteMediaGenreOnDeleteGenre)

	return mediaGenreService
}

// Create persists the given MediaGenre.
func (ser *MediaGenreService) Create(mg *models.MediaGenre, tx db.Tx) (int, error) {
	return tx.Database().Create(mg, ser, tx)
}

// Update rmglaces the value of the MediaGenre with the given ID.
func (ser *MediaGenreService) Update(mg *models.MediaGenre, tx db.Tx) error {
	return tx.Database().Update(mg, ser, tx)
}

// Delete deletes the MediaGenre with the given ID.
func (ser *MediaGenreService) Delete(id int, tx db.Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// DeleteByMedia deletes the MediaGenres with the given Media ID.
func (ser *MediaGenreService) DeleteByMedia(mID int, tx db.Tx) error {
	return tx.Database().DeleteFilter(ser, tx, func(m db.Model) bool {
		mg, err := ser.AssertType(m)
		if err != nil {
			return false
		}

		return mg.MediaID == mID
	})
}

// DeleteByGenre deletes the MediaGenres with the given Genre ID.
func (ser *MediaGenreService) DeleteByGenre(gID int, tx db.Tx) error {
	return tx.Database().DeleteFilter(ser, tx, func(m db.Model) bool {
		mg, err := ser.AssertType(m)
		if err != nil {
			return false
		}

		return mg.GenreID == gID
	})
}

// GetAll retrieves all persisted values of MediaGenre.
func (ser *MediaGenreService) GetAll(first *int, skip *int, tx db.Tx) ([]*models.MediaGenre, error) {
	vlist, err := tx.Database().GetAll(first, skip, ser, tx)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to MediaGenres: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of MediaGenre that pass the filter.
func (ser *MediaGenreService) GetFilter(
	first *int, skip *int, tx db.Tx, keep func(mg *models.MediaGenre) bool,
) ([]*models.MediaGenre, error) {
	vlist, err := tx.Database().GetFilter(first, skip, ser, tx,
		func(m db.Model) bool {
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
		return nil, fmt.Errorf("failed to map db.Models to MediaGenres: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted MediaGenre with the given ID.
func (ser *MediaGenreService) GetByID(id int, tx db.Tx) (*models.MediaGenre, error) {
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
	ids []int, tx db.Tx, keep func(mg *models.MediaGenre) bool,
) ([]*models.MediaGenre, error) {
	vlist, err := tx.Database().GetMultiple(ids, ser, tx,
		func(m db.Model) bool {
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
		return nil, fmt.Errorf("failed to map db.Models to MediaGenres: %w", err)
	}
	return list, nil
}

// GetByMedia retrieves a list of instances of MediaGenre with the given Media
// ID.
func (ser *MediaGenreService) GetByMedia(
	mID int, first *int, skip *int, tx db.Tx,
) ([]*models.MediaGenre, error) {
	return ser.GetFilter(first, skip, tx, func(mg *models.MediaGenre) bool {
		return mg.MediaID == mID
	})
}

// GetByGenre retrieves a list of instances of MediaGenre with the given Genre
// ID.
func (ser *MediaGenreService) GetByGenre(
	gID int, first *int, skip *int, tx db.Tx,
) ([]*models.MediaGenre, error) {
	return ser.GetFilter(first, skip, tx, func(mg *models.MediaGenre) bool {
		return mg.GenreID == gID
	})
}

// Bucket returns the name of the bucket for MediaGenre.
func (ser *MediaGenreService) Bucket() string {
	return "MediaGenre"
}

// Clean cleans the given MediaGenre for storage.
func (ser *MediaGenreService) Clean(_ db.Model, _ db.Tx) error {
	return nil
}

// Validate returns an error if the MediaGenre is not valid for the database.
func (ser *MediaGenreService) Validate(m db.Model, tx db.Tx) error {
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
func (ser *MediaGenreService) Initialize(_ db.Model, _ db.Tx) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing MediaGenre
// in updates.
func (ser *MediaGenreService) PersistOldProperties(_ db.Model, _ db.Model, _ db.Tx) error {
	return nil
}

// PersistHooks returns the persistence hook functions.
func (ser *MediaGenreService) PersistHooks() *db.PersistHooks {
	return &ser.Hooks
}

// Marshal transforms the given MediaGenre into JSON.
func (ser *MediaGenreService) Marshal(m db.Model) ([]byte, error) {
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
func (ser *MediaGenreService) Unmarshal(buf []byte) (db.Model, error) {
	var mg models.MediaGenre
	err := json.Unmarshal(buf, &mg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &mg, nil
}

// AssertType exposes the given db.Model as a MediaGenre.
func (ser *MediaGenreService) AssertType(m db.Model) (*models.MediaGenre, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	mg, ok := m.(*models.MediaGenre)
	if !ok {
		return nil, fmt.Errorf("model: %w", errors.New("not of MediaGenre type"))
	}
	return mg, nil
}

// mapfromModel returns a list of MediaGenre type asserted from the given list
// of db.Model.
func (ser *MediaGenreService) mapFromModel(vlist []db.Model) ([]*models.MediaGenre, error) {
	list := make([]*models.MediaGenre, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}
