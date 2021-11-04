package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Handles the websocket
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		w.WriteHeader(404)
		return
	}

	// get the channel
	ch := clients[id]
	log.Printf("%s connected!\n", id)

	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}

	defer func() {
		// Close connection gracefully
		conn.Close()
		clients_lock.Lock()
		log.Printf("Error sending message %s : %s", id, err)
		delete(clients, id)
		clients_lock.Unlock()
	}()

	// Quick cringe go funcs
	// One for reading from the websocket
	// One for writing to the websocket
	// TODO: Add db
	go func() {
		for {
			// Read message from client
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message %s : %s", id, err)
				return
			}
			clients_lock.Lock()
			for _, ch := range clients {
				// Add msg to channel for sending messages
				// Have to do it this way as websocket handler is seperate function
				select {
				case ch.ch <- msg:
				default:
					// the channel is blocking so we drop the data
				}
			}
			clients_lock.Unlock()
		}
	}()

	go func() {
		messages_lock.Lock()
		for _, ch := range clients {
			select {
			// Read message from channel
			case msg := <-ch.ch:
				// Add to main message bus. Basically allows us to bottleneck on purpose
				messages <- msg
			default:
				// dump the data
			}
		}
		messages_lock.Unlock()
	}()
	// Func above and below bottleneck the data and make sure all clients are on the same page
	go func() {
		for _, ch := range clients {
			select {
			case msg := <-messages:
				ch.ch <- msg
			default:
				// dump the data
			}
		}
	}()

	go func() {
		for {
			// Read message from channel
			msg := <-ch.ch
			// Send message to client
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("Error sending message %s : %s", id, err)
				return
			}
		}
	}()

}
