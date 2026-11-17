package database

import (
	"bandplan/src/models"
	"context"
	"errors"
	"fmt"
)

var ErrPushSubscriptionOwnedByAnotherUser = errors.New("push subscription belongs to another user")

func PushSubscriptionsTableCreateSubscription(ctx context.Context, userID string, endpoint string, p256dh string, auth string) error {

	query := `
		INSERT INTO push_subscriptions (user_id, endpoint, p256dh, auth)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (endpoint) DO UPDATE
		SET p256dh = EXCLUDED.p256dh,
			auth = EXCLUDED.auth
		WHERE push_subscriptions.user_id = EXCLUDED.user_id;
	`

	result, err := DB.ExecContext(
		ctx,
		query,
		userID,
		endpoint,
		p256dh,
		auth,
	)
	if err != nil {
		return fmt.Errorf("PushSubscriptionsTableCreateSubscription - create subscription: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("PushSubscriptionsTableCreateSubscription - count saved subscriptions: %w", err)
	}
	if affected == 0 {
		return ErrPushSubscriptionOwnedByAnotherUser
	}

	return nil
}

func PushSubscriptionsTableListPushSubscriptions(ctx context.Context, userID string) ([]models.PushSubscription, error) {
	fmt.Println("- PushSubscriptionsTableListPushSubscriptions")

	query := `
		SELECT 
			endpoint, 
			p256dh, 
			auth
		FROM push_subscriptions
		WHERE user_id = $1
	`

	rows, err := DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []models.PushSubscription

	for rows.Next() {
		var sub models.PushSubscription
		err := rows.Scan(
			&sub.Endpoint,
			&sub.P256DH,
			&sub.Auth,
		)

		if err != nil {
			return nil, err
		}

		subscriptions = append(subscriptions, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return subscriptions, nil
}
