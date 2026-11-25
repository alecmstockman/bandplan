package handlers

import (
	"bandplan/src/database"
	"bandplan/src/helpers"
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/gorilla/csrf"
)

func (h Handler) HandlerToDoListsPage(w http.ResponseWriter, r *http.Request) {

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

	primaryUserList, err := database.TodoListsTableGetPrimaryUserListByUserID(user.UserID)
	if err != nil {
		slog.Error(
			"unable to load primary user todo list",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to load todo lists", http.StatusInternalServerError)
		return
	}

	primaryBandList, err := database.TodoListsTableGetPrimaryBandListByBandID(band.BandID)
	if err != nil {
		slog.Error(
			"unable to load primary band todo list",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to load todo lists", http.StatusInternalServerError)
		return
	}

	lists, err := database.TodoListsTableGetListsByUserID(user.UserID, band.BandID)
	if err != nil {
		slog.Error(
			"unable to load user todo lists",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to load todo lists", http.StatusInternalServerError)
		return
	}

	// for _, list := range lists {
	// 	fmt.Printf("\n\n%+v", list)
	// }

	data := models.ToDoListsPageData{
		CSRFToken:       csrf.Token(r),
		User:            user,
		Band:            band,
		PrimaryUserList: primaryUserList,
		PrimaryBandList: primaryBandList,
		Items:           lists,
		ListCount:       len(lists) + 1,
		Notifications:   notifications,
	}

	err = h.Tmpl.ExecuteTemplate(w, "todo-lists.html", data)
	if err != nil {
		slog.Error(
			"unable to load todo-lists.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
	}
}

func (h Handler) HandlerToDoListCreatePage(w http.ResponseWriter, r *http.Request) {
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

	data := models.ToDoListCreatePageData{
		CSRFToken:     csrf.Token(r),
		User:          user,
		Band:          band,
		ToDoList:      models.ToDoList{},
		Songs:         songs,
		Setlists:      setlists,
		Events:        events,
		Members:       members,
		Notifications: notifications,
	}

	err = h.Tmpl.ExecuteTemplate(w, "todo-create.html", data)
	if err != nil {
		slog.Error(
			"unable to execute todo-create.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to load to do list creation page", http.StatusInternalServerError)
		return
	}
}

func (h Handler) HandlerTodoListAdd(w http.ResponseWriter, r *http.Request) {

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

	name := strings.TrimSpace(r.FormValue("todo-list-name"))
	if name == "" {
		http.Error(w, "To do list name is required", http.StatusBadRequest)
		return
	}
	if utf8.RuneCountInString(name) > 120 {
		http.Error(w, "To do list name is too long", http.StatusBadRequest)
		return
	}

	dueDate, dueTime, dueTimezone, err := helpers.ParseOptionalDueFields(
		strings.TrimSpace(r.FormValue("to-do-due-date")),
		strings.TrimSpace(r.FormValue("to-do-due-time")),
		strings.TrimSpace(r.FormValue("to-do-due-timezone")),
	)
	if err != nil {
		slog.Error(
			"unable to parse optional fields",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Invald time entry", http.StatusBadRequest)
		return
	}

	description := strings.TrimSpace(r.FormValue("todo-list-description"))

	newToDoList := models.ToDoList{
		ToDoListID:  uuid.New().String(),
		Name:        name,
		Description: &description,
		DueDate:     dueDate,
		DueTime:     dueTime,
		DueTimezone: dueTimezone,
		CreatedBy:   auth.User.UserID,
		UpdatedBy:   auth.User.UserID,
	}

	switch strings.TrimSpace(r.FormValue("to-do-list-owner")) {
	case "Personal:" + auth.User.UserID:
		newToDoList.UserID = &auth.User.UserID
	case "Band:" + auth.CurrentBand.BandID:
		newToDoList.BandID = &auth.CurrentBand.BandID
	default:
		http.Error(w, "Invalid to do list owner", http.StatusBadRequest)
		return
	}

	optionalValue := func(field string) *string {
		value := strings.TrimSpace(r.FormValue(field))
		if value == "" || value == "none" {
			return nil
		}
		return &value
	}

	newToDoList.SongID = optionalValue("to-do-song")
	newToDoList.SetlistID = optionalValue("to-do-setlist")
	newToDoList.EventID = optionalValue("to-do-event")

	assigneeIDs := []string{}
	if newToDoList.BandID != nil {
		seenAssigneeIDs := make(map[string]bool)
		for _, rawAssigneeID := range r.PostForm["todo-assignee-user-id"] {
			assigneeID := strings.TrimSpace(rawAssigneeID)
			if assigneeID == "" || seenAssigneeIDs[assigneeID] {
				continue
			}

			seenAssigneeIDs[assigneeID] = true
			assigneeIDs = append(assigneeIDs, assigneeID)
		}
	}

	err = database.TodoListsTableCreateTodoList(
		r.Context(),
		newToDoList,
		auth.User.UserID,
		auth.CurrentBand.BandID,
		assigneeIDs,
	)
	if errors.Is(err, database.ErrInvalidTodoListReferences) {
		slog.Error(
			"unable to assign todo list",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Invalid to do list selection", http.StatusBadRequest)
		return
	}
	if err != nil {
		slog.Error(
			"unable to create todo list",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to create to do list", http.StatusInternalServerError)
		return
	}

	notification := models.Notification{
		BandID:    &band.BandID,
		ItemID:    newToDoList.ToDoListID,
		Type:      models.NotificationTodoListAssigned,
		Body:      models.NotifBodyListAssigned,
		CreatedBy: user.UserID,
	}

	fmt.Println("assigneeIDs: ", assigneeIDs)

	if len(assigneeIDs) > 0 {
		recipientIDs, err := h.Services.CreateMemberNotifications(r.Context(), notification, assigneeIDs)
		if err != nil {
			slog.Error(
				"unable to create notification",
				"request_id", requestlog.GetRequestID(r.Context()),
				"type", notification.Type,
				"error", err,
			)
		} else {
			for _, recipientID := range recipientIDs {
				if h.Hub != nil {
					h.Hub.PublishNotification(band.BandID, recipientID)
				}
			}

			go func(userIDs []string) {
				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()

				for _, userID := range userIDs {
					if err := h.Services.SendPushToUser(ctx, userID, models.PushTodoAssigned, newToDoList.ToDoListID); err != nil {
						slog.Warn(
							"unable to send to-do assignment push",
							"user_id", userID,
							"error", err,
						)
					}
				}
			}(recipientIDs)
		}

	} else if newToDoList.BandID != nil {
		fmt.Println("newToDoList.BandID: ", newToDoList.BandID)

		notification.BandID = newToDoList.BandID

		recipientIDs, err := h.Services.CreateBandNotification(r.Context(), notification)
		if err != nil {
			slog.Error(
				"unable to create notification",
				"request_id", requestlog.GetRequestID(r.Context()),
				"type", notification.Type,
				"error", err,
			)
		}
		for _, recipientID := range recipientIDs {
			if h.Hub != nil {
				h.Hub.PublishNotification(band.BandID, recipientID)
			}
		}

		go func(userIDs []string) {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			for _, userID := range userIDs {
				if err := h.Services.SendPushToUser(ctx, userID, models.PushTodoAssigned, newToDoList.ToDoListID); err != nil {
					slog.Warn(
						"unable to send to-do assignment push",
						"user_id", userID,
						"error", err,
					)
				}
			}
		}(recipientIDs)
	}
	http.Redirect(w, r, "/todos", http.StatusSeeOther)
}

func (h Handler) HandlerToDoListEditPage(w http.ResponseWriter, r *http.Request) {
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

	listID := strings.TrimSpace(r.URL.Query().Get("todo-id"))
	if listID == "" {
		http.Error(w, "To do list is required", http.StatusBadRequest)
		return
	}

	user := auth.User
	band := auth.CurrentBand
	list, err := database.TodoListsTableGetListByID(listID, user.UserID, band.BandID)
	if errors.Is(err, database.ErrTodoListNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		slog.Error(
			"unable to get todo list for editing",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"list_id", listID,
			"error", err,
		)
		http.Error(w, "Unable to load to do list", http.StatusInternalServerError)
		return
	}

	notifications, err := database.NotifTableGetNotificationsByUserID(r.Context(), user.UserID, band.BandID)
	if err != nil {
		slog.Error("unable to load notifications", "request_id", requestlog.GetRequestID(r.Context()), "error", err)
		http.Error(w, "Unable to load notifications", http.StatusInternalServerError)
		return
	}

	songs, err := database.SongsTableGetSongNameAndID(user.UserID, band.BandID)
	if err != nil {
		slog.Error("unable to get song names and IDs", "request_id", requestlog.GetRequestID(r.Context()), "error", err)
	}
	setlists, err := database.SetlistsTableGetSetlistNamesAndIDs(band.BandID, user.UserID)
	if err != nil {
		slog.Error("unable to get setlist names and IDs", "request_id", requestlog.GetRequestID(r.Context()), "error", err)
	}
	events, err := database.EventsTableGetEventNameAndID(user.UserID, band.BandID)
	if err != nil {
		slog.Error("unable to get event names and IDs", "request_id", requestlog.GetRequestID(r.Context()), "error", err)
	}
	members, err := database.BandMembersGetMemberNameAndID(user.UserID, band.BandID)
	if err != nil {
		slog.Error("unable to get band member names and IDs", "request_id", requestlog.GetRequestID(r.Context()), "error", err)
	}

	assigneeIDs, err := database.ToDoListAssigneesTableGetAssigneeIDs(r.Context(), listID, band.BandID)
	if err != nil {
		slog.Error(
			"unable to get to-do list assignees",
			"request_id", requestlog.GetRequestID(r.Context()),
			"list_id", listID,
			"error", err,
		)
		http.Error(w, "Unable to load to-do list assignees", http.StatusInternalServerError)
		return
	}

	data := models.ToDoListCreatePageData{
		CSRFToken:     csrf.Token(r),
		User:          user,
		Band:          band,
		ToDoList:      list,
		AssigneeIDs:   assigneeIDs,
		Songs:         songs,
		Setlists:      setlists,
		Events:        events,
		Members:       members,
		Notifications: notifications,
	}

	if err := h.Tmpl.ExecuteTemplate(w, "todo-edit.html", data); err != nil {
		slog.Error(
			"unable to execute todo-edit.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to load to do list edit page", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerToDoListInfoPage(w http.ResponseWriter, r *http.Request) {
	auth, err := HelperGetAuthContext(r)
	if err != nil {
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	listID := strings.TrimSpace(r.URL.Query().Get("todo-id"))
	if listID == "" {
		http.Error(w, "To do list is required", http.StatusBadRequest)
		return
	}

	list, err := database.TodoListsTableGetListByID(
		listID,
		auth.User.UserID,
		auth.CurrentBand.BandID,
	)
	if errors.Is(err, database.ErrTodoListNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		slog.Error("unable to load to-do list info", "error", err)
		http.Error(w, "Unable to load to do list", http.StatusInternalServerError)
		return
	}

	songs, err := database.SongsTableGetSongNameAndID(auth.User.UserID, auth.CurrentBand.BandID)
	if err != nil {
		http.Error(w, "Unable to load song information", http.StatusInternalServerError)
		return
	}
	setlists, err := database.SetlistsTableGetSetlistNamesAndIDs(auth.CurrentBand.BandID, auth.User.UserID)
	if err != nil {
		http.Error(w, "Unable to load setlist information", http.StatusInternalServerError)
		return
	}
	events, err := database.EventsTableGetEventNameAndID(auth.User.UserID, auth.CurrentBand.BandID)
	if err != nil {
		http.Error(w, "Unable to load event information", http.StatusInternalServerError)
		return
	}
	members, err := database.BandMembersGetMemberNameAndID(auth.User.UserID, auth.CurrentBand.BandID)
	if err != nil {
		http.Error(w, "Unable to load member information", http.StatusInternalServerError)
		return
	}
	assigneeIDs, err := database.ToDoListAssigneesTableGetAssigneeIDs(
		r.Context(),
		listID,
		auth.CurrentBand.BandID,
	)
	if err != nil {
		http.Error(w, "Unable to load to-do list assignees", http.StatusInternalServerError)
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

	assigneeLabel := "Personal"
	if list.UserID == nil {
		assigneeLabel = todoAssigneeLabel(assigneeIDs, members)
	}

	returnTo := internalBackURL(r.URL.Query().Get("return"), "/todos")
	detailQuery := url.Values{"todo-id": {listID}, "return": {returnTo}}
	infoQuery := url.Values{"todo-id": {listID}, "return": {returnTo}}

	data := models.ToDoListInfoPageData{
		User:             auth.User,
		Band:             auth.CurrentBand,
		ToDoList:         list,
		AssigneeLabel:    assigneeLabel,
		Songs:            songs,
		Setlists:         setlists,
		Events:           events,
		BackURL:          "/todo?" + detailQuery.Encode(),
		RelatedReturnURL: "/todo/info?" + infoQuery.Encode(),
		Notifications:    notifications,
	}

	if err := h.Tmpl.ExecuteTemplate(w, "todo-info.html", data); err != nil {
		slog.Error("unable to render to-do list info", "error", err)
		http.Error(w, "Unable to load to do list info", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerTodoListUpdate(w http.ResponseWriter, r *http.Request) {

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

	listID := strings.TrimSpace(r.FormValue("todo-list-id"))
	if listID == "" {
		http.Error(w, "To do list is required", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("todo-list-name"))
	if name == "" {
		http.Error(w, "To do list name is required", http.StatusBadRequest)
		return
	}

	if utf8.RuneCountInString(name) > 120 {
		http.Error(w, "To do list name is too long", http.StatusBadRequest)
		return
	}

	dueDate, dueTime, dueTimezone, err := helpers.ParseOptionalDueFields(
		strings.TrimSpace(r.FormValue("to-do-due-date")),
		strings.TrimSpace(r.FormValue("to-do-due-time")),
		strings.TrimSpace(r.FormValue("to-do-due-timezone")),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	description := strings.TrimSpace(r.FormValue("todo-list-description"))
	updatedList := models.ToDoList{
		ToDoListID:  listID,
		Name:        name,
		Description: &description,
		DueDate:     dueDate,
		DueTime:     dueTime,
		DueTimezone: dueTimezone,
		UpdatedBy:   auth.User.UserID,
	}

	switch strings.TrimSpace(r.FormValue("to-do-list-owner")) {
	case "Personal:" + auth.User.UserID:
		updatedList.UserID = &auth.User.UserID
	case "Band:" + auth.CurrentBand.BandID:
		updatedList.BandID = &auth.CurrentBand.BandID
	default:
		http.Error(w, "Invalid to do list owner", http.StatusBadRequest)
		return
	}

	optionalValue := func(field string) *string {
		value := strings.TrimSpace(r.FormValue(field))
		if value == "" || value == "none" {
			return nil
		}
		return &value
	}

	updatedList.SongID = optionalValue("to-do-song")
	updatedList.SetlistID = optionalValue("to-do-setlist")
	updatedList.EventID = optionalValue("to-do-event")

	existingList, err := database.TodoListsTableGetListByID(
		listID,
		auth.User.UserID,
		auth.CurrentBand.BandID,
	)
	if errors.Is(err, database.ErrTodoListNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		slog.Error(
			"unable to load todo list for update",
			"request_id", requestlog.GetRequestID(r.Context()),
			"list_id", listID,
			"error", err,
		)
		http.Error(w, "Unable to update to do list", http.StatusInternalServerError)
		return
	}

	assigneeIDs := []string{}
	if updatedList.BandID != nil && !existingList.IsPrimaryBand {
		seenAssigneeIDs := make(map[string]bool)
		for _, rawAssigneeID := range r.PostForm["todo-assignee-user-id"] {
			assigneeID := strings.TrimSpace(rawAssigneeID)
			if assigneeID == "" || seenAssigneeIDs[assigneeID] {
				continue
			}

			seenAssigneeIDs[assigneeID] = true
			assigneeIDs = append(assigneeIDs, assigneeID)
		}
	}

	if existingList.IsPrimaryUser {
		err = database.TodoListsTableUpdatePrimaryUserTodoList(
			r.Context(),
			updatedList,
			auth.User.UserID,
			auth.CurrentBand.BandID,
		)
	} else {
		err = database.TodoListsTableUpdateTodoList(
			r.Context(),
			updatedList,
			auth.User.UserID,
			auth.CurrentBand.BandID,
			assigneeIDs,
		)
	}
	if errors.Is(err, database.ErrTodoListNotFound) {
		slog.Error(
			"todo list not found",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.NotFound(w, r)
		return
	}
	if errors.Is(err, database.ErrPrimaryTodoListOwnerChange) {
		slog.Error(
			"unable to chage primary todo list owner",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Primary to do list owner cannot be changed", http.StatusBadRequest)
		return
	}
	if errors.Is(err, database.ErrInvalidTodoListReferences) {
		slog.Error(
			"invalid todo list reference",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Invalid to do list selection", http.StatusBadRequest)
		return
	}
	if err != nil {
		slog.Error(
			"unable to update todo list",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"list_id", listID,
			"error", err,
		)
		http.Error(w, "Unable to update to do list", http.StatusInternalServerError)
		return
	}

	notification := models.Notification{
		BandID:    &band.BandID,
		ItemID:    updatedList.ToDoListID,
		Type:      models.NotificationTodoListUpdated,
		Body:      models.NotifBodyTodoListUpdated,
		CreatedBy: user.UserID,
	}

	if len(assigneeIDs) > 0 {
		recipientIDs, err := h.Services.CreateMemberNotifications(r.Context(), notification, assigneeIDs)
		if err != nil {
			slog.Error(
				"unable to create update todo list notification",
				"request_id", requestlog.GetRequestID(r.Context()),
				"type", notification.Type,
				"error", err,
			)
		} else {
			for _, recipientID := range recipientIDs {
				if h.Hub != nil {
					h.Hub.PublishNotification(band.BandID, recipientID)
				}
			}

			go func(userIDs []string) {
				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()

				for _, userID := range userIDs {
					if err := h.Services.SendPushToUser(ctx, userID, models.PushTodoAssigned, updatedList.ToDoListID); err != nil {
						slog.Warn(
							"unable to send to-do assignment push",
							"user_id", userID,
							"error", err,
						)
					}
				}
			}(recipientIDs)
		}

	} else if updatedList.BandID != nil {
		fmt.Println("newToDoList.BandID: ", updatedList.BandID)

		notification.BandID = updatedList.BandID

		recipientIDs, err := h.Services.CreateBandNotification(r.Context(), notification)
		if err != nil {
			slog.Error(
				"unable to create update todo list notification",
				"request_id", requestlog.GetRequestID(r.Context()),
				"type", notification.Type,
				"error", err,
			)
		}
		for _, recipientID := range recipientIDs {
			if h.Hub != nil {
				h.Hub.PublishNotification(band.BandID, recipientID)
			}
		}

		go func(userIDs []string) {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			for _, userID := range userIDs {
				if err := h.Services.SendPushToUser(ctx, userID, models.PushTodoListUpdated, updatedList.ToDoListID); err != nil {
					slog.Warn(
						"unable to send to-do assignment push",
						"user_id", userID,
						"error", err,
					)
				}
			}
		}(recipientIDs)
	}

	redirectURL := "/todo?todo-id=" + url.QueryEscape(listID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (h Handler) HandlerToDoListDelete(w http.ResponseWriter, r *http.Request) {
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

	r.Body = http.MaxBytesReader(w, r.Body, 4*1024)
	if err := r.ParseForm(); err != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "Invalid form", http.StatusBadRequest)
		}
		return
	}

	listID := strings.TrimSpace(r.FormValue("todo-list-id"))
	if listID == "" {
		http.Error(w, "To do list is required", http.StatusBadRequest)
		return
	}

	err = database.TodoListsTableDeleteList(
		listID,
		auth.User.UserID,
		auth.CurrentBand.BandID,
	)
	if errors.Is(err, database.ErrPrimaryTodoList) {
		http.Error(w, "Primary to do lists cannot be deleted", http.StatusForbidden)
		return
	}
	if errors.Is(err, database.ErrTodoListNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		slog.Error(
			"unable to delete todo list",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"list_id", listID,
			"error", err,
		)
		http.Error(w, "Unable to delete to do list", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/todos", http.StatusSeeOther)
}

func (h Handler) HandlerToDoListPageAll(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerToDoListPageAll")

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

	todoListID := r.URL.Query().Get("todo-id")

	list, err := database.TodoListsTableGetListByID(todoListID, user.UserID, band.BandID)
	if err != nil {
		slog.Error(
			"unable to get todo list",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"list_id", list.ToDoListID,
			"error", err,
		)
		http.Error(w, "Unable to load todo list", http.StatusInternalServerError)
		return
	}

	listItems, err := database.TodoItemsTableGetItemsByListID(list.ToDoListID)
	if err != nil {
		slog.Error(
			"unable to get todo list items",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"list_id", list.ToDoListID,
			"error", err,
		)
		http.Error(w, "Unable to load todo list", http.StatusInternalServerError)
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

	if err := loadToDoItemAssignees(r.Context(), list.ToDoListID, band.BandID, listItems, members); err != nil {
		slog.Error("unable to load to-do item assignees", "error", err)
		http.Error(w, "Unable to load to-do item assignees", http.StatusInternalServerError)
		return
	}
	list.Items = listItems

	data := models.ToDoListPage{
		CSRFToken: csrf.Token(r),
		BackURL:   internalBackURL(r.URL.Query().Get("return"), "/todos"),
		View:      models.TodoViewAll,
		User:      user,
		Band:      band,
		List:      list,
		Songs:     songs,
		Setlists:  setlists,
		Events:    events,
		Members:   members,
	}

	err = h.Tmpl.ExecuteTemplate(w, "todo.html", data)
	if err != nil {
		slog.Error(
			"unable to load todo.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
	}
}

func (h Handler) HandlerToDoListPageActive(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerToDoListPageActive")

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

	query := r.URL.Query()
	todoListID := query.Get("todo-id")

	var list models.ToDoList
	if todoListID == "" && query.Get("add-task") == "1" {
		list, err = database.TodoListsTableGetPrimaryBandListByBandID(band.BandID)
		if errors.Is(err, sql.ErrNoRows) {
			list, err = database.TodoListsTableGetPrimaryUserListByUserID(user.UserID)
		}
		if errors.Is(err, sql.ErrNoRows) {
			http.Redirect(w, r, "/todo/create", http.StatusSeeOther)
			return
		}
	} else {
		list, err = database.TodoListsTableGetListByID(todoListID, user.UserID, band.BandID)
	}
	if err != nil {
		slog.Error(
			"unable to get todo list",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"list_id", list.ToDoListID,
			"error", err,
		)
		http.Error(w, "Unable to load todo list", http.StatusInternalServerError)
		return
	}

	listItems, err := database.TodoItemsTableGetActiveItemsByListID(list.ToDoListID)
	if err != nil {
		slog.Error(
			"unable to get todo list items",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"list_id", list.ToDoListID,
			"error", err,
		)
		http.Error(w, "Unable to load todo list", http.StatusInternalServerError)
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

	if err := loadToDoItemAssignees(r.Context(), list.ToDoListID, band.BandID, listItems, members); err != nil {
		slog.Error("unable to load to-do item assignees", "error", err)
		http.Error(w, "Unable to load to-do item assignees", http.StatusInternalServerError)
		return
	}
	list.Items = listItems

	data := models.ToDoListPage{
		CSRFToken: csrf.Token(r),
		BackURL:   internalBackURL(r.URL.Query().Get("return"), "/todos"),
		View:      models.TodoViewActive,
		User:      user,
		Band:      band,
		List:      list,
		Songs:     songs,
		Setlists:  setlists,
		Events:    events,
		Members:   members,
	}

	err = h.Tmpl.ExecuteTemplate(w, "todo.html", data)
	if err != nil {
		slog.Error(
			"unable to load todo.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
	}
}

func (h Handler) HandlerToDoListPageCompleted(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerToDoListPageCompleted")

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

	todoListID := r.URL.Query().Get("todo-id")

	list, err := database.TodoListsTableGetListByID(todoListID, user.UserID, band.BandID)
	if err != nil {
		slog.Error(
			"unable to get todo list",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"list_id", list.ToDoListID,
			"error", err,
		)
		http.Error(w, "Unable to load todo list", http.StatusInternalServerError)
		return
	}

	listItems, err := database.TodoItemsTableGetCompletedItemsByListID(list.ToDoListID)
	if err != nil {
		slog.Error(
			"unable to get todo list items",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"list_id", list.ToDoListID,
			"error", err,
		)
		http.Error(w, "Unable to load todo list", http.StatusInternalServerError)
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

	if err := loadToDoItemAssignees(r.Context(), list.ToDoListID, band.BandID, listItems, members); err != nil {
		slog.Error("unable to load to-do item assignees", "error", err)
		http.Error(w, "Unable to load to-do item assignees", http.StatusInternalServerError)
		return
	}
	list.Items = listItems

	data := models.ToDoListPage{
		CSRFToken: csrf.Token(r),
		BackURL:   internalBackURL(r.URL.Query().Get("return"), "/todos"),
		View:      models.TodoViewCompleted,
		User:      user,
		Band:      band,
		List:      list,
		Songs:     songs,
		Setlists:  setlists,
		Events:    events,
		Members:   members,
	}

	err = h.Tmpl.ExecuteTemplate(w, "todo.html", data)
	if err != nil {
		slog.Error(
			"unable to load todo.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
	}
}

func (h Handler) HandlerToDoItemAdd(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerToDoItemAdd")

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

	listID := strings.TrimSpace(r.FormValue("todo-list-id"))
	if listID == "" {
		http.Error(w, "To do list is required", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("todo-item-name"))
	if name == "" {
		http.Error(w, "To do item name is required", http.StatusBadRequest)
		return
	}
	if utf8.RuneCountInString(name) > 120 {
		http.Error(w, "To do item name is too long", http.StatusBadRequest)
		return
	}

	dueDate, dueTime, dueTimezone, err := helpers.ParseOptionalDueFields(
		strings.TrimSpace(r.FormValue("todo-item-due-date")),
		strings.TrimSpace(r.FormValue("todo-item-due-time")),
		strings.TrimSpace(r.FormValue("todo-item-due-timezone")),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	optionalValue := func(field string) *string {
		value := strings.TrimSpace(r.FormValue(field))
		if value == "" || value == "none" {
			return nil
		}
		return &value
	}

	item := models.ToDoItem{
		ToDoListID:  listID,
		Name:        name,
		Body:        strings.TrimSpace(r.FormValue("todo-item-description")),
		SongID:      optionalValue("todo-item-song"),
		SetlistID:   optionalValue("todo-item-setlist"),
		EventID:     optionalValue("todo-item-event"),
		DueDate:     dueDate,
		DueTime:     dueTime,
		DueTimezone: dueTimezone,
		CreatedBy:   user.UserID,
		UpdatedBy:   user.UserID,
	}

	assigneeIDs := postedAssigneeIDs(r, "todo-item-assignee-user-id")
	err = database.TodoItemsTableCreateItem(
		r.Context(),
		item,
		user.UserID,
		band.BandID,
		assigneeIDs,
	)
	if errors.Is(err, database.ErrInvalidTodoItemReferences) {
		http.Error(w, "Invalid to do item selection", http.StatusBadRequest)
		return
	}
	if err != nil {
		slog.Error(
			"unable to create todo item",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"list_id", listID,
			"error", err,
		)
		http.Error(w, "Unable to create to do item", http.StatusInternalServerError)
		return
	}

	notif := models.Notification{
		BandID:    &band.BandID,
		UserID:    user.UserID,
		Type:      models.NotificationTaskAssigned,
		ItemID:    listID,
		Body:      models.NotifBodyTaskAssigned,
		CreatedBy: user.UserID,
	}

	fmt.Println("assigneeIDs: ", assigneeIDs)

	if len(assigneeIDs) > 0 {
		fmt.Println("len greate than 1")
		recipientIDs, err := h.Services.CreateMemberNotifications(r.Context(), notif, assigneeIDs)
		if err != nil {
			slog.Error(
				"unable to create notification",
				"request_id", requestlog.GetRequestID(r.Context()),
				"type", notif.Type,
				"error", err,
			)
		} else {
			fmt.Println("no error")
			for _, recipientID := range recipientIDs {
				if h.Hub != nil {
					h.Hub.PublishNotification(band.BandID, recipientID)
				}
			}

			go func(userIDs []string) {
				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()

				for _, userID := range userIDs {
					if err := h.Services.SendPushToUser(ctx, userID, models.PushTaskAssigned, listID); err != nil {
						slog.Warn(
							"unable to send task assignment push",
							"user_id", userID,
							"error", err,
						)
					}
				}
			}(recipientIDs)
		}
	} else {
		recipientIDs, err := h.Services.CreateBandNotification(r.Context(), notif)
		if err != nil {
			slog.Error(
				"unable to create notification",
				"request_id", requestlog.GetRequestID(r.Context()),
				"type", notif.Type,
				"error", err,
			)
		} else if h.Hub != nil {
			for _, recipientID := range recipientIDs {
				h.Hub.PublishNotification(band.BandID, recipientID)
			}
		}

		go func(userIDs []string) {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			for _, userID := range userIDs {
				if err := h.Services.SendPushToUser(ctx, userID, models.PushTaskAssigned, listID); err != nil {
					slog.Warn(
						"unable to send task assignment push",
						"user_id", userID,
						"error", err,
					)
				}
			}
		}(recipientIDs)
	}

	redirectURL := "/todo?todo-id=" + url.QueryEscape(listID) + "&toast=todo-task-created"
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (h Handler) HandlerToDoItemPage(w http.ResponseWriter, r *http.Request) {
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

	itemID := strings.TrimSpace(r.URL.Query().Get("item-id"))
	listID := strings.TrimSpace(r.URL.Query().Get("list-id"))
	if itemID == "" || listID == "" {
		http.Error(w, "To do item and list are required", http.StatusBadRequest)
		return
	}

	item, err := database.TodoItemsTableGetItemByID(
		itemID,
		listID,
		auth.User.UserID,
		auth.CurrentBand.BandID,
	)
	if errors.Is(err, database.ErrTodoItemNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		slog.Error(
			"unable to get todo item page",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"item_id", itemID,
			"list_id", listID,
			"error", err,
		)
		http.Error(w, "Unable to load todo item", http.StatusInternalServerError)
		return
	}

	data, err := loadToDoItemPageData(r, auth.User, auth.CurrentBand, item)
	if err != nil {
		slog.Error(
			"unable to load todo item page data",
			"request_id", requestlog.GetRequestID(r.Context()),
			"item_id", itemID,
			"list_id", listID,
			"error", err,
		)
		http.Error(w, "Unable to load todo item", http.StatusInternalServerError)
		return
	}

	data.Notifications, err = database.NotifTableGetNotificationsByUserID(
		r.Context(),
		auth.User.UserID,
		auth.CurrentBand.BandID,
	)
	if err != nil {
		slog.Error(
			"unable to load notifications",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load notifications", http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("return") == "calendar" {
		data.BackURL = "/calendar"
		data.ReturnTo = "calendar"
	}
	itemPageQuery := url.Values{}
	itemPageQuery.Set("item-id", item.ItemID)
	itemPageQuery.Set("list-id", item.ToDoListID)
	if data.ReturnTo == "calendar" {
		itemPageQuery.Set("return", "calendar")
	}
	data.RelatedReturnURL = "/todo/item/page?" + itemPageQuery.Encode()

	if err := h.Tmpl.ExecuteTemplate(w, "todo-item.html", data); err != nil {
		slog.Error(
			"unable to execute todo-item.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
	}
}

func (h Handler) HandlerToDoItemEdit(w http.ResponseWriter, r *http.Request) {
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

	listID := strings.TrimSpace(r.FormValue("todo-id"))
	itemID := strings.TrimSpace(r.FormValue("todo-item-id"))
	if listID == "" || itemID == "" {
		http.Error(w, "To do item and list are required", http.StatusBadRequest)
		return
	}

	item, err := database.TodoItemsTableGetItemByID(
		itemID,
		listID,
		user.UserID,
		band.BandID,
	)
	if errors.Is(err, database.ErrTodoItemNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		slog.Error(
			"unable to get todo item",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"item_id", itemID,
			"list_id", listID,
			"error", err,
		)
		http.Error(w, "Unable to load todo item", http.StatusInternalServerError)
		return
	}

	data, err := loadToDoItemPageData(r, user, band, item)
	if err != nil {
		slog.Error(
			"unable to load todo item edit data",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"item_id", itemID,
			"list_id", listID,
			"error", err,
		)
		http.Error(w, "Unable to load todo item", http.StatusInternalServerError)
		return
	}
	switch r.FormValue("return") {
	case "task-page", "task-page-calendar":
		data.ReturnTo = r.FormValue("return")
	}

	err = h.Tmpl.ExecuteTemplate(w, "todo_item_edit", data)
	if err != nil {
		slog.Error(
			"unable to execute todo_item_edit template",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
	}
}

func (h Handler) HandlerToDoItemUpdate(w http.ResponseWriter, r *http.Request) {

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

	itemID := strings.TrimSpace(r.FormValue("todo-item-id"))
	if itemID == "" {
		http.Error(w, "To do item is required", http.StatusBadRequest)
		return
	}

	listID := strings.TrimSpace(r.FormValue("todo-list-id"))
	if listID == "" {
		http.Error(w, "To do list is required", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("todo-item-name"))
	if name == "" {
		http.Error(w, "To do item name is required", http.StatusBadRequest)
		return
	}
	if utf8.RuneCountInString(name) > 120 {
		http.Error(w, "To do item name is too long", http.StatusBadRequest)
		return
	}

	var completion *bool
	switch r.PostForm.Get("todo-item-status") {
	case "":
	case "true":
		value := true
		completion = &value
	case "false":
		value := false
		completion = &value
	default:
		http.Error(w, "Invalid task status", http.StatusBadRequest)
		return
	}

	dueDate, dueTime, dueTimezone, err := helpers.ParseOptionalDueFields(
		strings.TrimSpace(r.FormValue("todo-item-due-date")),
		strings.TrimSpace(r.FormValue("todo-item-due-time")),
		strings.TrimSpace(r.FormValue("todo-item-due-timezone")),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	previousTask, err := database.TodoItemsTableGetItemByID(itemID, listID, auth.User.UserID, auth.CurrentBand.BandID)
	if err != nil {
		slog.Error(
			"unable to load previous task information",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}
	previousAssigneeIDs, err := database.ToDoItemAssigneesTableGetAssigneeIDs(
		r.Context(),
		itemID,
		band.BandID,
	)
	if err != nil {
		slog.Error(
			"unable to load previous task assignees",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load task assignees", http.StatusInternalServerError)
		return
	}
	previousTask.AssigneeIDs = previousAssigneeIDs

	optionalValue := func(field string) *string {
		value := strings.TrimSpace(r.FormValue(field))
		if value == "" || value == "none" {
			return nil
		}
		return &value
	}

	item := models.ToDoItem{
		ItemID:      itemID,
		ToDoListID:  listID,
		Name:        name,
		Body:        strings.TrimSpace(r.FormValue("todo-item-description")),
		SongID:      optionalValue("todo-item-song"),
		SetlistID:   optionalValue("todo-item-setlist"),
		EventID:     optionalValue("todo-item-event"),
		DueDate:     dueDate,
		DueTime:     dueTime,
		DueTimezone: dueTimezone,
		UpdatedBy:   auth.User.UserID,
	}

	assigneeIDs := postedAssigneeIDs(r, "todo-item-assignee-user-id")
	updatedItem, err := database.TodoItemsTableUpdateItem(
		r.Context(),
		item,
		auth.User.UserID,
		auth.CurrentBand.BandID,
		completion,
		assigneeIDs,
	)
	if errors.Is(err, database.ErrTodoItemNotFound) {
		http.NotFound(w, r)
		return
	}
	if errors.Is(err, database.ErrInvalidTodoItemReferences) {
		http.Error(w, "Invalid to do item selection", http.StatusBadRequest)
		return
	}
	if err != nil {
		slog.Error(
			"unable to update todo item",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"item_id", itemID,
			"list_id", listID,
			"error", err,
		)
		http.Error(w, "Unable to update to do item", http.StatusInternalServerError)
		return
	}

	var notifBody models.NotificationBody

	assignmentsChanged := !assigneeMapsEqual(previousTask.AssigneeIDs, updatedItem.AssigneeIDs)
	if !assignmentsChanged {
		notifBody = models.NotifBodyTaskUpdated
	} else {
		notifBody = models.NotifBodyTaskAssigned
	}

	notification := models.Notification{
		BandID:    &band.BandID,
		Type:      models.NotificationTaskAssigned,
		ItemID:    listID,
		Body:      notifBody,
		CreatedBy: user.UserID,
	}

	if len(assigneeIDs) > 0 {
		recipientIDs, err := h.Services.CreateMemberNotifications(r.Context(), notification, assigneeIDs)
		if err != nil {
			slog.Error(
				"unable to create notification",
				"request_id", requestlog.GetRequestID(r.Context()),
				"type", notification.Type,
				"error", err,
			)
		} else {
			for _, recipientID := range recipientIDs {
				if h.Hub != nil {
					h.Hub.PublishNotification(auth.CurrentBand.BandID, recipientID)
				}
			}

			if assignmentsChanged {
				go func(userIDs []string) {
					ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
					defer cancel()

					for _, userID := range userIDs {
						if err := h.Services.SendPushToUser(ctx, userID, models.PushTaskAssigned, listID); err != nil {
							slog.Warn(
								"unable to send task assignment push",
								"user_id", userID,
								"error", err,
							)
						}
					}
				}(recipientIDs)
			}
		}
	}

	switch r.PostForm.Get("return") {
	case "task-page", "task-page-calendar":
		query := url.Values{}
		query.Set("item-id", itemID)
		query.Set("list-id", listID)
		if r.PostForm.Get("return") == "task-page-calendar" {
			query.Set("return", "calendar")
		}
		redirectURL := "/todo/item/page?" + query.Encode()
		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Redirect", redirectURL)
			w.WriteHeader(http.StatusOK)
		} else {
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		}
		return
	}

	if completion != nil || r.Header.Get("HX-Request") != "true" {
		listPath := "/todo"
		if completion != nil && updatedItem.IsComplete {
			listPath = "/todo/completed"
		}
		redirectURL := listPath + "?todo-id=" + url.QueryEscape(listID)
		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Redirect", redirectURL)
			w.WriteHeader(http.StatusOK)
		} else {
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		}
		return
	}

	data, err := loadToDoItemPageData(r, auth.User, auth.CurrentBand, updatedItem)
	if err != nil {
		slog.Error(
			"unable to load updated todo item data",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"item_id", itemID,
			"list_id", listID,
			"error", err,
		)
		http.Error(w, "Unable to load updated to do item", http.StatusInternalServerError)
		return
	}

	if err := h.Tmpl.ExecuteTemplate(w, "todo_item_update", data); err != nil {
		slog.Error(
			"unable to execute todo_item_update template",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"item_id", itemID,
			"list_id", listID,
			"error", err,
		)
	}
}

func loadToDoItemPageData(r *http.Request, user models.User, band models.Band, item models.ToDoItem) (models.ToDoItemPage, error) {
	songs, err := database.SongsTableGetSongNameAndID(user.UserID, band.BandID)
	if err != nil {
		return models.ToDoItemPage{}, fmt.Errorf("load song options: %w", err)
	}

	setlists, err := database.SetlistsTableGetSetlistNamesAndIDs(band.BandID, user.UserID)
	if err != nil {
		return models.ToDoItemPage{}, fmt.Errorf("load setlist options: %w", err)
	}

	events, err := database.EventsTableGetEventNameAndID(user.UserID, band.BandID)
	if err != nil {
		return models.ToDoItemPage{}, fmt.Errorf("load event options: %w", err)
	}

	members, err := database.BandMembersGetMemberNameAndID(user.UserID, band.BandID)
	if err != nil {
		return models.ToDoItemPage{}, fmt.Errorf("load band member options: %w", err)
	}
	assigneeIDs, err := database.ToDoItemAssigneesTableGetAssigneeIDs(
		r.Context(),
		item.ItemID,
		band.BandID,
	)
	if err != nil {
		return models.ToDoItemPage{}, fmt.Errorf("load task assignees: %w", err)
	}
	item.AssigneeIDs = assigneeIDs
	item.AssigneeLabel = todoAssigneeLabel(assigneeIDs, members)

	return models.ToDoItemPage{
		CSRFToken:        csrf.Token(r),
		BackURL:          "/todo?todo-id=" + url.QueryEscape(item.ToDoListID),
		RelatedReturnURL: "/todo?todo-id=" + url.QueryEscape(item.ToDoListID),
		User:             user,
		Band:             band,
		Item:             item,
		AssigneeIDs:      assigneeIDs,
		Songs:            songs,
		Setlists:         setlists,
		Events:           events,
		Members:          members,
	}, nil
}

func postedAssigneeIDs(r *http.Request, field string) []string {
	assigneeIDs := []string{}
	seenAssigneeIDs := make(map[string]bool)
	for _, rawAssigneeID := range r.PostForm[field] {
		assigneeID := strings.TrimSpace(rawAssigneeID)
		if assigneeID == "" || seenAssigneeIDs[assigneeID] {
			continue
		}
		seenAssigneeIDs[assigneeID] = true
		assigneeIDs = append(assigneeIDs, assigneeID)
	}
	return assigneeIDs
}

func assigneeMapsEqual(left, right map[string]bool) bool {
	if len(left) != len(right) {
		return false
	}
	for assigneeID := range left {
		if !right[assigneeID] {
			return false
		}
	}
	return true
}

func loadToDoItemAssignees(
	ctx context.Context,
	listID, bandID string,
	items []models.ToDoItem,
	members []models.User,
) error {
	assigneesByItem, err := database.ToDoItemAssigneesTableGetByListID(ctx, listID, bandID)
	if err != nil {
		return err
	}
	for i := range items {
		items[i].AssigneeIDs = assigneesByItem[items[i].ItemID]
		items[i].AssigneeLabel = todoAssigneeLabel(items[i].AssigneeIDs, members)
	}
	return nil
}

func todoAssigneeLabel(assigneeIDs map[string]bool, members []models.User) string {
	if len(assigneeIDs) == 0 {
		return "Band"
	}
	if len(assigneeIDs) > 1 {
		return fmt.Sprintf("%d members", len(assigneeIDs))
	}
	for _, member := range members {
		if assigneeIDs[member.UserID] {
			return member.DisplayName
		}
	}
	return "Unavailable"
}

func (h Handler) HandlerToDoItemDelete(w http.ResponseWriter, r *http.Request) {
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

	r.Body = http.MaxBytesReader(w, r.Body, 4*1024)
	if err := r.ParseForm(); err != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "Invalid form", http.StatusBadRequest)
		}
		return
	}

	todoItemID := strings.TrimSpace(r.FormValue("item-id"))
	todoListID := strings.TrimSpace(r.FormValue("list-id"))
	if todoItemID == "" || todoListID == "" {
		http.Error(w, "Item ID and list ID are required", http.StatusBadRequest)
		return
	}

	err = database.TodoItemsTableDeleteItem(todoItemID, todoListID, user.UserID, band.BandID)
	if err != nil {
		if errors.Is(err, database.ErrTodoItemNotFound) {
			http.NotFound(w, r)
			return
		}
		slog.Error(
			"unable to delete to-do item",
			"request_id", requestlog.GetRequestID(r.Context()),
			"item_id", todoItemID,
			"error", err,
		)
		http.Error(w, "Unable to delete item", http.StatusInternalServerError)
		return
	}

	redirectURL := "/todo?todo-id=" + url.QueryEscape(todoListID)
	if r.FormValue("return") == "calendar" {
		redirectURL = "/calendar"
	}
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (h Handler) HandlerToDoItemComplete(w http.ResponseWriter, r *http.Request) {
	h.setToDoItemCompletion(w, r, true)
}

func (h Handler) HandlerToDoItemReopen(w http.ResponseWriter, r *http.Request) {
	h.setToDoItemCompletion(w, r, false)
}

func (h Handler) setToDoItemCompletion(w http.ResponseWriter, r *http.Request, complete bool) {
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

	todoItemID := strings.TrimSpace(r.URL.Query().Get("item-id"))
	todoListID := strings.TrimSpace(r.URL.Query().Get("list-id"))
	if todoItemID == "" || todoListID == "" {
		http.Error(w, "Item ID and list ID are required", http.StatusBadRequest)
		return
	}

	err = database.TodoItemsTableSetItemCompletion(todoItemID, todoListID, user.UserID, band.BandID, complete)
	if err != nil {
		if errors.Is(err, database.ErrTodoItemNotFound) {
			http.NotFound(w, r)
			return
		}
		slog.Error(
			"unable to update to-do item completion",
			"request_id", requestlog.GetRequestID(r.Context()),
			"item_id", todoItemID,
			"error", err,
		)
		http.Error(w, "Unable to update task status", http.StatusInternalServerError)
		return
	}

	redirectURL := "/todo?todo-id=" + url.QueryEscape(todoListID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
