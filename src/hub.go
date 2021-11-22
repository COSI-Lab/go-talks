package main

import "sync"

type Hub struct {
	sync.RWMutex

	// Map of clients TODO: Make a HashSet
	clients map[*Client]bool

	// Inbound messages from the clients
	broadcast chan []byte

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
			for client := range hub.clients {
				select {
				case client.send <- message:
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
