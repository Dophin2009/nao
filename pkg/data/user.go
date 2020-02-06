package data

import (
	"errors"
	"fmt"
	"strings"

	json "github.com/json-iterator/go"
	"gitlab.com/Dophin2009/nao/pkg/db"
	"golang.org/x/crypto/bcrypt"
)

// TODO: Delete all things related to User when deleting User

// User represents a single user.
type User struct {
	Username    string
	Email       string
	Password    []byte
	Permissions UserPermission
	Meta        db.ModelMetadata
	updatedPass bool `json:"-"`
}

// Metadata returns Meta.
func (u *User) Metadata() *db.ModelMetadata {
	return &u.Meta
}

// UserPermission contains a number of permissions for users for
// reading/writing data.
type UserPermission struct {
	// WriteMedia is the ability modify global Media.
	WriteMedia bool
	// WriteUsers is the ability to modify other Users.
	WriteUsers bool
}

// UserService performs operations on User.
type UserService struct {
	Hooks db.PersistHooks
}

// Create persists the given User.
func (ser *UserService) Create(u *User, tx db.Tx) (int, error) {
	return tx.Database().Create(u, ser, tx)
}

// Update rulaces the value of the User with the given ID.
func (ser *UserService) Update(u *User, tx db.Tx) error {
	return tx.Database().Update(u, ser, tx)
}

// Delete deletes the User with the given ID.
func (ser *UserService) Delete(id int, tx db.Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// GetAll retrieves all persisted values of User.
func (ser *UserService) GetAll(first *int, skip *int, tx db.Tx) ([]*User, error) {
	vlist, err := tx.Database().GetAll(first, skip, ser, tx)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to Users: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of User that pass the filter.
func (ser *UserService) GetFilter(
	first *int, skip *int, tx db.Tx, keep func(u *User) bool,
) ([]*User, error) {
	vlist, err := tx.Database().GetFilter(first, skip, ser, tx,
		func(m db.Model) bool {
			u, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(u)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to Users: %w", err)
	}
	return list, nil
}

// GetMultiple retrieves the persisted User values specified by the
// given IDs that pass the filter.
func (ser *UserService) GetMultiple(
	ids []int, first *int, tx db.Tx, keep func(u *User) bool,
) ([]*User, error) {
	vlist, err := tx.Database().GetMultiple(ids, first, ser, tx,
		func(m db.Model) bool {
			u, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(u)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to Users: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted User with the given ID.
func (ser *UserService) GetByID(id int, tx db.Tx) (*User, error) {
	m, err := tx.Database().GetByID(id, ser, tx)
	if err != nil {
		return nil, err
	}

	u, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return u, nil
}

// GetByUsername retrieves a single instance of User with the given username.
func (ser *UserService) GetByUsername(username string, tx db.Tx) (*User, error) {
	var e User
	_, err := tx.Database().FindFirst(ser, tx, func(m db.Model) (bool, error) {
		u, err := ser.AssertType(m)
		if err != nil {
			return false, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}

		if u.Username == username {
			e = *u
			return true, nil
		}

		return false, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate through keys: %w", err)
	}

	return &e, nil
}

// Authorize checks if the user with the given ID has permissions that meet
// the requirements.
func (ser *UserService) Authorize(userID int, req *UserPermission, tx db.Tx) (*User, error) {
	user, err := ser.GetByID(userID, tx)
	if err != nil {
		return nil, err
	}

	if !ser.RequirementsMet(&user.Permissions, req) {
		return nil, errors.New("insufficient permissions")
	}
	return user, nil
}

// RequirementsMet checks if the given permissions satisfy the required
// permissions.
func (ser *UserService) RequirementsMet(
	perm *UserPermission, req *UserPermission) bool {
	return !(req.WriteMedia && !perm.WriteMedia) &&
		!(req.WriteUsers && !perm.WriteUsers)
}

// AuthenticateWithPassword checks if the password for the User given by the
// username matches the provided password; returns nil if correct password,
// error if otherwise.
func (ser *UserService) AuthenticateWithPassword(
	username string, password string, tx db.Tx) error {
	u, err := ser.GetByUsername(username, tx)
	if err != nil {
		return fmt.Errorf("failed to get User by username %q: %w", username, err)
	}

	err = bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		return fmt.Errorf("failed to match passwords: %w", err)
	}

	return nil
}

// ChangePassword replaces the password of the User specified by the given ID
// with a new one.
func (ser *UserService) ChangePassword(userID int, password string, tx db.Tx) error {
	u, err := ser.GetByID(userID, tx)
	if err != nil {
		return fmt.Errorf("failed to get User by ID %d: %w", userID, err)
	}

	pass, err := ser.HashPassword([]byte(password))
	if err != nil {
		return fmt.Errorf("failed to generate password hash: %w", err)
	}
	u.Password = pass
	u.updatedPass = true

	err = tx.Database().Update(u, ser, tx)
	if err != nil {
		return err
	}

	return nil
}

// HashPassword hashes the given password and returns the result.
func (ser *UserService) HashPassword(pass []byte) ([]byte, error) {
	res, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password hash: %w", err)
	}

	return res, nil
}

// Bucket returns the name of the bucket for User.
func (ser *UserService) Bucket() string {
	return "User"
}

// Clean cleans the given User for storage.
func (ser *UserService) Clean(m db.Model, _ db.Tx) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	e.Username = strings.Trim(e.Username, " ")
	e.Email = strings.Trim(e.Email, " ")
	return nil
}

// Validate checks if the given User is valid for the database.
func (ser *UserService) Validate(m db.Model, tx db.Tx) error {
	u, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	// Check that username does not already exist
	sameUsername, err := ser.GetByUsername(u.Username, tx)
	if sameUsername != nil {
		return fmt.Errorf("username %q: %w", u.Username, errAlreadyExists)
	}

	return nil
}

// Initialize sets initial values for some properties.
func (ser *UserService) Initialize(m db.Model, _ db.Tx) error {
	u, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	pass, err := ser.HashPassword(u.Password)
	if err != nil {
		return fmt.Errorf("failed to generate password hash: %w", err)
	}
	u.Password = pass
	return nil
}

// PersistOldProperties maintains certain properties of the existing User in
// updates.
func (ser *UserService) PersistOldProperties(n db.Model, o db.Model, _ db.Tx) error {
	nu, err := ser.AssertType(n)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	ou, err := ser.AssertType(o)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	// Password may not be changed directly through update; must use
	// ChangePassword
	if !nu.updatedPass {
		nu.Password = ou.Password
	}
	return nil
}

// PersistHooks returns the persistence hook functions.
func (ser *UserService) PersistHooks() *db.PersistHooks {
	return &ser.Hooks
}

// Marshal transforms the given User into JSON.
func (ser *UserService) Marshal(m db.Model) ([]byte, error) {
	u, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	v, err := json.Marshal(u)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONMarshal, err)
	}

	return v, nil
}

// Unmarshal parses the given JSON into User.
func (ser *UserService) Unmarshal(buf []byte) (db.Model, error) {
	var u User
	err := json.Unmarshal(buf, &u)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &u, nil
}

// AssertType exposes the given db.Model as a User.
func (ser *UserService) AssertType(m db.Model) (*User, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	u, ok := m.(*User)
	if !ok {
		return nil, fmt.Errorf("model: %w", errors.New("not of User type"))
	}
	return u, nil
}

// mapfromModel returns a list of User type asserted from the given list of
// db.Model.
func (ser *UserService) mapFromModel(vlist []db.Model) ([]*User, error) {
	list := make([]*User, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}
