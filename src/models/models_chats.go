package models

import "time"

type MessageReaction struct {
	ID         string    `json:"_"`
	ReactionID string    `json:"reaction_id"`
	MessageID  string    `json:"-"`
	UserName   string    `json:"-"`
	UserID     string    `json:"user_id"`
	Reaction   string    `json:"reaction"`
	CreatedAt  time.Time `json:"created_at"`
}

type Message struct {
	ID               int
	MessageID        string
	BandID           string
	UserID           string
	ProfileImagePath string
	UserName         string
	ChatID           string
	Body             string
	Reactions        []MessageReaction
	CreatedAt        time.Time
}

type ChatPreview struct {
	ChatID            string
	UserID            string
	Name              string
	IsPrimary         string
	ImageID           string
	ImagePath         string
	LatestMessageID   string
	LatestSenderID    string
	LatestSenderName  string
	LatestMessage     string
	LatestMessageTime time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type ChatsPageData struct {
	User        User
	Band        Band
	Members     []User
	PrimaryChat ChatPreview
	Chats       []ChatPreview
}

type ChatPageData struct {
	User     User
	Band     Band
	Members  []User
	Chat     Chat
	Messages []Message
}

type Chat struct {
	ID                int
	ChatID            string
	BandID            string
	UserID            string
	Name              string
	Slug              string
	IsPrimary         bool
	ImageID           string
	ImagePath         string
	LatestMessage     string
	LatestMessageTime time.Time
	CreatedAt         time.Time
	CreatedBy         string
	UpdatedAt         time.Time
	UpdatedBy         string
	LastMessage       time.Time
}
