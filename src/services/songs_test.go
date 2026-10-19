package services

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"context"
	"database/sql"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func setupSongDataTest(t *testing.T) (sqlmock.Sqlmock, Service) {
	t.Helper()

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

	return mock, Service{}
}

func TestSongDataReturnsDatabaseError(t *testing.T) {
	mock, service := setupSongDataTest(t)

	user := models.User{UserID: "user-123"}
	band := models.Band{BandID: "band-456"}
	songID := "song-789"

	mock.ExpectQuery(`(?s)FROM songs.*WHERE song_id = \$1.*AND band_id = \$2`).
		WithArgs(songID, band.BandID).
		WillReturnError(sql.ErrNoRows)

	data, err := service.SongData(
		context.Background(),
		user,
		band,
		songID,
	)
	if err == nil {
		t.Fatal("SongData() error = nil; want database error")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("Songdata() error = %v; want sql.ErrNoRows", err)
	}
	if !strings.Contains(err.Error(), "get song song-789 from database") {
		t.Errorf("SongData() error = %q; missing song context", err)
	}
	if !reflect.DeepEqual(data, models.SongPageData{}) {
		t.Errorf("SongData() data = %#v; want empty page data", data)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}
