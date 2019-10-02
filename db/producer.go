package db

import (
	"errors"
	"strconv"
	"strings"

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
	Preload(db).Set("gorm:auto_preload", true).Find(&producers)
	return
}

// ProducerGetByID fetches a single Producer record by id
func ProducerGetByID(id uint, db *gorm.DB) (producer Producer, err error) {
	Preload(db).Set("gorm:auto_preload", true).First(&producer, id)
	if producer.ID == 0 {
		return producer, errors.New(strings.Join([]string{"producer with id", strconv.Itoa(int(id)), "not found"}, " "))
	}
	return producer, nil
}

// ProducerCreate persists a new record for the provided
// Producer instance
func ProducerCreate(producer *Producer, db *gorm.DB) error {
	if producer.ID != 0 {
		return errors.New("producer id must not be set")
	}

	if !db.NewRecord(producer) {
		return errors.New(strings.Join([]string{"producer with id", strconv.Itoa(int(producer.ID)), "already exists"}, " "))
	}

	Preload(db).Create(producer)
	return nil
}
