package models

import "time"

type SetlistPage struct {
	CSRFToken     string
	User          User
	Band          Band
	Setlist       Setlist
	BackURL       string
	Notifications []Notification
}

type SetlistSummary struct {
	SetlistID   string
	BandID      string
	Name        string
	Slug        string
	SongCount   string
	Length      int
	ArtworkID   string
	ArtworkPath string
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
	UpdatedBy   string
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

type SetlistItemType string

const (
	SetlistItemSong       SetlistItemType = "song"
	SetlistItemTransition SetlistItemType = "transition"
	SetlistItemBreak      SetlistItemType = "break"
)

type SetlistItemPage struct {
	User          User
	Band          Band
	Item          SetlistItem
	Title         string
	Length        int
	Type          SetlistItemType
	Notifications []Notification
}

type SetlistItem struct {
	ID        string
	SetlistID string
	ItemType  SetlistItemType
	ItemID    string
	Position  int

	PauseAfterSeconds int

	Song       *Song
	Transition *Transition
	Break      *Break

	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}

type SetlistPDFOptions struct {
	Numbered    bool
	Songs       bool
	Transitions bool
	Breaks      bool
	Length      bool
	Tempo       bool
	Key         bool
	Tuning      bool
	Capo        bool
}
