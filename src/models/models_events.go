package models

import "time"

type EventType string

const (
	EventTypeShow      EventType = "show"
	EventTypeGig       EventType = "gig"
	EventTypeFestival  EventType = "festival"
	EventTypeRehearsal EventType = "rehearsal"
	EventTypePractice  EventType = "practice"
	EventTypeWriting   EventType = "writing"
	EventTypeRecording EventType = "recording"
	EventTypeMeeting   EventType = "meeting"
	EventTypePhotos    EventType = "photos"
	EventTypePress     EventType = "press"
	EventtypeOther     EventType = "other"
)

type Event struct {
	ID      int
	EventID string
	BandID  string

	Name        string
	Slug        string
	ArtworkID   string
	ArtworkPath string

	EventType EventType
	Recurrig  bool

	Location  string
	Address   string
	StartTime time.Time
	EndTime   time.Time
	Timezone  string

	SetLocation      string
	LoadInTime       time.Time
	SetTime          time.Time
	SetLengthSeconds int

	VenueName  string
	AddressOne string
	AddressTwo string
	City       string
	State      string
	ZipCode    string

	PresaleTicketPrice string
	TicketPrice        string
	TicketLink         string

	SetlistID string
	Notes     string

	AddressLink string
	EventPage   string

	LinkOneName string
	LinkOne     string
	LinkTwoName string
	LinkTwo     string

	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}
