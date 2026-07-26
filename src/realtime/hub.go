package realtime

type BroadcastMessage struct {
	BandID  string
	Message []byte
}

type Hub struct {
	clients map[*Client]bool

	broadcast chan BroadcastMessage

	register chan *Client

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
