package data

import "strings"

// User represents a single user.
type User struct {
	ID       int
	Username string
	Email    string
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
