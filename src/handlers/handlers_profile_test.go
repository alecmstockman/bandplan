package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHandlerBandsPageMissingAuthContext(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/bands", nil)
	recorder := httptest.NewRecorder()

	handler := Handler{}
	handler.HandlerBandsPage(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusInternalServerError {
		t.Errorf("status code = %d; want %d", response.StatusCode, http.StatusInternalServerError)
	}
}

func TestHandlerBandPageMissingAuthContext(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/band?band-id=band-id", nil)
	recorder := httptest.NewRecorder()

	handler := Handler{}
	handler.HandlerBandPage(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusInternalServerError {
		t.Errorf("status code = %d; want %d", response.StatusCode, http.StatusInternalServerError)
	}
}

func TestHandlerBandPageMissingBandID(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/band", nil)
	request = request.WithContext(context.WithValue(request.Context(), AuthContextKey, AuthContext{}))
	recorder := httptest.NewRecorder()

	handler := Handler{}
	handler.HandlerBandPage(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("status code = %d; want %d", response.StatusCode, http.StatusBadRequest)
	}
}

func TestHandlerProfileFeedbackMissingAuthContext(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/profile/feedback", nil)
	recorder := httptest.NewRecorder()

	handler := Handler{}
	handler.HandlerProfileFeedback(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusInternalServerError {
		t.Errorf("status code = %d; want %d", response.StatusCode, http.StatusInternalServerError)
	}
}

func TestHandlerProfileFeedbackRejectsBlankFeedback(t *testing.T) {
	form := url.Values{"feedback": {"   "}}
	request := httptest.NewRequest(http.MethodPost, "/profile/feedback", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request = request.WithContext(context.WithValue(request.Context(), AuthContextKey, AuthContext{}))
	recorder := httptest.NewRecorder()

	handler := Handler{}
	handler.HandlerProfileFeedback(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("status code = %d; want %d", response.StatusCode, http.StatusBadRequest)
	}
}

func TestHandlerProfileFeedbackRejectsLongFeedback(t *testing.T) {
	form := url.Values{"feedback": {strings.Repeat("a", 4001)}}
	request := httptest.NewRequest(http.MethodPost, "/profile/feedback", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request = request.WithContext(context.WithValue(request.Context(), AuthContextKey, AuthContext{}))
	recorder := httptest.NewRecorder()

	handler := Handler{}
	handler.HandlerProfileFeedback(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("status code = %d; want %d", response.StatusCode, http.StatusBadRequest)
	}
}
