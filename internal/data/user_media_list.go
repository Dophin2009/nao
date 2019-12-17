package data

import (
	"errors"

	json "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
)

// UserMediaList represents a User-created list
// of UserMedia
type UserMediaList struct {
	ID           int
	UserID       int
	Names        []Info
	Descriptions []Info
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
func (ser *UserMediaListService) GetAll() ([]*UserMediaList, error) {
	vlist, err := GetAll(ser)
	if err != nil {
		return nil, err
	}

	return ser.mapFromModel(vlist)
}

// GetFilter retrieves all persisted values of UserMediaList that
// pass the filter.
func (ser *UserMediaListService) GetFilter(keep func(uml *UserMediaList) bool) ([]*UserMediaList, error) {
	vlist, err := GetFilter(ser, func(m Model) bool {
		uml, err := ser.assertType(m)
		if err != nil {
			return false
		}
		return keep(uml)
	})
	if err != nil {
		return nil, err
	}

	return ser.mapFromModel(vlist)
}

// GetByID retrieves the persisted UserMediaList with the given ID.
func (ser *UserMediaListService) GetByID(id int) (*UserMediaList, error) {
	m, err := GetByID(id, ser)
	if err != nil {
		return nil, err
	}

	uml, err := ser.assertType(m)
	if err != nil {
		return nil, err
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
	e, err := ser.assertType(m)
	if err != nil {
		return err
	}

	if err := infoListClean(e.Names); err != nil {
		return err
	}
	if err := infoListClean(e.Descriptions); err != nil {
		return err
	}
	return nil
}

// Validate returns an error if the UserMediaList is
// not valid for the database.
func (ser *UserMediaListService) Validate(m Model) error {
	e, err := ser.assertType(m)
	if err != nil {
		return err
	}

	return ser.DB.View(func(tx *bolt.Tx) error {
		// Check if User with ID specified in UserMediaList exists
		// Get User bucket, exit if error
		ub, err := Bucket(UserBucket, tx)
		if err != nil {
			return err
		}
		_, err = get(e.UserID, ub)
		if err != nil {
			return err
		}

		return nil
	})
}

// Initialize sets initial values for some properties.
func (ser *UserMediaListService) Initialize(m Model, id int) error {
	uml, err := ser.assertType(m)
	if err != nil {
		return err
	}
	uml.ID = id
	uml.Version = 0
	return nil
}

// PersistOldProperties maintains certain properties
// of the existing UserMediaList in updates.
func (ser *UserMediaListService) PersistOldProperties(n Model, o Model) error {
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

// Marshal transforms the given UserMediaList into JSON.
func (ser *UserMediaListService) Marshal(m Model) ([]byte, error) {
	uml, err := ser.assertType(m)
	if err != nil {
		return nil, err
	}

	v, err := json.Marshal(uml)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Unmarshal parses the given JSON into UserMediaList.
func (ser *UserMediaListService) Unmarshal(buf []byte) (Model, error) {
	var uml UserMediaList
	err := json.Unmarshal(buf, &uml)
	if err != nil {
		return nil, err
	}
	return &uml, nil
}

func (ser *UserMediaListService) assertType(m Model) (*UserMediaList, error) {
	if m == nil {
		return nil, errors.New("model must not be nil")
	}

	uml, ok := m.(*UserMediaList)
	if !ok {
		return nil, errors.New("model must be of UserMediaList type")
	}
	return uml, nil
}

// mapfromModel returns a list of UserMediaList type
// asserted from the given list of Model.
func (ser *UserMediaListService) mapFromModel(vlist []Model) ([]*UserMediaList, error) {
	list := make([]*UserMediaList, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.assertType(v)
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}
