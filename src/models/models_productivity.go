package models

import "time"

type ToDoItem struct {
	ID         string
	ToDoItemID string
	Name       string
	ToDoListID string
	IsComplete bool
	Body       string
	CreatedAt  time.Time
	CreatedBy  string
	UpdatedAt  time.Time
	UpdatedBy  string
}

type ToDoList struct {
	ID         string
	ToDoListID string
	Name       string
	BandID     string
	Items      []ToDoItem
	IsComplete bool
	CreatedAt  time.Time
	CreatedBy  string
	UpdatedAt  time.Time
	UpdatedBy  string
}
