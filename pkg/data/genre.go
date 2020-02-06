package data

import (
	"errors"
	"fmt"

	json "github.com/json-iterator/go"
	"gitlab.com/Dophin2009/nao/pkg/db"
)

// Genre represents a single instance of a genre.
type Genre struct {
	Names        []Title
	Descriptions []Title
	Meta         db.ModelMetadata
}

// Metadata returns Meta.
func (g *Genre) Metadata() *db.ModelMetadata {
	return &g.Meta
}

// GenreService performs operations on genre.
type GenreService struct {
	Hooks db.PersistHooks
}

// Create persists the given Genre.
func (ser *GenreService) Create(g *Genre, tx db.Tx) (int, error) {
	return tx.Database().Create(g, ser, tx)
}

// Update rglaces the value of the Genre with the given ID.
func (ser *GenreService) Update(g *Genre, tx db.Tx) error {
	return tx.Database().Update(g, ser, tx)
}

// Delete deletes the Genre with the given ID.
func (ser *GenreService) Delete(id int, tx db.Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// GetAll retrieves all persisted values of Genre.
func (ser *GenreService) GetAll(first *int, skip *int, tx db.Tx) ([]*Genre, error) {
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
	first *int, skip *int, tx db.Tx, keep func(g *Genre) bool,
) ([]*Genre, error) {
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
	ids []int, first *int, tx db.Tx, keep func(c *Genre) bool,
) ([]*Genre, error) {
	vlist, err := tx.Database().GetMultiple(ids, first, ser, tx,
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
func (ser *GenreService) GetByID(id int, tx db.Tx) (*Genre, error) {
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
	var g Genre
	err := json.Unmarshal(buf, &g)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &g, nil
}

// AssertType exposes the given Model as a Genre.
func (ser *GenreService) AssertType(m db.Model) (*Genre, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	g, ok := m.(*Genre)
	if !ok {
		return nil, fmt.Errorf("model: %w", errors.New("not of Genre type"))
	}
	return g, nil
}

// mapfromModel returns a list of Genre type asserted from the given list of
// Model.
func (ser *GenreService) mapFromModel(vlist []db.Model) ([]*Genre, error) {
	list := make([]*Genre, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}
