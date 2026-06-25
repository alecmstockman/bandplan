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
	CurrentPage string
}

type SongPageData struct {
	User User
	Band Band
	Song Song
}

type User struct {
	ID           int
	UserID       string
	Name         string
	Email        string
	PasswordHash string
	BandName     string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Band struct {
	ID        int
	BandID    string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Session struct {
	ID        int
	UsersID   string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type Message struct {
	ID        int
	MessageID string
	BandID    string
	UserID    string
	UserName  string
	Body      string
	CreatedAt time.Time
}

type Song struct {
	ID         int
	SongID     string
	Title      string
	AlbumTitle string
	BandID     string
	Genre      string

	MusicalKey    string
	Tuning        string
	RecordingBPM  int
	LiveBPM       int
	LengthSeconds int

	ReleaseDate    time.Time
	Lyrics         string
	SpotifyLink    string
	AppleMusicLink string
	YouTubeLink    string
	OtherLink      string

	Notes     string
	CreatedAt string
	UpdatedAt string
}
