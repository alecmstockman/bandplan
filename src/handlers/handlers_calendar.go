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

	data := models.MenuPageData{
		User: user,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "calendar.html", data)
}

func (h *Handler) HandlerCalendarEvents(w http.ResponseWriter, r *http.Request) {
	fmt.Println("- HandlerCalendarEvents")

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

	fmt.Println("Event one: ", events[0])

	var calendarEvents []models.CalendarEvent

	for _, event := range events {
		ce := models.CalendarEvent{
			ID:    event.EventID,
			Title: event.Name,
		}

		if event.StartTime != nil {
			ce.Start = event.StartTime.Format(time.RFC3339)
		}

		if event.EndTime != nil {
			ce.End = event.EndTime.Format(time.RFC3339)
		}

		ce.URL = fmt.Sprintf(
			"/event?event-id=%s&return-url=calendar",
			event.EventID,
		)

		calendarEvents = append(calendarEvents, ce)
	}

	fmt.Println("\nCalendar event: ", calendarEvents)

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(calendarEvents); err != nil {
		slog.Error(
			"unable to encode calendar events",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
	}
}
