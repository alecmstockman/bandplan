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

func ServeChatWebSocket(
	hub *Hub,
	bandID string,
	userID string,
	userName string,
	profileImagePath string,
	chatIDs map[string]bool,
	w http.ResponseWriter,
	r *http.Request,
) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Unable to upgrade chat WebSocket connection: %v", err)
		return
	}

	client := NewClient(
		hub,
		conn,
		chatClient,
		bandID,
		userID,
		userName,
		profileImagePath,
		chatIDs,
	)

	hub.register <- client

	go client.WritePump()

	client.ReadPump()
}

func ServeNotificationsWebSocket(
	hub *Hub,
	bandID string,
	userID string,
	w http.ResponseWriter,
	r *http.Request,
) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Unable to upgrade notifications WebSocket connection: %v", err)
		return
	}

	client := NewClient(
		hub,
		conn,
		notificationClient,
		bandID,
		userID,
		"",
		"",
		nil,
	)

	hub.register <- client

	go client.WritePump()

	client.ReadPump()
}
