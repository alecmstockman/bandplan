package database

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSongsTableUpdateSongTextField(t *testing.T) {
	tests := []struct {
		field  string
		column string
	}{
		{field: "lyrics", column: "lyrics"},
		{field: "description", column: "description"},
		{field: "notes", column: "notes"},
	}

	for _, test := range tests {
		t.Run(test.field, func(t *testing.T) {
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

			mock.ExpectExec(`UPDATE songs\s+SET `+test.column+` = \$1`).
				WithArgs("updated text", "user-id", "song-id", "band-id").
				WillReturnResult(sqlmock.NewResult(0, 1))

			if err := SongsTableUpdateSongTextField(
				t.Context(),
				"song-id",
				"band-id",
				"user-id",
				test.field,
				"updated text",
			); err != nil {
				t.Fatalf("update %s: %v", test.field, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("database expectations: %v", err)
			}
		})
	}
}

func TestSongsTableUpdateSongTextFieldRejectsInvalidField(t *testing.T) {
	err := SongsTableUpdateSongTextField(
		t.Context(),
		"song-id",
		"band-id",
		"user-id",
		"title",
		"updated text",
	)
	if err == nil {
		t.Fatal("update invalid field: expected error")
	}
}

func TestSongsTableUpdateSongTextFieldRequiresMatchingSong(t *testing.T) {
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

	mock.ExpectExec(`UPDATE songs\s+SET lyrics = \$1`).
		WithArgs("updated text", "user-id", "song-id", "band-id").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = SongsTableUpdateSongTextField(
		t.Context(),
		"song-id",
		"band-id",
		"user-id",
		"lyrics",
		"updated text",
	)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("update missing song error = %v; want sql.ErrNoRows", err)
	}
}
