package models

import "time"

type NotificationType string

const (
	NotificationSongCreated  NotificationType = "song_created"
	NotificationSongImported NotificationType = "song_imported"

	NotificationTaskAssigned NotificationType = "task_created"
	NotificationTaskUpdated  NotificationType = "task_updated"

	NotificationSetlistCreated NotificationType = "setlist_created"
	NotificationSetlistUpdated NotificationType = "setlist_updated"

	NotificationEventCreated NotificationType = "event_created"
	NotificationChatCreated  NotificationType = "chat_created"

	NotificationNewMember        NotificationType = "new_member"
	NotificationTodoListAssigned NotificationType = "list_assigned"
)

type NotificationBody string

const (
	NotifBodySongCreated  NotificationBody = "created a song!"
	NotifBodySongImported NotificationBody = "imported a song!"

	NotifBodyTaskAssigned NotificationBody = "assigned a task to you!"
	NotifBodyTaskUpdated  NotificationBody = "updated a task!"

	NotifBodySetlistCreated NotificationBody = "created a setlist!"
	NotifBodySetlistUpdated NotificationBody = "updated a setlist!"

	NotifBodyEventCreated NotificationBody = "created an event!"
	NotifBodyChatCreated  NotificationBody = "added you to a chat!"

	NotifBodyNewMember    NotificationType = "has joined the band!"
	NotifBodyListAssigned NotificationBody = "has assigned you a To-Do List!"
)

type Notification struct {
	ID             int
	NotificationID string
	BandID         *string
	UserID         string
	Recipients     map[string]bool
	Type           NotificationType
	ItemID         string
	Body           NotificationBody
	CreatedAt      time.Time
	CreatedBy      string
	CreatedByName  string
	ReadAt         *time.Time
}
