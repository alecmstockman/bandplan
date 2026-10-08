package handlers

import (
	"bandplan/src/database"
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

func (h Handler) HandlerToDoListsPage(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerToDoListsPage")

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

	primaryUserList, err := database.TodoListsTableGetPrimaryUserListByUserID(user.UserID)
	if err != nil {
		slog.Error(
			"unable to load primary user todo list",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
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
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	lists, err := database.TodoListsTableGetListsByUserID(user.UserID)
	if err != nil {
		slog.Error(
			"unable to load user todo lists",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	data := models.ToDoListsPageData{
		CSRFToken:       csrf.Token(r),
		User:            user,
		Band:            band,
		PrimaryUserList: primaryUserList,
		PrimaryBandList: primaryBandList,
		Items:           lists,
	}

	err = h.Tmpl.ExecuteTemplate(w, "todo_lists.html", data)
	if err != nil {
		slog.Error(
			"unable to load todo_lists.html",
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
		CSRFToken: csrf.Token(r),
		User:      auth.User,
		Band:      auth.CurrentBand,
		ToDoList:  models.ToDoList{},
		Songs:     songs,
		Setlists:  setlists,
		Events:    events,
		Members:   members,
	}

	err = h.Tmpl.ExecuteTemplate(w, "todo_create.html", data)
	if err != nil {
		slog.Error(
			"unable to execute todo_create.html",
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

	dueDate, dueTime, dueTimezone, err := parseOptionalDueFields(
		strings.TrimSpace(r.FormValue("to-do-due-date")),
		strings.TrimSpace(r.FormValue("to-do-due-time")),
		strings.TrimSpace(r.FormValue("to-do-due-timezone")),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	description := strings.TrimSpace(r.FormValue("todo-list-description"))

	newToDoList := models.ToDoList{
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
		newToDoList.UserID = auth.User.UserID
	case "Band:" + auth.CurrentBand.BandID:
		newToDoList.BandID = auth.CurrentBand.BandID
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

	newToDoList.AssignedTo = optionalValue("to-do-assigned-to")
	newToDoList.SongID = optionalValue("to-do-song")
	newToDoList.SetlistID = optionalValue("to-do-setlist")
	newToDoList.EventID = optionalValue("to-do-event")

	err = database.TodoListsTableCreateTodoList(
		newToDoList,
		auth.User.UserID,
		auth.CurrentBand.BandID,
	)
	if errors.Is(err, database.ErrInvalidTodoListReferences) {
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

	http.Redirect(w, r, "/todos", http.StatusSeeOther)
}

func (h Handler) HandlerToDoListPage(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerToDoListPage")

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

	list, err := database.TodoListsTableGetListByID(todoListID)
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

	list.Items = listItems

	data := models.ToDoListPage{
		CSRFToken: csrf.Token(r),
		BackURL:   "/todos",
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

	dueDate, dueTime, dueTimezone, err := parseOptionalDueFields(
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
		AssignedTo:  optionalValue("todo-item-assigned-to"),
		DueDate:     dueDate,
		DueTime:     dueTime,
		DueTimezone: dueTimezone,
		CreatedBy:   auth.User.UserID,
		UpdatedBy:   auth.User.UserID,
	}

	err = database.TodoItemsTableCreateItem(
		item,
		auth.User.UserID,
		auth.CurrentBand.BandID,
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

	redirectURL := "/todo?todo-id=" + url.QueryEscape(listID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (h Handler) HandlerToDoItemDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("- HandlerToDoItemDelete")

	_, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	todoItemID := r.URL.Query().Get("item-id")
	todoListID := r.URL.Query().Get("list-id")

	err = database.TodoItemsTableDeleteItem(todoItemID, todoListID)
	if err != nil {
		slog.Error(
			"unable to delete setlist item",
			"request_id", requestlog.GetRequestID(r.Context()),
			"item_id", todoItemID,
			"error", err,
		)
		http.Error(w, "Unable to delete item", http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("/todo?todo-id=%v", todoListID)

	http.Redirect(w, r, url, http.StatusSeeOther)
	return
}
