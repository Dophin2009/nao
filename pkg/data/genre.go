package data

// Genre represents a single instance of a genre.
type Genre struct {
	ID           int
	Names        []Info
	Descriptions []Info
	Version      int
}

// Validate returns an error if the Genre is
// not valid for the database.
func (ser *GenreService) Validate(e *Genre) (err error) {
	return nil
}
