package database

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestChecklistItemsTableResetChecklist(t *testing.T) {
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

	checklistID := "checklist-id"
	userID := "user-id"
	bandID := "band-id"

	mock.ExpectBegin()
	mock.ExpectQuery(`(?s)SELECT checklist_id.*FROM checklist_lists.*checklist_id = \$1`).
		WithArgs(checklistID, userID, bandID).
		WillReturnRows(sqlmock.NewRows([]string{"checklist_id"}).AddRow(checklistID))
	mock.ExpectExec(`(?s)UPDATE checklist_lists.*is_complete = FALSE.*completed_at = NULL.*completed_by = NULL`).
		WithArgs(checklistID, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`(?s)UPDATE checklist_items.*is_complete = FALSE.*completed_at = NULL.*completed_by = NULL`).
		WithArgs(checklistID, userID).
		WillReturnResult(sqlmock.NewResult(0, 3))
	mock.ExpectCommit()

	if err := ChecklistItemsTableResetChecklist(checklistID, userID, bandID); err != nil {
		t.Fatalf("reset checklist: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}

func TestChecklistItemsTableResetChecklistRejectsUnauthorizedUser(t *testing.T) {
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

	checklistID := "checklist-id"
	userID := "user-id"
	bandID := "band-id"

	mock.ExpectBegin()
	mock.ExpectQuery(`(?s)SELECT checklist_id.*FROM checklist_lists.*checklist_id = \$1`).
		WithArgs(checklistID, userID, bandID).
		WillReturnRows(sqlmock.NewRows([]string{"checklist_id"}))
	mock.ExpectRollback()

	err = ChecklistItemsTableResetChecklist(checklistID, userID, bandID)
	if !errors.Is(err, ErrChecklistNotFound) {
		t.Fatalf("reset checklist error = %v; want %v", err, ErrChecklistNotFound)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}
