package data

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	json "github.com/json-iterator/go"
	"gitlab.com/Dophin2009/nao/pkg/db"
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
	Meta            db.ModelMetadata
}

// Metadata returns Meta.
func (m *Media) Metadata() *db.ModelMetadata {
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
	// QuarterWinter is the first quarter of the year, encapsulating the months January,
	// February, and March.
	QuarterWinter Quarter = iota + 1

	// QuarterSpring is the second quarter of the year, encapsulating the months April,
	// May, and June.
	QuarterSpring

	// QuarterSummer is the third quarter of the year, encapsulating the months July,
	// August, and September.
	QuarterSummer

	// QuarterFall is the fouth quarter of the year, encapsulating the months October,
	// November, and December.
	QuarterFall
)

// IsValid checks if the Quarter has a value that is a valid one.
func (q Quarter) IsValid() bool {
	switch q {
	case QuarterWinter, QuarterSpring, QuarterSummer, QuarterFall:
		return true
	}
	return false
}

// String returns the written name of the Quarter.
func (q Quarter) String() string {
	switch q {
	case QuarterWinter:
		return "Winter"
	case QuarterSpring:
		return "Spring"
	case QuarterSummer:
		return "Summer"
	case QuarterFall:
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
		*q = QuarterWinter
	case "Spring":
		*q = QuarterSpring
	case "Summer":
		*q = QuarterSummer
	case "Fall":
		*q = QuarterFall
	default:
		return fmt.Errorf("%s: %w", str, errInvalid)
	}
	return nil
}

// MarshalGQL serializes the Quarter into a GraphQL readable form.
func (q Quarter) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(q.String()))
}

// MediaService performs operations on Media.
type MediaService struct {
	Hooks db.PersistHooks
}

// Create persists the given Media.
func (ser *MediaService) Create(md *Media, tx db.Tx) (int, error) {
	return tx.Database().Create(md, ser, tx)
}

// Update replaces the value of the Media with the given ID.
func (ser *MediaService) Update(md *Media, tx db.Tx) error {
	return tx.Database().Update(md, ser, tx)
}

// Delete deletes the Media with the given ID.
func (ser *MediaService) Delete(id int, tx db.Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// GetAll retrieves all persisted values of Media.
func (ser *MediaService) GetAll(first *int, skip *int, tx db.Tx) ([]*Media, error) {
	vlist, err := tx.Database().GetAll(first, skip, ser, tx)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to Media: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of Media that pass the filter.
func (ser *MediaService) GetFilter(
	first *int, skip *int, tx db.Tx, keep func(md *Media) bool,
) ([]*Media, error) {
	vlist, err := tx.Database().GetFilter(first, skip, ser, tx, func(m db.Model) bool {
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
		return nil, fmt.Errorf("failed to map db.Models to Media: %w", err)
	}
	return list, nil
}

// GetMultiple retrieves the persisted Media values specified by the given
// IDs that pass the filter.
func (ser *MediaService) GetMultiple(
	ids []int, first *int, skip *int, tx db.Tx, keep func(md *Media) bool,
) ([]*Media, error) {
	vlist, err := tx.Database().GetMultiple(ids, first, skip, ser, tx,
		func(m db.Model) bool {
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
		return nil, fmt.Errorf("failed to map db.Models to Media: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted Media with the given ID.
func (ser *MediaService) GetByID(id int, tx db.Tx) (*Media, error) {
	m, err := tx.Database().GetByID(id, ser, tx)
	if err != nil {
		return nil, err
	}

	md, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return md, nil
}

// Bucket returns the name of the bucket for Media.
func (ser *MediaService) Bucket() string {
	return "Media"
}

// Clean cleans the given Media for storage
func (ser *MediaService) Clean(m db.Model, _ db.Tx) error {
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
func (ser *MediaService) Validate(m db.Model, _ db.Tx) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Initialize sets initial values for some properties.
func (ser *MediaService) Initialize(_ db.Model, _ db.Tx) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing Media in
// updates.
func (ser *MediaService) PersistOldProperties(_ db.Model, _ db.Model, _ db.Tx) error {
	return nil
}

// PersistHooks returns the persistence hook functions.
func (ser *MediaService) PersistHooks() *db.PersistHooks {
	return &ser.Hooks
}

// Marshal transforms the given Media into JSON.
func (ser *MediaService) Marshal(m db.Model) ([]byte, error) {
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
func (ser *MediaService) Unmarshal(buf []byte) (db.Model, error) {
	var md Media
	err := json.Unmarshal(buf, &md)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &md, nil
}

// AssertType exposes the given db.Model as a Media.
func (ser *MediaService) AssertType(m db.Model) (*Media, error) {
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
// db.Model.
func (ser *MediaService) mapFromModel(vlist []db.Model) ([]*Media, error) {
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
