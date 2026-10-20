package handlers

import (
	"bandplan/src/database"
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"log/slog"
	"net/http"
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

	read, err := database.NotifTableMarkNotifAsRead(r.Context(), user.UserID, band.BandID, notificationID)
	if err != nil {
		slog.Error(
			"unable to mark notification as read",
			"request_id", requestlog.GetRequestID(r.Context()),
			"notification_id", notificationID,
			"error", err,
		)
		http.Error(w, "Unable to mark notification as read", http.StatusInternalServerError)
		return
	}
	if !read {
		http.Error(w, "Notification not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/notifications", http.StatusSeeOther)
}
