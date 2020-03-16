package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Dophin2009/nao/pkg/data/models"
	"github.com/Dophin2009/nao/pkg/db"
	json "github.com/json-iterator/go"
)

// MediaCharacterService performs operations on MediaCharacter.
type MediaCharacterService struct {
	MediaService     *MediaService
	CharacterService *CharacterService
	PersonService    *PersonService
	Hooks            db.PersistHooks
}

// NewMediaCharacterService returns a MediaCharacterService.
func NewMediaCharacterService(hooks db.PersistHooks, mediaService *MediaService,
	characterService *CharacterService, personService *PersonService) *MediaCharacterService {
	// Initialize MediaCharacterService
	mediaCharacterService := &MediaCharacterService{
		MediaService:     mediaService,
		CharacterService: characterService,
		PersonService:    personService,
	}

	// Add hook to delete MediaCharacter on Media deletion
	deleteMediaCharacterOnDeleteMedia := func(mdm db.Model, _ db.Service, tx db.Tx) error {
		mID := mdm.Metadata().ID
		err := mediaCharacterService.DeleteByMedia(mID, tx)
		if err != nil {
			return fmt.Errorf("failed to delete MediaCharacter by Media ID %d: %w",
				mID, err)
		}
		return nil
	}
	mdSerHooks := mediaService.PersistHooks()
	mdSerHooks.PreDeleteHooks =
		append(mdSerHooks.PreDeleteHooks, deleteMediaCharacterOnDeleteMedia)

	// Add hook to delete MediaCharacter on Character deletion
	deleteMediaCharacterOnDeleteCharacter := func(cm db.Model, _ db.Service, tx db.Tx) error {
		cID := cm.Metadata().ID
		err := mediaCharacterService.DeleteByCharacter(cID, tx)
		if err != nil {
			return fmt.Errorf("failed to delete MediaCharacter by Character ID %d: %w",
				cID, err)
		}
		return nil
	}
	cSerHooks := characterService.PersistHooks()
	cSerHooks.PreDeleteHooks =
		append(cSerHooks.PreDeleteHooks, deleteMediaCharacterOnDeleteCharacter)

	// Add hook to delete MediaCharaccter on Person deletion
	deleteMediaCharacterOnDeletePerson := func(pm db.Model, _ db.Service, tx db.Tx) error {
		pID := pm.Metadata().ID
		err := mediaCharacterService.DeleteByPerson(pID, tx)
		if err != nil {
			return fmt.Errorf("failed to delete MediaCharacter by Person ID %d: %w",
				pID, err)
		}
		return nil
	}
	pSerHooks := personService.PersistHooks()
	pSerHooks.PreDeleteHooks =
		append(pSerHooks.PreDeleteHooks, deleteMediaCharacterOnDeletePerson)

	return mediaCharacterService
}

// Create persists the given MediaCharacter.
func (ser *MediaCharacterService) Create(mc *models.MediaCharacter, tx db.Tx) (int, error) {
	return tx.Database().Create(mc, ser, tx)
}

// Update rmclaces the value of the MediaCharacter with the given ID.
func (ser *MediaCharacterService) Update(mc *models.MediaCharacter, tx db.Tx) error {
	return tx.Database().Update(mc, ser, tx)
}

// Delete deletes the MediaCharacter with the given ID.
func (ser *MediaCharacterService) Delete(id int, tx db.Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// DeleteByMedia deletes the MediaCharacters with the given Media ID.
func (ser *MediaCharacterService) DeleteByMedia(mID int, tx db.Tx) error {
	return tx.Database().DeleteFilter(ser, tx, func(m db.Model) bool {
		mc, err := ser.AssertType(m)
		if err != nil {
			return false
		}

		return mc.MediaID == mID
	})
}

// DeleteByCharacter deletes the MediaCharacters with the given Character ID.
func (ser *MediaCharacterService) DeleteByCharacter(cID int, tx db.Tx) error {
	return tx.Database().DeleteFilter(ser, tx, func(m db.Model) bool {
		mc, err := ser.AssertType(m)
		if err != nil {
			return false
		}

		if mc.CharacterID == nil {
			return false
		}

		return *mc.CharacterID == cID
	})
}

// DeleteByPerson deletes the MediaCharacters with the given Person ID.
func (ser *MediaCharacterService) DeleteByPerson(pID int, tx db.Tx) error {
	return tx.Database().DeleteFilter(ser, tx, func(m db.Model) bool {
		mc, err := ser.AssertType(m)
		if err != nil {
			return false
		}

		if mc.PersonID == nil {
			return false
		}

		return *mc.PersonID == pID
	})
}

// GetAll retrieves all persisted values of MediaCharacter.
func (ser *MediaCharacterService) GetAll(first *int, skip *int, tx db.Tx) ([]*models.MediaCharacter, error) {
	vlist, err := tx.Database().GetAll(first, skip, ser, tx)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to MediaCharacters: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of MediaCharacter that pass the
// filter.
func (ser *MediaCharacterService) GetFilter(
	first *int, skip *int, tx db.Tx, keep func(mc *models.MediaCharacter) bool,
) ([]*models.MediaCharacter, error) {
	vlist, err := tx.Database().GetFilter(first, skip, ser, tx,
		func(m db.Model) bool {
			mc, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(mc)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to MediaCharacters: %w", err)
	}
	return list, nil
}

// GetMultiple retrieves the persisted MediaCharacter values specified by the
// given IDs that pass the filter.
func (ser *MediaCharacterService) GetMultiple(
	ids []int, tx db.Tx, keep func(mc *models.MediaCharacter) bool,
) ([]*models.MediaCharacter, error) {
	vlist, err := tx.Database().GetMultiple(ids, ser, tx,
		func(m db.Model) bool {
			mc, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(mc)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to MediaCharacters: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted MediaCharacter with the given ID.
func (ser *MediaCharacterService) GetByID(id int, tx db.Tx) (*models.MediaCharacter, error) {
	m, err := tx.Database().GetByID(id, ser, tx)
	if err != nil {
		return nil, err
	}

	mc, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return mc, nil
}

// GetByMedia retrieves a list of instances of MediaCharacter with the given
// Media ID.
func (ser *MediaCharacterService) GetByMedia(
	mID int, first *int, skip *int, tx db.Tx,
) ([]*models.MediaCharacter, error) {
	return ser.GetFilter(first, skip, tx, func(mc *models.MediaCharacter) bool {
		return mc.MediaID == mID
	})
}

// GetByCharacter retrieves a list of instances of MediaCharacter with the
// given Character ID.
func (ser *MediaCharacterService) GetByCharacter(
	cID int, first *int, skip *int, tx db.Tx,
) ([]*models.MediaCharacter, error) {
	return ser.GetFilter(first, skip, tx, func(mc *models.MediaCharacter) bool {
		return *mc.CharacterID == cID
	})
}

// GetByPerson retrieves a list of instances of MediaCharacter with the given
// Person ID.
func (ser *MediaCharacterService) GetByPerson(
	pID int, first *int, skip *int, tx db.Tx,
) ([]*models.MediaCharacter, error) {
	return ser.GetFilter(first, skip, tx, func(mc *models.MediaCharacter) bool {
		return *mc.CharacterID == pID
	})
}

// Bucket returns the name of the bucket for MediaCharacter.
func (ser *MediaCharacterService) Bucket() string {
	return "MediaCharacter"
}

// Clean cleans the given MediaCharacter for storage.
func (ser *MediaCharacterService) Clean(m db.Model, _ db.Tx) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	if e.CharacterID != nil {
		*e.CharacterRole = strings.Trim(*e.CharacterRole, " ")
	}
	if e.PersonRole != nil {
		*e.PersonRole = strings.Trim(*e.PersonRole, " ")
	}
	return nil
}

// Validate returns an error if the MediaCharacter is not valid for the
// database.
func (ser *MediaCharacterService) Validate(m db.Model, tx db.Tx) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	db := tx.Database()

	// Check if Media with ID specified in MediaCharacter exists
	_, err = db.GetRawByID(e.MediaID, ser.MediaService, tx)
	if err != nil {
		return fmt.Errorf("failed to get Media with ID %d: %w", e.MediaID, err)
	}

	// Invalid if both Character and Person are not specified
	if e.CharacterID == nil && e.PersonID == nil {
		nsterr := fmt.Errorf("character ID and person ID: %w", errNil)
		return fmt.Errorf(
			"either character ID or person ID must be specified: %w", nsterr)
	}

	// Check if Character with ID specified in new MediaCharacter exists
	// CharacterID might be not specified
	if e.CharacterID != nil {
		// CharacterRole must be present if CharacterID is specified
		if e.CharacterRole == nil {
			nsterr := fmt.Errorf("character role: %w", errNil)
			return fmt.Errorf(
				"character role must not be nil if character ID is specified: %w",
				nsterr,
			)
		}

		cID := *e.CharacterID
		_, err = db.GetRawByID(cID, ser.CharacterService, tx)
		if err != nil {
			return fmt.Errorf("failed to get Character with ID %d: %w", cID, err)
		}
	} else {
		// CharacterRole must not be specified if CharacterID is not
		if e.CharacterRole != nil {
			nsterr := fmt.Errorf("character ID: %w", errNil)
			return fmt.Errorf(
				"character role must be nil if character ID is not specified: %w",
				nsterr,
			)
		}
	}

	// Check if Person with ID specified in new MediaCharacter exists
	// PersonID may be not specified
	if e.PersonID != nil {
		// PersonRole must be present if PersonID is specified
		if e.PersonRole == nil {
			nsterr := fmt.Errorf("person role: %w", errNil)
			return fmt.Errorf(
				"person role must not be nil if person ID is specified: %w", nsterr)
		}

		pID := *e.PersonID
		_, err = db.GetRawByID(pID, ser.PersonService, tx)
		if err != nil {
			return fmt.Errorf("failed to get Person with ID %d: %w", pID, err)
		}
	} else {
		// PersonRole must not be specified if PersonID is not
		if e.PersonRole != nil {
			nsterr := fmt.Errorf("person ID: %w", errNil)
			return fmt.Errorf(
				"person role must be nil if person ID is not specified: %w", nsterr)
		}
	}

	return nil
}

// Initialize sets initial values for some properties.
func (ser *MediaCharacterService) Initialize(_ db.Model, _ db.Tx) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing
// MediaCharacter in updates.
func (ser *MediaCharacterService) PersistOldProperties(_ db.Model, _ db.Model, _ db.Tx) error {
	return nil
}

// PersistHooks returns the persistence hook functions.
func (ser *MediaCharacterService) PersistHooks() *db.PersistHooks {
	return &ser.Hooks
}

// Marshal transforms the given MediaCharacter into JSON.
func (ser *MediaCharacterService) Marshal(m db.Model) ([]byte, error) {
	mc, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	v, err := json.Marshal(mc)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONMarshal, err)
	}

	return v, nil
}

// Unmarshal parses the given JSON into MediaCharacter.
func (ser *MediaCharacterService) Unmarshal(buf []byte) (db.Model, error) {
	var mc models.MediaCharacter
	err := json.Unmarshal(buf, &mc)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &mc, nil
}

// AssertType exposes the given db.Model as a MediaCharacter.
func (ser *MediaCharacterService) AssertType(m db.Model) (*models.MediaCharacter, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	mc, ok := m.(*models.MediaCharacter)
	if !ok {
		return nil, fmt.Errorf("model: %w", errors.New("not of MediaCharacter type"))
	}
	return mc, nil
}

// mapfromModel returns a list of MediaCharacter type asserted from the given
// list of db.Model.
func (ser *MediaCharacterService) mapFromModel(vlist []db.Model) ([]*models.MediaCharacter, error) {
	list := make([]*models.MediaCharacter, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}
