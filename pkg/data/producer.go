package data

import "strings"

// Producer represents a single studio, producer,
// licensor, etc.
type Producer struct {
	ID      int
	Titles  []Info
	Types   []string
	Version int
}

// Clean cleans the given Producer for storage
func (ser *ProducerService) Clean(e *Producer) (err error) {
	if err = infoListClean(e.Titles); err != nil {
		return err
	}
	for i, t := range e.Types {
		e.Types[i] = strings.Trim(t, " ")
	}
	return nil
}

// Validate returns an error if the Producer is
// not valid for the database
func (ser *ProducerService) Validate(e *Producer) (err error) {
	return nil
}

// persistOldProperties maintains certain properties
// of the existing entity in updates
func (ser *ProducerService) persistOldProperties(old *Producer, new *Producer) (err error) {
	new.Version = old.Version + 1
	return nil
}
