package data

import (
	"errors"
	"strconv"
	"strings"

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
func MediaProducerGetByID(id uint, db *gorm.DB) (mediaProducer MediaProducer, err error) {
	db.First(&mediaProducer, id)
	if mediaProducer.ID == 0 {
		return mediaProducer, errors.New(strings.Join([]string{"media-producer relation with id", strconv.Itoa(int(id)), "not found"}, " "))
	}
	return
}

// MediaProducerCreate persists a new record for the provided
// MediaProducer instance
func MediaProducerCreate(mediaProducer *MediaProducer, db *gorm.DB) error {
	if mediaProducer.ID != 0 {
		return errors.New("media-producer relation id must not be set")
	}

	if !db.NewRecord(mediaProducer) {
		return errors.New(strings.Join([]string{"media-producer relation with id", strconv.Itoa(int(mediaProducer.ID)), "already exists"}, " "))
	}

	if _, err := MediaGetByID(mediaProducer.MediaID, db); err != nil {
		return err
	}

	if _, err := ProducerGetByID(mediaProducer.ProducerID, db); err != nil {
		return err
	}

	db.Create(mediaProducer)
	return nil
}
