package models

type PushSubscription struct {
	Endpoint string
	P256DH   string
	Auth     string
}

type PushNotifType string

const (
	PushTaskAssigned PushNotifType = "task_assigned"
	PushTodoAssigned PushNotifType = "todo_assigned"
)

type PushNotifBody string

const (
	PushTaskAssignedBody PushNotifBody = "You have been assigned a task!"
	PushTodoAssignedBody PushNotifBody = "A To-Do list has been assigned to you!"

	PushDefaultBody PushNotifBody = "You have a notification!"
)
