package data

import (
	"encoding/json"
	"fmt"
	"strings"

	bolt "go.etcd.io/bbolt"
)

// MediaCharacter represents a relationship between single
// instances of Media and Producer
type MediaCharacter struct {
	ID            int
	MediaID       int
	CharacterID   int
	CharacterRole string
	PersonID      int
	PersonRole    string
	Version       int
}

const mediaCharacterBucketName = "MediaCharacter"

// MediaCharacterGet retrieves a single instance of MediaCharacter with
// the given ID
func MediaCharacterGet(ID int, db *bolt.DB) (mc MediaCharacter, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get MediaCharacter bucket, exit if error
		b, err := bucket(mediaCharacterBucketName, tx)
		if err != nil {
			return err
		}

		// Get MediaCharacter by ID, exit if error
		v, err := get(ID, b)
		if err != nil {
			return err
		}
		return json.Unmarshal(v, &mc)
	})

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
// with the given Producer ID
func MediaCharacterGetByCharacter(cID int, db *bolt.DB) (list []MediaCharacter, err error) {
	return MediaCharacterGetFilter(db, func(mc *MediaCharacter) bool {
		return mc.CharacterID == cID
	})
}

// MediaCharacterGetByVoiceActor retrieves a list of instances of MediaCharacter
// with the given Voice Actor ID
func MediaCharacterGetByVoiceActor(pID int, db *bolt.DB) (list []MediaCharacter, err error) {
	return MediaCharacterGetFilter(db, func(mc *MediaCharacter) bool {
		return mc.CharacterID == pID
	})
}

// MediaCharacterGetFilter retrieves all persisted MediaCharacter values
func MediaCharacterGetFilter(db *bolt.DB, filter func(mc *MediaCharacter) bool) (list []MediaCharacter, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get MediaCharacter bucket, exit if error
		b, err := bucket(mediaCharacterBucketName, tx)
		if err != nil {
			return err
		}

		// Unmarshal and add all MediaCharacters to slice,
		// exit if error
		b.ForEach(func(k, v []byte) error {
			mc := MediaCharacter{}
			err = json.Unmarshal(v, &mc)
			if err != nil {
				return err
			}

			if filter(&mc) {
				list = append(list, mc)
			}
			return err
		})

		return nil
	})

	return
}

// MediaCharacterCreate persists a new instance of MediaCharacter
func MediaCharacterCreate(mc *MediaCharacter, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Get MediaCharacter bucket, exit if error
		b, err := bucket(mediaCharacterBucketName, tx)
		if err != nil {
			return err
		}

		// Check if MediaCharacter properties are valid
		err = MediaCharacterCheckRelatedIDs(mc, tx)
		if err != nil {
			return err
		}

		// Get next ID in sequence and
		// assign to MediaCharacter
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		mc.ID = int(id)

		// Save MediaCharacter in bucket
		buf, err := json.Marshal(mc)
		if err != nil {
			return err
		}

		return b.Put(itob(mc.ID), buf)
	})
}

// MediaCharacterUpdate updates the properties of an existing
// persisted Producer instance
func MediaCharacterUpdate(mc *MediaCharacter, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Get MediaCharacter bucket, exit if error
		b, err := bucket(mediaCharacterBucketName, tx)
		if err != nil {
			return err
		}

		// Check if MediaCharacter with ID exists
		o, err := get(mc.ID, b)
		if err != nil {
			return err
		}

		// Check if MediaCharacter properties are valid
		err = MediaCharacterCheckRelatedIDs(mc, tx)
		if err != nil {
			return err
		}

		// Replace properties of new with immutable
		// ones of old
		old := MediaCharacter{}
		err = json.Unmarshal([]byte(o), &old)
		// Update version
		mc.Version = old.Version + 1

		// Save MediaCharacter
		buf, err := json.Marshal(mc)
		if err != nil {
			return err
		}

		return b.Put(itob(mc.ID), buf)
	})
}

// MediaCharacterCheckRelatedIDs checks if the entities specified
// by the related entity IDs exist for a MediaCharacter
func MediaCharacterCheckRelatedIDs(mc *MediaCharacter, tx *bolt.Tx) (err error) {
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
