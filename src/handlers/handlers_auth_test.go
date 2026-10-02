package handlers

import (
	"bandplan/src/database"
	"bandplan/src/services"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lib/pq"
)

func TestWriteRegistrationError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{name: "invalid", err: &services.RegistrationError{Kind: services.RegistrationInvalid, Err: errors.New("invalid")}, wantStatus: http.StatusBadRequest},
		{name: "forbidden", err: &services.RegistrationError{Kind: services.RegistrationForbidden, Err: errors.New("forbidden")}, wantStatus: http.StatusForbidden},
		{name: "not found", err: &services.RegistrationError{Kind: services.RegistrationNotFound, Err: errors.New("not found")}, wantStatus: http.StatusNotFound},
		{name: "conflict", err: &services.RegistrationError{Kind: services.RegistrationConflict, Err: errors.New("conflict")}, wantStatus: http.StatusConflict},
		{name: "expired", err: &services.RegistrationError{Kind: services.RegistrationExpired, Err: errors.New("expired")}, wantStatus: http.StatusGone},
		{name: "rate limited", err: &services.RegistrationError{Kind: services.RegistrationRateLimited, Err: errors.New("rate limited")}, wantStatus: http.StatusTooManyRequests},
		{name: "internal", err: &services.RegistrationError{Kind: services.RegistrationInternal, Err: errors.New("internal")}, wantStatus: http.StatusInternalServerError},
		{name: "database registration expired", err: database.ErrRegistrationExpired, wantStatus: http.StatusGone},
		{name: "database access code expired", err: database.ErrAccessCodeExpired, wantStatus: http.StatusGone},
		{name: "database not found", err: sql.ErrNoRows, wantStatus: http.StatusNotFound},
		{name: "database conflict", err: &pq.Error{Code: "23505"}, wantStatus: http.StatusConflict},
		{name: "unknown", err: errors.New("unexpected"), wantStatus: http.StatusInternalServerError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			writeRegistrationError(recorder, test.err)

			if recorder.Code != test.wantStatus {
				t.Fatalf("status = %d; want %d", recorder.Code, test.wantStatus)
			}
			if test.wantStatus == http.StatusTooManyRequests && recorder.Header().Get("Retry-After") != "10" {
				t.Fatalf("Retry-After = %q; want %q", recorder.Header().Get("Retry-After"), "10")
			}
		})
	}
}
