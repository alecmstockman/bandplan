package helpers

import (
	"bandplan/src/models"
	"net/url"
)

func PushNotificationsGetNotifBodyAndURL(notifType models.PushNotifType, itemID string) (models.PushNotifBody, string) {
	itemID = url.QueryEscape(itemID)
	switch string(notifType) {
	case string(models.PushSongCreated):
		return models.PushSongCreatedBody, "/song?id=" + itemID
	case string(models.PushSongImported):
		return models.PushSongImportedBody, "/song?id=" + itemID

	case string(models.PushSetlistCreated):
		return models.PushSetlistCreatedBody, "/setlist?id=" + itemID
	case string(models.PushSetlistUpdated):
		return models.PushSetlistUpdatedBody, "/setlist?id=" + itemID
	case string(models.PushSetlistDeleted):
		return models.PushSetlistDeletedBody, "/setlists"

	case string(models.PushTaskCreated):
		return models.PushTaskCreatedBody, "/todo?todo-id=" + itemID
	case string(models.PushTaskAssigned):
		return models.PushTaskAssignedBody, "/todo?todo-id=" + itemID
	case string(models.PushTaskUpdated):
		return models.PushTaskUpdatedBody, "/todo?todo-id=" + itemID
	case string(models.PushTaskDeleted):
		return models.PushTaskDeletedBody, "/todo?todo-id=" + itemID

	case string(models.PushTodoListCreated):
		return models.PushTodoListCreatedBody, "/todo?todo-id=" + itemID
	case string(models.PushTodoListUpdated):
		return models.PushTodoListUpdatedBody, "/todo?todo-id=" + itemID
	case string(models.PushTodoListAssigned), string(models.PushTodoAssigned):
		return models.PushTodoListAssignedBody, "/todo?todo-id=" + itemID
	case string(models.PushTodoListDeleted):
		return models.PushTodoListDeletedBody, "/todos"

	case string(models.PushChecklistCreated):
		return models.PushChecklistCreatedBody, "/checklist?checklist-id=" + itemID
	case string(models.PushChecklistUpdated):
		return models.PushChecklistUpdatedBody, "/checklist?checklist-id=" + itemID
	case string(models.PushChecklistAssigned):
		return models.PushChecklistAssignedBody, "/checklist?checklist-id=" + itemID
	case string(models.PushChecklistDeleted):
		return models.PushChecklistDeletedBody, "/checklists"

	case string(models.PushEventCreated):
		return models.PushEventCreatedBody, "/event?event-id=" + itemID
	case string(models.PushEventUpdated):
		return models.PushEventUpdatedBody, "/event?event-id=" + itemID
	case string(models.PushEventCancelled):
		return models.PushEventCancelledBody, "/event?event-id=" + itemID
	case string(models.PushEventMemberInvite):
		return models.PushEventMemberInviteBody, "/event?event-id=" + itemID
	case string(models.PushEventAllConfirmed):
		return models.PushEventAllConfirmedBody, "/event?event-id=" + itemID

	case string(models.PushChatCreated):
		return models.PushChatCreatedBody, "/chat?id=" + itemID
	case string(models.PushChatMention):
		return models.PushChatMentionBody, "/chat?id=" + itemID

	case string(models.PushNewMember):
		return models.PushNewMemberBody, "/band?band-id=" + itemID
	default:
		return models.PushDefaultBody, "/notifications"
	}
}
