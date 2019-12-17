package data

import (
	"errors"
	"fmt"
	"strings"

	json "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/crypto/bcrypt"
)

// User represents a single user.
type User struct {
	ID       int
	Username string
	Email    string
	Password []byte
	// Permissions states the permissions the user has
	// regarding shared/global data
	Permissions Permission
	Version     int
	Model
}

// Iden returns the ID.
func (u *User) Iden() int {
	return u.ID
}

// Permission contains a number of permissions for
// reading/writing data.
type Permission struct {
	// ReadMedia is the ability to query all Media.
	ReadMedia bool
	// ReadUsers is the ability to query all Users
	// and their data.
	ReadUsers bool
	// WriteMedia is the ability modify all Media.
	WriteMedia bool
	// WriteUsers is the ability to modify all Users.
	WriteUsers bool
}

// UserBucket is the name of the database bucket for
// User.
const UserBucket = "User"

// UserService performs operations on User.
type UserService struct {
	DB *bolt.DB
	Service
}

// Create persists the given User.
func (ser *UserService) Create(u *User) error {
	return Create(u, ser)
}

// Update rulaces the value of the User with the
// given ID.
func (ser *UserService) Update(u *User) error {
	return Update(u, ser)
}

// Delete deletes the User with the given ID.
func (ser *UserService) Delete(id int) error {
	return Delete(id, ser)
}

// GetAll retrieves all persisted values of User.
func (ser *UserService) GetAll() ([]*User, error) {
	vlist, err := GetAll(ser)
	if err != nil {
		return nil, err
	}

	return ser.mapFromModel(vlist)
}

// GetFilter retrieves all persisted values of User that
// pass the filter.
func (ser *UserService) GetFilter(keep func(u *User) bool) ([]*User, error) {
	vlist, err := GetFilter(ser, func(m Model) bool {
		u, err := ser.assertType(m)
		if err != nil {
			return false
		}
		return keep(u)
	})
	if err != nil {
		return nil, err
	}

	return ser.mapFromModel(vlist)
}

// GetByID retrieves the persisted User with the given ID.
func (ser *UserService) GetByID(id int) (*User, error) {
	m, err := GetByID(id, ser)
	if err != nil {
		return nil, err
	}

	u, err := ser.assertType(m)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// GetUserByUsername retrieves a single instance of User
// with the given username.
func (ser *UserService) GetByUsername(username string) (*User, error) {
	var e User
	err := ser.DB.View(func(tx *bolt.Tx) error {
		b, err := Bucket(ser.Bucket(), tx)
		if err != nil {
			return err
		}

		c := b.Cursor()
		for id, v := c.First(); id != nil; id, v = c.Next() {
			var u User
			err := json.Unmarshal(v, &u)
			if err != nil {
				return err
			}

			if u.Username == username {
				e = u
				return nil
			}
		}

		return fmt.Errorf("username %s not found", e.Username)
	})
	return &e, err
}

// AuthenticateWithPassword checks if the password for
// the User given by the username matches the provided
// password; returns nil if correct password, error if
// otherwise.
func (ser *UserService) AuthenticateWithPassword(username string, password string) error {
	u, err := ser.GetByUsername(username)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword(u.Password, []byte(password))
}

// ChangePassword replaces the password of the User specified
// by the given ID with a new one.
func (ser *UserService) ChangePassword(userID int, password string) error {
	u, err := ser.GetByID(userID)
	if err != nil {
		return err
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = pass

	return ser.DB.Update(func(tx *bolt.Tx) error {
		// Get bucket, exit if error
		b, err := Bucket(ser.Bucket(), tx)
		if err != nil {
			return err
		}

		// Save entity in bucket
		buf, err := json.Marshal(u)
		if err != nil {
			return err
		}

		return b.Put(itob(int(u.ID)), buf)
	})
}

// Database returns the database reference.
func (ser *UserService) Database() *bolt.DB {
	return ser.DB
}

// Bucket returns the name of the bucket for User.
func (ser *UserService) Bucket() string {
	return UserBucket
}

// Clean cleans the given User for storage.
func (ser *UserService) Clean(m Model) error {
	e, err := ser.assertType(m)
	if err != nil {
		return err
	}

	e.Username = strings.Trim(e.Username, " ")
	e.Email = strings.Trim(e.Email, " ")
	return nil
}

// Validate checks if the given User is valid
// for the database.
func (ser *UserService) Validate(m Model) error {
	u, err := ser.assertType(m)
	if err != nil {
		return err
	}

	// Check that username does not already exist
	err = ser.Database().View(func(tx *bolt.Tx) error {
		// Get bucket, exit if error
		b, err := Bucket(ser.Bucket(), tx)
		if err != nil {
			return err
		}

		// Check for duplicate username
		c := b.Cursor()
		for id, v := c.First(); id != nil; id, v = c.Next() {
			var w User
			err := json.Unmarshal(v, &w)
			if err != nil {
				return err
			}

			if w.Username == u.Username {
				return errors.New("username already exists")
			}
		}
		return nil
	})
	return err
}

// Initialize sets initial values for some properties.
func (ser *UserService) Initialize(m Model, id int) error {
	u, err := ser.assertType(m)
	if err != nil {
		return err
	}

	pass, err := bcrypt.GenerateFromPassword(u.Password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = pass

	u.ID = id
	u.Version = 0
	return nil
}

// PersistOldProperties maintains certain properties
// of the existing User in updates.
func (ser *UserService) PersistOldProperties(n Model, o Model) error {
	nu, err := ser.assertType(n)
	if err != nil {
		return err
	}
	ou, err := ser.assertType(o)
	if err != nil {
		return err
	}
	// Password may not be changed through update;
	// must use ChangePassword
	nu.Password = ou.Password
	nu.Version = ou.Version + 1
	return nil
}

// Marshal transforms the given User into JSON.
func (ser *UserService) Marshal(m Model) ([]byte, error) {
	u, err := ser.assertType(m)
	if err != nil {
		return nil, err
	}

	v, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Unmarshal parses the given JSON into User.
func (ser *UserService) Unmarshal(buf []byte) (Model, error) {
	var u User
	err := json.Unmarshal(buf, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (ser *UserService) assertType(m Model) (*User, error) {
	if m == nil {
		return nil, errors.New("model must not be nil")
	}

	u, ok := m.(*User)
	if !ok {
		return nil, errors.New("model must be of User type")
	}
	return u, nil
}

// mapfromModel returns a list of User type
// asserted from the given list of Model.
func (ser *UserService) mapFromModel(vlist []Model) ([]*User, error) {
	list := make([]*User, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.assertType(v)
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}
