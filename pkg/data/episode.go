package data

import (
	"errors"
	"fmt"
	"time"

	json "github.com/json-iterator/go"
	"gitlab.com/Dophin2009/nao/pkg/db"
)

// TODO: User rating/comments/etc. of Episodes

// Episode represents a single episode or chapter for some media.
type Episode struct {
	Titles   []Title
	Synopses []Title
	Date     *time.Time
	Duration *int
	Filler   bool
	Recap    bool
	Meta     db.ModelMetadata
}

// Metadata returns Meta.
func (ep *Episode) Metadata() *db.ModelMetadata {
	return &ep.Meta
}

// EpisodeSet is an ordered list of episodes.
type EpisodeSet struct {
	MediaID      int
	Descriptions []Title
	Episodes     []int
	Meta         db.ModelMetadata
}

// Metadata returns the Meta.
func (set *EpisodeSet) Metadata() *db.ModelMetadata {
	return &set.Meta
}

// EpisodeService performs operations on Episodes.
type EpisodeService struct {
	Hooks db.PersistHooks
}

// NewEpisodeService returns a EpisodeService.
func NewEpisodeService(hooks db.PersistHooks) *EpisodeService {
	return &EpisodeService{
		Hooks: hooks,
	}
}

// Create persists the given Episode.
func (ser *EpisodeService) Create(ep *Episode, tx db.Tx) (int, error) {
	return tx.Database().Create(ep, ser, tx)
}

// Update replaces the value of the Episode with the given ID.
func (ser *EpisodeService) Update(ep *Episode, tx db.Tx) error {
	return tx.Database().Update(ep, ser, tx)
}

// Delete deletes the Episode with the given ID.
func (ser *EpisodeService) Delete(id int, tx db.Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// GetAll retrieves all persisted values of Episode.
func (ser *EpisodeService) GetAll(first *int, skip *int, tx db.Tx) ([]*Episode, error) {
	vlist, err := tx.Database().GetAll(first, skip, ser, tx)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to Episodes: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of Episode that pass the filter.
func (ser *EpisodeService) GetFilter(
	first *int, skip *int, tx db.Tx, keep func(ep *Episode) bool,
) ([]*Episode, error) {
	vlist, err := tx.Database().GetFilter(first, skip, ser, tx,
		func(m db.Model) bool {
			ep, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(ep)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to Episodes: %w", err)
	}
	return list, nil
}

// GetMultiple retrieves the persisted Episode values specified by the given
// IDs that pass the filter.
func (ser *EpisodeService) GetMultiple(
	ids []int, tx db.Tx, keep func(ep *Episode) bool,
) ([]*Episode, error) {
	vlist, err := tx.Database().GetMultiple(ids, ser, tx,
		func(m db.Model) bool {
			ep, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(ep)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to Episodes: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted Episode with the given ID.
func (ser *EpisodeService) GetByID(id int, tx db.Tx) (*Episode, error) {
	m, err := tx.Database().GetByID(id, ser, tx)
	if err != nil {
		return nil, err
	}

	ep, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return ep, nil
}

// Bucket returns the name of the bucket for Episode.
func (ser *EpisodeService) Bucket() string {
	return "Episode"
}

// Clean cleans the given Episode for storage.
func (ser *EpisodeService) Clean(m db.Model, _ db.Tx) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Validate returns an error if the Episode is not valid for the database.
func (ser *EpisodeService) Validate(m db.Model, _ db.Tx) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Initialize sets initial values for some properties.
func (ser *EpisodeService) Initialize(_ db.Model, _ db.Tx) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing Episode in
// updates.
func (ser *EpisodeService) PersistOldProperties(_ db.Model, _ db.Model, _ db.Tx) error {
	return nil
}

// PersistHooks returns the persistence hook functions.
func (ser *EpisodeService) PersistHooks() *db.PersistHooks {
	return &ser.Hooks
}

// Marshal transforms the given Episode into JSON.
func (ser *EpisodeService) Marshal(m db.Model) ([]byte, error) {
	ep, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	v, err := json.Marshal(ep)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONMarshal, err)
	}

	return v, nil
}

// Unmarshal parses the given JSON into Episode.
func (ser *EpisodeService) Unmarshal(buf []byte) (db.Model, error) {
	var ep Episode
	err := json.Unmarshal(buf, &ep)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &ep, nil
}

// AssertType exposes the Model as an Episode.
func (ser *EpisodeService) AssertType(m db.Model) (*Episode, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	ep, ok := m.(*Episode)
	if !ok {
		return nil, fmt.Errorf("model: %w", errors.New("not of Episode type"))
	}
	return ep, nil
}

// mapfromModel returns a list of Episode type asserted from the given list of
// Model.
func (ser *EpisodeService) mapFromModel(vlist []db.Model) ([]*Episode, error) {
	list := make([]*Episode, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}

// EpisodeSetService performs operations on EpisodeSets.
type EpisodeSetService struct {
	EpisodeService *EpisodeService
	MediaService   *MediaService
	Hooks          db.PersistHooks
}

// NewEpisodeSetService returns an EpisodeSetService.
func NewEpisodeSetService(hooks db.PersistHooks, episodeService *EpisodeService,
	mediaService *MediaService) *EpisodeSetService {
	episodeSetService := &EpisodeSetService{
		EpisodeService: episodeService,
		MediaService:   mediaService,
		Hooks:          hooks,
	}

	// Add hook to update EpisodeSets' list of Episode IDs on Episode deletion
	updateEpisodeSetOnDeleteEpisode := func(epm db.Model, _ db.Service, tx db.Tx) error {
		epID := epm.Metadata().ID
		err := tx.Database().DoEach(nil, nil, episodeSetService, tx,
			func(m db.Model, _ db.Service, tx db.Tx) (exit bool, err error) {
				set, err := episodeSetService.AssertType(m)
				if err != nil {
					return true, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
				}

				// Find ID of Episode to be deleted in the list
				rmID := -1
				for _, id := range set.Episodes {
					if id == epID {
						rmID = id
						break
					}
				}
				// Episode ID not found, move onto next EpisodeSet
				if rmID < 0 {
					return false, nil
				}

				// Remove ID from Episodes
				set.Episodes = append(set.Episodes[:rmID], set.Episodes[rmID+1:]...)

				// Update persisted value
				err = tx.Database().Update(set, episodeSetService, tx)
				if err != nil {
					return true, fmt.Errorf("failed to update EpisodeSet: %w", err)
				}
				return false, nil
			}, nil,
		)
		if err != nil {
			return err
		}
		return nil
	}
	epSerHooks := episodeService.PersistHooks()
	epSerHooks.PreDeleteHooks =
		append(epSerHooks.PreDeleteHooks, updateEpisodeSetOnDeleteEpisode)

	deleteEpisodeSetOnDeleteMedia := func(mdm db.Model, ser db.Service, tx db.Tx) error {
		mID := mdm.Metadata().ID
		err := episodeSetService.DeleteByMedia(mID, tx)
		if err != nil {
			return fmt.Errorf("failed to delete EpisodeSets by Media ID %d: %w",
				mID, err)
		}
		return nil
	}
	mdSerHooks := mediaService.PersistHooks()
	mdSerHooks.PreDeleteHooks =
		append(mdSerHooks.PreDeleteHooks, deleteEpisodeSetOnDeleteMedia)

	return episodeSetService
}

// Create persists the given EpisodeSet.
func (ser *EpisodeSetService) Create(set *EpisodeSet, tx db.Tx) (int, error) {
	return tx.Database().Create(set, ser, tx)
}

// Update replaces the value of the EpisodeSet with the given ID.
func (ser *EpisodeSetService) Update(set *EpisodeSet, tx db.Tx) error {
	return tx.Database().Update(set, ser, tx)
}

// Delete deletes the EpisodeSet with the given ID.
func (ser *EpisodeSetService) Delete(id int, tx db.Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// DeleteByEpisode deletes the EpisodeSets who contain the Episode with the
// given ID.
func (ser *EpisodeSetService) DeleteByEpisode(epID int, tx db.Tx) error {
	return tx.Database().DeleteFilter(ser, tx, func(m db.Model) bool {
		set, err := ser.AssertType(m)
		if err != nil {
			return false
		}

		for _, id := range set.Episodes {
			if id == epID {
				return true
			}
		}

		return false
	})
}

// DeleteByMedia deletes the EpisodeSets with the given Media ID.
func (ser *EpisodeSetService) DeleteByMedia(mID int, tx db.Tx) error {
	return tx.Database().DeleteFilter(ser, tx, func(m db.Model) bool {
		set, err := ser.AssertType(m)
		if err != nil {
			return false
		}
		return set.MediaID == mID
	})
}

// GetAll retrieves all persisted values of EpisodeSet.
func (ser *EpisodeSetService) GetAll(first *int, skip *int, tx db.Tx) ([]*EpisodeSet, error) {
	vlist, err := tx.Database().GetAll(first, skip, ser, tx)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to EpisodeSets: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of EpisodeSet that pass the filter.
func (ser *EpisodeSetService) GetFilter(
	first *int, skip *int, tx db.Tx, keep func(*EpisodeSet) bool,
) ([]*EpisodeSet, error) {
	vlist, err := tx.Database().GetFilter(first, skip, ser, tx,
		func(m db.Model) bool {
			set, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(set)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to EpisodesSets: %w", err)
	}
	return list, nil
}

// GetMultiple retrieves the persisted EpisodeSet values specified by the given
// IDs that pass the filter.
func (ser *EpisodeSetService) GetMultiple(
	ids []int, tx db.Tx, keep func(set *EpisodeSet) bool,
) ([]*EpisodeSet, error) {
	vlist, err := tx.Database().GetMultiple(ids, ser, tx,
		func(m db.Model) bool {
			set, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(set)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map Models to EpisodeSets: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted EpisodeSet with the given ID.
func (ser *EpisodeSetService) GetByID(id int, tx db.Tx) (*EpisodeSet, error) {
	m, err := tx.Database().GetByID(id, ser, tx)
	if err != nil {
		return nil, err
	}

	set, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return set, nil
}

// GetByMedia retrieves a list of instances of EpisodeSet with the given Media
// ID.
func (ser *EpisodeSetService) GetByMedia(
	mID int, first *int, skip *int, tx db.Tx,
) ([]*EpisodeSet, error) {
	return ser.GetFilter(first, skip, tx, func(set *EpisodeSet) bool {
		return set.MediaID == mID
	})
}

// Bucket returns the name of the bucket for EpisodeSet.
func (ser *EpisodeSetService) Bucket() string {
	return "EpisodeSet"
}

// Clean cleans the given EpisodeSet for storage.
func (ser *EpisodeSetService) Clean(m db.Model, _ db.Tx) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return nil
}

// Validate returns an error if the EpisodeSet is not valid for the database.
func (ser *EpisodeSetService) Validate(m db.Model, tx db.Tx) error {
	set, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	for _, id := range set.Episodes {
		_, err := tx.Database().GetRawByID(id, ser, tx)
		if err != nil {
			return fmt.Errorf("failed to get Episode with ID %d: %w", id, err)
		}
	}
	return nil
}

// Initialize sets initial values for some properties.
func (ser *EpisodeSetService) Initialize(_ db.Model, _ db.Tx) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing EpisodeSet
// in updates.
func (ser *EpisodeSetService) PersistOldProperties(_ db.Model, _ db.Model, _ db.Tx) error {
	return nil
}

// PersistHooks returns the persistence hook functions.
func (ser *EpisodeSetService) PersistHooks() *db.PersistHooks {
	return &ser.Hooks
}

// Marshal transforms the given EpisodeSet into JSON.
func (ser *EpisodeSetService) Marshal(m db.Model) ([]byte, error) {
	set, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	v, err := json.Marshal(set)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONMarshal, err)
	}

	return v, nil
}

// Unmarshal parses the given JSON into EpisodeSet.
func (ser *EpisodeSetService) Unmarshal(buf []byte) (db.Model, error) {
	var set EpisodeSet
	err := json.Unmarshal(buf, &set)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &set, nil
}

// AssertType exposes the Model as an EpisodeSet.
func (ser *EpisodeSetService) AssertType(m db.Model) (*EpisodeSet, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	set, ok := m.(*EpisodeSet)
	if !ok {
		return nil, fmt.Errorf("model: %w", errors.New("not of EpisodeSet type"))
	}
	return set, nil
}

// mapfromModel returns a list of EpisodeSet type asserted from the given list
// of Model.
func (ser *EpisodeSetService) mapFromModel(vlist []db.Model) ([]*EpisodeSet, error) {
	list := make([]*EpisodeSet, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}
