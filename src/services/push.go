package services

import (
	"bandplan/src/database"
	"bandplan/src/helpers"
	"bandplan/src/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/SherClockHolmes/webpush-go"
)

func (s Service) SendPushToUser(ctx context.Context, userID string, pushType models.PushNotifType) error {
	fmt.Println("- SendPushToUser")

	subscriptions, err := database.PushSubscriptionsTableListPushSubscriptions(ctx, userID)
	if err != nil {
		slog.Info("loaded push subscriptions",
			"user_id", userID,
			"count", len(subscriptions),
		)
		return fmt.Errorf("list push subscriptions: %w", err)
	}

	body := helpers.PushNoticiationsGetNotifBody(pushType)

	var sendErr error
	for i, sub := range subscriptions {
		if err := ctx.Err(); err != nil {
			return errors.Join(sendErr, err)
		}

		if err := sendPush(ctx, sub.Endpoint, sub.P256DH, sub.Auth, body); err != nil {
			sendErr = errors.Join(sendErr, fmt.Errorf("send to subscription %d: %w", i+1, err))
		}
	}

	return sendErr
}

func sendPush(ctx context.Context, endpoint, p256dh, auth string, body models.PushNotifBody) error {
	fmt.Println("- sendPush")

	u, err := url.Parse(endpoint)
	host := strings.ToLower(u.Hostname())

	isFCM := host == "fcm.googleapis.com"
	isApplePush := host == "push.apple.com" ||
		strings.HasSuffix(host, ".push.apple.com")
	if err != nil ||
		u.Scheme != "https" ||
		u.User != nil ||
		(!isFCM && !isApplePush) {
		return fmt.Errorf("unsupported push endpoint")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	fmt.Println("body: ", body)

	payload, err := json.Marshal(struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}{
		Title: "BandPlan",
		Body:  string(body),
	})
	if err != nil {
		return fmt.Errorf("marshal push payload: %w", err)
	}

	resp, err := webpush.SendNotificationWithContext(
		ctx,
		payload,
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
		slog.Error(
			"unable to send push notification",
			"error", err,
		)
		return err
	}
	defer resp.Body.Close()

	responseBody, readErr := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
	if readErr != nil {
		return fmt.Errorf(
			"read push provider response: %w",
			readErr,
		)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf(
			"push provider returned HTTP %d: %s",
			resp.StatusCode,
			strings.TrimSpace(string(responseBody)),
		)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("push provider returned HTTP %d", resp.StatusCode)
	}
	return nil
}
