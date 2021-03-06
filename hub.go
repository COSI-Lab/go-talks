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

func processMessage(message *Message) bool {
	switch message.Type {
	case NEW:
		if message.Week == "" {
			message.Week = nextWednesday()
		}

		// You can not create a talk for a previous meeting
		if isPast(nextWednesday(), message.Week) {
			return false
		}

		// Validate talk type
		talk := &Talk{}
		if *message.Talktype > 4 {
			return false
		}
		talk.Type = *message.Talktype

		// Validate talk description
		if message.Description == "" {
			return false
		}
		talk.Description = message.Description

		// Update the message's description to be markdowned
		message.Description = string(markDowner(message.Description))

		// Validate talk name
		if message.Name == "" {
			return false
		}
		talk.Name = message.Name
		talk.Week = message.Week

		// TODO: Talk order
		talk.Order = 0

		message.Id = CreateTalk(talk)
	case HIDE:
		// Only hide talks during meetings, otherwise delete them
		if !DuringMeeting() {
			log.Println("[INFO] Delete talk {", message.Id, "}")
			DeleteTalk(message.Id)
		} else {
			log.Println("[INFO] Hide talk {", message.Id, "}")
			HideTalk(message.Id)
		}
	case DELETE:
		log.Println("[INFO] Delete talk {", message.Id, "}")
		DeleteTalk(message.Id)
	}

	return true
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
