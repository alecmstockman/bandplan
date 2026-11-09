package database

import (
	"bandplan/src/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func NotifTableCreateUserNotification(ctx context.Context, notif models.Notification) (string, error) {
	if notif.BandID == nil || *notif.BandID == "" {
		return "", fmt.Errorf("band ID is required")
	}
	if notif.UserID == "" {
		return "", fmt.Errorf("user ID is required")
	}
	if notif.CreatedBy == "" {
		return "", fmt.Errorf("created by ID is required")
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
			$5, $6, $7
		)
    `

	_, err := DB.ExecContext(
		ctx,
		query,
		notifID,
		*notif.BandID,
		notif.UserID,
		notif.Type,
		notif.ItemID,
		notif.Body,
		notif.CreatedBy,
	)
	if err != nil {
		return "", fmt.Errorf("NotifTableCreateUserNotification - create:  %w", err)
	}

	return notif.UserID, nil
}

func NotifTableCreateBandNotification(ctx context.Context, notif models.Notification) ([]string, error) {
	if notif.BandID == nil || *notif.BandID == "" {
		return nil, fmt.Errorf("band ID is required")
	}
	if notif.CreatedBy == "" {
		return nil, fmt.Errorf("created by ID is required")
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
		RETURNING user_id
    `

	rows, err := DB.QueryContext(
		ctx,
		query,
		*notif.BandID,
		notif.Type,
		notif.ItemID,
		notif.Body,
		notif.CreatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("NotifTableCreateBandNotification - create:  %w", err)
	}
	defer rows.Close()

	var userIDs []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("NotifTableCreateBandNotification - scan:  %w", err)
		}

		userIDs = append(userIDs, userID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("NotifTableCreateBandNotification - rows:  %w", err)
	}

	return userIDs, nil
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
		return []models.Notification{}, fmt.Errorf("NotifTableGetNotificationsByUserID - query:  %w", err)
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
			return []models.Notification{}, fmt.Errorf("NotifTableGetNotificationsByUserID - scan:  %w", err)
		}

		notifications = append(notifications, notif)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("NotifTableGetNotificationsByUserID - rows:  %w", err)
	}

	return notifications, nil
}

func NotifTableGetNotificationIDsByUserID(ctx context.Context, userID, bandID string) (map[string]bool, error) {
	query := `
		SELECT notification_id
		FROM notifications
		WHERE user_id = $1
			AND band_id = $2
	`

	rows, err := DB.QueryContext(ctx, query, userID, bandID)
	if err != nil {
		return nil, fmt.Errorf("NotifTableGetNotificationIDsByUserID - query:  %w", err)
	}
	defer rows.Close()

	notificationIDs := make(map[string]bool)

	for rows.Next() {
		var notificationID string

		if err := rows.Scan(&notificationID); err != nil {
			return nil, fmt.Errorf("NotifTableGetNotificationIDsByUserID - scan:  %w", err)
		}

		notificationIDs[notificationID] = true
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("NotifTableGetNotificationIDsByUserID - rows:  %w", err)
	}

	return notificationIDs, nil
}

func NotifTableMarkNotifAsRead(ctx context.Context, userID, bandID, notifID string) (models.Notification, error) {
	fmt.Println("- NotifTableMarkNotifAsRead")

	readTime := time.Now()

	query := `
		UPDATE notifications
		SET read_at = COALESCE(read_at, $4)
		WHERE notification_id = $3
			AND band_id = $2
			AND user_id = $1
		RETURNING
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
	`

	var notif models.Notification

	err := DB.QueryRowContext(ctx, query, userID, bandID, notifID, readTime).Scan(
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
		return models.Notification{}, fmt.Errorf("NotifTableMarkNotifAsRead - query row:  %w", err)
	}

	return notif, nil
}

func NotifTableMarkAllAsRead(ctx context.Context, userID, bandID string) error {
	fmt.Println("- NotifTableMarkAllAsRead")

	fmt.Println("userID: ", userID)

	query := `
		UPDATE notifications
		SET read_at = NOW()
		WHERE user_id = $1
			AND band_id = $2
	`

	_, err := DB.QueryContext(ctx, query, userID, bandID)
	if err != nil {
		return fmt.Errorf("NotifTableMarkAllAsRead - query: %w", err)
	}

	return nil
}
