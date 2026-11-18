package database

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestBandMembersTableRemoveMemberByIDRemovesChatMemberships(t *testing.T) {
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

	memberID := "member-id"
	bandID := "band-id"

	mock.ExpectBegin()
	mock.ExpectExec(`(?s)DELETE FROM chat_members cm.*USING chats c.*cm\.user_id = \$1.*c\.band_id = \$2`).
		WithArgs(memberID, bandID).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectExec(`(?s)DELETE FROM band_members.*user_id = \$1.*band_id = \$2`).
		WithArgs(memberID, bandID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	if err := BandMembersTableRemoveMemberByID(t.Context(), memberID, bandID); err != nil {
		t.Fatalf("remove band member: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}

func TestBandMembersTableRemoveMemberByIDRollsBackOnFailure(t *testing.T) {
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

	memberID := "member-id"
	bandID := "band-id"
	removeErr := errors.New("remove member failed")

	mock.ExpectBegin()
	mock.ExpectExec(`(?s)DELETE FROM chat_members cm.*USING chats c.*cm\.user_id = \$1.*c\.band_id = \$2`).
		WithArgs(memberID, bandID).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectExec(`(?s)DELETE FROM band_members.*user_id = \$1.*band_id = \$2`).
		WithArgs(memberID, bandID).
		WillReturnError(removeErr)
	mock.ExpectRollback()

	err = BandMembersTableRemoveMemberByID(t.Context(), memberID, bandID)
	if !errors.Is(err, removeErr) {
		t.Fatalf("remove band member error = %v; want %v", err, removeErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}
