package db

import (
	"fmt"
	"os"

	bolt "go.etcd.io/bbolt"
)

// BoltDatabase implements Database for boltDB.
type BoltDatabase struct {
	Bolt         *bolt.DB
	Buckets      []string
	ClearOnClose bool
}

// BoltTx implements Transaction for boltDB.
type BoltTx struct {
	DB *DatabaseService
	Tx *bolt.Tx
}

// Database returns the database of the transaction.
func (btx *BoltTx) Database() *DatabaseService {
	return btx.DB
}

// Unwrap returns the boltDB transaction object.
func (btx *BoltTx) Unwrap() interface{} {
	return btx.Tx
}

// BoltDatabaseConfig defines a set of options to be passed when opening a
// boltDB instance.
type BoltDatabaseConfig struct {
	Path         string
	FileMode     os.FileMode
	Buckets      []string
	ClearOnClose bool
}

// ConnectBoltDatabase connects to the database file at the given path and
// returns a new BoltDatabase pointer.
func ConnectBoltDatabase(conf *BoltDatabaseConfig) (*BoltDatabase, error) {
	// Open database connection
	bdb, err := bolt.Open(conf.Path, conf.FileMode, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Check buckets exist
	if len(conf.Buckets) > 0 {
		err = bdb.Update(func(tx *bolt.Tx) error {
			for _, bucket := range conf.Buckets {
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

	db := BoltDatabase{
		Bolt:         bdb,
		Buckets:      conf.Buckets,
		ClearOnClose: conf.ClearOnClose,
	}
	return &db, nil
}

// Close closes the database connection.
func (db *BoltDatabase) Close() error {
	if db.ClearOnClose {
		err := db.Clear()
		if err != nil {
			return fmt.Errorf("failed to clear database: %w", err)
		}
	}
	err := db.Bolt.Close()
	if err != nil {
		return fmt.Errorf("failed to close boltDB: %w", err)
	}
	return nil
}

// Clear removes all buckets in the given database.
func (db *BoltDatabase) Clear() error {
	err := db.Bolt.Update(func(tx *bolt.Tx) error {
		for _, bucket := range db.Buckets {
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
		DB: &DatabaseService{
			DatabaseDriver: db,
		},
		Tx: tx,
	}
	defer tx.Rollback()

	err = logic(btx)
	if err != nil {
		return err
	}

	if writable {
		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("failed to commit transaction, rolling back: %w", err)
		}
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
	err = CheckService(ser)
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
	meta := m.Metadata()
	meta.ID = int(id)

	// Save model in bucket
	buf, err := ser.Marshal(m)
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
	err = CheckService(ser)
	if err != nil {
		return err
	}

	// Get bucket, exit if error
	b, err := db.Bucket(ser.Bucket(), tx)
	if err != nil {
		return fmt.Errorf("%s %q: %w", errmsgBucketOpen, ser.Bucket(), err)
	}

	// Save model
	buf, err := ser.Marshal(m)
	if err != nil {
		return fmt.Errorf("%s: %w", errmsgModelMarshal, err)
	}

	err = b.Put(itob(m.Metadata().ID), buf)
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
	err = CheckService(ser)
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
	// Get raw value
	v, err := db.GetRawByID(id, ser, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to get by id %d: %w", id, err)
	}

	// Unmarshal and return
	m, err := ser.Unmarshal(v)
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
	err = CheckService(ser)
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
func (db *BoltDatabase) GetMultiple(ids []int, first *int, ser Service, tx Tx,
	keep func(m Model) bool) ([]Model, error) {
	list := []Model{}
	collect := func(m Model, _ Service, _ Tx) (exit bool, err error) {
		list = append(list, m)
		return false, nil
	}

	err := db.DoMultiple(ids, first, ser, tx, collect, keep)
	if err != nil {
		return nil, err
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
	list := []Model{}
	collect := func(m Model, ser Service, tx Tx) (exit bool, err error) {
		// Append element to list
		list = append(list, m)
		return false, nil
	}

	err := db.DoEach(first, skip, ser, tx, collect, keep)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// DoMultiple unmarshals and performs some function on the persisted elements
// that pass the given filter function specified by the given IDs.
func (db *BoltDatabase) DoMultiple(ids []int, first *int, ser Service, tx Tx,
	do func(Model, Service, Tx) (exit bool, err error), iff func(Model) bool) error {
	// Unwrap transaction
	_, err := db.unwrapTx(tx)
	if err != nil {
		return err
	}

	// Check service
	err = CheckService(ser)
	if err != nil {
		return err
	}

	// If filter function is nil, filter nothing
	if iff == nil {
		iff = func(_ Model) bool {
			return true
		}
	}

	// Calculate start and end numbers
	start, end := db.calculatePaginationBounds(first, nil)

	// Iterate through values
	i := start
	for _, id := range ids {
		if i >= end {
			break
		}

		m, err := db.GetByID(id, ser, tx)
		if err != nil {
			return fmt.Errorf("failed to get Model by id %d: %w", id, err)
		}

		// Check if pases filter
		if !iff(m) {
			continue
		}

		exit, err := do(m, ser, tx)
		if exit {
			return err
		}
		i++
	}

	return nil
}

// DoEach unmarshals and performs some function on each persisted element
// that passes the filter function.
func (db *BoltDatabase) DoEach(first *int, skip *int, ser Service, tx Tx,
	do func(Model, Service, Tx) (exit bool, err error), iff func(Model) bool) error {
	// Unwrap transaction
	_, err := db.unwrapTx(tx)
	if err != nil {
		return err
	}

	// Check service
	err = CheckService(ser)
	if err != nil {
		return err
	}

	// Get bucket, exit if error
	b, err := db.Bucket(ser.Bucket(), tx)
	if err != nil {
		return fmt.Errorf("%s %q: %w", errmsgBucketOpen, ser.Bucket(), err)
	}

	// If filter function is nil, filter nothing
	if iff == nil {
		iff = func(_ Model) bool {
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
		m, err := ser.Unmarshal(v)
		if err != nil {
			return fmt.Errorf("%s: %w", errmsgModelUnmarshal, err)
		}

		if iff(m) {
			i++
		}
	}

	// Iterate until end is reached
	for i := start; i < end && k != nil; k, v = c.Next() {
		// Unmarshal element
		m, err := ser.Unmarshal(v)
		if err != nil {
			return fmt.Errorf("%s: %w", errmsgModelUnmarshal, err)
		}

		// If element does not pass filter, continue to next
		if !iff(m) {
			continue
		}

		exit, err := do(m, ser, tx)
		if exit {
			return err
		}
		i++
	}

	return nil
}

// FindFirst returns the first element that matches the conditions in the
// given function. Elements are iterated through in key order.
func (db *BoltDatabase) FindFirst(
	ser Service, tx Tx, match func(Model) (bool, error)) (Model, error) {
	var found Model
	check := func(m Model, _ Service, _ Tx) (exit bool, err error) {
		t, err := match(m)
		if err != nil {
			return true, fmt.Errorf("failed to check if match was found: %w", err)
		}

		if t {
			found = m
			return true, nil
		}

		return false, nil
	}

	err := db.DoEach(nil, nil, ser, tx, check, func(m Model) bool {
		return true
	})
	if err != nil {
		return nil, err
	}

	return found, nil
}

// iterateKeys iterates through the keys of the given database bucket and
// passes the value at each key to some function.
//
// The loop exits only if do returns a true exit flag. If an error is returned
// but exit is false, the error will be ignored. If an error is returned and
// exit is true, the error will be returned.
// func (db *BoltDatabase) iterateKeys(bucketName string, tx Tx,
// do func(k, v []byte, tx Tx) (exit bool, err error)) error {
// // Unwrap transaction
// _, err := db.unwrapTx(tx)
// if err != nil {
// return err
// }

// // Get bucket, exit if error
// b, err := db.Bucket(bucketName, tx)
// if err != nil {
// return fmt.Errorf("%s %q: %w", errmsgBucketOpen, bucketName, err)
// }

// c := b.Cursor()
// for k, v := c.First(); k != nil; k, v = c.Next() {
// exit, err := do(k, v, tx)
// if !exit {
// return fmt.Errorf("iteration of values aborted: %w", err)
// }

// }

// return nil
// }

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
