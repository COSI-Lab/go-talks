package main

import "sync"

type Hub struct {
	sync.RWMutex

	// Map of clients TODO: Make a HashSet
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// registers or unregister a client from the hub.
	toggle chan *Client
}

func (hub *Hub) connections() int {
	hub.RLock()
	connections := len(hub.clients)
	hub.RUnlock()

	return connections
}

func (hub *Hub) run() {
	for {
		hub.Lock()
		select {
		case client := <-hub.toggle:
			// registers or unregister a client from the hub.
			if _, exists := hub.clients[client]; exists {
				delete(hub.clients, client)
				close(client.send)
			} else {
				hub.clients[client] = true
			}
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
