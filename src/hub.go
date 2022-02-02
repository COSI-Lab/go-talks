package main

import (
	"encoding/json"
	"log"
)

type Hub struct {
	// Map of clients TODO: Make a HashSet?
	clients map[*Client]bool

	// Inbound messages from the clients
	broadcast chan Message

	// registers a client from the hub
	register chan *Client

	// unregister a client from the hub
	unregister chan *Client
}

func (hub *Hub) countConnections() int {
	connections := len(hub.clients)

	return connections
}

func processMessage(message Message) bool {
	switch message.Type {
	case NEW:
		log.Println("[INFO] Create talk {", message.Name, message.Description, message.Talktype, "}")

		if message.Week == "" {
			message.Week = nextWednesday()
		}

		// You can not create a talk for a previous meeting
		if isPast(nextWednesday(), message.Week) {
			return false
		}

		talk := &Talk{}
		if *message.Talktype > 4 {
			return false
		}
		talk.Type = *message.Talktype

		if message.Description == "" {
			return false
		}
		talk.Description = message.Description

		if message.Name == "" {
			return false
		}
		talk.Name = message.Name
		talk.Week = message.Week
		talk.Order = 0

		CreateTalk(talk)
	case HIDE:
		log.Println("[INFO] Hide talk {", message.Id, "}")
		HideTalk(message.Id)
	}

	return false
}

func (hub *Hub) run() {
	for {
		select {
		case client := <-hub.register:
			// registers a client
			hub.clients[client] = true
			log.Println("[INFO] Registered client", client.conn.RemoteAddr())
		case client := <-hub.unregister:
			// unregister a client
			delete(hub.clients, client)
			close(client.send)
			log.Println("[INFO] Unregistered client", client.conn.RemoteAddr())
		case message := <-hub.broadcast:
			// broadcasts the message to all clients (including the one that sent the message)
			processMessage(message)

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
