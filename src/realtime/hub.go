package realtime

type BroadcastMessage struct {
	ClientType ClientType
	ItemID     string
	BandID     string
	UserID     string
	Message    []byte
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan BroadcastMessage
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan BroadcastMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case broadcast := <-h.broadcast:
			for client := range h.clients {
				if client.bandID != broadcast.BandID {
					continue
				}
				if client.clientType != broadcast.ClientType {
					continue
				}
				switch broadcast.ClientType {
				case chatClient:
					if !client.chatIDs[broadcast.ItemID] {
						continue
					}
				case notificationClient:
					if client.userID != broadcast.UserID {
						continue
					}
				default:
					continue
				}

				select {
				case client.send <- broadcast.Message:
				default:
					close(client.send)
					delete(h.clients, client)
				}

			}
		}
	}
}

func (h *Hub) PublishNotification(bandID, recipientUserID string) {
	h.broadcast <- BroadcastMessage{
		ClientType: notificationClient,
		BandID:     bandID,
		UserID:     recipientUserID,
		Message:    []byte(`{"type":"notification_message"}`),
	}
}
