package models

type PushSubscription struct {
	Endpoint string
	P256DH   string
	Auth     string
}

type PushNotifType string

const (
	PushChatMention PushNotifType = "chat_mention"

	PushTaskAssigned PushNotifType = "task_assigned"
	PushTodoAssigned PushNotifType = "todo_assigned"

	PushEventCreated PushNotifType = "event_created"
	PushEventUpdated PushNotifType = "event_updated"

	PushSetlistCreated PushNotifType = "setlist_created"
	PushSetlistUpdated PushNotifType = "setlist_updated"
)

type PushNotifBody string

const (
	PushChatMentionBody PushNotifBody = "You were mentioned in a chat!"

	PushTaskAssignedBody PushNotifBody = "You have been assigned a task!"
	PushTodoAssignedBody PushNotifBody = "A To-Do list has been assigned to you!"

	PushEventCreatedBody PushNotifBody = "An event has been created!"
	PushEventUpdatedBody PushNotifBody = "An event has been updated!"

	PushSetlistCreatedBody PushNotifBody = "A setlist has been created!"
	PushSetlistUpdatedBody PushNotifBody = "A setlist has been updated!"

	PushDefaultBody PushNotifBody = "You have a notification!"
)
