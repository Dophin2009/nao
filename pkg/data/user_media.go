package data

import (
	"errors"
	"fmt"
	"time"

	json "github.com/json-iterator/go"
	"gitlab.com/Dophin2009/nao/pkg/db"
)

// UserMedia represents a relationship between a User and a Media, containing
// information about the User's opinion on the Media.
type UserMedia struct {
	UserID           int
	MediaID          int
	Status           *WatchStatus
	Priority         *int
	Score            *int
	Recommended      *int
	WatchedInstances []WatchedInstance
	Comments         []Title
	Meta             db.ModelMetadata
}

// Metadata returns Meta
func (um *UserMedia) Metadata() *db.ModelMetadata {
	return &um.Meta
}

// WatchedInstance contains information about a single watch of some Media.
type WatchedInstance struct {
	Episodes  int
	Ongoing   bool
	StartDate *time.Time
	EndDate   *time.Time
	Comments  map[string]string
}

// WatchStatus is an enum that represents the status of a Media's consumption
// by a User.
type WatchStatus int

const (
	// WatchStatusCompleted means that the User has consumed the Media in its entirety at
	// least once.
	WatchStatusCompleted WatchStatus = iota

	// WatchStatusPlanning means that the User is planning to consume the Media sometime in
	// the future.
	WatchStatusPlanning

	// WatchStatusDropped means that the User has never consumed the Media in its entirety
	// and abandoned it in the middle somewhere.
	WatchStatusDropped

	// WatchStatusHold means the User has begun consuming the Media, but has placed it on
	// hold.
	WatchStatusHold
)

// UnmarshalJSON defines custom JSON deserialization for WatchStatus.
func (ws *WatchStatus) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}

	value, ok := map[string]WatchStatus{
		"Completed": WatchStatusCompleted,
		"Planning":  WatchStatusPlanning,
		"Dropped":   WatchStatusDropped,
		"Hold":      WatchStatusHold,
	}[s]
	if !ok {
		return fmt.Errorf("watch status value %q: %w", s, errInvalid)
	}
	*ws = value
	return nil
}

// MarshalJSON defines custom JSON serialization for WatchStatus.
func (ws *WatchStatus) MarshalJSON() (v []byte, err error) {
	value, ok := map[WatchStatus]string{
		WatchStatusCompleted: "Completed",
		WatchStatusPlanning:  "Planning",
		WatchStatusDropped:   "Dropped",
		WatchStatusHold:      "Hold",
	}[*ws]
	if !ok {
		return nil, fmt.Errorf("watch status value %d: %w", *ws, err)
	}

	v, err = json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONMarshal, err)
	}

	return v, nil
}

// UserMediaService performs operations on UserMedia.
type UserMediaService struct {
	UserService  *UserService
	MediaService *MediaService
	Hooks        db.PersistHooks
}

// NewUserMediaService returns a UserMediaService.
func NewUserMediaService(hooks db.PersistHooks, userService *UserService,
	mediaService *MediaService) *UserMediaService {
	// Initialize UserMediaService
	userMediaService := &UserMediaService{
		UserService:  userService,
		MediaService: mediaService,
		Hooks:        hooks,
	}

	// Add hook to delete UserMedia on User deletion
	deleteUserMediaOnDeleteUser := func(um db.Model, _ db.Service, tx db.Tx) error {
		uID := um.Metadata().ID
		err := userMediaService.DeleteByUser(uID, tx)
		if err != nil {
			return fmt.Errorf("failed to delete UserMedia by User ID %d: %w",
				uID, err)
		}
		return nil
	}
	uSerHooks := userService.PersistHooks()
	uSerHooks.PreDeleteHooks =
		append(uSerHooks.PreDeleteHooks, deleteUserMediaOnDeleteUser)

	// Add hook to delete UserMedia on Media deletion
	deleteUserMediaOnDeleteMedia := func(mdm db.Model, _ db.Service, tx db.Tx) error {
		mID := mdm.Metadata().ID
		err := userMediaService.DeleteByMedia(mID, tx)
		if err != nil {
			return fmt.Errorf("failed to delete UserMedia by Media ID %d: %w",
				mID, err)
		}
		return nil
	}
	mdSerHooks := mediaService.PersistHooks()
	mdSerHooks.PreDeleteHooks =
		append(mdSerHooks.PreDeleteHooks, deleteUserMediaOnDeleteMedia)

	return userMediaService
}

// Create persists the given UserMedia.
func (ser *UserMediaService) Create(um *UserMedia, tx db.Tx) (int, error) {
	return tx.Database().Create(um, ser, tx)
}

// Update rumlaces the value of the UserMedia with the given ID.
func (ser *UserMediaService) Update(um *UserMedia, tx db.Tx) error {
	return tx.Database().Update(um, ser, tx)
}

// Delete deletes the UserMedia with the given ID.
func (ser *UserMediaService) Delete(id int, tx db.Tx) error {
	return tx.Database().Delete(id, ser, tx)
}

// DeleteByUser deletes the UserMedia with the given User ID.
func (ser *UserMediaService) DeleteByUser(uID int, tx db.Tx) error {
	return tx.Database().DeleteFilter(ser, tx, func(m db.Model) bool {
		um, err := ser.AssertType(m)
		if err != nil {
			return false
		}
		return um.UserID == uID
	})
}

// DeleteByMedia deletes the UserMedia with the given Media ID.
func (ser *UserMediaService) DeleteByMedia(mID int, tx db.Tx) error {
	return tx.Database().DeleteFilter(ser, tx, func(m db.Model) bool {
		um, err := ser.AssertType(m)
		if err != nil {
			return false
		}
		return um.MediaID == mID
	})
}

// GetAll retrieves all persisted values of UserMedia.
func (ser *UserMediaService) GetAll(first *int, skip *int, tx db.Tx) ([]*UserMedia, error) {
	vlist, err := tx.Database().GetAll(first, skip, ser, tx)
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to UserMedia: %w", err)
	}
	return list, nil
}

// GetFilter retrieves all persisted values of UserMedia that pass the filter.
func (ser *UserMediaService) GetFilter(
	first *int, skip *int, tx db.Tx, keep func(um *UserMedia) bool,
) ([]*UserMedia, error) {
	vlist, err := tx.Database().GetFilter(first, skip, ser, tx,
		func(m db.Model) bool {
			um, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(um)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to UserMedia: %w", err)
	}
	return list, nil
}

// GetMultiple retrieves the persisted UserMedia values specified by the
// given IDs that pass the filter.
func (ser *UserMediaService) GetMultiple(
	ids []int, tx db.Tx, keep func(um *UserMedia) bool,
) ([]*UserMedia, error) {
	vlist, err := tx.Database().GetMultiple(ids, ser, tx,
		func(m db.Model) bool {
			um, err := ser.AssertType(m)
			if err != nil {
				return false
			}
			return keep(um)
		})
	if err != nil {
		return nil, err
	}

	list, err := ser.mapFromModel(vlist)
	if err != nil {
		return nil, fmt.Errorf("failed to map db.Models to UserMedias: %w", err)
	}
	return list, nil
}

// GetByID retrieves the persisted UserMedia with the given ID.
func (ser *UserMediaService) GetByID(id int, tx db.Tx) (*UserMedia, error) {
	m, err := tx.Database().GetByID(id, ser, tx)
	if err != nil {
		return nil, err
	}

	um, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}
	return um, nil
}

// GetByUser retrieves the persisted UserMedia with the given User ID.
func (ser *UserMediaService) GetByUser(
	uID int, first *int, skip *int, tx db.Tx,
) ([]*UserMedia, error) {
	return ser.GetFilter(first, skip, tx, func(um *UserMedia) bool {
		return um.UserID == uID
	})
}

// GetByMedia retrieves the persisted UserMedia with the given Media ID.
func (ser *UserMediaService) GetByMedia(
	mID int, first *int, skip *int, tx db.Tx,
) ([]*UserMedia, error) {
	return ser.GetFilter(first, skip, tx, func(um *UserMedia) bool {
		return um.MediaID == mID
	})
}

// Bucket returns the name of the bucket for UserMedia.
func (ser *UserMediaService) Bucket() string {
	return "UserMedia"
}

// Clean cleans the given UserMedia for storage.
func (ser *UserMediaService) Clean(m db.Model, _ db.Tx) error {
	_, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s :%w", errmsgModelAssertType, err)
	}
	return nil
}

// Validate returns an error if the UserMedia is not valid for the database.
func (ser *UserMediaService) Validate(m db.Model, tx db.Tx) error {
	e, err := ser.AssertType(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	db := tx.Database()

	// Check if User with ID specified in UserMedia exists
	_, err = db.GetRawByID(e.UserID, ser.UserService, tx)
	if err != nil {
		return fmt.Errorf("failed to get User with ID %d: %w", e.UserID, err)
	}

	// Check if Media with ID specified in MediaCharacter exists
	_, err = db.GetRawByID(e.MediaID, ser.MediaService, tx)
	if err != nil {
		return fmt.Errorf("failed to get Media with ID %d: %w", e.MediaID, err)
	}

	return nil
}

// Initialize sets initial values for some properties.
func (ser *UserMediaService) Initialize(_ db.Model, _ db.Tx) error {
	return nil
}

// PersistOldProperties maintains certain properties of the existing UserMedia
// in updates.
func (ser *UserMediaService) PersistOldProperties(_ db.Model, _ db.Model, _ db.Tx) error {
	return nil
}

// PersistHooks returns the persistence hook functions.
func (ser *UserMediaService) PersistHooks() *db.PersistHooks {
	return &ser.Hooks
}

// Marshal transforms the given UserMedia into JSON.
func (ser *UserMediaService) Marshal(m db.Model) ([]byte, error) {
	um, err := ser.AssertType(m)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
	}

	v, err := json.Marshal(um)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONMarshal, err)
	}

	return v, nil
}

// Unmarshal parses the given JSON into UserMedia.
func (ser *UserMediaService) Unmarshal(buf []byte) (db.Model, error) {
	var um UserMedia
	err := json.Unmarshal(buf, &um)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgJSONUnmarshal, err)
	}
	return &um, nil
}

// AssertType exposes the given db.Model as a UserMedia.
func (ser *UserMediaService) AssertType(m db.Model) (*UserMedia, error) {
	if m == nil {
		return nil, fmt.Errorf("model: %w", errNil)
	}

	um, ok := m.(*UserMedia)
	if !ok {
		return nil, fmt.Errorf("model: %w", errors.New("not of UserMedia type"))
	}
	return um, nil
}

// mapfromModel returns a list of UserMedia type asserted from the given list
// of db.Model.
func (ser *UserMediaService) mapFromModel(vlist []db.Model) ([]*UserMedia, error) {
	list := make([]*UserMedia, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.AssertType(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelAssertType, err)
		}
	}
	return list, nil
}
