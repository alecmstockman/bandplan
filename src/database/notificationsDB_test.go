package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNotifTableMarkNotifAsRead(t *testing.T) {
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

	mock.ExpectExec(`UPDATE notifications`).
		WithArgs("user-id", "band-id", "notification-id").
		WillReturnResult(sqlmock.NewResult(0, 1))

	updated, err := NotifTableMarkNotifAsRead(
		t.Context(),
		"user-id",
		"band-id",
		"notification-id",
	)
	if err != nil {
		t.Fatalf("mark notification as read: %v", err)
	}
	if !updated {
		t.Error("updated = false; want true")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}

func TestNotifTableMarkNotifAsReadNotFound(t *testing.T) {
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

	mock.ExpectExec(`UPDATE notifications`).
		WithArgs("user-id", "band-id", "notification-id").
		WillReturnResult(sqlmock.NewResult(0, 0))

	updated, err := NotifTableMarkNotifAsRead(
		t.Context(),
		"user-id",
		"band-id",
		"notification-id",
	)
	if err != nil {
		t.Fatalf("mark notification as read: %v", err)
	}
	if updated {
		t.Error("updated = true; want false")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}
