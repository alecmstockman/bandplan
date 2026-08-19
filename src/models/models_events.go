package models

import "time"

type Event struct {
	ID      int
	EventID string
	BandID  string

	Title       string
	Slug        string
	Description string

	EventType string

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
	LinkOneName string
	LinkOne     string

	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}
