package models

import (
	"time"
)

type HomePageData struct {
	User          User
	Band          Band
	ChatID        string
	Event         Event
	ToDoItems     []ToDoItem
	Messages      []Message
	Notifications []Notification
}

type SettingsPageData struct {
	CSRFToken     string
	User          User
	CurrentBand   Band
	Notifications []Notification
}

type BandsPageData struct {
	CSRFToken     string
	User          User
	Band          Band
	Bands         []Band
	Notifications []Notification
}

type BandPageData struct {
	CSRFToken     string
	User          User
	Band          Band
	SelectedBand  Band
	CanEdit       bool
	Notifications []Notification
}

type UserPageData struct {
	CSRFToken     string
	User          User
	Band          Band
	Member        User
	Notifications []Notification
	CurrentPage   string
	BackURL       string
}

type MenuPageData struct {
	CSRFToken     string
	User          User
	Band          Band
	Songs         []Song
	Setlists      []Setlist
	Notifications []Notification
	CurrentPage   string
}

type AdminPageData struct {
	User          User
	Band          Band
	Users         []User
	Code          string
	Notifications []Notification
}

type UserPermissionsPageData struct {
	CSRFToken string
	User      User
	Band      Band
	Member    UserPermissions
}

type User struct {
	ID     int
	UserID string

	Name        string
	FirstName   string
	LastName    string
	DisplayName string

	Email            string
	Slug             string
	PasswordHash     string
	IsAdmin          bool
	ProfileImageID   string
	ProfileImagePath string
	TimeZone         string
	IsEmailVerified  bool
	Permissions      Permissions
	LegalAccepted    bool
	LegalAcceptedAt  *time.Time

	TermsVersion   string
	PrivacyVersion string
	LastLogin      time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Permissions struct {
	EditBand      bool
	ManageMembers bool

	CreateSongs bool
	EditSongs   bool
	DeleteSongs bool

	CreateSetlists bool
	EditSetlists   bool
	DeleteSetlists bool

	CreateEvents bool
	EditEvents   bool
	DeleteEvents bool

	CreateTodos bool
	EditTodos   bool
	DeleteTodos bool
}

type UserPermissions struct {
	ID     int
	UserID string
	BandID string

	DisplayName      string
	FirstName        string
	LastName         string
	Email            string
	ProfileImagePath string

	IsAdmin       bool
	EditBand      bool
	ManageMembers bool

	CreateSongs bool
	EditSongs   bool
	DeleteSongs bool

	CreateSetlists bool
	EditSetlists   bool
	DeleteSetlists bool

	CreateEvents bool
	EditEvents   bool
	DeleteEvents bool

	CreateTodos bool
	EditTodos   bool
	DeleteTodos bool

	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}

type Band struct {
	ID     int
	BandID string
	Name   string
	Slug   string

	ProfileImageID   string
	ProfileImagePath string

	EmailAddress    string
	IsEmailVerified bool
	EmailVerifiedAt time.Time
	Timezone        string

	City    string
	State   string
	Country string

	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}

type AuthContext struct {
	User        User
	CurrentBand Band
	Session     Session
}

type CreateSessionParams struct {
	UserID    string
	BandID    *string
	Token     string
	TokenHash string
}

type Session struct {
	ID        int
	UsersID   string
	BandID    *string
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type ToastData struct {
	Body string
}

type CalendarPageData struct {
	CSRFToken     string
	User          User
	Band          Band
	Events        []Event
	Notifications []Notification
}
