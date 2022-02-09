package main

import (
	"encoding/json"
	"log"
	"time"

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
	NEW MessageType = iota
	HIDE
	DELETE
	MOVE
	AUTH
)

type Message struct {
	Type        MessageType `json:"type"`
	Status      uint32      `json:"status,omitempty"`
	Id          uint32      `json:"id,omitempty"`
	Password    string      `json:"password,omitempty"`
	Name        string      `json:"name,omitempty"`
	Talktype    *TalkType   `json:"talktype,omitempty"`
	Description string      `json:"description,omitempty"`
	Week        string      `json:"week,omitempty"`
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
			log.Printf("[ERROR] %v", err)
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
			var resp []byte
			if message.Password == talks_password {
				resp, err = json.Marshal(Message{Type: AUTH})
				if err != nil {
					log.Println("[WARN] Marshalling password response failed! Should never happen.", err)
				}

				c.auth = true
			} else {
				resp, err = json.Marshal(Message{Type: AUTH})
				if err != nil {
					log.Println("[WARN] Marshalling password response failed! Should never happen.", err)
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

	ticker := time.NewTicker(30 * time.Second)

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				return
			}

			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("[ERROR] %v", err)
				return
			}
		case <-ticker.C:
			err := c.conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Printf("[ERROR] %v", err)
				return
			}
		}
	}
}
