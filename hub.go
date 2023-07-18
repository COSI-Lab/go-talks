package main

import (
	"encoding/json"
	"log"
)

// Hub maintains a set of active clients and broadcasts messages to the clients
type Hub struct {
	// Map of clients
	clients map[*Client]struct{}

	// Inbound messages from the clients
	broadcast chan Message

	// registers a client from the hub
	register chan *Client

	// unregister a client from the hub
	unregister chan *Client
}

func (hub *Hub) countConnections() int {
	return len(hub.clients)
}

func processMessage(message *Message) bool {
	switch message.Type {
	case NEW:
		// You can not create a talk for a previous meeting
		wednesday := nextWednesday()
		if message.New.Week == "" {
			message.New.Week = wednesday
		} else if isPast(wednesday, message.New.Week) {
			return false
		}

		// Validate talk type
		if message.New.Talktype > 4 {
			return false
		}

		// Validate talk description
		if message.New.Description == "" {
			return false
		}
		// Validate talk name
		if message.New.Name == "" {
			return false
		}

		message.New.ID = talks.Create(message.New.Name, message.New.Talktype, message.New.Description, message.New.Week)

		// Update the message's description to be parsed as markdown
		message.New.Description = string(markDownerSafe(message.New.Description))

		return true
	case HIDE:
		// During meetings we hide talks instead of deleting them
		if duringMeeting() {
			log.Println("[INFO] Hide talk {", message.Hide.ID, "}")
			talks.Hide(message.Hide.ID)
		} else {
			log.Println("[INFO] Delete talk {", message.Hide.ID, "}")
			talks.Delete(message.Hide.ID)
		}
		return true
	case DELETE:
		log.Println("[INFO] Delete talk {", message.Hide.ID, "}")
		talks.Delete(message.Hide.ID)
		return true
	default:
		return false
	}
}

func (hub *Hub) run() {
	for {
		select {
		case client := <-hub.register:
			// registers a client
			hub.clients[client] = struct{}{}
		case client := <-hub.unregister:
			// unregister a client
			delete(hub.clients, client)
			close(client.send)
		case message := <-hub.broadcast:
			log.Println("[INFO] Broadcast message:", message)

			// broadcasts the message to all clients (including the one that sent the message)
			if !processMessage(&message) {
				log.Println("[WARN] Invalid message")
			}

			// Serialize message into a byte slice
			bytes, err := json.Marshal(message)
			if err != nil {
				log.Println("[WARN] failed to marshal", err)
			}

			for client := range hub.clients {
				select {
				case client.send <- bytes:
				default:
					// if sending to a client blocks we drop the client
					close(client.send)
					delete(hub.clients, client)
				}
			}
		}
	}
}
