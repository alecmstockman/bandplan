package handlers

import (
	"bandplan/src/models"
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSongTextPageData(t *testing.T) {
	song := models.Song{
		SongID:      "song-id",
		Title:       "Song Title",
		Lyrics:      "song lyrics",
		Description: "song description",
		Notes:       "song notes",
	}

	tests := []struct {
		field string
		label string
		value string
	}{
		{field: "lyrics", label: "Lyrics", value: song.Lyrics},
		{field: "description", label: "Description", value: song.Description},
		{field: "notes", label: "Notes", value: song.Notes},
	}

	for _, test := range tests {
		t.Run(test.field, func(t *testing.T) {
			data, ok := songTextPageData(song, test.field, "csrf-token")
			if !ok {
				t.Fatalf("songTextPageData(%q) returned invalid", test.field)
			}
			if data.Label != test.label || data.Value != test.value {
				t.Fatalf("data = %#v; want label %q and value %q", data, test.label, test.value)
			}
		})
	}

	if _, ok := songTextPageData(song, "title", "csrf-token"); ok {
		t.Fatal("songTextPageData accepted invalid field")
	}
}

func TestHandlerSongsAdd_MissingAuthContext(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/song/add",
		nil,
	)
	recorder := httptest.NewRecorder()

	handler := Handler{}

	handler.HandlerSongSave(recorder, request)

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
		"/song/create",
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

	handler.HandlerSongSave(recorder, request)

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

	if location != "/song/create" {
		t.Errorf(
			"redirect location = %q; want %q",
			location,
			"/song/create",
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
		"/song/create",
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

	handler.HandlerSongSave(recorder, request)

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
