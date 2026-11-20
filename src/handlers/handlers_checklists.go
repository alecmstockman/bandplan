package handlers

import (
	"bandplan/src/database"
	"bandplan/src/helpers"
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
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

	checklists, err := database.ChecklistListsTableGetChecklists(user.UserID, band.BandID)
	if err != nil {
		slog.Error(
			"unable to load checklists",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load checklists", http.StatusInternalServerError)
		return
	}

	data := models.ChecklistsPageData{
		CSRFToken:     csrf.Token(r),
		User:          user,
		Band:          band,
		Checklists:    checklists,
		Notifications: notifications,
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
		CSRFToken:     csrf.Token(r),
		User:          user,
		Band:          band,
		Checklist:     models.Checklist{},
		Songs:         songs,
		Setlists:      setlists,
		Events:        events,
		Members:       members,
		Notifications: notifications,
	}

	err = h.Tmpl.ExecuteTemplate(w, "checklist-create.html", data)
	if err != nil {
		slog.Error(
			"unable to execute checklist-create.html",
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

	dueDate, dueTime, dueTimezone, err := helpers.ParseOptionalDueFields(
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

func (h Handler) HandlerChecklistPage(w http.ResponseWriter, r *http.Request) {
	auth, err := HelperGetAuthContext(r)
	if err != nil {
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	checklistID := strings.TrimSpace(r.URL.Query().Get("checklist-id"))
	if checklistID == "" {
		http.Error(w, "Checklist ID is required", http.StatusBadRequest)
		return
	}

	checklist, err := database.ChecklistListsTableGetChecklistByID(
		checklistID,
		auth.User.UserID,
		auth.CurrentBand.BandID,
	)
	if errors.Is(err, database.ErrChecklistNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		slog.Error("unable to load checklist", "error", err)
		http.Error(w, "Unable to load checklist", http.StatusInternalServerError)
		return
	}

	items, err := database.ChecklistItemsTableGetItems(checklistID)
	if err != nil {
		slog.Error("unable to load checklist items", "error", err)
		http.Error(w, "Unable to load checklist", http.StatusInternalServerError)
		return
	}
	checklist.Items = items

	data, err := loadChecklistPageData(r, auth.User, auth.CurrentBand, checklist)
	if err != nil {
		slog.Error("unable to load checklist options", "error", err)
		http.Error(w, "Unable to load checklist", http.StatusInternalServerError)
		return
	}

	if err := h.Tmpl.ExecuteTemplate(w, "checklist.html", data); err != nil {
		slog.Error("unable to render checklist", "error", err)
		http.Error(w, "Unable to load checklist", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerChecklistEditPage(w http.ResponseWriter, r *http.Request) {
	auth, err := HelperGetAuthContext(r)
	if err != nil {
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	checklistID := strings.TrimSpace(r.URL.Query().Get("checklist-id"))
	if checklistID == "" {
		http.Error(w, "Checklist ID is required", http.StatusBadRequest)
		return
	}
	checklist, err := database.ChecklistListsTableGetChecklistByID(
		checklistID,
		auth.User.UserID,
		auth.CurrentBand.BandID,
	)
	if errors.Is(err, database.ErrChecklistNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		slog.Error("unable to load checklist for editing", "error", err)
		http.Error(w, "Unable to load checklist", http.StatusInternalServerError)
		return
	}

	songs, setlists, events, members, err := loadChecklistOptions(auth.User, auth.CurrentBand)
	if err != nil {
		http.Error(w, "Unable to load checklist options", http.StatusInternalServerError)
		return
	}
	notifications, err := database.NotifTableGetNotificationsByUserID(
		r.Context(),
		auth.User.UserID,
		auth.CurrentBand.BandID,
	)
	if err != nil {
		http.Error(w, "Unable to load notifications", http.StatusInternalServerError)
		return
	}

	data := models.ChecklistCreatePageData{
		CSRFToken:     csrf.Token(r),
		User:          auth.User,
		Band:          auth.CurrentBand,
		Checklist:     checklist,
		Songs:         songs,
		Setlists:      setlists,
		Events:        events,
		Members:       members,
		Notifications: notifications,
	}
	if err := h.Tmpl.ExecuteTemplate(w, "checklist-create.html", data); err != nil {
		slog.Error("unable to render checklist edit page", "error", err)
		http.Error(w, "Unable to load checklist edit page", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerChecklistUpdate(w http.ResponseWriter, r *http.Request) {
	auth, err := HelperGetAuthContext(r)
	if err != nil {
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 16*1024)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	checklistID := strings.TrimSpace(r.FormValue("checklist-id"))
	name := strings.TrimSpace(r.FormValue("checklist-name"))
	if checklistID == "" || name == "" {
		http.Error(w, "Checklist ID and name are required", http.StatusBadRequest)
		return
	}
	if utf8.RuneCountInString(name) > 120 {
		http.Error(w, "Checklist name is too long", http.StatusBadRequest)
		return
	}

	dueDate, dueTime, dueTimezone, err := helpers.ParseOptionalDueFields(
		strings.TrimSpace(r.FormValue("checklist-due-date")),
		strings.TrimSpace(r.FormValue("checklist-due-time")),
		strings.TrimSpace(r.FormValue("checklist-due-timezone")),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	checklist := models.Checklist{
		ChecklistID: checklistID,
		Name:        name,
		Description: strings.TrimSpace(r.FormValue("checklist-description")),
		AssignedTo:  optionalChecklistValue(r, "checklist-assigned-to"),
		SongID:      optionalChecklistValue(r, "checklist-song"),
		SetlistID:   optionalChecklistValue(r, "checklist-setlist"),
		EventID:     optionalChecklistValue(r, "checklist-event"),
		DueDate:     dueDate,
		DueTime:     dueTime,
		DueTimezone: dueTimezone,
		UpdatedBy:   auth.User.UserID,
	}
	switch strings.TrimSpace(r.FormValue("checklist-owner")) {
	case "Personal:" + auth.User.UserID:
		userID := auth.User.UserID
		checklist.UserID = &userID
	case "Band:" + auth.CurrentBand.BandID:
		bandID := auth.CurrentBand.BandID
		checklist.BandID = &bandID
	default:
		http.Error(w, "Invalid checklist owner", http.StatusBadRequest)
		return
	}

	err = database.ChecklistListsTableUpdateChecklist(
		checklist,
		auth.User.UserID,
		auth.CurrentBand.BandID,
	)
	if errors.Is(err, database.ErrChecklistNotFound) {
		http.NotFound(w, r)
		return
	}
	if errors.Is(err, database.ErrInvalidChecklistReferences) {
		http.Error(w, "Invalid checklist selection", http.StatusBadRequest)
		return
	}
	if err != nil {
		slog.Error("unable to update checklist", "error", err)
		http.Error(w, "Unable to update checklist", http.StatusInternalServerError)
		return
	}

	http.Redirect(
		w,
		r,
		"/checklist?checklist-id="+url.QueryEscape(checklistID),
		http.StatusSeeOther,
	)
}

func (h Handler) HandlerChecklistDuplicate(w http.ResponseWriter, r *http.Request) {
	auth, err := HelperGetAuthContext(r)
	if err != nil {
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	checklistID := strings.TrimSpace(r.FormValue("checklist-id"))
	if checklistID == "" {
		http.Error(w, "Checklist ID is required", http.StatusBadRequest)
		return
	}
	newChecklistID, err := database.ChecklistListsTableDuplicateChecklist(
		checklistID,
		auth.User.UserID,
		auth.CurrentBand.BandID,
	)
	if errors.Is(err, database.ErrChecklistNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		slog.Error("unable to duplicate checklist", "error", err)
		http.Error(w, "Unable to duplicate checklist", http.StatusInternalServerError)
		return
	}

	http.Redirect(
		w,
		r,
		"/checklist/edit?checklist-id="+url.QueryEscape(newChecklistID),
		http.StatusSeeOther,
	)
}

func (h Handler) HandlerChecklistDelete(w http.ResponseWriter, r *http.Request) {
	auth, err := HelperGetAuthContext(r)
	if err != nil {
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	checklistID := strings.TrimSpace(r.FormValue("checklist-id"))
	if checklistID == "" {
		http.Error(w, "Checklist ID is required", http.StatusBadRequest)
		return
	}
	err = database.ChecklistListsTableDeleteChecklist(
		checklistID,
		auth.User.UserID,
		auth.CurrentBand.BandID,
	)
	if errors.Is(err, database.ErrChecklistNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		slog.Error("unable to delete checklist", "error", err)
		http.Error(w, "Unable to delete checklist", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/checklists", http.StatusSeeOther)
}

func (h Handler) HandlerChecklistItemAdd(w http.ResponseWriter, r *http.Request) {
	auth, err := HelperGetAuthContext(r)
	if err != nil {
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 16*1024)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	checklistID := strings.TrimSpace(r.FormValue("checklist-id"))
	name := strings.TrimSpace(r.FormValue("checklist-item-name"))
	if checklistID == "" || name == "" {
		http.Error(w, "Checklist and item name are required", http.StatusBadRequest)
		return
	}
	if utf8.RuneCountInString(name) > 120 {
		http.Error(w, "Checklist item name is too long", http.StatusBadRequest)
		return
	}

	dueDate, dueTime, dueTimezone, err := helpers.ParseOptionalDueFields(
		strings.TrimSpace(r.FormValue("checklist-item-due-date")),
		strings.TrimSpace(r.FormValue("checklist-item-due-time")),
		strings.TrimSpace(r.FormValue("checklist-item-due-timezone")),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item := models.ChecklistItem{
		ChecklistID: checklistID,
		Name:        name,
		Body:        strings.TrimSpace(r.FormValue("checklist-item-description")),
		SongID:      optionalChecklistValue(r, "checklist-item-song"),
		SetlistID:   optionalChecklistValue(r, "checklist-item-setlist"),
		EventID:     optionalChecklistValue(r, "checklist-item-event"),
		AssignedTo:  optionalChecklistValue(r, "checklist-item-assigned-to"),
		DueDate:     dueDate,
		DueTime:     dueTime,
		DueTimezone: dueTimezone,
		CreatedBy:   auth.User.UserID,
		UpdatedBy:   auth.User.UserID,
	}

	err = database.ChecklistItemsTableCreateItem(
		item,
		auth.User.UserID,
		auth.CurrentBand.BandID,
	)
	if errors.Is(err, database.ErrInvalidChecklistItemReferences) {
		http.Error(w, "Invalid checklist item selection", http.StatusBadRequest)
		return
	}
	if err != nil {
		slog.Error("unable to create checklist item", "error", err)
		http.Error(w, "Unable to create checklist item", http.StatusInternalServerError)
		return
	}

	http.Redirect(
		w,
		r,
		"/checklist?checklist-id="+url.QueryEscape(checklistID),
		http.StatusSeeOther,
	)
}

func (h Handler) HandlerChecklistItemEdit(w http.ResponseWriter, r *http.Request) {
	auth, err := HelperGetAuthContext(r)
	if err != nil {
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	checklistID := strings.TrimSpace(r.FormValue("checklist-id"))
	itemID := strings.TrimSpace(r.FormValue("checklist-item-id"))
	if checklistID == "" || itemID == "" {
		http.Error(w, "Checklist and item are required", http.StatusBadRequest)
		return
	}

	item, err := database.ChecklistItemsTableGetItemByID(
		itemID,
		checklistID,
		auth.User.UserID,
		auth.CurrentBand.BandID,
	)
	if errors.Is(err, database.ErrChecklistItemNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Unable to load checklist item", http.StatusInternalServerError)
		return
	}

	data, err := loadChecklistItemPageData(r, auth.User, auth.CurrentBand, item)
	if err != nil {
		http.Error(w, "Unable to load checklist item", http.StatusInternalServerError)
		return
	}
	if err := h.Tmpl.ExecuteTemplate(w, "checklist_item_edit", data); err != nil {
		slog.Error("unable to render checklist item edit", "error", err)
	}
}

func (h Handler) HandlerChecklistItemUpdate(w http.ResponseWriter, r *http.Request) {
	auth, err := HelperGetAuthContext(r)
	if err != nil {
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 16*1024)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	checklistID := strings.TrimSpace(r.FormValue("checklist-id"))
	itemID := strings.TrimSpace(r.FormValue("checklist-item-id"))
	name := strings.TrimSpace(r.FormValue("checklist-item-name"))
	if checklistID == "" || itemID == "" || name == "" {
		http.Error(w, "Checklist, item, and name are required", http.StatusBadRequest)
		return
	}
	if utf8.RuneCountInString(name) > 120 {
		http.Error(w, "Checklist item name is too long", http.StatusBadRequest)
		return
	}

	var completion *bool
	switch r.PostForm.Get("checklist-item-status") {
	case "":
	case "true":
		value := true
		completion = &value
	case "false":
		value := false
		completion = &value
	default:
		http.Error(w, "Invalid checklist item status", http.StatusBadRequest)
		return
	}

	dueDate, dueTime, dueTimezone, err := helpers.ParseOptionalDueFields(
		strings.TrimSpace(r.FormValue("checklist-item-due-date")),
		strings.TrimSpace(r.FormValue("checklist-item-due-time")),
		strings.TrimSpace(r.FormValue("checklist-item-due-timezone")),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item := models.ChecklistItem{
		ItemID:      itemID,
		ChecklistID: checklistID,
		Name:        name,
		Body:        strings.TrimSpace(r.FormValue("checklist-item-description")),
		SongID:      optionalChecklistValue(r, "checklist-item-song"),
		SetlistID:   optionalChecklistValue(r, "checklist-item-setlist"),
		EventID:     optionalChecklistValue(r, "checklist-item-event"),
		AssignedTo:  optionalChecklistValue(r, "checklist-item-assigned-to"),
		DueDate:     dueDate,
		DueTime:     dueTime,
		DueTimezone: dueTimezone,
		UpdatedBy:   auth.User.UserID,
	}

	updated, err := database.ChecklistItemsTableUpdateItem(
		item,
		auth.User.UserID,
		auth.CurrentBand.BandID,
		completion,
	)
	if errors.Is(err, database.ErrChecklistItemNotFound) {
		http.NotFound(w, r)
		return
	}
	if errors.Is(err, database.ErrInvalidChecklistItemReferences) {
		http.Error(w, "Invalid checklist item selection", http.StatusBadRequest)
		return
	}
	if err != nil {
		slog.Error("unable to update checklist item", "error", err)
		http.Error(w, "Unable to update checklist item", http.StatusInternalServerError)
		return
	}

	data, err := loadChecklistItemPageData(r, auth.User, auth.CurrentBand, updated)
	if err != nil {
		http.Error(w, "Unable to load updated checklist item", http.StatusInternalServerError)
		return
	}
	if err := h.Tmpl.ExecuteTemplate(w, "checklist_item_update", data); err != nil {
		slog.Error("unable to render updated checklist item", "error", err)
	}
}

func (h Handler) HandlerChecklistItemComplete(w http.ResponseWriter, r *http.Request) {
	h.setChecklistItemCompletion(w, r, true)
}

func (h Handler) HandlerChecklistItemReopen(w http.ResponseWriter, r *http.Request) {
	h.setChecklistItemCompletion(w, r, false)
}

func (h Handler) HandlerChecklistReset(w http.ResponseWriter, r *http.Request) {
	auth, err := HelperGetAuthContext(r)
	if err != nil {
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	checklistID := strings.TrimSpace(r.FormValue("checklist-id"))
	if checklistID == "" {
		http.Error(w, "Checklist ID is required", http.StatusBadRequest)
		return
	}

	err = database.ChecklistItemsTableResetChecklist(
		checklistID,
		auth.User.UserID,
		auth.CurrentBand.BandID,
	)
	if errors.Is(err, database.ErrChecklistNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		slog.Error(
			"unable to reset checklist",
			"request_id", requestlog.GetRequestID(r.Context()),
			"checklist_id", checklistID,
			"error", err,
		)
		http.Error(w, "Unable to reset checklist", http.StatusInternalServerError)
		return
	}

	http.Redirect(
		w,
		r,
		"/checklist?checklist-id="+url.QueryEscape(checklistID),
		http.StatusSeeOther,
	)
}

func (h Handler) setChecklistItemCompletion(w http.ResponseWriter, r *http.Request, complete bool) {
	auth, err := HelperGetAuthContext(r)
	if err != nil {
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	itemID := strings.TrimSpace(r.URL.Query().Get("item-id"))
	checklistID := strings.TrimSpace(r.URL.Query().Get("checklist-id"))
	if itemID == "" || checklistID == "" {
		http.Error(w, "Checklist and item are required", http.StatusBadRequest)
		return
	}

	err = database.ChecklistItemsTableSetCompletion(
		itemID,
		checklistID,
		auth.User.UserID,
		auth.CurrentBand.BandID,
		complete,
	)
	if errors.Is(err, database.ErrChecklistItemNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Unable to update checklist item", http.StatusInternalServerError)
		return
	}

	http.Redirect(
		w,
		r,
		"/checklist?checklist-id="+url.QueryEscape(checklistID),
		http.StatusSeeOther,
	)
}

func loadChecklistPageData(
	r *http.Request,
	user models.User,
	band models.Band,
	checklist models.Checklist,
) (models.ChecklistPageData, error) {
	songs, setlists, events, members, err := loadChecklistOptions(user, band)
	if err != nil {
		return models.ChecklistPageData{}, err
	}
	return models.ChecklistPageData{
		CSRFToken: csrf.Token(r),
		User:      user,
		Band:      band,
		Checklist: checklist,
		Songs:     songs,
		Setlists:  setlists,
		Events:    events,
		Members:   members,
	}, nil
}

func loadChecklistItemPageData(
	r *http.Request,
	user models.User,
	band models.Band,
	item models.ChecklistItem,
) (models.ChecklistItemPage, error) {
	songs, setlists, events, members, err := loadChecklistOptions(user, band)
	if err != nil {
		return models.ChecklistItemPage{}, err
	}
	return models.ChecklistItemPage{
		CSRFToken: csrf.Token(r),
		User:      user,
		Band:      band,
		Item:      item,
		Songs:     songs,
		Setlists:  setlists,
		Events:    events,
		Members:   members,
	}, nil
}

func loadChecklistOptions(
	user models.User,
	band models.Band,
) ([]models.Song, []models.Setlist, []models.Event, []models.User, error) {
	songs, err := database.SongsTableGetSongNameAndID(user.UserID, band.BandID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("load songs: %w", err)
	}
	setlists, err := database.SetlistsTableGetSetlistNamesAndIDs(band.BandID, user.UserID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("load setlists: %w", err)
	}
	events, err := database.EventsTableGetEventNameAndID(user.UserID, band.BandID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("load events: %w", err)
	}
	members, err := database.BandMembersGetMemberNameAndID(user.UserID, band.BandID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("load members: %w", err)
	}
	return songs, setlists, events, members, nil
}

func optionalChecklistValue(r *http.Request, field string) *string {
	value := strings.TrimSpace(r.FormValue(field))
	if value == "" || value == "none" {
		return nil
	}
	return &value
}
