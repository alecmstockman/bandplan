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

	mock.ExpectQuery(`FROM sessions s`).
		WithArgs(tokenHash).
		WillReturnRows(authUserAndBandRows(userID, false, ""))

	mock.ExpectQuery(`FROM bands`).
		WithArgs(userID).
		WillReturnRows(authBandRows(bandID))
	mock.ExpectExec(`UPDATE sessions`).
		WithArgs(bandID, tokenHash, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(`FROM sessions s`).
		WithArgs(tokenHash).
		WillReturnRows(authUserAndBandRows(userID, true, bandID))

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

func authUserAndBandRows(userID string, editBand bool, bandID string) *sqlmock.Rows {
	now := time.Now()
	columns := []string{
		"id", "user_id", "name", "first_name", "last_name", "display_name",
		"email", "slug", "password_hash", "is_admin", "edit_band", "manage_members",
		"create_songs", "edit_songs", "delete_songs", "create_setlists", "edit_setlists",
		"delete_setlists", "create_events", "edit_events", "delete_events", "create_todos",
		"edit_todos", "delete_todos", "profile_image_id", "profile_image_path", "timezone",
		"is_email_verified", "last_login", "created_at", "updated_at",
	}
	columns = append(columns, authBandColumns()...)

	rows := sqlmock.NewRows(columns)
	if bandID == "" {
		return rows.AddRow(
			1, userID, "Test User", "Test", "User", "Test User", "test@example.com", "test-user",
			"password-hash", true, editBand, false, false, false, false, false, false, false, false,
			false, false, false, false, false, "", "", "UTC", true, now, now, now,
			nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		)
	}

	return rows.AddRow(
		1, userID, "Test User", "Test", "User", "Test User", "test@example.com", "test-user",
		"password-hash", true, editBand, false, false, false, false, false, false, false, false,
		false, false, false, false, false, "", "", "UTC", true, now, now, now,
		1, bandID, "Test Band", "test-band", "", "", now, userID, now, userID,
	)
}

func authBandRows(bandID string) *sqlmock.Rows {
	now := time.Now()
	return sqlmock.NewRows(authBandColumns()).
		AddRow(1, bandID, "Test Band", "test-band", "", "", now, "user-id", now, "user-id")
}

func authBandColumns() []string {
	return []string{
		"id", "band_id", "name", "slug", "profile_image_id", "profile_image_path",
		"created_at", "created_by", "updated_at", "updated_by",
	}
}
