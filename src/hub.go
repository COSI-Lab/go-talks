package main

import (
	"encoding/json"
	"log"
	"time"
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

// Returns the next Wednesday in YYYYMMDD format
// If today is a wednesday today's date is returned
func nextWednesday() string {
	const format = "20060102"
	now := time.Now()

	daysUntilWednesday := time.Wednesday - now.Weekday()

	if daysUntilWednesday == 0 {
		return now.Format(format)
	} else if daysUntilWednesday > 0 {
		return now.AddDate(0, 0, int(daysUntilWednesday)).Format(format)
	} else {
		return now.AddDate(0, 0, int(daysUntilWednesday)+7).Format(format)
	}
}

func processMessage(message Message) bool {
	switch message.Type {
	case NEW:
		log.Println("Create talk {", message.Name, message.Description, message.Talktype, "}")

		if message.Week == "" {
			message.Week = nextWednesday()
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
		talk.Order = 0

		CreateTalk(talk)
	case HIDE:
		log.Println("Hide talk {", message.Id, "}")
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
			log.Println("Registered client", client.conn.RemoteAddr())
		case client := <-hub.unregister:
			// unregister a client
			delete(hub.clients, client)
			close(client.send)
			log.Println("Unregistered client", client.conn.RemoteAddr())
		case message := <-hub.broadcast:
			// broadcasts the message to all clients (including the one that sent the message)
			processMessage(message)

			// Serialize message into a byte slice
			bytes, err := json.Marshal(message)
			if err != nil {
				log.Println("failed to marshal", err)
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
