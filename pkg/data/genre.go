package data

// Genre represents a single instance of a genre.
type Genre struct {
	ID           int
	Names        []Info
	Descriptions []Info
	Version      int
}

// Clean cleans the given Genre for storage
func (ser *GenreService) Clean(e *Genre) (err error) {
	if err = infoListClean(e.Names); err != nil {
		return err
	}
	if err = infoListClean(e.Descriptions); err != nil {
		return err
	}
	return nil
}

// Validate returns an error if the Genre is
// not valid for the database.
func (ser *GenreService) Validate(e *Genre) (err error) {
	return nil
}
