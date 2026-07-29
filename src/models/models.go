package models

import "time"

type HomePageData struct {
	User     User
	Band     Band
	Messages []Message
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

type Message struct {
	ID               int
	MessageID        string
	BandID           string
	UserID           string
	ProfileImagePath string
	UserName         string
	Body             string
	CreatedAt        time.Time
}
