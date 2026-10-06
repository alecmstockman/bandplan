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

type ToDoList struct {
	ID         int
	ToDoListID string
	UserID     string
	BandID     string

	Name        string
	Description string
	Items       []ToDoItem

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

type ToDoListCreatePageData struct {
	CSRFToken string
	User      User
	Band      Band
	ToDoList  ToDoList
	Songs     []Song
}

type ChecklistItem struct {
	ID          int
	ItemID      string
	ChecklistID string

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

type Checklist struct {
	ID          int
	ChecklistID string
	UserID      *string
	BandID      *string

	Name        string
	Description string
	Items       []ChecklistItem

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

type ChecklistCreatePageData struct {
	CSRFToken string
	User      User
	Band      Band
	Checklist Checklist
}
