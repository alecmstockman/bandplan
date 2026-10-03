package models

import "time"

type Goal struct {
	ID          int
	ItemID      string
	GoalsListID string
	UserID      string
	BandID      string

	Name     string
	Position int
	Body     string

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

	AssignedTo  *string
	IsComplete  bool
	CompletedAt *time.Time
	CompletedBy *string

	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
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

	AssignedTo  *string
	IsComplete  bool
	CompletedAt *time.Time
	CompletedBy *string

	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}

type ChecklistItem struct {
	ID          int
	ItemID      string
	ChecklistID string
	UserID      string
	BandID      string

	Name     string
	Position int
	Body     string

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
	UserID      string
	BandID      string

	Name        string
	Description string
	Items       []ChecklistItem

	AssignedTo  *string
	IsComplete  bool
	CompletedAt *time.Time
	CompletedBy *string

	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}
