package services

import (
	"bandplan/src/database"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/SherClockHolmes/webpush-go"
)

func (s Service) SendPushToUser(ctx context.Context, userID string) error {
	fmt.Println("- SendPushToUser")

	subscriptions, err := database.PushSubscriptionsTableListPushSubscriptions(ctx, userID)
	if err != nil {
		slog.Info("loaded push subscriptions",
			"user_id", userID,
			"count", len(subscriptions),
		)
		return fmt.Errorf("list push subscriptions: %w", err)
	}

	var sendErr error
	for i, sub := range subscriptions {
		if err := ctx.Err(); err != nil {

			return errors.Join(sendErr, err)
		}

		if err := sendPush(ctx, sub.Endpoint, sub.P256DH, sub.Auth); err != nil {
			sendErr = errors.Join(sendErr, fmt.Errorf("send to subscription %d: %w", i+1, err))
		}
	}

	return sendErr
}

func sendPush(ctx context.Context, endpoint, p256dh, auth string) error {
	fmt.Println("- sendPush")

	u, err := url.Parse(endpoint)
	if err != nil || u.Scheme != "https" ||
		u.Host != "fcm.googleapis.com" || u.User != nil {
		return fmt.Errorf("expected a Chrome push endpoint at fcm.googleapis.com")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := webpush.SendNotificationWithContext(
		ctx,
		[]byte("BandPlan test push"),
		&webpush.Subscription{
			Endpoint: endpoint,
			Keys:     webpush.Keys{P256dh: p256dh, Auth: auth},
		},
		&webpush.Options{
			HTTPClient:      client,
			Subscriber:      os.Getenv("VAPID_CONTACT_EMAIL"),
			VAPIDPublicKey:  os.Getenv("VAPID_PUBLIC_KEY"),
			VAPIDPrivateKey: os.Getenv("VAPID_PRIVATE_KEY"),
			TTL:             60,
		},
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("push provider returned HTTP %d", resp.StatusCode)
	}
	return nil
}
