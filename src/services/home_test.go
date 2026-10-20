package services

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestHomePageWithoutChats(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create database mock: %v", err)
	}

	originalDB := database.DB
	database.DB = db

	t.Cleanup(func() {
		database.DB = originalDB
		db.Close()
	})

	user := models.User{UserID: "user-id"}
	band := models.Band{BandID: "band-id"}

	mock.ExpectQuery(`FROM messages m`).
		WithArgs(user.UserID, band.BandID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	mock.ExpectQuery(
		`(?s)COALESCE\(e\.set_location, ''\).*FROM events e.*e\.event_date >= CURRENT_DATE`).
		WithArgs(band.BandID, user.UserID).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectQuery(`FROM todo_items i`).
		WithArgs(user.UserID, band.BandID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectQuery(`FROM notifications`).
		WithArgs(user.UserID, band.BandID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	data, err := (Service{}).HomePage(t.Context(), user, band)
	if err != nil {
		t.Fatalf("load home page: %v", err)
	}
	if len(data.Messages) != 0 {
		t.Errorf("messages count = %d; want 0", len(data.Messages))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}

func TestHomePageWithoutUpcomingEvent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create database mock: %v", err)
	}

	originalDB := database.DB
	database.DB = db
	t.Cleanup(func() {
		database.DB = originalDB
		db.Close()
	})

	user := models.User{UserID: "user-id"}
	band := models.Band{BandID: "band-id"}

	mock.ExpectQuery(`FROM messages m`).
		WithArgs(user.UserID, band.BandID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectQuery(`(?s)COALESCE\(e\.set_location, ''\).*FROM events e.*e\.event_date >= CURRENT_DATE`).
		WithArgs(band.BandID, user.UserID).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery(`FROM todo_items i`).
		WithArgs(user.UserID, band.BandID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectQuery(`FROM notifications`).
		WithArgs(user.UserID, band.BandID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	data, err := (Service{}).HomePage(t.Context(), user, band)
	if err != nil {
		t.Fatalf("load home page: %v", err)
	}
	if data.User.UserID != user.UserID {
		t.Errorf("user ID = %q; want %q", data.User.UserID, user.UserID)
	}
	if data.Band.BandID != band.BandID {
		t.Errorf("band ID = %q; want %q", data.Band.BandID, band.BandID)
	}
	if data.Event.Name != "" {
		t.Errorf("event name = %q; want empty event", data.Event.Name)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}
