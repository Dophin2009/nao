package data

import (
	"errors"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
)

// Media represents a single media type
type Media struct {
	gorm.Model
	Synopsis  string
	Titles    []Title `gorm:"polymorphic:Owner;save_associations:true"`
	Producers []MediaProducer
	Episodes  []Episode
	Related   []MediaRelation `gorm:"foreignKey:Owner"`
	RelatedTo []MediaRelation `gorm:"foreignKey:Related"`
}

// MediaGetAll fetches all Media records
func MediaGetAll(db *gorm.DB) (media []Media) {
	Preload(db).Find(&media)
	return
}

// MediaGetByID fetches a single Media record by id
func MediaGetByID(id uint, db *gorm.DB) (media Media, err error) {
	Preload(db).First(&media, id)
	if media.ID == 0 {
		return media, errors.New(strings.Join([]string{"media with id", strconv.Itoa(int(id)), "not found"}, " "))
	}
	return media, nil
}

// MediaCreate persists a new record for the provided
// Media instance
func MediaCreate(media *Media, db *gorm.DB) error {
	if media.ID != 0 {
		return errors.New("media id must not be set")
	}

	if !db.NewRecord(media) {
		return errors.New(strings.Join([]string{"media with id", strconv.Itoa(int(media.ID)), "already exists"}, " "))
	}

	Preload(db).Create(media)
	return nil
}
