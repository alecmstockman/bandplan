package models

import "time"

type EventType string
type EventRecurrence string
type EventAges string
type EventAttendanceStatus string

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
	EventTypeOther       EventType = "other"
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

const (
	EventAttendanceUnseen    EventAttendanceStatus = "unseen"
	EventAttendanceMaybe     EventAttendanceStatus = "maybe"
	EventAttendanceConfirmed EventAttendanceStatus = "confirmed"
	EventAttendanceCannot    EventAttendanceStatus = "cannot"
)

type EventStatus string

const (
	EventStatusActive    EventStatus = "active"
	EventStatusCompleted EventStatus = "completed"
	EventStatusCancelled EventStatus = "cancelled"
)

type EventsPageData struct {
	CSRFToken      string
	User           User
	Band           Band
	UpcomingEvents []Event
	PastEvents     []Event
	Notifications  []Notification
}

type EventEmailData struct {
	User   User
	Band   Band
	Emails []string
	Event  Event
}

type EventEmailShareData struct {
	CSRFToken string
	User      User
	Band      Band
	Members   []User
	Event     Event
}

type EventPageData struct {
	CSRFToken         string
	BackURL           string
	User              User
	Band              Band
	Event             Event
	Setlists          []Setlist
	TodoLists         []ToDoList
	Checklists        []Checklist
	Time              time.Time
	Notifications     []Notification
	GoogleCalendarURL string
	Attendance        []EventAttendance
	CurrentAttendance *EventAttendance
	ConfirmedCount    int
	AttendanceCount   int
}

type EventCreatePageData struct {
	CSRFToken     string
	BackURL       string
	User          User
	Band          Band
	Setlists      []Setlist
	TodoLists     []ToDoList
	Checklists    []Checklist
	Notifications []Notification
	Invitees      EventInviteFields
}

type EventEditPageData struct {
	EventPageData
	Invitees EventInviteFields
}

type EventInviteFields struct {
	BandName string
	Members  []EventInvitee
}

type EventInvitee struct {
	User
	Selected bool
}

func (fields EventInviteFields) AllSelected() bool {
	for _, member := range fields.Members {
		if !member.Selected {
			return false
		}
	}
	return len(fields.Members) > 0
}

type EventAttendance struct {
	ID                int
	EventAttendanceID string
	EventID           string
	UserID            string
	DisplayName       string
	Status            EventAttendanceStatus
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type CalendarExtendedProps struct {
	MultidayPriority string      `json:"multidayPriority"`
	EventType        EventType   `json:"eventType"`
	EventStatus      EventStatus `json:"eventStatus"`
}

type CalendarEvent struct {
	ID              string                `json:"id"`
	Title           string                `json:"title"`
	Start           string                `json:"start"`
	End             string                `json:"end,omitempty"`
	URL             string                `json:"url"`
	AllDay          bool                  `json:"allDay"`
	ExtendedProps   CalendarExtendedProps `json:"extendedProps"`
	BackgroundColor string                `json:"backgroundColor,omitempty"`
	BorderColor     string                `json:"borderColor,omitempty"`
	TextColor       string                `json:"textColor,omitempty"`
}

type Event struct {
	ID      int    `json:"-"`
	EventID string `json:"id"`
	BandID  string `json:"band_id"`

	Name      string `json:"title"`
	Slug      string `json:"slug"`
	ImageID   string `json:"image_id"`
	ImagePath string `json:"image_path"`

	EventDate   time.Time       `json:"start"`
	EndDate     *time.Time      `json:"end_date"`
	EventType   EventType       `json:"event_type"`
	Ages        EventAges       `json:"ages"`
	Recurrence  EventRecurrence `json:"recurrence"`
	EventStatus EventStatus     `json:"event_status"`

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

	SetlistID     string `json:"setlist_id"`
	SetlistName   string `json:"setlist_name"`
	TodoListID    string `json:"todo_list_id"`
	TodoListName  string `json:"todo_list_name"`
	ChecklistID   string `json:"checklist_id"`
	ChecklistName string `json:"checklist_name"`
	Notes         string `json:"notes"`

	LinkOneName string `json:"link_one_name"`
	LinkOne     string `json:"link_one"`
	LinkTwoName string `json:"link_two_name"`
	LinkTwo     string `json:"link_two"`

	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
}
