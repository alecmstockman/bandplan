package handlers

import (
	"bandplan/src/database"
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/csrf"
)

func (h Handler) HandlerNotifications(w http.ResponseWriter, r *http.Request) {
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

	err = h.Tmpl.ExecuteTemplate(w, "notifications.html", data)
	if err != nil {
		slog.Error(
			"unable to load notifications.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to load notifications", http.StatusInternalServerError)
		return
	}
}

func (h Handler) HandlerNotification(w http.ResponseWriter, r *http.Request) {
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

	notificationID := strings.TrimSpace(r.URL.Query().Get("notification-id"))
	if notificationID == "" {
		http.Error(w, "Notification ID is required", http.StatusBadRequest)
		return
	}

	newNotif, err := database.NotifTableMarkNotifAsRead(r.Context(), user.UserID, band.BandID, notificationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Notification not found", http.StatusNotFound)
			return
		}

		slog.Error(
			"unable to mark notification as read",
			"request_id", requestlog.GetRequestID(r.Context()),
			"notification_id", notificationID,
			"error", err,
		)
		http.Error(w, "Unable to mark notification as read", http.StatusInternalServerError)
		return
	}

	itemID := url.QueryEscape(newNotif.ItemID)
	destination := "/notifications"
	switch newNotif.Type {
	case models.NotificationSongCreated, models.NotificationSongImported:
		destination = "/song?id=" + itemID
	case models.NotificationTaskAssigned:
		_, err := database.TodoListsTableGetListByID(newNotif.ItemID, user.UserID, band.BandID)
		if errors.Is(err, database.ErrTodoListNotFound) {
			destination = "/todos"
			break
		}
		if err != nil {
			slog.Error(
				"unable to load notification to do list",
				"request_id", requestlog.GetRequestID(r.Context()),
				"notification_id", notificationID,
				"list_id", newNotif.ItemID,
				"error", err,
			)
			http.Error(w, "Unable to load notification destination", http.StatusInternalServerError)
			return
		}
		destination = "/todo?todo-id=" + itemID
	case models.NotificationSetlistCreated:
		destination = "/setlist?id=" + itemID
	case models.NotificationEventCreated:
		destination = "/event?event-id=" + itemID
	case models.NotificationChatCreated:
		destination = "/chat?id=" + itemID
	case models.NotificationNewMember:
		destination = "/band?band-id=" + itemID
	}

	http.Redirect(w, r, destination, http.StatusSeeOther)
}
