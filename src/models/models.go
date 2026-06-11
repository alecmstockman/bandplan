package models

import "time"

type User struct {
	ID           int
	UserID       string
	Name         string
	Band         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

type Session struct {
	ID        int
	UsersID   string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type Message struct {
	ID int
}
