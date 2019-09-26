package db

import (
	"github.com/jinzhu/gorm"
	// Required to connect to an SQLite db
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Tables provides a slice of zero-value interfaces
// for the types to be persisted
func Tables() (tables []interface{}) {
	tables = append(tables, &Title{}, &Media{})
	return
}

// Connect establishes a connection to or creates an
// SQLite3 database at the provided path
func Connect(dbPath string) (database *gorm.DB) {
	database, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		panic("could not connect to database: " + dbPath)
	}

	return
}

// ConnectWithMigrations performs table migrations after
// establishing a connection to an SQLite3 database at the
// provided path
func ConnectWithMigrations(dbPath string) (database *gorm.DB) {
	database = Connect(dbPath)

	// Perform migrations
	for _, table := range Tables() {
		database.DropTableIfExists(table)
		database.CreateTable(table)
	}

	return
}
