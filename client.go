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
	// SYNC requests a sync of the talks
	SYNC
)

// Message is the format of messages sent between the client and server
// Since go doesn't have the strongest type system we pack every message
// into a single struct
type Message struct {
	Type MessageType    `json:"type"`
	New  *NewMessage    `json:"new,omitempty"`
	Hide *HideMessage   `json:"hide,omitempty"`
	Del  *DeleteMessage `json:"delete,omitempty"`
	Auth *AuthMessage   `json:"auth,omitempty"`
	Sync *SyncMessage   `json:"sync,omitempty"`
}

// NewMessage gives the client all the information needed to add a new talk to
// the page
type NewMessage struct {
	ID          uint32   `json:"id"`
	Name        string   `json:"name"`
	Talktype    TalkType `json:"talktype"`
	Description string   `json:"description"`
	Week        string   `json:"week"`
}

// HideMessage gives the client the ID of the talk to hide
type HideMessage struct {
	ID uint32 `json:"id"`
}

// DeleteMessage gives the client the ID of the talk to delete
type DeleteMessage struct {
	ID uint32 `json:"id"`
}

// AuthMessage gives the client the password to authenticate
type AuthMessage struct {
	Password string `json:"password"`
}

// SyncMessage starts and caps off a sync
type SyncMessage struct {
	Week string `json:"week"`
}

func authenticatedMessage(b bool) []byte {
	if b {
		return []byte("{\"type\": 3, \"auth\": {\"status\": true}}")
	}

	return []byte("{\"type\": 3, \"auth\": {\"status\": false}}")
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
			// Print the message and continue
			log.Printf("[WARN] %v", err)
			continue
		}

		switch message.Type {
		case NEW, HIDE, DELETE:
			// NEW, HIDE, and DELETE need to be serialized through the hub
			hub.broadcast <- message
		case AUTH:
			// AUTH is handled without having to contact the hub
			log.Printf("[INFO] Client %v is trying to authenticate", c.conn.RemoteAddr())
			c.send <- authenticatedMessage(message.Auth.Password == config.Password)
		case SYNC:
			// SYNC messages don't need to be serialized, go straight to the db
			log.Printf("[INFO] Client %v is requesting a sync", c.conn.RemoteAddr())
			for _, talk := range talks.AllTalks(message.Sync.Week) {
				var msg Message
				if talk.Hidden {
					// Send a hide message
					msg = Message{
						Type: HIDE,
						Hide: &HideMessage{
							ID: talk.ID,
						},
					}
				} else {
					// Send a create message
					msg = Message{
						Type: NEW,
						New: &NewMessage{
							ID:          talk.ID,
							Name:        talk.Name,
							Talktype:    talk.Type,
							Description: talk.Description,
							Week:        talk.Week,
						},
					}
				}
				// Send the message
				raw, _ := json.Marshal(msg)
				c.send <- raw
			}
			raw, _ := json.Marshal(message)
			c.send <- raw
		default:
			log.Printf("[WARN] Client %v sent an invalid message type", c.conn.RemoteAddr())
		}
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
				log.Println("[INFO] Closing connection")
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
