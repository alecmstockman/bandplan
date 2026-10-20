package models

import "time"

type NotificationType string

const (
	NotificationSongCreated    NotificationType = "song_created"
	NotificationSetlistCreated NotificationType = "setlist_created"
	NotificationEventCreated   NotificationType = "event_created"
	NotificationChatCreated    NotificationType = "chat_created"
)

type Notification struct {
	ID             int
	NotificationID string
	BandID         *string
	UserID         string
	Type           NotificationType
	ItemID         string
	Body           string
	CreatedAt      time.Time
	CreatedBy      string
	ReadAt         *time.Time
}
