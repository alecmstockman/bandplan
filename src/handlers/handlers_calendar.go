package handlers

import (
	"bandplan/src/database"
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/csrf"
)

func (h Handler) HandlerCalendar(w http.ResponseWriter, r *http.Request) {

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	notifications, err := database.NotifTableGetNotificationsByUserID(r.Context(), user.UserID, band.BandID)
	if err != nil {
		slog.Error(
			"unable to load notifications",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load notifications", http.StatusInternalServerError)
		return
	}

	data := models.MenuPageData{
		CSRFToken:     csrf.Token(r),
		User:          user,
		Band:          band,
		Notifications: notifications,
	}

	err = h.Tmpl.ExecuteTemplate(w, "calendar.html", data)
}

func (h *Handler) HandlerCalendarEvents(w http.ResponseWriter, r *http.Request) {
	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	events, err := database.EventsTableGetAllCalendarEventsByBandIDAndUserID(band.BandID, user.UserID)
	if err != nil {
		slog.Error(
			"unable to load calendar events",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load calendar events", http.StatusInternalServerError)
		return
	}

	tasks, err := database.TodoItemsTableGetItemsByDueDate(user.UserID, band.BandID)
	if err != nil {
		slog.Error(
			"unable to load calendar tasks",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load calendar tasks", http.StatusInternalServerError)
		return
	}

	calendarEvents := make([]models.CalendarEvent, 0, len(events)+len(tasks))

	for _, event := range events {
		ce := models.CalendarEvent{
			ID:    event.EventID,
			Title: event.Name,
		}

		location, err := time.LoadLocation(event.Timezone)
		if err != nil {
			location = time.UTC
		}

		if event.StartTime != nil {
			ce.Start = event.StartTime.Format(time.RFC3339)
		} else if event.EndTime != nil {
			startOfDay, _ := time.ParseInLocation("2006-01-02", event.EventDate.Format("2006-01-02"), location)
			ce.Start = startOfDay.Format(time.RFC3339)
		} else {
			ce.Start = event.EventDate.Format("2006-01-02")
			ce.AllDay = true
		}

		if event.EndTime != nil {
			ce.End = event.EndTime.Format(time.RFC3339)
		} else if event.EndDate != nil {
			exclusiveEnd := event.EndDate.AddDate(0, 0, 1)
			if ce.AllDay {
				ce.End = exclusiveEnd.Format("2006-01-02")
			} else {
				endOfRange, _ := time.ParseInLocation("2006-01-02", exclusiveEnd.Format("2006-01-02"), location)
				ce.End = endOfRange.Format(time.RFC3339)
			}
		}

		switch event.EventType {
		case models.EventTypeShow, models.EventTypeGig, models.EventTypeFestival:
			ce.BackgroundColor = "#166534"
		case models.EventTypeRehearsal, models.EventTypePractice:
			ce.BackgroundColor = "#0F766E"
		case models.EventTypeUnavailable:
			ce.BackgroundColor = "#B91C1C"
		case models.EventTypeOther:
			ce.BackgroundColor = "#854D0E"
		case models.EventTypeWriting:
			ce.BackgroundColor = "#6D28D9"
		case models.EventTypeRecording:
			ce.BackgroundColor = "#C2410C"
		case models.EventTypePhotos, models.EventTypeMeeting, models.EventTypePress:
			ce.BackgroundColor = "#A21CAF"
		default:
			ce.BackgroundColor = "#475569"
		}
		ce.BorderColor = ce.BackgroundColor
		ce.TextColor = "#fcfcfc"

		ce.URL = fmt.Sprintf(
			"/event?event-id=%s&return=calendar",
			event.EventID,
		)
		ce.ExtendedProps = models.CalendarExtendedProps{
			MultidayPriority: "1",
			EventType:        event.EventType,
			EventStatus:      event.EventStatus,
		}

		calendarEvents = append(calendarEvents, ce)

	}

	for _, task := range tasks {
		calendarEvents = append(calendarEvents, calendarEventFromTodoItem(task))
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(calendarEvents); err != nil {
		slog.Error(
			"unable to encode calendar events",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
	}
}

func calendarEventFromTodoItem(task models.ToDoItem) models.CalendarEvent {
	ce := models.CalendarEvent{
		ID:              task.ItemID,
		Title:           task.Name,
		BackgroundColor: "#1D4ED8",
		BorderColor:     "#1D4ED8",
		TextColor:       "#fcfcfc",
		ExtendedProps: models.CalendarExtendedProps{
			MultidayPriority: "1",
			EventType:        "task",
			EventStatus:      "active",
		},
	}

	location := time.UTC
	if task.DueTimezone != nil {
		if taskLocation, err := time.LoadLocation(*task.DueTimezone); err == nil {
			location = taskLocation
		}
	}

	if task.DueDate != nil {
		if task.DueTime != nil {
			dueAt := time.Date(
				task.DueDate.Year(),
				task.DueDate.Month(),
				task.DueDate.Day(),
				task.DueTime.Hour(),
				task.DueTime.Minute(),
				task.DueTime.Second(),
				task.DueTime.Nanosecond(),
				location,
			)
			ce.Start = dueAt.Format(time.RFC3339)
		} else {
			ce.Start = task.DueDate.Format("2006-01-02")
			ce.AllDay = true
		}
	}

	query := url.Values{}
	query.Set("item-id", task.ItemID)
	query.Set("list-id", task.ToDoListID)
	query.Set("return", "calendar")
	ce.URL = "/todo/item/page?" + query.Encode()

	return ce
}
