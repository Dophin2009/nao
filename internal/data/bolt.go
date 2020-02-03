package data

import (
	"errors"
	"fmt"
	"os"
	"time"

	bolt "go.etcd.io/bbolt"
)

// BoltDatabase implements Database for boltDB.
type BoltDatabase struct {
	Bolt *bolt.DB
}

// BoltTx implements Transaction for boltDB.
type BoltTx struct {
	DB *BoltDatabase
	Tx *bolt.Tx
}

// Database returns the database of the transaction.
func (btx *BoltTx) Database() Database {
	return btx.DB
}

// Unwrap returns the boltDB transaction object.
func (btx *BoltTx) Unwrap() interface{} {
	return btx.Tx
}

// ConnectBoltDatabase connects to the database file at the given path and
// returns a new BoltDatabase pointer.
func ConnectBoltDatabase(path string, mode os.FileMode, buckets []string) (*BoltDatabase, error) {
	// Open database connection
	bdb, err := bolt.Open(path, mode, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Check buckets exist
	if len(buckets) > 0 {
		err = bdb.Update(func(tx *bolt.Tx) error {
			for _, bucket := range Buckets() {
				_, err = tx.CreateBucketIfNotExists([]byte(bucket))
				if err != nil {
					return fmt.Errorf("failed to create bucket: %w", err)
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	db := BoltDatabase{bdb}
	return &db, nil
}

// Clear removes all buckets in the given database.
func (db *BoltDatabase) Clear() error {
	err := db.Bolt.Update(func(tx *bolt.Tx) error {
		for _, bucket := range Buckets() {
			err := tx.DeleteBucket([]byte(bucket))
			if err != nil {
				return fmt.Errorf("failed to delete bucket: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// Bucket returns the bucket with the given name
func (db *BoltDatabase) Bucket(name string, tx Tx) (*bolt.Bucket, error) {
	// Unwrap transaction
	btx, err := db.unwrapTx(tx)
	if err != nil {
		return nil, err
	}

	// Return bucket
	bucket := btx.Bucket([]byte(name))
	if bucket == nil {
		return nil, fmt.Errorf("bucket: %w", errNotFound)
	}
	return bucket, nil

}

// Transaction is a wrapper method that begins a transaction and passes it to
// the given function.
func (db *BoltDatabase) Transaction(writable bool, logic func(Tx) error) error {
	tx, err := db.Bolt.Begin(writable)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	btx := &BoltTx{
		DB: db,
		Tx: tx,
	}
	defer tx.Rollback()

	err = logic(btx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction, rolling back: %w", err)
	}

	return nil
}

// Create persists the given Model.
func (db *BoltDatabase) Create(m Model, ser Service, tx Tx) (int, error) {
	// Unwrap transaction
	btx, err := db.unwrapTx(tx)
	if err != nil {
		return 0, err
	}

	if !btx.Writable() {
		return 0, errUnwritableTx
	}

	// Check service
	err = checkService(ser)
	if err != nil {
		return 0, err
	}

	// Get bucket, exit if error
	b, err := db.Bucket(ser.Bucket(), tx)
	if err != nil {
		return 0, fmt.Errorf("%s %q: %w", errmsgBucketOpen, ser.Bucket(), err)
	}

	// Get next ID in sequence and assign to
	// model
	id, err := b.NextSequence()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errmsgBucketNextSeq, err)
	}
	m.Metadata().ID = int(id)

	// Verify validity of model
	err = ser.Validate(m, tx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errmsgModelValidation, err)
	}

	// Clean model
	err = ser.Clean(m, tx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errmsgModelCleaning, err)
	}

	// Initialize metadata
	meta := m.Metadata()
	meta.CreatedAt = time.Now()
	meta.UpdatedAt = time.Now()
	meta.Version = 0
	err = ser.Initialize(m, tx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errmsgModelInitialize, err)
	}

	// Save model in bucket
	buf, err := ser.Marshal(m, tx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errmsgModelMarshal, err)
	}

	err = b.Put(itob(meta.ID), buf)
	if err != nil {
		return 0, fmt.Errorf("%s %q: %w", errmsgBucketPut, ser.Bucket(), err)
	}

	// Return new ID
	return meta.ID, nil
}

// Update replaces the value of the model with the given ID.
func (db *BoltDatabase) Update(m Model, ser Service, tx Tx) error {
	// Unwrap transaction
	btx, err := db.unwrapTx(tx)
	if err != nil {
		return err
	}

	// Ensure transaction allows updates
	if !btx.Writable() {
		return errUnwritableTx
	}

	// Check service
	err = checkService(ser)
	if err != nil {
		return err
	}

	// Get bucket, exit if error
	b, err := db.Bucket(ser.Bucket(), tx)
	if err != nil {
		return fmt.Errorf("%s %q: %w", errmsgBucketOpen, ser.Bucket(), err)
	}

	// Verify validity of model
	err = ser.Validate(m, tx)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelValidation, err)
	}

	// Prepare
	err = ser.Clean(m, tx)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelCleaning, err)
	}

	// Check if entity with ID exists
	v, err := db.GetRawByID(m.Metadata().ID, ser, tx)
	if err != nil {
		return fmt.Errorf("failed to get by id %d: %w", m.Metadata().ID, err)
	}

	// Unmarshal old
	o, err := ser.Unmarshal(v, tx)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelUnmarshal, err)
	}

	// Replace properties of updated with certain frozen
	// ones of old
	meta := m.Metadata()
	meta.UpdatedAt = time.Now()
	meta.Version = meta.Version + 1
	err = ser.PersistOldProperties(m, o, tx)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelPersistOld, err)
	}

	// Save model
	buf, err := ser.Marshal(m, tx)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelMarshal, err)
	}

	err = b.Put(itob(meta.ID), buf)
	if err != nil {
		return fmt.Errorf("%s %q: %w", errmsgBucketPut, ser.Bucket(), err)
	}

	return nil
}

// Delete deletes the model with the given ID.
func (db *BoltDatabase) Delete(id int, ser Service, tx Tx) error {
	// Unwrap transaction
	btx, err := db.unwrapTx(tx)
	if err != nil {
		return err
	}

	// Ensure transaction allows updates
	if !btx.Writable() {
		return errUnwritableTx
	}

	// Check service
	err = checkService(ser)
	if err != nil {
		return err
	}

	// Get bucket, exit if error
	b, err := db.Bucket(ser.Bucket(), tx)
	if err != nil {
		return fmt.Errorf("%s %q: %w", errmsgBucketOpen, ser.Bucket(), err)
	}

	// Store existing model to return
	err = b.Delete(itob(id))
	if err != nil {
		return fmt.Errorf("failed to delete by id %d: %w", id, err)
	}

	return nil
}

// GetByID retrieves the persisted Model with the given ID. The given service
// and its DB should not be nil.
func (db *BoltDatabase) GetByID(id int, ser Service, tx Tx) (Model, error) {
	// Unwrap transaction
	_, err := db.unwrapTx(tx)
	if err != nil {
		return nil, err
	}

	// Check service
	err = checkService(ser)
	if err != nil {
		return nil, err
	}

	v, err := db.GetRawByID(id, ser, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to get by id %d: %w", id, err)
	}

	// Unmarshal and return
	m, err := ser.Unmarshal(v, tx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errmsgModelMarshal, err)
	}

	return m, nil
}

// GetRawByID is a generic function that queries the given bucket in the given
// database for an entity of the given ID. The given DB pointer should not be
// nil.
func (db *BoltDatabase) GetRawByID(id int, ser Service, tx Tx) ([]byte, error) {
	// Unwrap transaction
	_, err := db.unwrapTx(tx)
	if err != nil {
		return nil, err
	}

	// Check service
	err = checkService(ser)
	if err != nil {
		return nil, err
	}

	// Get bucket, exit if error
	b, err := db.Bucket(ser.Bucket(), tx)
	if err != nil {
		return nil, fmt.Errorf("%s %q: %w", errmsgBucketOpen, ser.Bucket(), err)
	}

	// Get entity by ID, exit if error
	v := b.Get(itob(id))
	if v == nil {
		return nil, fmt.Errorf("model with id %d: %w", id, errNotFound)
	}

	return v, nil
}

// GetMultiple retrieves the persisted instances of a Model type with the given
// IDs.
//
// See GetFilter for details on `first` and `skip`.
func (db *BoltDatabase) GetMultiple(ids []int, first *int, skip *int,
	ser Service, tx Tx, keep func(m Model) bool) ([]Model, error) {
	// Unwrap transaction
	_, err := db.unwrapTx(tx)
	if err != nil {
		return nil, err
	}

	// Check service
	err = checkService(ser)
	if err != nil {
		return nil, err
	}

	// If filter function is nil, filter nothing
	if keep == nil {
		keep = func(_ Model) bool {
			return true
		}
	}

	// Calculate start and end numbers
	start, end := db.calculatePaginationBounds(first, skip)

	// List to return
	list := []Model{}

	// Iterate through values
	i := start
	for _, id := range ids {
		if i >= end {
			break
		}

		m, err := db.GetByID(id, ser, tx)
		if err != nil {
			return nil, fmt.Errorf("failed to get Model by id %d: %w", id, err)
		}

		// Check if pases filter
		if keep(m) {
			list = append(list, m)
			i++
		}
	}

	return list, nil
}

// GetAll retrieves all persisted instances of a Model type with the given data
// layer service.
//
// See GetFilter for details on `first` and `skip`.
func (db *BoltDatabase) GetAll(first *int, skip *int, ser Service, tx Tx) ([]Model, error) {
	return db.GetFilter(first, skip, ser, tx, func(m Model) bool { return true })
}

// GetFilter retrieves all persisted instances of a Model type that pass the
// filter.
//
// Collection begins on the first valid element after skipping the `skip` valid
// elements and continues for `first` valid elements that pass the filter. If
// `skip` is given as nil, collection begins with the first valid element. If
// `first` is given as nil, collection continues until the last persisted
// element is queried. The given service and its DB should not be nil. A nil
// filter function passes all.
func (db *BoltDatabase) GetFilter(first *int, skip *int, ser Service, tx Tx,
	keep func(m Model) bool) ([]Model, error) {
	// Unwrap transaction
	_, err := db.unwrapTx(tx)
	if err != nil {
		return nil, err
	}

	// Check service
	err = checkService(ser)
	if err != nil {
		return nil, err
	}

	// Get bucket, exit if error
	b, err := db.Bucket(ser.Bucket(), tx)
	if err != nil {
		return nil, fmt.Errorf("%s %q: %w", errmsgBucketOpen, ser.Bucket(), err)
	}

	// If filter function is nil, filter nothing
	if keep == nil {
		keep = func(_ Model) bool {
			return true
		}
	}

	// Calculate start and end numbers
	start, end := db.calculatePaginationBounds(first, skip)

	// Get cursor for bucket
	c := b.Cursor()

	// Move cursor to starting element
	var k, v []byte
	c.First()
	for i := 0; i < start; k, v = c.Next() {
		if k == nil {
			continue
		}
		i++
	}

	// Iterate until end is reached
	list := []Model{}
	for i := start; i < end && k != nil; k, v = c.Next() {
		// Unmarshal element
		m, err := ser.Unmarshal(v, tx)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errmsgModelUnmarshal, err)
		}

		// If element does not pass filter, continue to next
		if !keep(m) {
			continue
		}

		list = append(list, m)
		i++
	}

	return list, nil
}

// iterateKeys iterates through the keys of the given database bucket and
// passes the value at each key to some function.
//
// The loop exits only if do returns a true exit flag. If an error is returned
// but exit is false, the error will be ignored. If an error is returned and
// exit is true, the error will be returned.
func (db *BoltDatabase) iterateKeys(bucketName string, tx Tx,
	do func(v []byte, tx Tx) (exit bool, err error)) error {
	// Unwrap transaction
	_, err := db.unwrapTx(tx)
	if err != nil {
		return err
	}

	// Get bucket, exit if error
	b, err := db.Bucket(bucketName, tx)
	if err != nil {
		return fmt.Errorf("%s %q: %w", errmsgBucketOpen, bucketName, err)
	}

	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		exit, err := do(v, tx)
		if !exit {
			return fmt.Errorf("iteration of values aborted: %w", err)
		}

	}

	return nil
}

func (db *BoltDatabase) assertTx(tx Tx) (*BoltTx, error) {
	btx, ok := tx.(*BoltTx)
	if !ok {
		return nil, fmt.Errorf("transaction type %T: %w", tx, errInvalid)
	}
	return btx, nil
}

func (db *BoltDatabase) unwrapTx(tx Tx) (*bolt.Tx, error) {
	btx, err := db.assertTx(tx)
	if err != nil {
		return nil, err
	}

	unwrapped := btx.Unwrap()
	inner, ok := unwrapped.(*bolt.Tx)
	if !ok {
		return nil,
			fmt.Errorf("wrapped transaction type %T: %w", unwrapped, errInvalid)
	}

	return inner, nil
}

func (db *BoltDatabase) calculatePaginationBounds(first *int, skip *int) (int, int) {
	// The number of elements to skip
	var start int
	if skip == nil || *skip <= 0 {
		start = 0
	} else {
		start = *skip
	}

	// When iterator reaches this number, stop
	var end int
	if first == nil || *first < 0 {
		// Return all elements if `first` is nil
		end = -1
	} else if *first == 0 {
		end = start
	} else {
		end = start + *first
	}

	return start, end
}

var (
	// errNil is an error returned when some pointer is nil.
	errNil = errors.New("is nil")
	// errNotFound is an error returned when the requested object is not found.
	errNotFound = errors.New("not found")
	// errAlreadyExists is an error returned when a unique value already exists.
	errAlreadyExists = errors.New("already exists")
	// errInvalid is an error returned when some value is invalid.
	errInvalid = errors.New("invalid")
	// errUnwritableTx is an error returned when an update attempt was made with
	// a transaction object that does now allow updates.
	errUnwritableTx = errors.New("read-only transaction")
)

const (
	errmsgModelCleaning   = "failed to clean model"
	errmsgModelValidation = "failed to validate model"
	errmsgModelInitialize = "failed to initialize model values"
	errmsgModelPersistOld = "failed to persist old model values"
	errmsgModelMarshal    = "failed to marshal model"
	errmsgModelUnmarshal  = "failed to unmarshal model"
	errmsgModelAssertType = "failed to assert type of model"
	errmsgBucketOpen      = "failed to open bucket"
	errmsgBucketNextSeq   = "failed to generate next sequence ID"
	errmsgBucketPut       = "failed to put value in bucket"
	errmsgBucketDelete    = "failed to delete value in bucket"

	errmsgJSONMarshal   = "failed to marshal to JSON"
	errmsgJSONUnmarshal = "failed to unmarshal from JSON"
)
