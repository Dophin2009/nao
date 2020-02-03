package data

import (
	"errors"
	"fmt"

	json "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
)

// UserPerson represents a relationship between a User and a Person,
// containing information about the User's opinion on the Person.
type UserPerson struct {
	UserID   int
	PersonID int
	Score    *int
	Comments []Title
	Meta     ModelMetadata
}

// Metadata returns Meta.
func (up *UserPerson) Metadata() *ModelMetadata {
	return &up.Meta
}

// UserPersonBucket is the name of the database bucket for UserPerson.
const UserPersonBucket = "UserPerson"

// UserPersonService performs operations on UserPerson.
type UserPersonService struct {
	DB *bolt.DB
	Service
}

// Create persists the given UserPerson.
func (ser *UserPersonService) Create(up *UserPerson) error {
	return Create(up, ser)
}

// Update ruplaces the value of the UserPerson with the given ID.
func (ser *UserPersonService) Update(up *UserPerson) error {
	return Update(up, ser)
}

// Delete deletes the UserPerson with the given ID.
func (ser *UserPersonService) Delete(id int) error {
	return Delete(id, ser)
}

// GetAll retrieves all persisted values of UserPerson.
func (ser *UserPersonService) GetAll(first *int, skip *int) ([]*UserPerson, error) {
	vlist, err := GetAll(ser, first, skip)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to UserPerson: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of UserPerson that pass the
// filter.
func (ser *UserPersonService) GetFilter(
	first *int, skip *int, keep func(up *UserPerson) bool,
) ([]*UserPerson, error) {
	vlist, err := GetFilter(ser, first, skip, func(m Model) bool {
		up, err := ser.AssertType(m)
		if err != nil {
			return false
		}
		return keep(up)
	})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to UserPerson: %w", err)
	}
	return list, nil
}

// GetMultiple retrieves the persisted UserPerson values specified by the
// given IDs that pass the filter.
func (ser *UserPersonService) GetMultiple(
	ids []int, first *int, skip *int, keep func(up *UserPerson) bool,
) ([]*UserPerson, error) {
	vlist, err := GetMultiple(ser, ids, first, skip, func(m Model) bool {
		up, err := ser.AssertType(m)
		if err != nil {
			return false
		}
		return keep(up)
	})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to UserPersons: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted UserPerson with the given ID.
func (ser *UserPersonService) GetByID(id int) (*UserPerson, error) {
	m, err := GetByID(id, ser)
	if err != nil {
		return nil, err
	}

	up, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return up, nil
}

// GetByUser retrieves the persisted UserPerson with the given User ID.
func (ser *UserPersonService) GetByUser(
	uID int, first *int, skip *int,
) ([]*UserPerson, error) {
	return ser.GetFilter(first, skip, func(up *UserPerson) bool {
		return up.UserID == uID
	})
}

// GetByPerson retrieves the persisted UserPerson with the given Person ID.
func (ser *UserPersonService) GetByPerson(
	cID int, first *int, skip *int,
) ([]*UserPerson, error) {
	return ser.GetFilter(first, skip, func(up *UserPerson) bool {
		return up.PersonID == cID
	})
}

// Database returns the database reference.
func (ser *UserPersonService) Database() *bolt.DB {
	return ser.DB
}

// Bucket returns the name of the bucket for UserPerson.
func (ser *UserPersonService) Bucket() string {
	return UserPersonBucket
}

// Clean cleans the given UserPerson for storage.
func (ser *UserPersonService) Clean(m Model) error {
	_, err := ser.AssertType(m)
	return err
}

// Validate returns an error if the UserPerson is not valid for the database.
func (ser *UserPersonService) Validate(m Model) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	return ser.DB.View(func(tx *bolt.Tx) error {
		// Check if User with ID specified in UserPerson exists
		// Get User bucket, exit if error
		ub, err := Bucket(UserBucket, tx)
		if err != nil {
			return fmt.Errorf("%s %q: %w", errmsgBucketOpen, UserBucket, err)
		}
		_, err = get(e.UserID, ub)
		if err != nil {
			return fmt.Errorf("failed to get User with ID %d: %w", e.UserID, err)
		}

		// Check if Person with ID specified in UserPerson exists
		// Get Person bucket, exit if error
		cb, err := Bucket(PersonBucket, tx)
		if err != nil {
			return fmt.Errorf("%s %q: %w", errmsgBucketOpen, PersonBucket, err)
		}
		_, err = get(e.PersonID, cb)
		if err != nil {
			return fmt.Errorf(
				"failed to get Person with ID %d: %w", e.PersonID, err)
		}

		return nil
	})
}

// Initialize sets initial values for some properties.
func (ser *UserPersonService) Initialize(m Model) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing
// UserPerson in updates.
func (ser *UserPersonService) PersistOldProperties(n Model, o Model) error {
	return nil
}

// Marshal transforms the given UserPerson into JSON.
func (ser *UserPersonService) Marshal(m Model) ([]byte, error) {
	up, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	v, err := json.Marshal(up)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONMarshal, err)
	}

	return v, nil
}

// Unmarshal parses the given JSON into UserPerson.
func (ser *UserPersonService) Unmarshal(buf []byte) (Model, error) {
	var up UserPerson
	err := json.Unmarshal(buf, &up)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &up, nil
}

// AssertType exposes the given Model as a UserPerson.
func (ser *UserPersonService) AssertType(m Model) (*UserPerson, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	up, ok := m.(*UserPerson)
	if !ok {
		return nil,
			fmt.Errorf("model: %w", errors.New("not of UserPerson type"))
	}
	return up, nil
}

// mapfromModel returns a list of UserPerson type asserted from the given
// list of Model.
func (ser *UserPersonService) mapFromModel(vlist []Model) ([]*UserPerson, error) {
	list := make([]*UserPerson, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}
