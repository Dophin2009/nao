package data

import (
	"fmt"
	"strings"

	bolt "go.etcd.io/bbolt"
)

// MediaCharacter represents a relationship between single
// instances of Media and Character
type MediaCharacter struct {
	ID            int
	MediaID       int
	CharacterID   int
	CharacterRole string
	PersonID      int
	PersonRole    string
	Version       int
}

// Identifier returns the ID of the MediaCharacter
func (mc *MediaCharacter) Identifier() int {
	return mc.ID
}

// SetIdentifier sets the ID of the MediaCharacter
func (mc *MediaCharacter) SetIdentifier(ID int) {
	mc.ID = ID
}

// Ver returns the verison of the MediaCharacter
func (mc *MediaCharacter) Ver() int {
	return mc.Version
}

// UpdateVer increments the version of the
// MediaCharacter by one
func (mc *MediaCharacter) UpdateVer() {
	mc.Version++
}

// Validate returns an error if the MediaCharacter is
// not valid for the database
func (mc *MediaCharacter) Validate(tx *bolt.Tx) (err error) {
	// Check if Media with ID specified in new MediaCharacter exists
	// Get Media bucket, exit if error
	mb, err := bucket(mediaBucketName, tx)
	if err != nil {
		return err
	}
	_, err = get(mc.MediaID, mb)
	if err != nil {
		return err
	}

	if mc.CharacterID == 0 && mc.PersonID == 0 {
		return fmt.Errorf("either character id or person id must be specified")
	}

	// Check if Character with ID specified in new MediaCharacter exists
	// CharacterID may be not specified
	if mc.CharacterID != 0 {
		// Get Character bucket, exit if error
		cb, err := bucket(characterBucketName, tx)
		if err != nil {
			return err
		}
		_, err = get(mc.CharacterID, cb)
		if err != nil {
			return err
		}
	} else {
		if strings.Trim(mc.CharacterRole, " ") != "" {
			return fmt.Errorf("character role must be blank if character id is not specified")
		}
	}

	// Check if Person with ID specified in new MediaCharacter exists
	// PersonID may be not specified
	if mc.PersonID != 0 {
		// Get Person bucket, exit if error
		pb, err := bucket(personBucketName, tx)
		if err != nil {
			return err
		}
		_, err = get(mc.PersonID, pb)
		if err != nil {
			return err
		}
	} else {
		if strings.Trim(mc.PersonRole, " ") != "" {
			return fmt.Errorf("person role must be blank if person id is not specified")
		}
	}

	return nil
}

const mediaCharacterBucketName = "MediaCharacter"

// MediaCharacterGet retrieves a single instance of MediaCharacter with
// the given ID
func MediaCharacterGet(ID int, db *bolt.DB) (mc MediaCharacter, err error) {
	err = getByID(ID, &mc, mediaCharacterBucketName, db)
	return
}

// MediaCharacterGetAll retrieves all persisted MediaCharacter values
func MediaCharacterGetAll(db *bolt.DB) (list []MediaCharacter, err error) {
	return MediaCharacterGetFilter(db, func(mc *MediaCharacter) bool { return true })
}

// MediaCharacterGetByMedia retrieves a list of instances of MediaCharacter
// with the given Media ID
func MediaCharacterGetByMedia(mID int, db *bolt.DB) (list []MediaCharacter, err error) {
	return MediaCharacterGetFilter(db, func(mc *MediaCharacter) bool {
		return mc.MediaID == mID
	})
}

// MediaCharacterGetByCharacter retrieves a list of instances of MediaCharacter
// with the given Character ID
func MediaCharacterGetByCharacter(cID int, db *bolt.DB) (list []MediaCharacter, err error) {
	return MediaCharacterGetFilter(db, func(mc *MediaCharacter) bool {
		return mc.CharacterID == cID
	})
}

// MediaCharacterGetByPerson retrieves a list of instances of MediaCharacter
// with the given Person ID
func MediaCharacterGetByPerson(pID int, db *bolt.DB) (list []MediaCharacter, err error) {
	return MediaCharacterGetFilter(db, func(mc *MediaCharacter) bool {
		return mc.CharacterID == pID
	})
}

// MediaCharacterGetFilter retrieves all persisted MediaCharacter values
func MediaCharacterGetFilter(db *bolt.DB, filter func(mc *MediaCharacter) bool) (list []MediaCharacter, err error) {
	ilist, err := getFilter(&Media{}, func(entity Idenitifiable) (bool, error) {
		mc, ok := entity.(*MediaCharacter)
		if !ok {
			return false, fmt.Errorf("type assertion failed: entity is not a MediaCharacter")
		}
		return filter(mc), nil
	}, mediaCharacterBucketName, db)

	list = make([]MediaCharacter, len(ilist))
	for i, mc := range ilist {
		list[i] = *mc.(*MediaCharacter)
	}

	return
}

// MediaCharacterCreate persists a new instance of MediaCharacter
func MediaCharacterCreate(mc *MediaCharacter, db *bolt.DB) error {
	return create(mc, mediaCharacterBucketName, db)
}

// MediaCharacterUpdate updates the properties of an existing
// persisted Producer instance
func MediaCharacterUpdate(mc *MediaCharacter, db *bolt.DB) error {
	return update(mc, mediaCharacterBucketName, db)
}
