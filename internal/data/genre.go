package data

import (
	"errors"
	"fmt"

	"github.com/Dophin2009/nao/pkg/models"
	"github.com/Dophin2009/nao/pkg/db"
	json "github.com/json-iterator/go"
)

// GenreService performs operations on genre.
type GenreService struct {
	Hooks db.PersistHooks
}

// NewGenreService returns a GenreService.
func NewGenreService(hooks db.PersistHooks) *GenreService {
	return &GenreService{
		Hooks: hooks,
	}
}

// Create persists the given Genre.
func (ser *GenreService) Create(g *models.Genre, tx db.Tx) (int, error) {
	return tx.Database().Create(g, ser, tx)
}

// Update rglaces the value of the Genre with the given ID.
func (ser *GenreService) Update(g *models.Genre, tx db.Tx) error {
	return tx.Database().Update(g, ser, tx)
}

// Delete deletes the Genre with the given ID.
func (ser *GenreService) Delete(id int, tx db.Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// GetAll retrieves all persisted values of Genre.
func (ser *GenreService) GetAll(first *int, skip *int, tx db.Tx) ([]*models.Genre, error) {
	vlist, err := tx.Database().GetAll(first, skip, ser, tx)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to Genres: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of Genre that pass the filter.
func (ser *GenreService) GetFilter(
	first *int, skip *int, tx db.Tx, keep func(g *models.Genre) bool,
) ([]*models.Genre, error) {
	vlist, err := tx.Database().GetFilter(first, skip, ser, tx,
		func(m db.Model) bool {
			g, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(g)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to Genres: %w", err)
	}
	return list, nil
}

// GetMultiple retrieves the persisted Genre values specified by the given
// IDs that pass the filter.
func (ser *GenreService) GetMultiple(
	ids []int, tx db.Tx, keep func(c *models.Genre) bool,
) ([]*models.Genre, error) {
	vlist, err := tx.Database().GetMultiple(ids, ser, tx,
		func(m db.Model) bool {
			g, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(g)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to Genres: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted Genre with the given ID.
func (ser *GenreService) GetByID(id int, tx db.Tx) (*models.Genre, error) {
	m, err := tx.Database().GetByID(id, ser, tx)
	if err != nil {
		return nil, err
	}

	g, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return g, nil
}

// Bucket returns the name of the bucket for Genre.
func (ser *GenreService) Bucket() string {
	return "Genre"
}

// Clean cleans the given Genre for storage
func (ser *GenreService) Clean(m db.Model, _ db.Tx) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Validate returns an error if the Genre is not valid for the database.
func (ser *GenreService) Validate(m db.Model, _ db.Tx) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Initialize sets initial values for some properties.
func (ser *GenreService) Initialize(_ db.Model, _ db.Tx) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing Genre in
// updates.
func (ser *GenreService) PersistOldProperties(_ db.Model, _ db.Model, _ db.Tx) error {
	return nil
}

// PersistHooks returns the persistence hook functions.
func (ser *GenreService) PersistHooks() *db.PersistHooks {
	return &ser.Hooks
}

// Marshal transforms the given Genre into JSON.
func (ser *GenreService) Marshal(m db.Model) ([]byte, error) {
	g, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	v, err := json.Marshal(g)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONMarshal, err)
	}

	return v, nil
}

// Unmarshal parses the given JSON into Genre.
func (ser *GenreService) Unmarshal(buf []byte) (db.Model, error) {
	var g models.Genre
	err := json.Unmarshal(buf, &g)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &g, nil
}

// AssertType exposes the given Model as a Genre.
func (ser *GenreService) AssertType(m db.Model) (*models.Genre, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	g, ok := m.(*models.Genre)
	if !ok {
		return nil, fmt.Errorf("model: %w", errors.New("not of Genre type"))
	}
	return g, nil
}

// mapfromModel returns a list of Genre type asserted from the given list of
// Model.
func (ser *GenreService) mapFromModel(vlist []db.Model) ([]*models.Genre, error) {
	list := make([]*models.Genre, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}
