package services

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func FormatNotif(notifType models.NotificationType, creatorID string) string {
	switch notifType {
	case models.NotificationSongCreated:
		return fmt.Sprintf("")
	case models.NotificationSetlistCreated:
		return fmt.Sprintf("")
	case models.NotificationEventCreated:
		return fmt.Sprintf("")
	case models.NotificationChatCreated:
		return fmt.Sprintf("")
	default:
		return ""
	}
}

func (s Service) CreateNotification(ctx context.Context, notification models.Notification) error {

	notID := uuid.New().String()

	notification.NotificationID = notID

	err := database.NotifTableCreateNotification(ctx, notification)
	if err != nil {
		return err
	}

	return nil
}
