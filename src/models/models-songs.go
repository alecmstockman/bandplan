package models

import "time"

type SongDownloadData struct {
	User User
	Band Band
}

type SongPageData struct {
	BackURL  string
	User     User
	Band     Band
	Song     Song
	Setlists []Setlist
}

type LyricsPageData struct {
	SongID string
	BandID string
	Title  string
	Lyrics string
}

type SetlistData struct {
	User     User
	Band     Band
	Setlists []Setlist
}

type SetlistPage struct {
	User    User
	Band    Band
	Setlist Setlist
}

type SetlistSummary struct {
	SetlistID string
	BandID    string
	Name      string
	SongCount string
	Length    string
	UpdatedAt time.Time
	UpdatedBy string
}

type Setlist struct {
	ID          int
	SetlistID   string
	BandID      string
	Name        string
	Slug        string
	Explicit    bool
	Notes       string
	ArtworkID   string
	ArtworkPath string
	Songs       []SetlistItem
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
	UpdatedBy   string
}

type SetlistItem struct {
	ID        string
	SetlistID string
	Type      string
	SongID    string
	Position  int
	Song      Song
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
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

type Transition struct {
	ID           int
	TransitionID string
	BandID       string

	Title     string
	TitleSlug string

	BPM           int
	TimeSignature string
	Key           string
	Tuning        string
	Capo          string
	LengthSeconds int

	Explicit  bool
	Chords    string
	ChartLink string

	LinkOne   string
	LinkTwo   string
	LinkThree string

	Lyrics      string
	Description string
	Notes       string

	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}

type Break struct {
	ID        int
	BreakID   string
	BandID    string
	Title     string
	TitleSlug string

	Notes         string
	Description   string
	LengthSeconds int

	LinkOne   string
	LinkTwo   string
	LinkThree string

	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}
