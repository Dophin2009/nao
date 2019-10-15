package data

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

// Character represents a single  character
type Character struct {
	ID          int
	Names       []Info
	Information []Info
	Version     int
}

// Identifier returns the ID of the Character
func (c *Character) Identifier() int {
	return c.ID
}

// SetIdentifier sets the ID of the Character
func (c *Character) SetIdentifier(ID int) {
	c.ID = ID
}

// Ver returns the verison of the Character
func (c *Character) Ver() int {
	return c.Version
}

// UpdateVer increments the version of the
// Character by one
func (c *Character) UpdateVer() {
	c.Version++
}

// Validate returns an error if the Character is
// not valid for the database
func (c *Character) Validate(tx *bolt.Tx) (err error) {
	return nil
}

// CharacterBucketName provides the database bucket name
// for the Character entity
const characterBucketName = "Character"

// CharacterGet retrieves a single instance of Character with
// the given ID
func CharacterGet(ID int, db *bolt.DB) (c Character, err error) {
	err = getByID(ID, &c, characterBucketName, db)
	return
}

// CharacterGetAll retrieves all persisted Character values
func CharacterGetAll(db *bolt.DB) (list []Character, err error) {
	return CharacterGetFilter(db, func(p *Character) bool { return true })
}

// CharacterGetFilter retrieves all persisted Character values
// that pass the filter
func CharacterGetFilter(db *bolt.DB, filter func(c *Character) bool) (list []Character, err error) {
	ilist, err := getFilter(&Character{}, func(entity Idenitifiable) (bool, error) {
		c, ok := entity.(*Character)
		if !ok {
			return false, fmt.Errorf("type assertion failed: entity is not a Character")
		}
		return filter(c), nil
	}, characterBucketName, db)

	list = make([]Character, len(ilist))
	for i, c := range ilist {
		list[i] = *c.(*Character)
	}

	return
}

// CharacterCreate persists a new instance of Character
func CharacterCreate(c *Character, db *bolt.DB) error {
	return create(c, characterBucketName, db)
}

// CharacterUpdate updates the properties of an existing
// persisted Character instance
func CharacterUpdate(c *Character, db *bolt.DB) error {
	return update(c, characterBucketName, db)
}
