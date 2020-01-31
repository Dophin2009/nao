package data

import (
	"errors"
	"fmt"

	json "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
)

// UserMediaList represents a User-created list
// of UserMedia
type UserMediaList struct {
	ID           int
	UserID       int
	Names        []Title
	Descriptions []Title
	Version      int
	Model
}

// Iden returns the ID.
func (uml *UserMediaList) Iden() int {
	return uml.ID
}

// UserMediaListBucket is the name of the database bucket for
// UserMediaList.
const UserMediaListBucket = "UserMediaList"

// UserMediaListService performs operations on UserMediaList.
type UserMediaListService struct {
	DB *bolt.DB
	Service
}

// Create persists the given UserMediaList.
func (ser *UserMediaListService) Create(uml *UserMediaList) error {
	return Create(uml, ser)
}

// Update rumllaces the value of the UserMediaList with the
// given ID.
func (ser *UserMediaListService) Update(uml *UserMediaList) error {
	return Update(uml, ser)
}

// Delete deletes the UserMediaList with the given ID.
func (ser *UserMediaListService) Delete(id int) error {
	return Delete(id, ser)
}

// GetAll retrieves all persisted values of UserMediaList.
func (ser *UserMediaListService) GetAll(first int, prefixID *int) ([]*UserMediaList, error) {
	vlist, err := GetAll(ser, first, prefixID)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to UserMediaLists: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of UserMediaList that
// pass the filter.
func (ser *UserMediaListService) GetFilter(first int, prefixID *int, keep func(uml *UserMediaList) bool) ([]*UserMediaList, error) {
	vlist, err := GetFilter(ser, first, prefixID, func(m Model) bool {
		uml, err := ser.AssertType(m)
		if err != nil {
			return false
		}
		return keep(uml)
	})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to UserMediaLists: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted UserMediaList with the given ID.
func (ser *UserMediaListService) GetByID(id int) (*UserMediaList, error) {
	m, err := GetByID(id, ser)
	if err != nil {
		return nil, err
	}

	uml, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return uml, nil
}

// Database returns the database reference.
func (ser *UserMediaListService) Database() *bolt.DB {
	return ser.DB
}

// Bucket returns the name of the bucket for UserMediaList.
func (ser *UserMediaListService) Bucket() string {
	return UserMediaListBucket
}

// Clean cleans the given UserMediaList for storage
func (ser *UserMediaListService) Clean(m Model) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Validate returns an error if the UserMediaList is
// not valid for the database.
func (ser *UserMediaListService) Validate(m Model) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	return ser.DB.View(func(tx *bolt.Tx) error {
		// Check if User with ID specified in UserMediaList exists
		// Get User bucket, exit if error
		ub, err := Bucket(UserBucket, tx)
		if err != nil {
			return fmt.Errorf("%s: %w", errmsgBucketOpen, err)
		}
		_, err = get(e.UserID, ub)
		if err != nil {
			return fmt.Errorf("failed to get User with ID %d: %w", e.UserID, err)
		}

		return nil
	})
}

// Initialize sets initial values for some properties.
func (ser *UserMediaListService) Initialize(m Model, id int) error {
	uml, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	uml.ID = id
	uml.Version = 0
	return nil
}

// PersistOldProperties maintains certain properties
// of the existing UserMediaList in updates.
func (ser *UserMediaListService) PersistOldProperties(n Model, o Model) error {
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

// Marshal transforms the given UserMediaList into JSON.
func (ser *UserMediaListService) Marshal(m Model) ([]byte, error) {
	uml, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	v, err := json.Marshal(uml)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONMarshal, err)
	}

	return v, nil
}

// Unmarshal parses the given JSON into UserMediaList.
func (ser *UserMediaListService) Unmarshal(buf []byte) (Model, error) {
	var uml UserMediaList
	err := json.Unmarshal(buf, &uml)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &uml, nil
}

// AssertType exposes the given Model as a UserMediaList.
func (ser *UserMediaListService) AssertType(m Model) (*UserMediaList, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	uml, ok := m.(*UserMediaList)
	if !ok {
		return nil, fmt.Errorf("model: %w", errors.New("not of UserMediaList type"))
	}
	return uml, nil
}

// mapfromModel returns a list of UserMediaList type
// asserted from the given list of Model.
func (ser *UserMediaListService) mapFromModel(vlist []Model) ([]*UserMediaList, error) {
	list := make([]*UserMediaList, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}
