package data

// Character represents a single character.
type Character struct {
	ID          int
	Names       []Info
	Information []Info
	Version     int
}

// Validate returns an error if the Character is
// not valid for the database.
func (ser *CharacterService) Validate(e *Character) (err error) {
	return nil
}
