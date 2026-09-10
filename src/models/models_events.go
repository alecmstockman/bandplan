package models

import "time"

type EventType string

const (
	EventTypeShow      EventType = "show"
	EventTypeGig       EventType = "gig"
	EventTypeRehearsal EventType = "rehearsal"
	EventTypePractice  EventType = "practice"
	EventTypeFestival  EventType = "festival"
	EventTypeMeeting   EventType = "festival"
	EventtypeOther     EventType = "other"
)

type Event struct {
	ID      int
	EventID string
	BandID  string

	Title       string
	Slug        string
	Description string

	EventType EventType
	Recurrig  bool

	StartTime time.Time
	EndTime   time.Time
	Timezone  string

	VenueName string
	Address   string
	City      string
	State     string
	ZipCode   string

	SetlistID string

	Notes string

	AddressLink string
	TicketLink  string
	EventPage   string

	LinkOneName   string
	LinkOne       string
	LinkTwoName   string
	LinkTwo       string
	LinkThreeName string
	LinkThree     string

	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}
