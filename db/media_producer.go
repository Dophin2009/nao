package db

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// MediaProducer represents a relationship between single
// instances of Media and Producer
type MediaProducer struct {
	gorm.Model
	MediaID    uint
	ProducerID uint
	Role       string
}

// MediaProducerGetAll fetches all MediaProducer records
func MediaProducerGetAll(db *gorm.DB) (mediaProducer []MediaProducer) {
	db.Find(&mediaProducer)
	return
}

// MediaProducerGetByID fetches a single MediaProducer record by id
func MediaProducerGetByID(id uint, db *gorm.DB) (mediaProducer MediaProducer) {
	db.First(&mediaProducer, id)
	return
}

// MediaProducerCreate persists a new record for the provided
// MediaProducer instance
func MediaProducerCreate(mediaProducer *MediaProducer, db *gorm.DB) error {
	if !db.NewRecord(mediaProducer) {
		return errors.New("database insertion failed")
	}

	db.Create(&mediaProducer)
	return nil
}
