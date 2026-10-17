package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerToDoListDeleteMissingAuthContext(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/todo/delete", nil)
	recorder := httptest.NewRecorder()

	handler := Handler{}
	handler.HandlerToDoListDelete(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusInternalServerError {
		t.Errorf("status code = %d; want %d", response.StatusCode, http.StatusInternalServerError)
	}
}

func TestHandlerToDoListDeleteMissingListID(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/todo/delete", nil)
	request = request.WithContext(context.WithValue(request.Context(), AuthContextKey, AuthContext{}))
	recorder := httptest.NewRecorder()

	handler := Handler{}
	handler.HandlerToDoListDelete(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("status code = %d; want %d", response.StatusCode, http.StatusBadRequest)
	}
}
