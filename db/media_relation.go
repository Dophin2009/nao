package db

import (
	"errors"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
)

// MediaRelation represents a single relationship between two media
type MediaRelation struct {
	gorm.Model
	Relation string
	Owner    uint
	Related  uint
}

// MediaRelationGetAll fetches all MediaRelation records
func MediaRelationGetAll(db *gorm.DB) (relation []MediaRelation) {
	db.Set("gorm:auto_preload", true).Find(&relation)
	return
}

// MediaRelationGetByID fetches a single MediaRelation record by id
func MediaRelationGetByID(id uint, db *gorm.DB) (relation MediaRelation, err error) {
	db.Set("gorm:auto_preload", true).First(&relation, id)
	if relation.ID == 0 {
		return relation, errors.New(strings.Join([]string{"media relation with id", strconv.Itoa(int(id)), "not found"}, " "))
	}
	return
}

// MediaRelationCreate persists a new record for the provided
// MediaRelation instance
func MediaRelationCreate(relation *MediaRelation, db *gorm.DB) error {
	if relation.ID != 0 {
		return errors.New("media relation id must not be set")
	}

	if !db.NewRecord(relation) {
		return errors.New(strings.Join([]string{"media relation with id", strconv.Itoa(int(relation.ID)), "already exists"}, " "))
	}

	if _, err := MediaGetByID(relation.Owner, db); err != nil {
		return err
	}

	if _, err := MediaGetByID(relation.Related, db); err != nil {
		return err
	}

	db.Create(relation)
	return nil
}
