package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	// The websocket connection
	conn *websocket.Conn

	// Outbound messages
	send chan []byte

	// Marks the client as authenticated
	auth bool
}

// MessageType are the types of messages that can be sent/received
type MessageType uint

const (
	// NEW creates a new talk
	NEW MessageType = iota
	// HIDE hides a talk
	HIDE
	// DELETE deletes a talk
	DELETE
	// AUTH communicates authentication request/response
	AUTH
)

// Message is the format of messages sent between the client and server
// Since go doesn't have the strongest type system we pack every message
// into a single struct
type Message struct {
	Type        MessageType `json:"type"`
	ID          uint32      `json:"id,omitempty"`
	Password    string      `json:"password,omitempty"`
	Name        string      `json:"name,omitempty"`
	Talktype    *TalkType   `json:"talktype,omitempty"`
	Description string      `json:"description,omitempty"`
	Week        string      `json:"week,omitempty"`
}

func authenticatedMessage(b bool) []byte {
	if b {
		return []byte("{\"type\": 3, \"status\": true}")
	}

	return []byte("{\"type\": 3, \"status\": false}")
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
			log.Printf("[WARN] %v", err)
			break
		}

		// Format the message
		var message Message
		err = json.Unmarshal(raw, &message)
		if err != nil {
			continue
		}

		// Handle authentication without consulting the hub
		if message.Type == AUTH {
			log.Printf("[INFO] Client %v is trying to authenticate", c.conn.RemoteAddr())

			c.auth = message.Password == config.Password
			c.send <- authenticatedMessage(c.auth)

			continue
		}

		// Forward all other message to be processed and broadcasted to other client
		hub.broadcast <- message
	}
}

func (c *Client) write() {
	defer func() {
		c.conn.WriteMessage(websocket.CloseMessage, []byte{})
		c.conn.Close()
	}()

	ticker := time.NewTicker(30 * time.Second)

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				return
			}

			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("[WARN] %v", err)
				return
			}
		case <-ticker.C:
			err := c.conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Printf("[WARN] %v", err)
				return
			}
		}
	}
}
