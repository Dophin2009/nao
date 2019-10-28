package data

import bolt "go.etcd.io/bbolt"

// UserMediaList represents a User-created list
// of UserMedia
type UserMediaList struct {
	ID           int
	UserID       int
	Names        []Info
	Descriptions []Info
	Version      int
}

// Clean cleans the given UserMediaList for storage
func (ser *UserMediaListService) Clean(e *UserMediaList) (err error) {
	if err := infoListClean(e.Names); err != nil {
		return err
	}
	if err := infoListClean(e.Descriptions); err != nil {
		return err
	}
	return nil
}

// Validate returns an error if the UserMediaList is
// not valid for the database.
func (ser *UserMediaListService) Validate(e *UserMediaList) (err error) {
	return ser.DB.View(func(tx *bolt.Tx) error {
		// Check if User with ID specified in UserMediaList exists
		// Get User bucket, exit if error
		ub, err := Bucket(UserBucketName, tx)
		if err != nil {
			return err
		}
		_, err = get(e.UserID, ub)
		if err != nil {
			return err
		}

		return nil
	})
}

// persistOldProperties maintains certain properties
// of the existing entity in updates
func (ser *UserMediaListService) persistOldProperties(old *UserMediaList, new *UserMediaList) (err error) {
	new.Version = old.Version + 1
	return nil
}
