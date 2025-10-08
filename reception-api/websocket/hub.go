package websocket

import (
	"encoding/json"
	"log"
)

const (
	MessageTypePatientCreated     = "patient_created"
	MessageTypePatientDeleted     = "patient_deleted"
	MessageTypePatientHISIDUpdate = "patient_his_id_update"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan Message, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("WebSocket client connected. Total: %d", len(h.clients))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("WebSocket client disconnected. Total: %d", len(h.clients))
			}

		case message := <-h.broadcast:
			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling WebSocket message: %v", err)
				continue
			}

			log.Printf("Broadcasting message type: %s to %d clients", message.Type, len(h.clients))

			for client := range h.clients {
				select {
				case client.send <- data:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) BroadcastPatientCreated(patient interface{}) {
	h.broadcast <- Message{
		Type: MessageTypePatientCreated,
		Data: patient,
	}
}

func (h *Hub) BroadcastPatientDeleted(id int) {
	h.broadcast <- Message{
		Type: MessageTypePatientDeleted,
		Data: map[string]int{"id": id},
	}
}

func (h *Hub) BroadcastPatientHISIDUpdate(patientID int, hisID string) {
	h.broadcast <- Message{
		Type: MessageTypePatientHISIDUpdate,
		Data: map[string]interface{}{
			"id":             patientID,
			"his_patient_id": hisID,
		},
	}
}
