package handlers

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerSongsAdd_MissingAuthContext(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/songs/add",
		nil,
	)
	recorder := httptest.NewRecorder()

	handler := Handler{}

	handler.HandlerSongsAdd(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusInternalServerError {
		t.Errorf(
			"status code = %d; want %d",
			response.StatusCode,
			http.StatusInternalServerError,
		)
	}
}

func TestHandlerSongsAdd_MissingSongTitle(t *testing.T) {

	var body bytes.Buffer

	form := multipart.NewWriter(&body)

	if err := form.Close(); err != nil {
		t.Fatalf("unable to close multipart form: %v", err)
	}

	request := httptest.NewRequest(
		http.MethodPost,
		"/songs/add",
		&body,
	)

	request.Header.Set("Content-Type", form.FormDataContentType())

	ctx := context.WithValue(
		request.Context(),
		AuthContextKey,
		AuthContext{},
	)

	request = request.WithContext(ctx)

	recorder := httptest.NewRecorder()
	handler := Handler{}

	handler.HandlerSongsAdd(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusSeeOther {
		t.Errorf(
			"status code = %d; want %d",
			response.StatusCode,
			http.StatusSeeOther,
		)
	}

	location := response.Header.Get("Location")

	if location != "/songs/add" {
		t.Errorf(
			"redirect location = %q; want %q",
			location,
			"/songs/add",
		)
	}
}

func TestHandlerSongsAdd_InvalidReleaseDate(t *testing.T) {
	var body bytes.Buffer
	form := multipart.NewWriter(&body)

	if err := form.WriteField("song-title", "Test Song"); err != nil {
		t.Fatalf("unable to add song title: %v", err)
	}

	if err := form.WriteField("artist-name", "Test Artist"); err != nil {
		t.Fatalf("unable to add artist name: %v", err)
	}

	if err := form.WriteField("release-date", "not-a-date"); err != nil {
		t.Fatalf("unable to add release date: %v", err)
	}

	if err := form.WriteField("minutes", "0"); err != nil {
		t.Fatalf("unable to add minutes: %v", err)
	}

	if err := form.WriteField("seconds", "0"); err != nil {
		t.Fatalf("unable to add seconds: %v", err)
	}

	if err := form.Close(); err != nil {
		t.Fatalf("unable to close multipart form: %v", err)
	}

	request := httptest.NewRequest(
		http.MethodPost,
		"/songs/add",
		&body,
	)

	request.Header.Set("Content-Type", form.FormDataContentType())

	ctx := context.WithValue(
		request.Context(),
		AuthContextKey,
		AuthContext{},
	)

	request = request.WithContext(ctx)

	recorder := httptest.NewRecorder()
	handler := Handler{}

	handler.HandlerSongsAdd(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf(
			"status code = %d; want %d",
			response.StatusCode,
			http.StatusBadRequest,
		)
	}
}
