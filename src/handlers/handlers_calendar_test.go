package handlers

import (
	"bandplan/src/models"
	"net/url"
	"testing"
	"time"
)

func TestCalendarEventFromTodoItemTimed(t *testing.T) {
	dueDate := time.Date(2026, time.October, 5, 0, 0, 0, 0, time.UTC)
	dueTime := time.Date(1, time.January, 1, 14, 30, 0, 0, time.UTC)
	timezone := "America/Chicago"

	event := calendarEventFromTodoItem(models.ToDoItem{
		ItemID:      "task-1",
		ToDoListID:  "list-1",
		Name:        "Load gear",
		DueDate:     &dueDate,
		DueTime:     &dueTime,
		DueTimezone: &timezone,
	})

	if event.Start != "2026-10-05T14:30:00-05:00" {
		t.Errorf("start = %q; want %q", event.Start, "2026-10-05T14:30:00-05:00")
	}
	if event.AllDay {
		t.Error("timed task unexpectedly marked all-day")
	}

	taskURL, err := url.Parse(event.URL)
	if err != nil {
		t.Fatalf("parse task URL: %v", err)
	}
	if taskURL.Path != "/todo/item/page" {
		t.Errorf("URL path = %q; want %q", taskURL.Path, "/todo/item/page")
	}
	if got := taskURL.Query().Get("item-id"); got != "task-1" {
		t.Errorf("item-id = %q; want %q", got, "task-1")
	}
	if got := taskURL.Query().Get("list-id"); got != "list-1" {
		t.Errorf("list-id = %q; want %q", got, "list-1")
	}
	if got := taskURL.Query().Get("return"); got != "calendar" {
		t.Errorf("return = %q; want %q", got, "calendar")
	}
}

func TestCalendarEventFromTodoItemAllDay(t *testing.T) {
	dueDate := time.Date(2026, time.November, 24, 0, 0, 0, 0, time.UTC)

	event := calendarEventFromTodoItem(models.ToDoItem{
		ItemID:     "task-2",
		ToDoListID: "list-2",
		Name:       "Confirm venue",
		DueDate:    &dueDate,
	})

	if event.Start != "2026-11-24" {
		t.Errorf("start = %q; want %q", event.Start, "2026-11-24")
	}
	if !event.AllDay {
		t.Error("date-only task was not marked all-day")
	}
}

func TestCalendarEventFromTodoItemInvalidTimezoneUsesUTC(t *testing.T) {
	dueDate := time.Date(2026, time.November, 24, 0, 0, 0, 0, time.UTC)
	dueTime := time.Date(1, time.January, 1, 9, 15, 0, 0, time.UTC)
	timezone := "invalid"

	event := calendarEventFromTodoItem(models.ToDoItem{
		ItemID:      "task-3",
		ToDoListID:  "list-3",
		Name:        "Post flyer",
		DueDate:     &dueDate,
		DueTime:     &dueTime,
		DueTimezone: &timezone,
	})

	if event.Start != "2026-11-24T09:15:00Z" {
		t.Errorf("start = %q; want %q", event.Start, "2026-11-24T09:15:00Z")
	}
}
