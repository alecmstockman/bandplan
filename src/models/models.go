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

type SetlistData struct {
	User     User
	Band     Band
	Setlists []Setlist
}

type SetlistSong struct {
	ID        string
	SetlistID string
	SongID    string
	Position  int
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy time.Time
}

type AdminPageData struct {
	User  User
	Band  Band
	Users []User
	Code  string
}

type SongPageData struct {
	User User
	Band Band
	Song Song
}

type SongDownloadData struct {
	User User
	Band Band
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

type Song struct {
	ID     int
	SongID string
	BandID string

	Title     string
	TitleSlug string

	AlbumTitle string
	AlbumID    string
	AlbumSlug  string

	ArtistName string
	ArtistID   string
	ArtistSlug string

	ArtworkID   string
	ArtworkPath string
	ReleaseDate time.Time
	Genre       string

	RecordingBPM  int
	LiveBPM       int
	TimeSignature string
	OriginalKey   string
	LiveKey       string
	Tuning        string
	Capo          string
	LengthSeconds int

	Status    string
	Explicit  bool
	IsCover   bool
	Chords    string
	ChartLink string

	SpotifyLink     string
	AppleMusicLink  string
	YouTubeLink     string
	AmazonMusicLink string
	PandoraLink     string
	DeezerLink      string
	TidalLink       string
	OtherLink       string

	Lyrics      string
	Description string
	Notes       string

	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}

type Setlist struct {
	ID          int
	SetlistID   string
	BandID      string
	Name        string
	Notes       string
	ArtworkPath string
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
	UpdatedBy   string
}
