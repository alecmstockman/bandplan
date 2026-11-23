package models

import "time"

type NotificationType string

const (
	NotificationSongCreated  NotificationType = "song_created"
	NotificationSongImported NotificationType = "song_imported"

	NotificationSetlistCreated NotificationType = "setlist_created"
	NotificationSetlistUpdated NotificationType = "setlist_updated"
	NotificationSetlistDeleted NotificationType = "setlist_deleted"

	NotificationTaskCreated  NotificationType = "task_created"
	NotificationTaskAssigned NotificationType = "task_assigned"
	NotificationTaskUpdated  NotificationType = "task_updated"
	NotificationTaskDeleted  NotificationType = "task_deleted"

	NotificationTodoListCreated  NotificationType = "todo_list_created"
	NotificationTodoListUpdated  NotificationType = "todo_list_updated"
	NotificationTodoListAssigned NotificationType = "todo_list_assigned"
	NotificationTodoListDeleted  NotificationType = "todo_list_deleted"

	NotificationChecklistCreated  NotificationType = "checklist_created"
	NotificationChecklistUpdated  NotificationType = "checklist_updated"
	NotificationChecklistAssigned NotificationType = "checklist_assigned"
	NotificationChecklistDeleted  NotificationType = "checklist_deleted"

	NotificationEventCreated      NotificationType = "event_created"
	NotificationEventUpdated      NotificationType = "event_updated"
	NotificationEventCancelled    NotificationType = "event_cancelled"
	NotificationEventMemberInvite NotificationType = "event_invite"
	NotificationEventAllConfirmed NotificationType = "event_all_confirmed"

	NotificationChatCreated NotificationType = "chat_created"
	NotificationChatMention NotificationType = "chat_mention"

	NotificationNewMember NotificationType = "new_member"
)

type NotificationBody string

const (
	NotifBodyChatCreated NotificationBody = "added you to a chat!"
	NotifBodyChatMention NotificationBody = "mentioned you in a chat!"

	NotifBodySongCreated  NotificationBody = "created a song!"
	NotifBodySongImported NotificationBody = "imported a song!"

	NotifBodyTaskCreated  NotificationBody = "created a task for you!"
	NotifBodyTaskAssigned NotificationBody = "assigned a task to you!"
	NotifBodyTaskUpdated  NotificationBody = "updated a task!"
	NotifBodyTaskDeleted  NotificationBody = "deleted a task!"

	NotifBodyTodoListCreated  NotificationBody = "created a todo list!"
	NotifBodyTodoListAssigned NotificationBody = "assigned you a todo list!"
	NotifBodyTodoListUpdated  NotificationBody = "updated your todo list!"
	NotifBodyTodoListDeleted  NotificationBody = "deleted a todo list!"

	NotifBodyChecklistCreated  NotificationBody = "created a checklist!"
	NotifBodyChecklistAssigned NotificationBody = "assigned you a checklist!"
	NotifBodyChecklistUpdated  NotificationBody = "updated a checklist!"
	NotifBodyChecklistDeleted  NotificationBody = "deleted a checklist!"

	NotifBodySetlistCreated NotificationBody = "created a setlist!"
	NotifBodySetlistUpdated NotificationBody = "updated a setlist!"
	NotifBodySetlistDeleted NotificationBody = "deleted a setlist!"

	NotifBodyEventCreated      NotificationBody = "created an event!"
	NotifBodyEventUpdated      NotificationBody = "updated an event!"
	NotifBodyEventCancelled    NotificationBody = "cancelled an event!"
	NotifBodyEventMemberInvite NotificationBody = "invited you to an event!"
	NotifBodyEventAllConfirmed NotificationBody = "confirmed! All members confirmed!"

	NotifBodyNewMember    NotificationBody = "has joined the band!"
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
