package ws

import (
	"encoding/json"
	"log"
	"slices"

	"github.com/google/uuid"
)

type Hub struct {
	Clients      map[*Client]bool
	Broadcast    chan []byte
	BroadcastFor chan []byte
	register     chan *Client
	unregister   chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Clients:      map[*Client]bool{},
		Broadcast:    make(chan []byte),
		BroadcastFor: make(chan []byte),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
	}
}

type broadcastForRequest struct {
	Clients []uuid.UUID
	Message []byte
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.Clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.send)
			}
			h.Clients[client] = false
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.Clients, client)
				}
			}
		case broadcastForReq := <-h.BroadcastFor:
			var broadcastFor broadcastForRequest
			err := json.Unmarshal(broadcastForReq, &broadcastFor)
			if err != nil {
				log.Printf("coudlnt unmarshal broadcast request: %v", err)
			}

			clientsToBroadcast := []*Client{}
			for client := range h.Clients {
				if slices.Contains(broadcastFor.Clients, client.UserID) {
					clientsToBroadcast = append(clientsToBroadcast, client)
				}
			}

			for _, client := range clientsToBroadcast {
				select {
				case client.send <- broadcastFor.Message:
				default:
					close(client.send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
