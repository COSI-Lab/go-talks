package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	// The websocket connection
	conn *websocket.Conn

	// Outbound messages
	send chan []byte

	// Marks the client as autheniticated
	auth bool
}

type MessageType uint

const (
	Move MessageType = iota
	New
	Hide
	Delete
	Auth
)

type Message struct {
	Type        MessageType `json:"type"`
	Status      uint32      `json:"status,omitempty"`
	Id          uint32      `json:"id,omitempty"`
	Password    string      `json:"password,omitempty"`
	Name        string      `json:"name,omitempty`
	Talktype    TalkType    `json:"talktype,omitempty`
	Description string      `json:"description,omitempty`
}

func (c *Client) read() {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	for {
		// Read from the connection
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		// Format the message
		var message Message
		err = json.Unmarshal(raw, &message)
		if err != nil {
			continue
		}

		// Handle authentication without consulting the hub
		if message.Type == Auth {
			var resp []byte
			if message.Password == talks_password {
				resp, err = json.Marshal(Message{Type: Auth, Status: 200})
				if err != nil {
					log.Println("Marshalling password response failed! Should never happen.", err)
				}

				c.auth = true
			} else {
				resp, err = json.Marshal(Message{Type: Auth, Status: 403})
				if err != nil {
					log.Println("Marshalling password response failed! Should never happen.", err)
				}

				c.auth = false
			}

			c.send <- resp
			continue
		}

		// Foward all other message to be processed and broadcasted to other client
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
			break
		}

		w.Write(message)
	}
}

var upgrader = websocket.Upgrader{} // TODO: Don't use default options

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the connection to a websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{conn: conn, send: make(chan []byte, 256)}
	hub.register <- client

	// Run send and recieve in goroutines
	go client.write()
	go client.read()
}
