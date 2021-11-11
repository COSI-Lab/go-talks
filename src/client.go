package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Client struct {
	// The websocket connection
	conn *websocket.Conn

	// Outbound messages
	send chan []byte
}

func (c *Client) read() {
	defer func() {
		hub.toggle <- c
		c.conn.Close()
	}()

	for {
		// Read from the connection
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		// TODO Check that the message is in a good format
		hub.broadcast <- message
	}
}

func (c *Client) write() {
	defer func() {
		c.conn.WriteMessage(websocket.CloseMessage, []byte{})
		c.conn.Close()
	}()

	for message := range c.send {
		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}

		w.Write(message)
	}
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Handles the websocket
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		w.WriteHeader(404)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{conn: conn, send: make(chan []byte, 256)}
	hub.toggle <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.write()
	go client.read()
}
