package data

import (
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

// Clean cleans the given User for storage
func (ser *UserService) Clean(e *User) (err error) {
	e.Username = strings.Trim(e.Username, " ")
	e.Email = strings.Trim(e.Email, " ")
	return nil
}

// Validate checks if the given User is valid
// for the database.
func (ser *UserService) Validate(e *User) (err error) {
	return nil
}

// persistOldProperties maintains certain properties
// of the existing entity in updates
func (ser *UserService) persistOldProperties(old *User, new *User) (err error) {
	new.Version = old.Version + 1
	new.Password = old.Password
	return nil
}

// Create persists a new instance of User
// to the database. Passwords are hashed
// and stored in the given User.
func (ser *UserService) Create(e *User) (err error) {
	e.Password, err = bcrypt.GenerateFromPassword(e.Password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = ser.Clean(e)
	if err != nil {
		return err
	}

	// Verify validity of struct
	err = ser.Validate(e)
	if err != nil {
		return err
	}

	return ser.DB.Update(func(tx *bolt.Tx) error {
		// Get bucket, exit if error
		b, err := Bucket(UserBucketName, tx)
		if err != nil {
			return err
		}

		// Get next ID in sequence and
		// assign to entity
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		e.ID = int(id)
		e.Version = 0

		// Save entity in bucket
		buf, err := json.Marshal(e)
		if err != nil {
			return err
		}

		return b.Put(itob(int(id)), buf)
	})
}

// GetByUsername retrieves a single instance of User
// with the given username.
func (ser *UserService) GetByUsername(e *User) (err error) {
	return ser.DB.View(func(tx *bolt.Tx) error {
		b, err := Bucket(UserBucketName, tx)
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

			if u.Username == e.Username {
				e = &u
				return nil
			}
		}

		return fmt.Errorf("entity with username %s not found", e.Username)
	})
}

// AuthenticateWithPassword checks if the password for
// the User given by the userID matches the provided
// password; returns nil if correct password, error if
// otherwise.
func (ser *UserService) AuthenticateWithPassword(username string, password string) (err error) {
	var e User
	err = ser.GetByUsername(&e)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword(e.Password, []byte(password))
}

// ChangePassword replaces the User specified
// by the given ID with a new one.
func (ser *UserService) ChangePassword(userID int, password string) (err error) {
	var e User
	err = ser.GetByID(&e)
	if err != nil {
		return err
	}

	e.Password, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return ser.DB.Update(func(tx *bolt.Tx) error {
		// Get bucket, exit if error
		b, err := Bucket(UserBucketName, tx)
		if err != nil {
			return err
		}

		// Save entity in bucket
		buf, err := json.Marshal(e)
		if err != nil {
			return err
		}

		return b.Put(itob(int(e.ID)), buf)
	})
}
