package data

// Producer represents a single studio, producer,
// licensor, etc.
type Producer struct {
	ID      int
	Titles  []Info
	Type    string
	Version int
}

// Validate returns an error if the Producer is
// not valid for the database
func (ser *ProducerService) Validate(e *Producer) (err error) {
	return nil
}
