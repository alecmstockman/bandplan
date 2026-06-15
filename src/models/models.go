package models

import "time"

type HomePageData struct {
	User     User
	messages []string
}

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
	ID        int
	MessageID string
	UserID    string
	UserName  string
	Body      string
	CreatedAt time.Time
}
