package handlers

import (
	"bandplan/src/database"
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"bandplan/src/realtime"
	"context"
	"log"
	"log/slog"
	"net/http"
)

func (h Handler) HandlerChatWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerChatWebSocket")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusUnauthorized)
		return
	}

	if h.Hub == nil {
		http.Error(w, "Websocket service unavailable", http.StatusServiceUnavailable)
		return
	}

	userChatIDs, err := database.ChatMembersTableGetChatIDsByUserID(
		auth.User.UserID,
		auth.CurrentBand.BandID,
	)
	if err != nil {
		slog.Error(
			"unable to load chat messages",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to get user chat history", http.StatusInternalServerError)
		return
	}

	var sendPush func(context.Context, string, models.PushNotifType) error
	if h.Services != nil {
		sendPush = h.Services.SendPushToUser
	}

	realtime.ServeChatWebSocket(
		h.Hub,
		auth.CurrentBand.BandID,
		auth.User.UserID,
		auth.User.Name,
		auth.User.ProfileImagePath,
		userChatIDs,
		sendPush,
		w,
		r,
	)

}

func (h Handler) HandlerNotificationsWebSocket(w http.ResponseWriter, r *http.Request) {

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusUnauthorized)
		return
	}

	if h.Hub == nil {
		http.Error(w, "Websocket service unavailable", http.StatusServiceUnavailable)
		return
	}

	realtime.ServeNotificationsWebSocket(
		h.Hub,
		auth.CurrentBand.BandID,
		auth.User.UserID,
		w,
		r,
	)
}
