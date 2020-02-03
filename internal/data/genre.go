package data

import (
	"errors"
	"fmt"

	json "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
)

// Genre represents a single instance of a genre.
type Genre struct {
	Names        []Title
	Descriptions []Title
	Meta         ModelMetadata
}

// Metadata returns Meta.
func (g *Genre) Metadata() *ModelMetadata {
	return &g.Meta
}

// GenreBucket is the name of the database bucket for Genre.
const GenreBucket = "Genre"

// GenreService performs operations on genre.
type GenreService struct {
	DB *bolt.DB
}

// Create persists the given Genre.
func (ser *GenreService) Create(g *Genre) error {
	return Create(g, ser)
}

// Update rglaces the value of the Genre with the given ID.
func (ser *GenreService) Update(g *Genre) error {
	return Update(g, ser)
}

// Delete deletes the Genre with the given ID.
func (ser *GenreService) Delete(id int) error {
	return Delete(id, ser)
}

// GetAll retrieves all persisted values of Genre.
func (ser *GenreService) GetAll(first *int, skip *int) ([]*Genre, error) {
	vlist, err := GetAll(ser, first, skip)
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
	first *int, skip *int, keep func(g *Genre) bool,
) ([]*Genre, error) {
	vlist, err := GetFilter(ser, first, skip, func(m Model) bool {
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
	ids []int, first *int, skip *int, keep func(c *Genre) bool,
) ([]*Genre, error) {
	vlist, err := GetMultiple(ser, ids, first, skip, func(m Model) bool {
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
func (ser *GenreService) GetByID(id int) (*Genre, error) {
	m, err := GetByID(id, ser)
	if err != nil {
		return nil, err
	}

	g, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return g, nil
}

// Database returns the database reference.
func (ser *GenreService) Database() *bolt.DB {
	return ser.DB
}

// Bucket returns the name of the bucket for Genre.
func (ser *GenreService) Bucket() string {
	return GenreBucket
}

// Clean cleans the given Genre for storage
func (ser *GenreService) Clean(m Model) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Validate returns an error if the Genre is not valid for the database.
func (ser *GenreService) Validate(m Model) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Initialize sets initial values for some properties.
func (ser *GenreService) Initialize(m Model) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing Genre in
// updates.
func (ser *GenreService) PersistOldProperties(n Model, o Model) error {
	return nil
}

// Marshal transforms the given Genre into JSON.
func (ser *GenreService) Marshal(m Model) ([]byte, error) {
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
func (ser *GenreService) Unmarshal(buf []byte) (Model, error) {
	var g Genre
	err := json.Unmarshal(buf, &g)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &g, nil
}

// AssertType exposes the given Model as a Genre.
func (ser *GenreService) AssertType(m Model) (*Genre, error) {
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
func (ser *GenreService) mapFromModel(vlist []Model) ([]*Genre, error) {
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
