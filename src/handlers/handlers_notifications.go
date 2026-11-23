package handlers

import (
	"bandplan/src/database"
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/csrf"
)

func (h Handler) HandlerNotificationsPage(w http.ResponseWriter, r *http.Request) {
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
	backURL := strings.TrimSpace(r.URL.Query().Get("back"))
	if backURL == "" || !strings.HasPrefix(backURL, "/") || strings.HasPrefix(backURL, "//") {
		backURL = "/"
	}

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

	subscriptions, err := database.PushSubscriptionsTableListPushSubscriptions(r.Context(), user.UserID)
	if err != nil {
		slog.Error(
			"unable to load subscriptions on notification page",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
	}

	for _, sub := range subscriptions {
		fmt.Printf("\nSub: \n%+v", sub)
	}

	data := struct {
		models.MenuPageData
		VAPIDPublicKey string
		Subsciptions   []models.PushSubscription
	}{
		MenuPageData: models.MenuPageData{
			CSRFToken:     csrf.Token(r),
			BackURL:       backURL,
			User:          user,
			Band:          band,
			Notifications: notifications,
		},
		VAPIDPublicKey: os.Getenv("VAPID_PUBLIC_KEY"),
		Subsciptions:   subscriptions,
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
	fmt.Println("- HandlerNotification")

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

	fmt.Println("\nnotificationID: ", notificationID)

	newNotif, err := database.NotifTableMarkNotifAsRead(r.Context(), user.UserID, band.BandID, notificationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error(
				"notification not found",
				"request_id", requestlog.GetRequestID(r.Context()),
				"error", err,
			)
		}
		slog.Error(
			"unable to mark notification as read",
			"request_id", requestlog.GetRequestID(r.Context()),
			"notification_id", notificationID,
			"error", err,
		)
		// http.Error(w, "Unable to mark notification as read", http.StatusInternalServerError)
		// return
	}

	itemID := url.QueryEscape(newNotif.ItemID)
	destination := "/notifications"

	switch newNotif.Type {
	case models.NotificationSongCreated, models.NotificationSongImported:
		destination = "/song?id=" + itemID
		fmt.Println("song created: ", destination)

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
	case models.NotificationEventCreated, models.NotificationEventUpdated:
		destination = "/event?event-id=" + itemID
	case models.NotificationChatCreated, models.NotificationChatMention:
		destination = "/chat?id=" + itemID
	case models.NotificationNewMember:
		destination = "/band?band-id=" + itemID
	case models.NotificationTodoListAssigned:
		destination = "/todo?todo-id=" + itemID
	}

	http.Redirect(w, r, destination, http.StatusSeeOther)
}

func (h Handler) HandlerNotificationsMarkAllRead(w http.ResponseWriter, r *http.Request) {
	fmt.Println("- HandlerNotification")

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

	err = database.NotifTableMarkAllAsRead(r.Context(), user.UserID, band.BandID)
	if err != nil {
		slog.Error(
			"unable to mark notifications as read",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
	}

	http.Redirect(w, r, "/notifications", http.StatusSeeOther)
}

type pushSubscriptionRequest struct {
	Endpoint string `json:"endpoint"`
	Keys     struct {
		P256DH string `json:"p256dh"`
		Auth   string `json:"auth"`
	} `json:"keys"`
}

func (h Handler) HandlerPushSubscription(w http.ResponseWriter, r *http.Request) {
	fmt.Println("- HandlerPushSubscription")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"error getting auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Sign in to enable push notifications", http.StatusUnauthorized)
		return
	}

	var subscription pushSubscriptionRequest
	body := http.MaxBytesReader(w, r.Body, 8<<10)
	decoder := json.NewDecoder(body)

	if err := decoder.Decode(&subscription); err != nil {
		slog.Error(
			"error with json decoder",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Invalid push subscription JSON", http.StatusBadRequest)
		return
	}

	if err := decoder.Decode(new(json.RawMessage)); err != io.EOF {
		http.Error(w, "Expected one push subscription", http.StatusBadRequest)
		return
	}

	endpoint, err := url.Parse(subscription.Endpoint)
	if err != nil || endpoint.Scheme != "https" || endpoint.Hostname() == "" ||
		endpoint.User != nil || endpoint.Fragment != "" ||
		subscription.Keys.P256DH == "" || subscription.Keys.Auth == "" {
		http.Error(w, "Invalid push subscription", http.StatusBadRequest)
		return
	}

	err = database.PushSubscriptionsTableCreateSubscription(
		r.Context(),
		auth.User.UserID,
		subscription.Endpoint,
		subscription.Keys.P256DH,
		subscription.Keys.Auth,
	)
	if err != nil {
		if errors.Is(err, database.ErrPushSubscriptionOwnedByAnotherUser) {
			http.Error(w, "Subscription belongs to another account", http.StatusConflict)
			return
		}

		slog.Error("unable to save push subscription", "error", err)
		http.Error(w, "Unable to save push subscription", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
