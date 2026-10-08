package models

import "time"

type HomePageData struct {
	User     User
	Band     Band
	ChatID   string
	Event    Event
	Messages []Message
}

type SettingsPageData struct {
	CSRFToken   string
	User        User
	CurrentBand Band
}

type BandsPageData struct {
	CSRFToken string
	User      User
	Band      Band
	Bands     []Band
}

type BandPageData struct {
	User         User
	Band         Band
	SelectedBand Band
}

type MenuPageData struct {
	CSRFToken   string
	User        User
	Band        Band
	Songs       []Song
	Setlists    []Setlist
	CurrentPage string
}

type AdminPageData struct {
	User  User
	Band  Band
	Users []User
	Code  string
}

type User struct {
	ID               int
	UserID           string
	Name             string
	FirstName        string
	LastName         string
	DisplayName      string
	Email            string
	Slug             string
	PasswordHash     string
	IsAdmin          bool
	ProfileImageID   string
	ProfileImagePath string
	TimeZone         string
	IsEmailVerified  bool
	LegalAccepted    bool
	LegalAcceptedAt  *time.Time
	TermsVersion     string
	PrivacyVersion   string
	LastLogin        time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type UserPermissions struct {
	ID             int
	UserID         string
	BandID         string
	UpdateBand     bool
	AddSongs       bool
	EditSongs      bool
	DeleteSongs    bool
	AddSetlists    bool
	EditSetlists   bool
	DeleteSetlists bool
	AddEvent       bool
	UpdateEvent    bool
	DeleteEvent    bool

	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}

type Band struct {
	ID        int
	BandID    string
	Name      string
	Slug      string
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
