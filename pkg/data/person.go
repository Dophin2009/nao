package data

// Person represents a single person
type Person struct {
	ID          int
	Names       []Info
	Information []Info
	Version     int
}

// Validate returns an error if the Person is
// not valid for the database.
func (ser *PersonService) Validate(e *Person) (err error) {
	return nil
}
