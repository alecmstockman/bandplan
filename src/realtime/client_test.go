package realtime

import (
	"bandplan/src/database"
	"bandplan/src/helpers"
	"bandplan/src/models"
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateMentionNotificationsDispatchesPush(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	previousDB := database.DB
	database.DB = db
	t.Cleanup(func() {
		database.DB = previousDB
		db.Close()
	})

	mock.ExpectQuery("FROM chat_members").
		WithArgs("chat").
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "display_name", "profile_image_path"}).
			AddRow("sender", "Sender", "").
			AddRow("recipient", "Recipient", ""))
	mock.ExpectExec("INSERT INTO notifications").
		WithArgs(sqlmock.AnyArg(), "band", "recipient", models.NotificationChatMention, "chat", models.NotifBodyChatMention, "sender").
		WillReturnResult(sqlmock.NewResult(0, 1))

	type push struct {
		userID string
		body   models.PushNotifBody
		url    string
	}
	pushes := make(chan push, 4)
	hub := &Hub{broadcast: make(chan BroadcastMessage, 4)}
	client := &Client{
		hub:    hub,
		bandID: "band",
		userID: "sender",
		sendPush: func(_ context.Context, userID string, pushType models.PushNotifType, itemID string) error {
			body, notificationURL := helpers.PushNotificationsGetNotifBodyAndURL(pushType, itemID)
			pushes <- push{userID, body, notificationURL}
			return nil
		},
	}

	client.createMentionNotifications(context.Background(), IncomingMessage{
		ChatID:           "chat",
		MentionedUserIDs: []string{"sender", "outsider", "recipient", "recipient"},
	})

	select {
	case got := <-pushes:
		want := push{"recipient", models.PushChatMentionBody, "/chat?id=chat"}
		if got != want {
			t.Fatalf("push = %+v, want %+v", got, want)
		}
	case <-time.After(time.Second):
		t.Fatal("no push dispatched")
	}

	if len(hub.broadcast) != 1 {
		t.Fatalf("got %d in-app notifications, want 1", len(hub.broadcast))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
