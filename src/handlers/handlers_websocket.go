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
			"request started",
			"request_id", requestlog.GetRequestID(r.Context()),
			"user_id", auth.User.UserID,
			"band_id", auth.CurrentBand.BandID,
			"method", r.Method,
			"path", r.URL.Path,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusUnauthorized)
		return
	}

	if h.Hub == nil {
		log.Println("   WebSocket hub is not initialized")
		http.Error(w, "Websocket service unavailable", http.StatusServiceUnavailable)
		return
	}

	userChatIDs, err := database.ChatMembersTableGetChatIDsByUserID(auth.User.UserID)
	if err != nil {
		log.Println("   Unable to get user chat ids from chat members table: ", err)
		http.Error(w, "Unable to get user chat history", http.StatusInternalServerError)
		return
	}

	realtime.ServeWebSocket(
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
