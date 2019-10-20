package data

import (
	"fmt"
	"strings"

	bolt "go.etcd.io/bbolt"
)

// MediaCharacter represents a relationship between single
// instances of Media and Character.
type MediaCharacter struct {
	ID            int
	MediaID       int
	CharacterID   int
	CharacterRole string
	PersonID      int
	PersonRole    string
	Version       int
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

		if e.CharacterID == 0 && e.PersonID == 0 {
			return fmt.Errorf("either character id or person id must be specified")
		}

		// Check if Character with ID specified in new MediaCharacter exists
		// CharacterID may be not specified
		if e.CharacterID != 0 {
			// Get Character bucket, exit if error
			cb, err := Bucket(CharacterBucketName, tx)
			if err != nil {
				return err
			}
			_, err = get(e.CharacterID, cb)
			if err != nil {
				return err
			}
		} else {
			if strings.Trim(e.CharacterRole, " ") != "" {
				return fmt.Errorf("character role must be blank if character id is not specified")
			}
		}

		// Check if Person with ID specified in new MediaCharacter exists
		// PersonID may be not specified
		if e.PersonID != 0 {
			// Get Person bucket, exit if error
			pb, err := Bucket(PersonBucketName, tx)
			if err != nil {
				return err
			}
			_, err = get(e.PersonID, pb)
			if err != nil {
				return err
			}
		} else {
			if strings.Trim(e.PersonRole, " ") != "" {
				return fmt.Errorf("person role must be blank if person id is not specified")
			}
		}

		return nil
	})
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
		return mc.CharacterID == cID
	})
}

// GetByPerson retrieves a list of instances of MediaCharacter
// with the given Person ID.
func (ser *MediaCharacterService) GetByPerson(pID int, db *bolt.DB) (list []MediaCharacter, err error) {
	return ser.GetFilter(func(mc *MediaCharacter) bool {
		return mc.CharacterID == pID
	})
}
