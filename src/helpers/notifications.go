package helpers

import (
	"bandplan/src/models"
)

func PushNoticiationsGetNotifBody(notifType models.PushNotifType) models.PushNotifBody {

	switch notifType {
	case models.PushTaskAssigned:
		return models.PushTaskAssignedBody
	case models.PushTodoAssigned:
		return models.PushTodoAssignedBody
	default:
		return models.PushDefaultBody
	}
}
