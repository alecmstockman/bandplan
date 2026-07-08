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

type SongDownloadData struct {
	User User
	Band Band
}

type User struct {
	ID               int
	UserID           string
	Name             string
	Email            string
	PasswordHash     string
	BandName         string
	Role             string
	ProfileImagePath string
	CreatedAt        time.Time
	UpdatedAt        time.Time
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
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CreatedBy   string
	UpdatedBy   string
}
