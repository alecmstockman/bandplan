package handlers

import (
	"bandplan/src/database"
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
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

	calendarEvents := make([]models.CalendarEvent, 0, len(events))
	calendarExtendedProps := models.CalendarExtendedProps{
		MultidayPriority: "1",
	}

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

		ce.URL = fmt.Sprintf(
			"/event?event-id=%s&return-url=calendar",
			event.EventID,
		)

		calendarEvents = append(calendarEvents, ce)
		ce.ExtendedProps = calendarExtendedProps
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
