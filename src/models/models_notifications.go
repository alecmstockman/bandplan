package models

import "time"

type NotificationType string

const (
	NotificationSongCreated  NotificationType = "song_created"
	NotificationSongImported NotificationType = "song_imported"

	NotificationTaskAssigned NotificationType = "task_created"
	NotificationTaskUpdated  NotificationType = "task_updated"

	NotificationChecklistCreated  NotificationType = "checklist_created"
	NotificationChecklistUpdated  NotificationType = "checklist_updated"
	NotificationChecklistAssigned NotificationType = "checklist_assigned"
	NotificationChecklistDeleted  NotificationType = "checklist_deleted"

	NotificationSetlistCreated NotificationType = "setlist_created"
	NotificationSetlistUpdated NotificationType = "setlist_updated"

	NotificationEventCreated      NotificationType = "event_created"
	NotificationEventUpdated      NotificationType = "event_updated"
	NotificationEventAllConfirmed NotificationType = "event_all_confirmed"

	NotificationChatCreated NotificationType = "chat_created"
	NotificationChatMention NotificationType = "chat_mention"

	NotificationNewMember        NotificationType = "new_member"
	NotificationTodoListAssigned NotificationType = "list_assigned"
)

type NotificationBody string

const (
	NotifBodyChatCreated NotificationBody = "added you to a chat!"
	NotifBodyChatMention NotificationBody = "mentioned you in a chat!"

	NotifBodySongCreated  NotificationBody = "created a song!"
	NotifBodySongImported NotificationBody = "imported a song!"

	NotifBodyTaskAssigned NotificationBody = "assigned a task to you!"
	NotifBodyTaskUpdated  NotificationBody = "updated a task!"

	NotifBodyChecklistCreated  NotificationBody = "created a checklist!"
	NotifBodyChecklistAssigned NotificationBody = "assigned you a checklist!"
	NotifBodyChecklistUpdated  NotificationBody = "updated a checklist!"
	NotifBodyChecklistDeleted  NotificationBody = "deleted a checklist!"

	NotifBodySetlistCreated NotificationBody = "created a setlist!"
	NotifBodySetlistUpdated NotificationBody = "updated a setlist!"

	NotifBodyEventCreated      NotificationBody = "created an event!"
	NotifBodyEventUpdated      NotificationBody = "updated an event!"
	NotifBodyEventAllConfirmed NotificationBody = "confirmed! All members confirmed!"

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
