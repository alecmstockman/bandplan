package realtime

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ServeWebSocket(
	hub *Hub,
	bandID string,
	userID string,
	userName string,
	profileImagePath string,
	w http.ResponseWriter,
	r *http.Request,
) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Unable to upgrade WebSocket connection: %v", err)
		return
	}

	client := NewClient(
		hub,
		conn,
		bandID,
		userID,
		userName,
		profileImagePath,
	)

	hub.register <- client

	go client.WritePump()

	client.ReadPump()
}
