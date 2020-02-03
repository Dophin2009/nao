package data

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	json "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
)

// TODO: Move UnmarshalGQL and MarshalGQL of Quarter to graphql package; they
// do not belong in this package
// TODO: Fuzzy search of models

// Media represents a single instance of a media
type Media struct {
	Titles          []Title
	Synopses        []Title
	Background      []Title
	StartDate       *time.Time
	EndDate         *time.Time
	SeasonPremiered Season
	Type            *string
	Source          *string
	Meta            ModelMetadata
}

// Metadata returns Meta.
func (m *Media) Metadata() *ModelMetadata {
	return &m.Meta
}

// Season contains information about the quarter and year.
type Season struct {
	Quarter *Quarter
	Year    *int
}

// Quarter represents the quarter of the year by integer.
type Quarter int

const (
	// Winter is the first quarter of the year, encapsulating the months January,
	// February, and March.
	Winter Quarter = iota + 1

	// Spring is the second quarter of the year, encapsulating the months April,
	// May, and June.
	Spring

	// Summer is the third quarter of the year, encapsulating the months July,
	// August, and September.
	Summer

	// Fall is the fouth quarter of the year, encapsulating the months October,
	// November, and December.
	Fall
)

// IsValid checks if the Quarter has a value that is a valid one.
func (q Quarter) IsValid() bool {
	switch q {
	case Winter, Spring, Summer, Fall:
		return true
	}
	return false
}

// String returns the written name of the Quarter.
func (q Quarter) String() string {
	switch q {
	case Winter:
		return "Winter"
	case Spring:
		return "Spring"
	case Summer:
		return "Summer"
	case Fall:
		return "Fall"
	}
	return fmt.Sprintf("%d", int(q))
}

// UnmarshalGQL casts the type of the given value to a Quarter.
func (q *Quarter) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("%v: %w", v, errInvalid)
	}

	switch str {
	case "Winter":
		*q = Winter
	case "Spring":
		*q = Spring
	case "Summer":
		*q = Summer
	case "Fall":
		*q = Fall
	default:
		return fmt.Errorf("%s: %w", str, errInvalid)
	}
	return nil
}

// MarshalGQL serializes the Quarter into a GraphQL readable form.
func (q Quarter) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(q.String()))
}

// MediaBucket is the name of the database bucket for Media.
const MediaBucket = "Media"

// MediaService performs operations on Media.
type MediaService struct {
	DB *bolt.DB
	Service
}

// Create persists the given Media.
func (ser *MediaService) Create(md *Media) error {
	return Create(md, ser)
}

// Update replaces the value of the Media with the given ID.
func (ser *MediaService) Update(md *Media) error {
	return Update(md, ser)
}

// Delete deletes the Media with the given ID.
func (ser *MediaService) Delete(id int) error {
	return Delete(id, ser)
}

// GetAll retrieves all persisted values of Media.
func (ser *MediaService) GetAll(first *int, skip *int) ([]*Media, error) {
	vlist, err := GetAll(ser, first, skip)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to Media: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of Media that pass the filter.
func (ser *MediaService) GetFilter(
	first *int, skip *int, keep func(md *Media) bool,
) ([]*Media, error) {
	vlist, err := GetFilter(ser, first, skip, func(m Model) bool {
		md, err := ser.AssertType(m)
		if err != nil {
			return false
		}
		return keep(md)
	})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to Media: %w", err)
	}
	return list, nil
}

// GetMultiple retrieves the persisted Media values specified by the given
// IDs that pass the filter.
func (ser *MediaService) GetMultiple(
	ids []int, first *int, skip *int, keep func(md *Media) bool,
) ([]*Media, error) {
	vlist, err := GetMultiple(ser, ids, first, skip, func(m Model) bool {
		md, err := ser.AssertType(m)
		if err != nil {
			return false
		}
		return keep(md)
	})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to Media: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted Media with the given ID.
func (ser *MediaService) GetByID(id int) (*Media, error) {
	m, err := GetByID(id, ser)
	if err != nil {
		return nil, err
	}

	md, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return md, nil
}

// Database returns the database reference.
func (ser *MediaService) Database() *bolt.DB {
	return ser.DB
}

// Bucket returns the name of the bucket for Media.
func (ser *MediaService) Bucket() string {
	return MediaBucket
}

// Clean cleans the given Media for storage
func (ser *MediaService) Clean(m Model) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	if e.Type != nil {
		*e.Type = strings.Trim(*e.Type, " ")
	}
	if e.Source != nil {
		*e.Source = strings.Trim(*e.Source, " ")
	}

	if e.SeasonPremiered.Quarter != nil && *e.SeasonPremiered.Quarter > 4 {
		*e.SeasonPremiered.Quarter = 0
	}
	return nil
}

// Validate checks if the given Media is valid.
func (ser *MediaService) Validate(m Model) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Initialize sets initial values for some properties.
func (ser *MediaService) Initialize(m Model) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing Media in
// updates.
func (ser *MediaService) PersistOldProperties(n Model, o Model) error {
	return nil
}

// Marshal transforms the given Media into JSON.
func (ser *MediaService) Marshal(m Model) ([]byte, error) {
	md, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	v, err := json.Marshal(md)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONMarshal, err)
	}

	return v, nil
}

// Unmarshal parses the given JSON into Media.
func (ser *MediaService) Unmarshal(buf []byte) (Model, error) {
	var md Media
	err := json.Unmarshal(buf, &md)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &md, nil
}

// AssertType exposes the given Model as a Media.
func (ser *MediaService) AssertType(m Model) (*Media, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	md, ok := m.(*Media)
	if !ok {
		return nil, fmt.Errorf("model: %w", errors.New("not of Media type"))
	}
	return md, nil
}

// mapFromModel returns a list of Media type asserted from the given list of
// Model.
func (ser *MediaService) mapFromModel(vlist []Model) ([]*Media, error) {
	list := make([]*Media, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}
