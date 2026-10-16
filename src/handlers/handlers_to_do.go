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
	log.Print("\n- HandlerToDoListsPage")

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

	fmt.Println("TEST")

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
		ListCount:       len(lists) + 1,
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
	fmt.Println("todoListID: ", todoListID)

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
		ItemID:      itemID,
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
		UpdatedBy:   auth.User.UserID,
	}

	updatedItem, err := database.TodoItemsTableUpdateItem(
		item,
		auth.User.UserID,
		auth.CurrentBand.BandID,
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

	if r.Header.Get("HX-Request") != "true" {
		redirectURL := "/todo?todo-id=" + url.QueryEscape(listID)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
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

	return models.ToDoItemPage{
		CSRFToken: csrf.Token(r),
		BackURL:   "/todo?todo-id=" + url.QueryEscape(item.ToDoListID),
		User:      user,
		Band:      band,
		Item:      item,
		Songs:     songs,
		Setlists:  setlists,
		Events:    events,
		Members:   members,
	}, nil
}

func (h Handler) HandlerToDoItemDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("- HandlerToDoItemDelete")

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
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	return
}
