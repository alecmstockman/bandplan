package helpers

import (
	"bandplan/src/database"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestHelperGetAuthenticatedUserAndBandSetsMissingSessionBand(t *testing.T) {
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

	const (
		token  = "session-token"
		userID = "user-id"
		bandID = "band-id"
	)
	tokenHash := HashSessionToken(token)

	mock.ExpectBegin()
	mock.ExpectQuery(`FROM users u`).
		WithArgs(tokenHash).
		WillReturnRows(authUserRows(userID, false))
	mock.ExpectQuery(`FROM sessions`).
		WithArgs(tokenHash, userID).
		WillReturnRows(emptyAuthBandRows())
	mock.ExpectRollback()

	mock.ExpectQuery(`FROM bands`).
		WithArgs(userID).
		WillReturnRows(authBandRows(bandID))
	mock.ExpectExec(`UPDATE sessions`).
		WithArgs(bandID, tokenHash, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectBegin()
	mock.ExpectQuery(`FROM users u`).
		WithArgs(tokenHash).
		WillReturnRows(authUserRows(userID, true))
	mock.ExpectQuery(`FROM sessions`).
		WithArgs(tokenHash, userID).
		WillReturnRows(authBandRows(bandID))
	mock.ExpectRollback()

	request := httptest.NewRequest("GET", "/", nil)
	request.AddCookie(&http.Cookie{Name: "session_token", Value: token})

	user, band, err := HelperGetAuthenticatedUserAndBand(request)
	if err != nil {
		t.Fatalf("authenticate user: %v", err)
	}
	if user.UserID != userID {
		t.Errorf("user ID = %q; want %q", user.UserID, userID)
	}
	if band.BandID != bandID {
		t.Errorf("band ID = %q; want %q", band.BandID, bandID)
	}
	if !user.Permissions.EditBand {
		t.Error("EditBand permission = false; want refreshed permissions")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}

func authUserRows(userID string, editBand bool) *sqlmock.Rows {
	now := time.Now()
	return sqlmock.NewRows([]string{
		"id", "user_id", "name", "first_name", "last_name", "display_name",
		"email", "slug", "password_hash", "is_admin", "edit_band", "manage_members",
		"create_songs", "edit_songs", "delete_songs", "create_setlists", "edit_setlists",
		"delete_setlists", "create_events", "edit_events", "delete_events", "create_todos",
		"edit_todos", "delete_todos", "profile_image_id", "profile_image_path", "timezone",
		"is_email_verified", "last_login", "created_at", "updated_at",
	}).AddRow(
		1, userID, "Test User", "Test", "User", "Test User", "test@example.com", "test-user",
		"password-hash", true, editBand, false, false, false, false, false, false, false, false,
		false, false, false, false, false, "", "", "UTC", true, now, now, now,
	)
}

func authBandRows(bandID string) *sqlmock.Rows {
	now := time.Now()
	return sqlmock.NewRows(authBandColumns()).
		AddRow(1, bandID, "Test Band", "test-band", "", "", now, "user-id", now, "user-id")
}

func emptyAuthBandRows() *sqlmock.Rows {
	return sqlmock.NewRows(authBandColumns())
}

func authBandColumns() []string {
	return []string{
		"id", "band_id", "name", "slug", "profile_image_id", "profile_image_path",
		"created_at", "created_by", "updated_at", "updated_by",
	}
}
