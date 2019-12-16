package data

import (
	"errors"
	"fmt"
	"strings"

	bolt "go.etcd.io/bbolt"
)

// MediaCharacter represents a relationship between single
// instances of Media and Character.
type MediaCharacter struct {
	ID            int
	MediaID       int
	CharacterID   *int
	CharacterRole *string
	PersonID      *int
	PersonRole    *string
	Version       int
}

// Clean cleans the given MediaCharacter for storage
func (ser *MediaCharacterService) Clean(e *MediaCharacter) (err error) {
	if e.CharacterID != nil {
		*e.CharacterRole = strings.Trim(*e.CharacterRole, " ")
	}
	if e.PersonRole != nil {
		*e.PersonRole = strings.Trim(*e.PersonRole, " ")
	}
	return nil
}

// Validate returns an error if the MediaCharacter is
// not valid for the database.
func (ser *MediaCharacterService) Validate(e *MediaCharacter) (err error) {
	return ser.DB.View(func(tx *bolt.Tx) error {
		// Check if Media with ID specified in new MediaCharacter exists
		// Get Media bucket, exit if error
		mb, err := Bucket(MediaBucketName, tx)
		if err != nil {
			return err
		}
		_, err = get(e.MediaID, mb)
		if err != nil {
			return err
		}

		if e.CharacterID == nil && e.PersonID == nil {
			return fmt.Errorf("either character id or person id must be specified")
		}

		// Check if Character with ID specified in new MediaCharacter exists
		// CharacterID may be not specified
		if e.CharacterID != nil {
			if e.CharacterRole == nil {
				return errors.New("character role must not be nil if character id is specified")
			}
			// Get Character bucket, exit if error
			cb, err := Bucket(CharacterBucketName, tx)
			if err != nil {
				return err
			}
			_, err = get(*e.CharacterID, cb)
			if err != nil {
				return err
			}
		} else {
			if e.CharacterRole != nil {
				return fmt.Errorf("character role must be nil if character id is not specified")
			}
		}

		// Check if Person with ID specified in new MediaCharacter exists
		// PersonID may be not specified
		if e.PersonID != nil {
			if e.PersonRole == nil {
				return errors.New("person role must not be nil if person id is specified")
			}
			// Get Person bucket, exit if error
			pb, err := Bucket(PersonBucketName, tx)
			if err != nil {
				return err
			}
			_, err = get(*e.PersonID, pb)
			if err != nil {
				return err
			}
		} else {
			if e.PersonRole != nil {
				return fmt.Errorf("person role must be nil if person id is not specified")
			}
		}

		return nil
	})
}

// persistOldProperties maintains certain properties
// of the existing entity in updates
func (ser *MediaCharacterService) persistOldProperties(old *MediaCharacter, new *MediaCharacter) (err error) {
	new.Version = old.Version + 1
	return nil
}

// GetByMedia retrieves a list of instances of MediaCharacter
// with the given Media ID.
func (ser *MediaCharacterService) GetByMedia(mID int, db *bolt.DB) (list []MediaCharacter, err error) {
	return ser.GetFilter(func(mc *MediaCharacter) bool {
		return mc.MediaID == mID
	})
}

// GetByCharacter retrieves a list of instances of MediaCharacter
// with the given Character ID.
func (ser *MediaCharacterService) GetByCharacter(cID int, db *bolt.DB) (list []MediaCharacter, err error) {
	return ser.GetFilter(func(mc *MediaCharacter) bool {
		return *mc.CharacterID == cID
	})
}

// GetByPerson retrieves a list of instances of MediaCharacter
// with the given Person ID.
func (ser *MediaCharacterService) GetByPerson(pID int, db *bolt.DB) (list []MediaCharacter, err error) {
	return ser.GetFilter(func(mc *MediaCharacter) bool {
		return *mc.CharacterID == pID
	})
}
