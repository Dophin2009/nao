package data

import (
	"strings"
	"time"
)

// Media represents a single instance of a media
type Media struct {
	ID              int
	Titles          []Info
	Synopses        []Info
	Background      []Info
	StartDate       *time.Time
	EndDate         *time.Time
	SeasonPremiered *Season
	Type            *string
	Source          *string
	Version         int
}

// Season contains information about the quarter
// and year
type Season struct {
	Quarter Quarter
	Year    int
}

// Quarter represents the quarter of the year
// by integer
type Quarter int

const (
	// Winter is the first quarter of the year,
	// encapsulating the months January, February,
	// and March
	Winter Quarter = iota

	// Spring is the second quarter of the year,
	// encapsulating the months April, May, and June
	Spring

	// Summer is the third quarter of the year,
	// encapsulating the months July, August, and
	// September
	Summer

	// Fall is the fouth quarter of the year,
	// encapsulating the months October,
	// November, and December
	Fall
)

// Clean cleans the given Media for storage
func (ser *MediaService) Clean(e *Media) (err error) {
	if err := infoListClean(e.Titles); err != nil {
		return err
	}
	if err := infoListClean(e.Synopses); err != nil {
		return err
	}
	if err := infoListClean(e.Background); err != nil {
		return err
	}
	if e.Type != nil {
		*e.Type = strings.Trim(*e.Type, " ")
	}
	if e.Source != nil {
		*e.Source = strings.Trim(*e.Source, " ")
	}
	return nil
}

// Validate checks if the given Media is valid
func (ser *MediaService) Validate(e *Media) (err error) {
	return nil
}
