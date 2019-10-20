package data

// User represents a single user.
type User struct {
	ID       int
	Username string
	Email    string
	Version  int
}

// Validate checks if the given Media is valid
// for the database.
func (ser *UserService) Validate(e *User) (err error) {
	return nil
}
