package realtime

import (
	"bandplan/src/database"
	"bytes"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 4096
)

type IncomingMessage struct {
	Type   string `json:"type"`
	ChatID string `json:"chat_id"`
	Body   string `json:"body"`
}

type OutgoingMessage struct {
	Type             string `json:"type"`
	ChatID           string `json:"chat_id"`
	MessageID        string `json:"message_id"`
	UserID           string `json:"user_id"`
	UserName         string `json:"user_name"`
	ProfileImagePath string `json:"profile_image_path"`
	Body             string `json:"body"`
	CreatedAt        string `json:"created_at"`
	DisplayTime      string `json:"display_time"`
}

type Client struct {
	hub              *Hub
	conn             *websocket.Conn
	send             chan []byte
	bandID           string
	userID           string
	userName         string
	profileImagePath string
	chatIDs          map[string]bool
}

func NewClient(
	hub *Hub,
	conn *websocket.Conn,
	bandID string,
	userID string,
	userName string,
	profileImagePath string,
	chatIDs map[string]bool,
) *Client {
	return &Client{
		hub:              hub,
		conn:             conn,
		send:             make(chan []byte, 256),
		bandID:           bandID,
		userID:           userID,
		userName:         userName,
		profileImagePath: profileImagePath,
		chatIDs:          chatIDs,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)

	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Printf("Unable to set WebSocket read deadline: %v", err)
	}

	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseNormalClosure,
			) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		if messageType != websocket.TextMessage {
			continue
		}

		message = bytes.TrimSpace(message)

		if len(message) == 0 {
			continue
		}

		var incoming IncomingMessage

		err = json.Unmarshal(message, &incoming)
		if err != nil {
			log.Printf("Unable to decode WebSocket message: %v", err)
			continue
		}

		incoming.Body = strings.TrimSpace(incoming.Body)

		if incoming.Type != "chat_message" {
			continue
		}

		if incoming.Body == "" {
			continue
		}

		savedMessage, err := database.MessagesTableCreateMessage(
			c.bandID,
			c.userID,
			c.userName,
			incoming.ChatID,
			incoming.Body,
		)
		if err != nil {
			log.Printf("Unable to save WebSocket chat messages: %v", err)
			continue
		}

		outgoing := OutgoingMessage{
			Type:             "chat_message",
			ChatID:           "chat_id",
			MessageID:        savedMessage.MessageID,
			UserID:           savedMessage.UserID,
			UserName:         c.userName,
			ProfileImagePath: c.profileImagePath,
			Body:             savedMessage.Body,
			CreatedAt:        savedMessage.CreatedAt.Format(time.RFC3339),
			DisplayTime:      savedMessage.CreatedAt.Format("3:04 PM"),
		}

		encodedMessage, err := json.Marshal(outgoing)
		if err != nil {
			log.Printf("   Unable to encode WebSocket message %v", err)
			continue
		}

		c.hub.broadcast <- BroadcastMessage{
			BandID:  c.bandID,
			ChatID:  incoming.ChatID,
			Message: encodedMessage,
		}
	}
}

func (c *Client) WritePump() {
	log.Println("- WritePump")
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}

			if !ok {
				c.conn.WriteMessage(
					websocket.CloseMessage,
					[]byte{},
				)
				return
			}

			err = c.conn.WriteMessage(
				websocket.TextMessage,
				message,
			)
			if err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}

		case <-ticker.C:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}

			err = c.conn.WriteMessage(
				websocket.PingMessage,
				nil,
			)
			if err != nil {
				return
			}
		}
	}
}
