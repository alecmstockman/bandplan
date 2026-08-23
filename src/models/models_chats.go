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
	ChatID            string
	UserID            string
	Name              string
	IsPrimary         string
	LatestMessageID   string
	LatestSenderID    string
	LatestSenderName  string
	LatestMessage     string
	LatestMessageTime time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type ChatsDataPage struct {
	User    User
	Band    Band
	Members []User
	Chats   []ChatPreview
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
