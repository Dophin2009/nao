package db

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// Media represents a single media type
type Media struct {
	gorm.Model
	Synopsis  string
	Titles    []Title `gorm:"polymorphic:Owner"`
	Producers []MediaProducer
}

// MediaGetAll fetches all Media records
func MediaGetAll(db *gorm.DB) (media []Media) {
	db.Set("gorm:auto_preload", true).Find(&media)
	return
}

// MediaGetByID fetches a single Media record by id
func MediaGetByID(id uint, db *gorm.DB) (media Media) {
	db.Set("gorm:auto_preload", true).First(&media, id)
	return
}

// MediaCreate persists a new record for the provided
// Media instance
func MediaCreate(media *Media, db *gorm.DB) error {
	if !db.NewRecord(media) {
		return errors.New("database insertion failed")
	}

	db.Create(&media)
	return nil
}
