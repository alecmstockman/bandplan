package models

import "time"

type EventType string
type EventRecurrence string
type EventAges string

const (
	EventTypeNone        EventType = "none"
	EventTypeShow        EventType = "show"
	EventTypeGig         EventType = "gig"
	EventTypeFestival    EventType = "festival"
	EventTypeRehearsal   EventType = "rehearsal"
	EventTypePractice    EventType = "practice"
	EventTypeWriting     EventType = "writing"
	EventTypeRecording   EventType = "recording"
	EventTypeMeeting     EventType = "meeting"
	EventTypePhotos      EventType = "photos"
	EventTypePress       EventType = "press"
	EventTypeUnavailable EventType = "unavailable"
	EventtypeOther       EventType = "other"
)

const (
	EventRecurrenceNone     EventRecurrence = "none"
	EventRecurrenceDaily    EventRecurrence = "daily"
	EventRecurrenceWeekly   EventRecurrence = "weekly"
	EventRecurrenceBiweekly EventRecurrence = "biweekly"
	EventRecurrenceMonthly  EventRecurrence = "monthly"
	EventRecurrenceYearly   EventRecurrence = "yearly"
)

const (
	EventNA      EventAges = "N/A"
	EventAllAges EventAges = "allages"
	Event21Plus  EventAges = "21plus"
	Event18Plus  EventAges = "18plus"
)

type EventsPageData struct {
	User          User
	Band          Band
	Events        []Event
	UpcomingLen   int
	PastLen       int
	Notifications []Notification
}

type EventPageData struct {
	CSRFToken     string
	BackURL       string
	User          User
	Band          Band
	Event         Event
	Setlists      []Setlist
	Time          time.Time
	Notifications []Notification
}

type EventCreatePageData struct {
	CSRFToken     string
	User          User
	Band          Band
	Setlists      []Setlist
	Notifications []Notification
}

type CalendarExtendedProps struct {
	MultidayPriority string
}

type CalendarEvent struct {
	ID            string                `json:"id"`
	Title         string                `json:"title"`
	Start         string                `json:"start"`
	End           string                `json:"end,omitempty"`
	URL           string                `json:"url"`
	AllDay        bool                  `json:"allDay"`
	ExtendedProps CalendarExtendedProps `json:"extendedProps"`
}

type Event struct {
	ID      int    `json:"-"`
	EventID string `json:"id"`
	BandID  string `json:"band_id"`

	Name      string `json:"title"`
	Slug      string `json:"slug"`
	ImageID   string `json:"image_id"`
	ImagePath string `json:"image_path"`

	EventDate  time.Time       `json:"start"`
	EndDate    *time.Time      `json:"end_date"`
	EventType  EventType       `json:"event_type"`
	Ages       EventAges       `json:"ages"`
	Recurrence EventRecurrence `json:"recurrence"`

	Location  string     `json:"location"`
	Address   string     `json:"address"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end"`
	Timezone  string     `json:"timezone"`

	SetLocation        string     `json:"set_location"`
	LoadInTime         *time.Time `json:"load_in_time"`
	LoadInInstructions string     `json:"load_in_instructions"`
	SoundCheckTime     *time.Time `json:"sound_check_time"`
	SetTime            *time.Time `json:"set_time"`
	SetLengthSeconds   int        `json:"set_length_seconds"`

	VenueName  string `json:"venue_name"`
	AddressOne string `json:"address_one"`
	AddressTwo string `json:"address_two"`
	City       string `json:"city"`
	State      string `json:"state"`
	ZipCode    string `json:"zip_code"`

	PresaleTicketPrice string `json:"presale_ticket_price"`
	TicketPrice        string `json:"ticket_price"`
	TicketLink         string `json:"ticket_link"`

	SetlistID   string `json:"setlist_id"`
	SetlistName string `json:"setlist_name"`
	Notes       string `json:"notes"`

	LinkOneName string `json:"link_one_name"`
	LinkOne     string `json:"link_one"`
	LinkTwoName string `json:"link_two_name"`
	LinkTwo     string `json:"link_two"`

	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
}
