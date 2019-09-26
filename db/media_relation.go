package db

import (
	"errors"

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
func MediaRelationGetByID(id uint, db *gorm.DB) (relation MediaRelation) {
	db.Set("gorm:auto_preload", true).First(&relation, id)
	return
}

// MediaRelationCreate persists a new record for the provided
// MediaRelation instance
func MediaRelationCreate(relation *MediaRelation, db *gorm.DB) error {
	if !db.NewRecord(relation) {
		return errors.New("database insertion failed")
	}

	db.Create(relation)
	return nil
}
