package handlers

import (
	"bandplan/src/database"
	requestlog "bandplan/src/logging"
	"bandplan/src/realtime"
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

	realtime.ServeChatWebSocket(
		h.Hub,
		auth.CurrentBand.BandID,
		auth.User.UserID,
		auth.User.Name,
		auth.User.ProfileImagePath,
		userChatIDs,
		w,
		r,
	)

}

func (h Handler) HandlerNotificationsWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerNotificationsWebSocket")

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
