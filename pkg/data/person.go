package data

// Person represents a single person
type Person struct {
	ID          int
	Names       []Info
	Information []Info
	Version     int
}

// Clean cleans the given Person for storage
func (ser *PersonService) Clean(e *Person) (err error) {
	if err := infoListClean(e.Names); err != nil {
		return err
	}
	if err := infoListClean(e.Information); err != nil {
		return err
	}
	return nil
}

// Validate returns an error if the Person is
// not valid for the database.
func (ser *PersonService) Validate(e *Person) (err error) {
	return nil
}
