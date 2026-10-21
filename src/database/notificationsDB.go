package database

import (
	"bandplan/src/models"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func NotifTableCreateUserNotification(ctx context.Context, notif models.Notification) error {
	if notif.UserID == "" {
		return fmt.Errorf("user ID is required")
	}
	if notif.CreatedBy == "" {
		return fmt.Errorf("created by ID is required")
	}

	notifID := uuid.New().String()

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
			$1, $2, $3, $4, 
			$5, %6, $7
		)
    `

	_, err := DB.ExecContext(
		ctx,
		query,
		notifID,
		notif.UserID,
		notif.Type,
		notif.ItemID,
		notif.Body,
		notif.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("NotifTableCreateUserNotification - exec: %w", err)
	}

	return nil
}

func NotifTableCreateBandNotification(ctx context.Context, notif models.Notification) error {
	if notif.BandID == nil || *notif.BandID == "" {
		return fmt.Errorf("band ID is required")
	}
	if notif.CreatedBy == "" {
		return fmt.Errorf("created by ID is required")
	}

	query := `
        INSERT INTO notifications (
            notification_id,
            band_id,
            user_id,
            type,
            item_id,
            body,
            created_by
        )
        SELECT
            gen_random_uuid()::text,
            bm.band_id,
            bm.user_id,
            $2,
            $3,
            $4,
            $5
        FROM band_members bm
        WHERE bm.band_id = $1
			AND bm.user_id <> $5
    `

	_, err := DB.ExecContext(
		ctx,
		query,
		*notif.BandID,
		notif.Type,
		notif.ItemID,
		notif.Body,
		notif.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("create band notifications: %w", err)
	}

	return nil
}

func NotifTableGetNotificationsByUserID(ctx context.Context, userID, bandID string) ([]models.Notification, error) {

	query := `
		SELECT
			n.id,
			n.notification_id,
			n.band_id,
			n.user_id,
			n.type,
			n.item_id,
			n.body,
			n.created_at,
			n.created_by,
			u.display_name,
			n.read_at
		FROM notifications n
		LEFT JOIN users u
			ON n.created_by = u.user_Id
		WHERE n.user_id = $1
			AND band_id = $2
		ORDER BY 
			(n.read_at IS NULL) DESC,
			n.created_at DESC
		LIMIT 25
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
			&notif.CreatedByName,
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
