package data

import (
	"encoding/json"
	"strings"

	bolt "go.etcd.io/bbolt"
	"golang.org/x/crypto/bcrypt"
)

// User represents a single user.
type User struct {
	ID       int
	Username string
	Email    string
	Password []byte
	Version  int
}

// Clean cleans the given User for storage
func (ser *UserService) Clean(e *User) (err error) {
	e.Username = strings.Trim(e.Username, " ")
	e.Email = strings.Trim(e.Email, " ")
	return nil
}

// Validate checks if the given Media is valid
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

// Authenticate checks if the password for
// the User given by the userID matches the
// provided password; returns nil if correct
// password, error if otherwise.
func (ser *UserService) Authenticate(userID int, password string) (err error) {
	var e User
	err = ser.GetByID(&e)
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
