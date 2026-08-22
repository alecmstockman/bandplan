package models

import "time"

type Message struct {
	ID               int
	MessageID        string
	BandID           string
	UserID           string
	ProfileImagePath string
	UserName         string
	ChatID           string
	Body             string
	CreatedAt        time.Time
}

type ChatPreview struct {
	Body      string
	CreatedAt time.Time
}

type ChatsDataPage struct {
	User  User
	Band  Band
	Chats []Chat
}

type Chat struct {
	ID                int
	ChatID            string
	BandID            string
	UserID            string
	Name              string
	Slug              string
	IsPrimary         bool
	LatestMessage     string
	LatestMessageTime time.Time
	CreatedAt         time.Time
	CreatedBy         string
	UpdatedAt         time.Time
	UpdatedBy         string
	LastMessage       time.Time
}
