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
	User        User
	CurrentBand Band
}

type MenuPageData struct {
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
	DisplayName      string
	Email            string
	Slug             string
	PasswordHash     string
	IsAdmin          bool
	ProfileImageID   string
	ProfileImagePath string
	TimeZone         string
	IsEmailVerified  bool
	LastLogin        time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type UserPermissions struct {
 ID int
 UserID string
 BandID string
 UpdateBand bool
 AddSongs bool
 EditSongs bool
 DeleteSongs bool
 AddSetlists bool
 EditSetlists bool
 DeleteSetlists bool
 AddEvent bool
 UpdateEvent bool
 DeleteEvent bool

 CreatedAt        time.Time
 CreatedBy string
	UpdatedAt        time.Time
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
	UserID string
	BandID *string
	Token  string
}

type Session struct {
	ID        int
	UsersID   string
	BandID    *string
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
}
