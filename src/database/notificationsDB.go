package database

import (
	"bandplan/src/models"
	"context"
	"fmt"
)

func NotifTableCreateNotification(ctx context.Context, notif models.Notification) error {
	fmt.Println("\n - NotifTableCreateNotification")

	fmt.Printf("\nnotif: \n%+v", notif)

	query := `
		INSERT INTO notifications (
			notification_id,
			band_id, 
			user_id,
			type,
			item_id,
			body,
			created_by
		) VALUES (
			$1, NULLIF($2, ''), $3, $4,
			$5, $6, $7
		)
	`

	_, err := DB.ExecContext(
		ctx,
		query,
		notif.NotificationID,
		notif.BandID,
		notif.UserID,
		notif.Type,
		notif.ItemID,
		notif.Body,
		notif.CreatedBy,
	)

	if err != nil {
		return fmt.Errorf("create notification %s error: %w", notif.NotificationID, err)
	}
	return nil
}

func NotifTableGetNotificationsByUserID(ctx context.Context, userID, bandID string) ([]models.Notification, error) {

	query := `
		SELECT
			id,
			notification_id,
			band_id,
			user_id,
			type,
			item_id,
			body,
			created_at,
			created_by,
			read_at
		FROM notifications
		WHERE user_id = $1
			AND band_id = $2
		ORDER BY created_at DESC
		LIMIT 20
	`

	rows, err := DB.QueryContext(ctx, query, userID, bandID)
	if err != nil {
		return []models.Notification{}, fmt.Errorf("NotifTableGetNotificationsByUserID - get notifications by user: %w", err)
	}

	defer rows.Close()

	var notifications []models.Notification

	for rows.Next() {
		var notif models.Notification

		err = rows.Scan(
			&notif.ID,
			&notif.NotificationID,
			&notif.BandID,
			&notif.UserID,
			&notif.Type,
			&notif.ItemID,
			&notif.Body,
			&notif.CreatedAt,
			&notif.CreatedBy,
			&notif.ReadAt,
		)
		if err != nil {
			return []models.Notification{}, fmt.Errorf("NotifTableGetNotificationsByUserID - scan notifications: %w", err)
		}

		notifications = append(notifications, notif)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("NotifTableGetNotificationsByUserID - err: %w", err)
	}

	return notifications, nil
}

func NotifTableMarkNotifAsRead(ctx context.Context, userID, bandID, notifID string) (bool, error) {

	query := `
		UPDATE notifications
		SET read_at = NOW()
		WHERE user_id = $1
			AND band_id = $2
			AND notification_id = $3
	`

	result, err := DB.ExecContext(ctx, query, userID, bandID, notifID)
	if err != nil {
		return false, fmt.Errorf("NotifTableMarkNotifAsRead - err: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("NotifTableMarkNotifAsRead - rows affected: %w", err)
	}

	return rowsAffected == 1, nil
}
