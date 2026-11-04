package database

import (
	"bandplan/src/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSetlistsTableCreateSetlistMakesFirstSetlistPrimary(t *testing.T) {
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

	setlist := models.Setlist{
		BandID:      "band-id",
		Name:        "First Setlist",
		Slug:        "first-setlist",
		CreatedBy:   "user-id",
		UpdatedBy:   "user-id",
		ArtworkID:   "",
		ArtworkPath: "",
	}

	mock.ExpectExec(`(?s)INSERT INTO setlists.*is_primary.*NOT EXISTS`).
		WithArgs(
			sqlmock.AnyArg(),
			setlist.BandID,
			setlist.Name,
			setlist.Slug,
			setlist.Explicit,
			setlist.Notes,
			setlist.ArtworkID,
			setlist.ArtworkPath,
			setlist.CreatedBy,
			setlist.UpdatedBy,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := SetlistsTableCreateSetlist(setlist); err != nil {
		t.Fatalf("create setlist: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}
