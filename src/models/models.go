package models

import "time"

type HomePageData struct {
	User     User
	Band     Band
	Messages []Message
}

type MenuPageData struct {
	User User
	Band Band
}

type User struct {
	ID           int
	UserID       string
	Name         string
	Email        string
	PasswordHash string
	BandName     string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
