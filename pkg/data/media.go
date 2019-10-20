package data

import (
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
	Type            string
	Source          string
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

// Validate checks if the given Media is valid
func (ser *MediaService) Validate(e *Media) (err error) {
	return nil
}
