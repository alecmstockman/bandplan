package handlers

import (
	"net/http"
	"net/http/httptest"
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
