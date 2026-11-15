package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	webpush "github.com/SherClockHolmes/webpush-go"
)

func main() {
	endpoint := requiredEnv("PUSH_ENDPOINT")
	p256dh := requiredEnv("PUSH_P256DH")
	auth := requiredEnv("PUSH_AUTH")
	requiredEnv("VAPID_CONTACT_EMAIL")
	requiredEnv("VAPID_PUBLIC_KEY")
	requiredEnv("VAPID_PRIVATE_KEY")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := sendTestPush(ctx, endpoint, p256dh, auth); err != nil {
		log.Fatal("Unable to send test push: ", err)
	}

	log.Println("Test push sent successfully")
}

func requiredEnv(name string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		log.Fatalf("%s environment variable is required", name)
	}

	return value
}

func sendTestPush(ctx context.Context, endpoint, p256dh, auth string) error {
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
