package data

// Character represents a single character.
type Character struct {
	ID          int
	Names       []Info
	Information []Info
	Version     int
}

// Clean cleans the given Character for storage
func (ser *CharacterService) Clean(e *Character) (err error) {
	if err = infoListClean(e.Names); err != nil {
		return err
	}
	if err = infoListClean(e.Information); err != nil {
		return err
	}
	return nil
}

// Validate returns an error if the Character is
// not valid for the database.
func (ser *CharacterService) Validate(e *Character) (err error) {
	return nil
}

// persistOldProperties maintains certain properties
// of the existing entity in updates
func (ser *CharacterService) persistOldProperties(old *Character, new *Character) (err error) {
	new.Version = old.Version + 1
	return nil
}
