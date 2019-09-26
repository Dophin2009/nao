package db

import (
	"github.com/jinzhu/gorm"
)

// Title represents a name for some other object
type Title struct {
	gorm.Model
	Name      string
	Language  string
	OwnerID   int
	OwnerType string
}
