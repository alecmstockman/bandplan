package models

import "time"

type NotificationType string

const (
	NotificationSongCreated    NotificationType = "song_created"
	NotificationTaskAssigned   NotificationType = "task_created"
	NotificationSongImported   NotificationType = "song_imported"
	NotificationSetlistCreated NotificationType = "setlist_created"
	NotificationEventCreated   NotificationType = "event_created"
	NotificationChatCreated    NotificationType = "chat_created"
	NotificationNewMember      NotificationType = "new_member"
)

type NotificationBody string

const (
	NotifBodySongCreated    NotificationBody = "created a song!"
	NotifBodySongImported   NotificationBody = "imported a song!"
	NotifTaskAssigned       NotificationBody = "assigned a task to you!"
	NotifBodySetlistCreated NotificationBody = "created a setlist!"
	NotifBodyEventCreated   NotificationBody = "created an event!"
	NotifBodyChatCreated    NotificationBody = "added you to a chat!"
	NotifBodyNewMember      NotificationType = "has joined the band!"
)

type Notification struct {
	ID             int
	NotificationID string
	BandID         *string
	UserID         string
	Type           NotificationType
	ItemID         string
	Body           NotificationBody
	CreatedAt      time.Time
	CreatedBy      string
	CreatedByName  string
	ReadAt         *time.Time
}
