package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHandlerBandEditPageMissingAuthContext(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/band/edit?band-id=band-id", nil)
	recorder := httptest.NewRecorder()

	Handler{}.HandlerBandEditPage(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusInternalServerError {
		t.Errorf("status code = %d; want %d", response.StatusCode, http.StatusInternalServerError)
	}
}

func TestHandlerBandEditPageMissingBandID(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/band/edit", nil)
	request = request.WithContext(context.WithValue(request.Context(), AuthContextKey, AuthContext{}))
	recorder := httptest.NewRecorder()

	Handler{}.HandlerBandEditPage(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("status code = %d; want %d", response.StatusCode, http.StatusBadRequest)
	}
}

func TestHandlerBandUpdateMissingAuthContext(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/band/update", nil)
	recorder := httptest.NewRecorder()

	Handler{}.HandlerBandUpdate(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusInternalServerError {
		t.Errorf("status code = %d; want %d", response.StatusCode, http.StatusInternalServerError)
	}
}

func TestHandlerBandUpdateMissingBandID(t *testing.T) {
	form := url.Values{"band-name": {"Band"}, "timezone": {"UTC"}}
	request := httptest.NewRequest(http.MethodPost, "/band/update", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request = request.WithContext(context.WithValue(request.Context(), AuthContextKey, AuthContext{}))
	recorder := httptest.NewRecorder()

	Handler{}.HandlerBandUpdate(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("status code = %d; want %d", response.StatusCode, http.StatusBadRequest)
	}
}

func TestHandlerBandEmailVerificationRejectsInvalidToken(t *testing.T) {
	form := url.Values{"token": {"invalid"}}
	request := httptest.NewRequest(http.MethodPost, "/band/email/verification", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recorder := httptest.NewRecorder()

	Handler{}.HandlerBandEmailVerification(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("status code = %d; want %d", response.StatusCode, http.StatusBadRequest)
	}
}
