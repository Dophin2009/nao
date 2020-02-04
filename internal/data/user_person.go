package data

import (
	"errors"
	"fmt"

	json "github.com/json-iterator/go"
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
	UserService   *UserService
	PersonService *PersonService
}

// Create persists the given UserPerson.
func (ser *UserPersonService) Create(up *UserPerson, tx Tx) (int, error) {
	return tx.Database().Create(up, ser, tx)
}

// Update ruplaces the value of the UserPerson with the given ID.
func (ser *UserPersonService) Update(up *UserPerson, tx Tx) error {
	return tx.Database().Update(up, ser, tx)
}

// Delete deletes the UserPerson with the given ID.
func (ser *UserPersonService) Delete(id int, tx Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// GetAll retrieves all persisted values of UserPerson.
func (ser *UserPersonService) GetAll(first *int, skip *int, tx Tx) ([]*UserPerson, error) {
	vlist, err := tx.Database().GetAll(first, skip, ser, tx)
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
	first *int, skip *int, tx Tx, keep func(up *UserPerson) bool,
) ([]*UserPerson, error) {
	vlist, err := tx.Database().GetFilter(first, skip, ser, tx,
		func(m Model) bool {
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
	ids []int, first *int, skip *int, tx Tx, keep func(up *UserPerson) bool,
) ([]*UserPerson, error) {
	vlist, err := tx.Database().GetMultiple(ids, first, skip, ser, tx,
		func(m Model) bool {
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
func (ser *UserPersonService) GetByID(id int, tx Tx) (*UserPerson, error) {
	m, err := tx.Database().GetByID(id, ser, tx)
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
	uID int, first *int, skip *int, tx Tx,
) ([]*UserPerson, error) {
	return ser.GetFilter(first, skip, tx, func(up *UserPerson) bool {
		return up.UserID == uID
	})
}

// GetByPerson retrieves the persisted UserPerson with the given Person ID.
func (ser *UserPersonService) GetByPerson(
	cID int, first *int, skip *int, tx Tx,
) ([]*UserPerson, error) {
	return ser.GetFilter(first, skip, tx, func(up *UserPerson) bool {
		return up.PersonID == cID
	})
}

// Bucket returns the name of the bucket for UserPerson.
func (ser *UserPersonService) Bucket() string {
	return UserPersonBucket
}

// Clean cleans the given UserPerson for storage.
func (ser *UserPersonService) Clean(m Model, _ Tx) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Validate returns an error if the UserPerson is not valid for the database.
func (ser *UserPersonService) Validate(m Model, tx Tx) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	db := tx.Database()

	// Check if User with ID specified in UserPerson exists
	_, err = db.GetRawByID(e.UserID, ser.UserService, tx)
	if err != nil {
		return fmt.Errorf("failed to get User with ID %d: %w", e.UserID, err)
	}

	// Check if Person with ID specified in UserPerson exists
	_, err = db.GetRawByID(e.PersonID, ser.PersonService, tx)
	if err != nil {
		return fmt.Errorf(
			"failed to get Person with ID %d: %w", e.PersonID, err)
	}

	return nil
}

// Initialize sets initial values for some properties.
func (ser *UserPersonService) Initialize(_ Model, _ Tx) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing
// UserPerson in updates.
func (ser *UserPersonService) PersistOldProperties(_ Model, _ Model, _ Tx) error {
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
