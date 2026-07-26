package realtime

import (
	"bytes"
	"encoding/json"
	"log"
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
	Type string `json:"type"`
	Body string `json:"body"`
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte

	bandID string
	userID string
}

func NewClient(
	hub *Hub,
	conn *websocket.Conn,
	bandID string,
	userID string,
) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		bandID: bandID,
		userID: userID,
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
			log.Println("Unable to decode WebSocket message: %v", err)
			continue
		}

		if incoming.Type == "" || incoming.Body == "" {
			continue
		}

		encodedMessage, err := json.Marshal(incoming)
		if err != nil {
			log.Printf("Unable to encode WebSocket message: %v", err)
			continue
		}

		c.hub.broadcast <- BroadcastMessage{
			BandID:  c.bandID,
			Message: encodedMessage,
		}

	}
}

func (c *Client) WritePump() {
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
