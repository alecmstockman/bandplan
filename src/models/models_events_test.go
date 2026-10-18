package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestEventJSONUsesFullCalendarFields(t *testing.T) {
	eventDate := time.Date(2026, time.October, 18, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2026, time.October, 18, 22, 30, 0, 0, time.UTC)
	event := Event{
		ID:        42,
		EventID:   "event-123",
		Name:      "Band Practice",
		EventDate: eventDate,
		EndTime:   &endTime,
	}

	encoded, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("marshal event: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(encoded, &result); err != nil {
		t.Fatalf("unmarshal event JSON: %v", err)
	}

	if result["id"] != "event-123" {
		t.Errorf("id = %v, want event-123", result["id"])
	}
	if result["title"] != "Band Practice" {
		t.Errorf("title = %v, want Band Practice", result["title"])
	}
	if result["start"] != eventDate.Format(time.RFC3339) {
		t.Errorf("start = %v, want %s", result["start"], eventDate.Format(time.RFC3339))
	}
	if result["end"] != endTime.Format(time.RFC3339) {
		t.Errorf("end = %v, want %s", result["end"], endTime.Format(time.RFC3339))
	}
	if _, exists := result["ID"]; exists {
		t.Error("internal database ID must not be included in JSON")
	}
}
