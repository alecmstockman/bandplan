package models

type PushSubscription struct {
	Endpoint string
	P256DH   string
	Auth     string
}

type PushNotifType string

const (
	PushSongCreated  PushNotifType = PushNotifType(NotificationSongCreated)
	PushSongImported PushNotifType = PushNotifType(NotificationSongImported)

	PushSetlistCreated PushNotifType = PushNotifType(NotificationSetlistCreated)
	PushSetlistUpdated PushNotifType = PushNotifType(NotificationSetlistUpdated)
	PushSetlistDeleted PushNotifType = PushNotifType(NotificationSetlistDeleted)

	PushTaskCreated  PushNotifType = PushNotifType(NotificationTaskCreated)
	PushTaskAssigned PushNotifType = PushNotifType(NotificationTaskAssigned)
	PushTaskUpdated  PushNotifType = PushNotifType(NotificationTaskUpdated)
	PushTaskDeleted  PushNotifType = PushNotifType(NotificationTaskDeleted)

	PushTodoListCreated  PushNotifType = PushNotifType(NotificationTodoListCreated)
	PushTodoListUpdated  PushNotifType = PushNotifType(NotificationTodoListUpdated)
	PushTodoListAssigned PushNotifType = PushNotifType(NotificationTodoListAssigned)
	PushTodoListDeleted  PushNotifType = PushNotifType(NotificationTodoListDeleted)

	// PushTodoAssigned is retained for existing callers using the legacy value.
	PushTodoAssigned PushNotifType = "todo_assigned"

	PushChecklistCreated  PushNotifType = PushNotifType(NotificationChecklistCreated)
	PushChecklistUpdated  PushNotifType = PushNotifType(NotificationChecklistUpdated)
	PushChecklistAssigned PushNotifType = PushNotifType(NotificationChecklistAssigned)
	PushChecklistDeleted  PushNotifType = PushNotifType(NotificationChecklistDeleted)

	PushEventCreated      PushNotifType = PushNotifType(NotificationEventCreated)
	PushEventUpdated      PushNotifType = PushNotifType(NotificationEventUpdated)
	PushEventCancelled    PushNotifType = PushNotifType(NotificationEventCancelled)
	PushEventMemberInvite PushNotifType = PushNotifType(NotificationEventMemberInvite)
	PushEventAllConfirmed PushNotifType = PushNotifType(NotificationEventAllConfirmed)

	PushChatCreated PushNotifType = PushNotifType(NotificationChatCreated)
	PushChatMention PushNotifType = PushNotifType(NotificationChatMention)

	PushNewMember PushNotifType = PushNotifType(NotificationNewMember)
)

type PushNotifBody string

const (
	PushSongCreatedBody  PushNotifBody = "A song has been created!"
	PushSongImportedBody PushNotifBody = "A song has been imported!"

	PushSetlistCreatedBody PushNotifBody = "A setlist has been created!"
	PushSetlistUpdatedBody PushNotifBody = "A setlist has been updated!"
	PushSetlistDeletedBody PushNotifBody = "A setlist has been deleted!"

	PushTaskCreatedBody  PushNotifBody = "A task has been created!"
	PushTaskAssignedBody PushNotifBody = "You have been assigned a task!"
	PushTaskUpdatedBody  PushNotifBody = "A task has been updated!"
	PushTaskDeletedBody  PushNotifBody = "A task has been deleted!"

	PushTodoListCreatedBody  PushNotifBody = "A To-Do list has been created!"
	PushTodoListUpdatedBody  PushNotifBody = "One of your To-Do lists has been updated!"
	PushTodoListAssignedBody PushNotifBody = "A To-Do list has been assigned to you!"
	PushTodoListDeletedBody  PushNotifBody = "A To-Do list has been deleted!"

	// PushTodoAssignedBody is retained for existing callers.
	PushTodoAssignedBody PushNotifBody = PushTodoListAssignedBody

	PushChecklistCreatedBody  PushNotifBody = "A checklist has been created!"
	PushChecklistUpdatedBody  PushNotifBody = "A checklist has been updated!"
	PushChecklistAssignedBody PushNotifBody = "A checklist has been assigned to you!"
	PushChecklistDeletedBody  PushNotifBody = "A checklist has been deleted!"

	PushEventCreatedBody      PushNotifBody = "An event has been created!"
	PushEventUpdatedBody      PushNotifBody = "An event has been updated!"
	PushEventCancelledBody    PushNotifBody = "An event has been cancelled!"
	PushEventMemberInviteBody PushNotifBody = "You have been invited to an event!"
	PushEventAllConfirmedBody PushNotifBody = "All attendees confirmed!"

	PushChatCreatedBody PushNotifBody = "You have been added to a chat!"
	PushChatMentionBody PushNotifBody = "You were mentioned in a chat!"

	PushNewMemberBody PushNotifBody = "A new member has joined the band!"

	PushDefaultBody PushNotifBody = "You have a notification!"
)
