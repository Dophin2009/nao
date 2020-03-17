package models

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/Dophin2009/nao/pkg/db"
)

// Media represents a single instance of a media
type Media struct {
	Titles          []Title
	Synopses        []Title
	Background      []Title
	StartDate       *time.Time
	EndDate         *time.Time
	SeasonPremiered Season
	Type            *string
	Source          *string
	Meta            db.ModelMetadata
}

// Metadata returns Meta.
func (m *Media) Metadata() *db.ModelMetadata {
	return &m.Meta
}

// Season contains information about the quarter and year.
type Season struct {
	Quarter *Quarter
	Year    *int
}

// Quarter represents the quarter of the year by integer.
type Quarter int

const (
	// QuarterWinter is the first quarter of the year, encapsulating the months January,
	// February, and March.
	QuarterWinter Quarter = iota + 1

	// QuarterSpring is the second quarter of the year, encapsulating the months April,
	// May, and June.
	QuarterSpring

	// QuarterSummer is the third quarter of the year, encapsulating the months July,
	// August, and September.
	QuarterSummer

	// QuarterFall is the fouth quarter of the year, encapsulating the months October,
	// November, and December.
	QuarterFall
)

// IsValid checks if the Quarter has a value that is a valid one.
func (q Quarter) IsValid() bool {
	switch q {
	case QuarterWinter, QuarterSpring, QuarterSummer, QuarterFall:
		return true
	}
	return false
}

// String returns the written name of the Quarter.
func (q Quarter) String() string {
	switch q {
	case QuarterWinter:
		return "Winter"
	case QuarterSpring:
		return "Spring"
	case QuarterSummer:
		return "Summer"
	case QuarterFall:
		return "Fall"
	}
	return fmt.Sprintf("%d", int(q))
}

// UnmarshalGQL casts the type of the given value to a Quarter.
func (q *Quarter) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid value: %v", v)
	}

	switch str {
	case "Winter":
		*q = QuarterWinter
	case "Spring":
		*q = QuarterSpring
	case "Summer":
		*q = QuarterSummer
	case "Fall":
		*q = QuarterFall
	default:
		return fmt.Errorf("invalid value: %s", str)
	}
	return nil
}

// MarshalGQL serializes the Quarter into a GraphQL readable form.
func (q Quarter) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(q.String()))
}

// Character represents a single character.
type Character struct {
	Names       []Title
	Information []Title
	Meta        db.ModelMetadata
}

// Metadata returns Meta.
func (c *Character) Metadata() *db.ModelMetadata {
	return &c.Meta
}

// Episode represents a single episode or chapter for some media.
type Episode struct {
	Titles   []Title
	Synopses []Title
	Date     *time.Time
	Duration *int
	Filler   bool
	Recap    bool
	Meta     db.ModelMetadata
}

// Metadata returns Meta.
func (ep *Episode) Metadata() *db.ModelMetadata {
	return &ep.Meta
}

// EpisodeSet is an ordered list of episodes.
type EpisodeSet struct {
	MediaID      int
	Descriptions []Title
	Episodes     []int
	Meta         db.ModelMetadata
}

// Metadata returns the Meta.
func (set *EpisodeSet) Metadata() *db.ModelMetadata {
	return &set.Meta
}

// Genre represents a single instance of a genre.
type Genre struct {
	Names        []Title
	Descriptions []Title
	Meta         db.ModelMetadata
}

// Metadata returns Meta.
func (g *Genre) Metadata() *db.ModelMetadata {
	return &g.Meta
}

// Person represents a single person
type Person struct {
	Names       []Title
	Information []Title
	Meta        db.ModelMetadata
}

// Metadata returns Meta.
func (p *Person) Metadata() *db.ModelMetadata {
	return &p.Meta
}

// Producer represents a single studio, producer, licensor, etc.
type Producer struct {
	Titles []Title
	Types  []string
	Meta   db.ModelMetadata
}

// Metadata return Meta.
func (p *Producer) Metadata() *db.ModelMetadata {
	return &p.Meta
}

// MediaCharacter represents a relationship between single instances of Media
// and Character.
type MediaCharacter struct {
	MediaID       int
	CharacterID   *int
	CharacterRole *string
	PersonID      *int
	PersonRole    *string
	Meta          db.ModelMetadata
}

// Metadata returns Meta.
func (mc *MediaCharacter) Metadata() *db.ModelMetadata {
	return &mc.Meta
}

// MediaGenre represents a relationship between single instances of Media and
// Genre.
type MediaGenre struct {
	MediaID int
	GenreID int
	Meta    db.ModelMetadata
}

// Metadata returns Meta.
func (mg *MediaGenre) Metadata() *db.ModelMetadata {
	return &mg.Meta
}

// MediaProducer represents a relationship between single instances of Media
// and Producer.
type MediaProducer struct {
	MediaID    int
	ProducerID int
	Role       string
	Meta       db.ModelMetadata
}

// Metadata returns Meta.
func (mp *MediaProducer) Metadata() *db.ModelMetadata {
	return &mp.Meta
}

// MediaRelation represents a relationship between single instances of Media
// and Producer.
type MediaRelation struct {
	OwnerID      int
	RelatedID    int
	Relationship string
	Meta         db.ModelMetadata
}

// Metadata returns Meta.
func (mr *MediaRelation) Metadata() *db.ModelMetadata {
	return &mr.Meta
}

// User represents a single user.
type User struct {
	Username    string
	Email       string
	Password    []byte
	Permissions UserPermission
	Meta        db.ModelMetadata
}

// Metadata returns Meta.
func (u *User) Metadata() *db.ModelMetadata {
	return &u.Meta
}

// UserPermission contains a number of permissions for users for
// reading/writing data.
type UserPermission struct {
	// WriteMedia is the ability modify global Media.
	WriteMedia bool
	// WriteUsers is the ability to modify other Users.
	WriteUsers bool
}

// UserCharacter represents a relationship between a User and a Character,
// containing information about the User's opinion on the Character.
type UserCharacter struct {
	UserID      int
	CharacterID int
	Score       *int
	Comments    []Title
	Meta        db.ModelMetadata
}

// Metadata returns Meta.
func (uc *UserCharacter) Metadata() *db.ModelMetadata {
	return &uc.Meta
}

// UserEpisode represents a relationship between a User and an Episode,
// containing information about the User's opinion on the Episode.
type UserEpisode struct {
	UserID    int
	EpisodeID int
	Score     *int
	Comments  []Title
	Meta      db.ModelMetadata
}

// Metadata returns Meta.
func (uep *UserEpisode) Metadata() *db.ModelMetadata {
	return &uep.Meta
}

// UserMedia represents a relationship between a User and a Media, containing
// information about the User's opinion on the Media.
type UserMedia struct {
	UserID         int
	MediaID        int
	Priority       *int
	Score          *int
	Recommended    *int
	Status         *WatchStatus
	WatchInstances []WatchedInstance
	Comments       []Title
	Meta           db.ModelMetadata
}

// Metadata returns Meta
func (um *UserMedia) Metadata() *db.ModelMetadata {
	return &um.Meta
}

// WatchedInstance contains information about a single watch of some Media.
type WatchedInstance struct {
	Episodes  int
	Ongoing   bool
	StartDate *time.Time
	EndDate   *time.Time
	Comments  []Title
}

// WatchStatus is an enum that represents the status of a Media's consumption
// by a User.
type WatchStatus int

const (
	// WatchStatusCurrent means the User is currently consuming the Media.
	WatchStatusCurrent WatchStatus = iota

	// WatchStatusCompleted means that the User has consumed the Media in its entirety at
	// least once.
	WatchStatusCompleted

	// WatchStatusPlanning means that the User is planning to consume the Media sometime in
	// the future.
	WatchStatusPlanning

	// WatchStatusDropped means that the User has never consumed the Media in its entirety
	// and abandoned it in the middle somewhere.
	WatchStatusDropped

	// WatchStatusHold means the User has begun consuming the Media, but has placed it on
	// hold.
	WatchStatusHold
)

// UnmarshalJSON defines custom JSON deserialization for WatchStatus.
func (ws *WatchStatus) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	value, ok := map[string]WatchStatus{
		"Completed": WatchStatusCompleted,
		"Planning":  WatchStatusPlanning,
		"Dropped":   WatchStatusDropped,
		"Hold":      WatchStatusHold,
	}[s]
	if !ok {
		return fmt.Errorf("invalid value: %q", s)
	}
	*ws = value
	return nil
}

// MarshalJSON defines custom JSON serialization for WatchStatus.
func (ws *WatchStatus) MarshalJSON() ([]byte, error) {
	value, ok := map[WatchStatus]string{
		WatchStatusCompleted: "Completed",
		WatchStatusPlanning:  "Planning",
		WatchStatusDropped:   "Dropped",
		WatchStatusHold:      "Hold",
	}[*ws]
	if !ok {
		return nil, fmt.Errorf("invalid value: %d", *ws)
	}

	v, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return v, nil
}

// UserMediaList represents a User-created list of UserMedia.
type UserMediaList struct {
	UserID       int
	Names        []Title
	Descriptions []Title
	UserMedia    []int
	Meta         db.ModelMetadata
}

// Metadata returns Meta.
func (uml *UserMediaList) Metadata() *db.ModelMetadata {
	return &uml.Meta
}

// UserPerson represents a relationship between a User and a Person,
// containing information about the User's opinion on the Person.
type UserPerson struct {
	UserID   int
	PersonID int
	Score    *int
	Comments []Title
	Meta     db.ModelMetadata
}

// Metadata returns Meta.
func (up *UserPerson) Metadata() *db.ModelMetadata {
	return &up.Meta
}
