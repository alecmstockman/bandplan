package helpers

import (
	"bandplan/src/models"
)

func PushNoticiationsGetNotifBody(notifType models.PushNotifType) models.PushNotifBody {

	switch notifType {

	case models.PushTodoAssigned:
		return models.PushTodoAssignedBody
	case models.PushTaskAssigned:
		return models.PushTaskAssignedBody

	case models.PushEventCreated:
		return models.PushEventCreatedBody
	case models.PushEventUpdated:
		return models.PushEventUpdatedBody

	case models.PushSetlistCreated:
		return models.PushSetlistCreatedBody
	case models.PushSetlistUpdated:
		return models.PushSetlistUpdatedBody
	default:
		return models.PushDefaultBody
	}
}
