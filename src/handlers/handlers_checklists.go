package handlers

import (
	"bandplan/src/database"
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gorilla/csrf"
)

func (h Handler) HandlerChecklistsPage(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerChecklistsPage")

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

	err = h.Tmpl.ExecuteTemplate(w, "checklists.html", data)
}

func (h Handler) HandlerChecklistCreatePage(w http.ResponseWriter, r *http.Request) {

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

	songs, err := database.SongsTableGetSongNameAndID(user.UserID, band.BandID)
	if err != nil {
		slog.Error(
			"unable to get song names and IDs",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
	}

	setlists, err := database.SetlistsTableGetSetlistNamesAndIDs(band.BandID, user.UserID)
	if err != nil {
		slog.Error(
			"unable to get setlist names and IDs",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
	}

	events, err := database.EventsTableGetEventNameAndID(user.UserID, band.BandID)
	if err != nil {
		slog.Error(
			"unable to get event names and IDs",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
	}

	members, err := database.BandMembersGetMemberNameAndID(user.UserID, band.BandID)
	if err != nil {
		slog.Error(
			"unable to get band member names and IDs",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
	}

	data := models.ChecklistCreatePageData{
		CSRFToken: csrf.Token(r),
		User:      auth.User,
		Band:      auth.CurrentBand,
		Checklist: models.Checklist{},
		Songs:     songs,
		Setlists:  setlists,
		Events:    events,
		Members:   members,
	}

	err = h.Tmpl.ExecuteTemplate(w, "checklist_create.html", data)
	if err != nil {
		slog.Error(
			"unable to execute checklist_create.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to load checklist creation page", http.StatusInternalServerError)
		return
	}
}

func (h Handler) HandlerChecklistAdd(w http.ResponseWriter, r *http.Request) {

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

	r.Body = http.MaxBytesReader(w, r.Body, 16*1024)

	if err := r.ParseForm(); err != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "Invalid form", http.StatusBadRequest)
		}
		return
	}

	name := strings.TrimSpace(r.FormValue("checklist-name"))
	if name == "" {
		http.Error(w, "Checklist name is required", http.StatusBadRequest)
		return
	}
	if utf8.RuneCountInString(name) > 120 {
		http.Error(w, "Checklist name is too long", http.StatusBadRequest)
		return
	}

	dueDate, dueTime, dueTimezone, err := parseOptionalDueFields(
		strings.TrimSpace(r.FormValue("checklist-due-date")),
		strings.TrimSpace(r.FormValue("checklist-due-time")),
		strings.TrimSpace(r.FormValue("checklist-due-timezone")),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newChecklist := models.Checklist{
		Name:        name,
		Description: strings.TrimSpace(r.FormValue("checklist-description")),
		DueDate:     dueDate,
		DueTime:     dueTime,
		DueTimezone: dueTimezone,
		CreatedBy:   auth.User.UserID,
		UpdatedBy:   auth.User.UserID,
	}

	switch strings.TrimSpace(r.FormValue("checklist-owner")) {
	case "Personal:" + auth.User.UserID:
		userID := auth.User.UserID
		newChecklist.UserID = &userID
	case "Band:" + auth.CurrentBand.BandID:
		bandID := auth.CurrentBand.BandID
		newChecklist.BandID = &bandID
	default:
		http.Error(w, "Invalid checklist owner", http.StatusBadRequest)
		return
	}

	optionalValue := func(field string) *string {
		value := strings.TrimSpace(r.FormValue(field))
		if value == "" || value == "none" {
			return nil
		}
		return &value
	}

	newChecklist.AssignedTo = optionalValue("checklist-assigned-to")
	newChecklist.SongID = optionalValue("checklist-song")
	newChecklist.SetlistID = optionalValue("checklist-setlist")
	newChecklist.EventID = optionalValue("checklist-event")

	err = database.ChecklistListsTableCreateChecklist(
		newChecklist,
		auth.User.UserID,
		auth.CurrentBand.BandID,
	)
	if errors.Is(err, database.ErrInvalidChecklistReferences) {
		http.Error(w, "Invalid checklist selection", http.StatusBadRequest)
		return
	}
	if err != nil {
		slog.Error(
			"unable to create checklist",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to create checklist", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/checklists", http.StatusSeeOther)
}
