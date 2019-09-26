package db

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// Producer represents a single studio, producer,
// licensor, etc.
type Producer struct {
	gorm.Model
	Titles []Title `gorm:"polymorphic:Owner"`
	Media  []MediaProducer
}

// ProducerGetAll fetches all Producer records
func ProducerGetAll(db *gorm.DB) (producers []Producer) {
	db.Set("gorm:auto_preload", true).Find(&producers)
	return
}

// ProducerGetByID fetches a single Producer record by id
func ProducerGetByID(id uint, db *gorm.DB) (producer Producer) {
	db.Set("gorm:auto_preload", true).First(&producer, id)
	return
}

// ProducerCreate persists a new record for the provided
// Producer instance
func ProducerCreate(producer *Producer, db *gorm.DB) error {
	if !db.NewRecord(producer) {
		return errors.New("database insertion failed")
	}

	db.Create(&producer)
	return nil
}
