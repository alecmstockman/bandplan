package handlers

import (
	"bandplan/src/realtime"
	"log"
	"net/http"
)

func (h Handler) HandlerChatWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerChatWebSocket")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusUnauthorized)
		return
	}

	if h.Hub == nil {
		log.Println("   WebSocket hub is not initialized")
		http.Error(w, "Websocket service unavailable", http.StatusServiceUnavailable)
		return
	}

	realtime.ServeWebSocket(
		h.Hub,
		auth.CurrentBand.BandID,
		auth.User.UserID,
		auth.User.Name,
		auth.User.ProfileImagePath,
		w,
		r,
	)

}
