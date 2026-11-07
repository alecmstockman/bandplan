package database

import (
	"bandplan/src/models"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNotifTableCreateUserNotificationIncludesBandAndBody(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create database mock: %v", err)
	}

	originalDB := DB
	DB = db
	t.Cleanup(func() {
		DB = originalDB
		db.Close()
	})

	bandID := "band-id"
	notification := models.Notification{
		BandID:    &bandID,
		UserID:    "assigned-user-id",
		Type:      models.NotificationTodoListAssigned,
		ItemID:    "todo-list-id",
		Body:      models.NotifBodyListAssigned,
		CreatedBy: "assigning-user-id",
	}

	mock.ExpectExec(`INSERT INTO notifications`).
		WithArgs(
			sqlmock.AnyArg(),
			bandID,
			notification.UserID,
			notification.Type,
			notification.ItemID,
			notification.Body,
			notification.CreatedBy,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if _, err := NotifTableCreateUserNotification(t.Context(), notification); err != nil {
		t.Fatalf("create user notification: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}

func TestNotifTableCreateUserNotificationRequiresBand(t *testing.T) {
	notification := models.Notification{
		UserID:    "assigned-user-id",
		CreatedBy: "assigning-user-id",
	}

	_, err := NotifTableCreateUserNotification(t.Context(), notification)
	if err == nil {
		t.Fatal("create user notification without band: expected error")
	}
}

func TestNotifTableGetNotificationIDsByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create database mock: %v", err)
	}

	originalDB := DB
	DB = db
	t.Cleanup(func() {
		DB = originalDB
		db.Close()
	})

	query := `
		SELECT notification_id
		FROM notifications
		WHERE user_id = $1
			AND band_id = $2
	`
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("user-id", "band-id").
		WillReturnRows(sqlmock.NewRows([]string{"notification_id"}).
			AddRow("notification-1").
			AddRow("notification-2"))

	notificationIDs, err := NotifTableGetNotificationIDsByUserID(
		t.Context(),
		"user-id",
		"band-id",
	)
	if err != nil {
		t.Fatalf("get notification IDs: %v", err)
	}
	if !notificationIDs["notification-1"] || !notificationIDs["notification-2"] {
		t.Fatalf("unexpected notification IDs: %v", notificationIDs)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}
