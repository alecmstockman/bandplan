package models

import (
	"time"
)

type Goal struct {
	ID          int
	ItemID      string
	GoalsListID string
	UserID      string
	BandID      string

	Name     string
	Position int
	Body     string

	song_id    *string
	setlist_id *string
	event_id   *string

	AssignedTo  *string
	IsComplete  bool
	CompletedAt *time.Time
	CompletedBy *string

	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}

type GoalsList struct {
	ID          int
	GoalsListID string
	UserID      string
	BandID      string

	Name        string
	Description string
	Items       []Goal

	song_id    *string
	setlist_id *string
	event_id   *string

	AssignedTo  *string
	IsComplete  bool
	CompletedAt *time.Time
	CompletedBy *string

	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}

type GoalsListCreatePageData struct {
	CSRFToken string
	User      User
	Band      Band
	GoalsList GoalsList
}

type ToDoItem struct {
	ID         int
	ItemID     string
	ToDoListID string
	UserID     string
	BandID     string

	Name     string
	Position int
	Body     string

	SongID    *string
	SetlistID *string
	EventID   *string

	AssignedTo  *string
	DueDate     *time.Time
	DueTime     *time.Time
	DueTimezone *string
	IsComplete  bool
	CompletedAt *time.Time
	CompletedBy *string

	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}

type ToDoListPage struct {
	CSRFToken string
	BackURL   string
	User      User
	Band      Band
	List      ToDoList
	Songs     []Song
	Setlists  []Setlist
	Events    []Event
	Members   []User
}

type ToDoList struct {
	ID         int
	ToDoListID string
	UserID     string
	BandID     string

	Name          string
	Description   *string
	IsPrimaryBand bool
	IsPrimaryUser bool
	Items         []ToDoItem

	SongID    *string
	SetlistID *string
	EventID   *string

	AssignedTo  *string
	DueDate     *time.Time
	DueTime     *time.Time
	DueTimezone *string
	IsComplete  bool
	CompletedAt *time.Time
	CompletedBy *string

	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}

type ToDoListsPageData struct {
	CSRFToken       string
	User            User
	Band            Band
	PrimaryUserList ToDoList
	PrimaryBandList ToDoList
	Items           []ToDoList
	ListCount       int
}

type ToDoListCreatePageData struct {
	CSRFToken string
	User      User
	Band      Band
	ToDoList  ToDoList
	Songs     []Song
	Setlists  []Setlist
	Events    []Event
	Members   []User
}

type ChecklistItem struct {
	ID          int
	ItemID      string
	ChecklistID string

	Name     string
	Position int
	Body     string

	SongID    *string
	SetlistID *string
	EventID   *string

	AssignedTo  *string
	DueDate     *time.Time
	DueTime     *time.Time
	DueTimezone *string
	IsComplete  bool
	CompletedAt *time.Time
	CompletedBy *string

	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}

type Checklist struct {
	ID          int
	ChecklistID string
	UserID      *string
	BandID      *string

	Name        string
	Description string
	Items       []ChecklistItem

	SongID    *string
	SetlistID *string
	EventID   *string

	AssignedTo  *string
	DueDate     *time.Time
	DueTime     *time.Time
	DueTimezone *string
	IsComplete  bool
	CompletedAt *time.Time
	CompletedBy *string

	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}

type ChecklistCreatePageData struct {
	CSRFToken string
	User      User
	Band      Band
	Checklist Checklist
	Songs     []Song
	Setlists  []Setlist
	Events    []Event
	Members   []User
}
