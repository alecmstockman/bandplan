package helpers

import (
	"bandplan/src/models"
	"net/url"
)

func PushNotificationsGetNotifBodyAndURL(notifType models.PushNotifType, itemID string) (models.PushNotifBody, string) {
	itemID = url.QueryEscape(itemID)
	switch notifType {
	case models.PushChatMention:
		return models.PushChatMentionBody, "/chat?id=" + itemID
	case models.PushTodoAssigned:
		return models.PushTodoAssignedBody, "/todo?todo-id=" + itemID
	case models.PushTaskAssigned:
		return models.PushTaskAssignedBody, "/todo?todo-id=" + itemID
	case models.PushEventCreated, models.PushEventAllConfirmed:
		return models.PushEventCreatedBody, "/event?event-id=" + itemID
	case models.PushEventUpdated:
		return models.PushEventUpdatedBody, "/event?event-id=" + itemID
	case models.PushSetlistCreated:
		return models.PushSetlistCreatedBody, "/setlist?id=" + itemID
	case models.PushSetlistUpdated:
		return models.PushSetlistUpdatedBody, "/setlist?id=" + itemID
	default:
		return models.PushDefaultBody, "/notifications"
	}
}
