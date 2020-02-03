package data

import (
	"errors"
	"fmt"
	"strings"

	json "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/crypto/bcrypt"
)

// TODO: Delete all things related to User when deleting User

// User represents a single user.
type User struct {
	Username string
	Email    string
	Password []byte
	// Permissions states the permissions the user has regarding shared/global
	// data.
	Permissions Permission
	Meta        ModelMetadata
}

// Metadata returns Meta.
func (u *User) Metadata() *ModelMetadata {
	return &u.Meta
}

// Permission contains a number of permissions for reading/writing data.
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

// UserBucket is the name of the database bucket for User.
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

// Update rulaces the value of the User with the given ID.
func (ser *UserService) Update(u *User) error {
	return Update(u, ser)
}

// Delete deletes the User with the given ID.
func (ser *UserService) Delete(id int) error {
	return Delete(id, ser)
}

// GetAll retrieves all persisted values of User.
func (ser *UserService) GetAll(first *int, skip *int) ([]*User, error) {
	vlist, err := GetAll(ser, first, skip)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to Users: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of User that pass the filter.
func (ser *UserService) GetFilter(
	first *int, skip *int, keep func(u *User) bool,
) ([]*User, error) {
	vlist, err := GetFilter(ser, first, skip, func(m Model) bool {
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
		return nil, fmt.Errorf("failed to map Models to Users: %w", err)
	}
	return list, nil
}

// GetMultiple retrieves the persisted User values specified by the
// given IDs that pass the filter.
func (ser *UserService) GetMultiple(
	ids []int, first *int, skip *int, keep func(u *User) bool,
) ([]*User, error) {
	vlist, err := GetMultiple(ser, ids, first, skip, func(m Model) bool {
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
		return nil, fmt.Errorf("failed to map Models to Users: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted User with the given ID.
func (ser *UserService) GetByID(id int) (*User, error) {
	m, err := GetByID(id, ser)
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
func (ser *UserService) GetByUsername(username string) (*User, error) {
	var e User
	err := ser.DB.View(func(tx *bolt.Tx) error {
		// Open User bucket
		b, err := Bucket(ser.Bucket(), tx)
		if err != nil {
			return fmt.Errorf("%s %q: %w", errmsgBucketOpen, ser.Bucket(), err)
		}

		// Iterate through values until username matches
		c := b.Cursor()
		for id, v := c.First(); id != nil; id, v = c.Next() {
			m, err := ser.Unmarshal(v)
			if err != nil {
				return fmt.Errorf("%s: %w", errmsgModelUnmarshal, err)
			}

			u, err := ser.AssertType(m)
			if err != nil {
				return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
			}

			if u.Username == username {
				e = *u
				return nil
			}
		}

		return fmt.Errorf("username %q: %w", username, errNotFound)
	})
	if err != nil {
		return nil, err
	}

	return &e, nil
}

// AuthenticateWithPassword checks if the password for the User given by the
// username matches the provided password; returns nil if correct password,
// error if otherwise.
func (ser *UserService) AuthenticateWithPassword(username string, password string) error {
	u, err := ser.GetByUsername(username)
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
func (ser *UserService) ChangePassword(userID int, password string) error {
	u, err := ser.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get User by ID %d: %w", userID, err)
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to generate password hash: %w", err)
	}
	u.Password = pass

	return ser.DB.Update(func(tx *bolt.Tx) error {
		// Get bucket, exit if error
		b, err := Bucket(ser.Bucket(), tx)
		if err != nil {
			return fmt.Errorf("%s %q: %w", errmsgBucketOpen, ser.Bucket(), err)
		}

		// Save entity in bucket
		buf, err := json.Marshal(u)
		if err != nil {
			return fmt.Errorf("%s: %w", errmsgJSONMarshal, err)
		}

		err = b.Put(itob(int(u.Meta.ID)), buf)
		if err != nil {
			return fmt.Errorf("%s: %w", errmsgBucketPut, err)
		}

		return nil
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
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	e.Username = strings.Trim(e.Username, " ")
	e.Email = strings.Trim(e.Email, " ")
	return nil
}

// Validate checks if the given User is valid for the database.
func (ser *UserService) Validate(m Model) error {
	u, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	// Check that username does not already exist
	err = ser.Database().View(func(tx *bolt.Tx) error {
		// Get bucket, exit if error
		b, err := Bucket(ser.Bucket(), tx)
		if err != nil {
			return fmt.Errorf("%s %q: %w", errmsgBucketOpen, ser.Bucket(), err)
		}

		// Check for duplicate username
		c := b.Cursor()
		for id, v := c.First(); id != nil; id, v = c.Next() {
			m, err := ser.Unmarshal(v)
			if err != nil {
				return fmt.Errorf("%s: %w", errmsgModelUnmarshal, err)
			}

			w, err := ser.AssertType(m)
			if err != nil {
				return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
			}

			if w.Username == u.Username {
				return fmt.Errorf("username %q: %w", u.Username, errAlreadyExists)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to verify uniqueness of username: %w", err)
	}

	return nil
}

// Initialize sets initial values for some properties.
func (ser *UserService) Initialize(m Model) error {
	u, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	pass, err := bcrypt.GenerateFromPassword(u.Password, bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to generate password hash: %w", err)
	}
	u.Password = pass
	return nil
}

// PersistOldProperties maintains certain properties of the existing User in
// updates.
func (ser *UserService) PersistOldProperties(n Model, o Model) error {
	nu, err := ser.AssertType(n)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	ou, err := ser.AssertType(o)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	// Password may not be changed through update; must use ChangePassword
	nu.Password = ou.Password
	return nil
}

// Marshal transforms the given User into JSON.
func (ser *UserService) Marshal(m Model) ([]byte, error) {
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
func (ser *UserService) Unmarshal(buf []byte) (Model, error) {
	var u User
	err := json.Unmarshal(buf, &u)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &u, nil
}

// AssertType exposes the given Model as a User.
func (ser *UserService) AssertType(m Model) (*User, error) {
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
// Model.
func (ser *UserService) mapFromModel(vlist []Model) ([]*User, error) {
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
