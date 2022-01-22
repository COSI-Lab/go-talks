package main

import (
	"encoding/json"
	"log"
	"sync"
)

type Hub struct {
	sync.RWMutex

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
	hub.RLock()
	connections := len(hub.clients)
	hub.RUnlock()

	return connections
}

func processMessage(message Message) bool {
	switch message.Type {
	case New:
		DB.Create(&Talk{})
	}

	return false
}

func (hub *Hub) run() {
	for {
		hub.Lock()
		select {
		case client := <-hub.register:
			// registers a client
			hub.clients[client] = true
		case client := <-hub.unregister:
			// unregister a client
			delete(hub.clients, client)
			close(client.send)
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
		hub.Unlock()
	}
}
