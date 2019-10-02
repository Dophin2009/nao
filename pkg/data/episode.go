package data

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// Episode represents a single episode or chapter
// for some media
type Episode struct {
	gorm.Model
	MediaID  uint
	Titles   []Title    `gorm:"polymorphic:Owner;save_associations:true"`
	Date     *time.Time `gorm:"type:datetime"`
	Synopsis string
	Duration uint
	Filler   bool
	Recap    bool
}

// EpisodeGetAll fetches all Episode records
func EpisodeGetAll(db *gorm.DB) (episodes []Episode) {
	Preload(db).Find(&episodes)
	return
}

// EpisodeGetByID fetches a single Episode record by id
func EpisodeGetByID(id uint, db *gorm.DB) (episode Episode) {
	Preload(db).First(&episode, id)
	return
}

// EpisodeCreate persists a new record for the provided
// Episode instance
func EpisodeCreate(episode *Episode, db *gorm.DB) error {
	if episode.ID != 0 {
		return errors.New("episode id must not be set")
	}

	if !db.NewRecord(episode) {
		return errors.New(strings.Join([]string{"episode with id", strconv.Itoa(int(episode.ID)), "already exists"}, " "))
	}

	Preload(db).Create(episode)
	return nil
}
