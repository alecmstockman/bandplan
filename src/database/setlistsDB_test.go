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

func TestSetlistsTableUpdateSetlistSwitchesPrimaryInTransaction(t *testing.T) {
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
		SetlistID:   "setlist-id",
		BandID:      "band-id",
		Name:        "Primary Setlist",
		Slug:        "primary-setlist",
		Notes:       "notes",
		ArtworkID:   "artwork-id",
		ArtworkPath: "artwork-path",
		IsPrimary:   true,
		UpdatedBy:   "user-id",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`(?s)UPDATE setlists.*SET is_primary = FALSE`).
		WithArgs(setlist.BandID, setlist.SetlistID, "user-id").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`(?s)UPDATE setlists.*is_primary = \$7`).
		WithArgs(
			setlist.Name,
			setlist.Slug,
			setlist.Explicit,
			setlist.Notes,
			setlist.ArtworkID,
			setlist.ArtworkPath,
			setlist.IsPrimary,
			setlist.UpdatedBy,
			setlist.SetlistID,
			setlist.BandID,
			"user-id",
		).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	updated, err := SetlistsTableUpdateSetlist(setlist, "user-id")
	if err != nil {
		t.Fatalf("update setlist: %v", err)
	}
	if !updated {
		t.Fatal("update setlist returned false; want true")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}
