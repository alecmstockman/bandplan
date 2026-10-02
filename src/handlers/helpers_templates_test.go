package handlers

import (
	"bandplan/src/models"
	"html/template"
	"os"
	"strings"
	"testing"
	"time"
)

func TestFormatTimeValue(t *testing.T) {
	formatTimeValue := funcMap["formatTimeValue"].(func(*time.Time, string) string)
	value := time.Date(2026, time.September, 17, 15, 30, 0, 0, time.UTC)

	tests := []struct {
		name     string
		value    *time.Time
		timeZone string
		want     string
	}{
		{name: "local time", value: &value, timeZone: "America/Chicago", want: "10:30"},
		{name: "invalid timezone fallback", value: &value, timeZone: "invalid", want: "15:30"},
		{name: "nil time", value: nil, timeZone: "UTC", want: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := formatTimeValue(test.value, test.timeZone); got != test.want {
				t.Errorf("formatTimeValue() = %q; want %q", got, test.want)
			}
		})
	}
}

func TestFormatEventAges(t *testing.T) {
	formatEventAges := funcMap["formatEventAges"].(func(models.EventAges) string)

	tests := []struct {
		ages models.EventAges
		want string
	}{
		{ages: models.EventNA, want: "N/A"},
		{ages: models.EventAllAges, want: "All Ages"},
		{ages: models.Event18Plus, want: "18+"},
		{ages: models.Event21Plus, want: "21+"},
		{ages: models.EventAges("unknown"), want: "N/A"},
	}

	for _, test := range tests {
		if got := formatEventAges(test.ages); got != test.want {
			t.Errorf("formatEventAges(%q) = %q; want %q", test.ages, got, test.want)
		}
	}
}

func TestEventTemplatesParse(t *testing.T) {
	_, err := template.New("").Funcs(funcMap).ParseFiles(
		"../../templates/events/event_create.html",
		"../../templates/events/event-edit.html",
		"../../templates/events/event.html",
	)
	if err != nil {
		t.Fatalf("parse event templates: %v", err)
	}
}

func TestNativeFallbackFormsIncludeCSRFToken(t *testing.T) {
	tests := []struct {
		path   string
		action string
	}{
		{path: "../../templates/chats/chat_create.html", action: `action="/chats/create"`},
		{path: "../../templates/setlists/setlist-add.html", action: `action="/setlists/create"`},
		{path: "../../templates/setlists/setlist-edit.html", action: `action="/setlist/update"`},
		{path: "../../templates/profile/settings.html", action: `action="/logout"`},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			contents, err := os.ReadFile(test.path)
			if err != nil {
				t.Fatalf("read template: %v", err)
			}

			templateContents := string(contents)
			if !strings.Contains(templateContents, test.action) {
				t.Errorf("template does not contain native fallback %s", test.action)
			}
			if !strings.Contains(templateContents, `name="gorilla.csrf.Token"`) {
				t.Error("template does not contain a CSRF form field")
			}
		})
	}
}
