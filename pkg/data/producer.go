package data

import "strings"

// Producer represents a single studio, producer,
// licensor, etc.
type Producer struct {
	ID      int
	Titles  []Info
	Type    string
	Version int
}

// Clean cleans the given Producer for storage
func (ser *ProducerService) Clean(e *Producer) (err error) {
	if err = infoListClean(e.Titles); err != nil {
		return err
	}
	e.Type = strings.Trim(e.Type, " ")
	return nil
}

// Validate returns an error if the Producer is
// not valid for the database
func (ser *ProducerService) Validate(e *Producer) (err error) {
	return nil
}
