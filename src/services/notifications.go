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

	fmt.Println("\n- CreateNotification")
	fmt.Println("id: ", notification.ID)
	fmt.Println("notif id: ", notification.NotificationID)
	fmt.Println("band id: ", notification.BandID)
	fmt.Println("user id: ", notification.UserID)
	fmt.Println("type: ", notification.Type)
	fmt.Println("item id: ", notification.ItemID)
	fmt.Println("body: ", notification.Body)
	fmt.Println("created at: ", notification.CreatedAt)
	fmt.Println("created by: ", notification.CreatedBy)
	fmt.Println("read at: ", notification.ReadAt)

	notID := uuid.New().String()

	notification.NotificationID = notID

	err := database.NotifTableCreateBandNotification(ctx, notification)
	if err != nil {
		return err
	}

	return nil
}
