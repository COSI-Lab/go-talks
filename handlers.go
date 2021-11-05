package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/thanhpk/randstr"
)

type client struct {
	ch       chan []byte
	isAuthed bool
}

var clients map[string]client       // map of clients and channels for each client
var clients_lock sync.RWMutex       // Lock for clients map
var upgrader = websocket.Upgrader{} // use default options
var messages chan []byte            // Main bus for all messages
var messages_lock sync.RWMutex      // Lock for messages

func indexHandler(w http.ResponseWriter, r *http.Request) {

}

func allHandler(w http.ResponseWriter, r *http.Request) {

}

func talksHandler(w http.ResponseWriter, r *http.Request) {

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Return list of active clients
	// Mostly for diagnostic purposes
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprint(len(clients))))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// Create UUID but badly
	// Should work as we arent serving enough clients were psuedo random will mess us up
	id := randstr.Hex(16)

	// Check if client is within our subnet and so should be auto authed
	strIp := r.Header.Get("X-Forwarded-For")
	ip := net.ParseIP(strIp)

	authed := false
	if ip != nil {
		authed = isInSubnet(ip)
	} else {
		log.Println("Could not parse ip address \"X-Forwarded-For\":", strIp)
	}

	// Create client object
	client := client{make(chan []byte), authed}
	clients_lock.Lock()
	clients[id] = client
	clients_lock.Unlock()
	log.Printf("new connection registered: %s\n", id)

	// Send id to client
	w.WriteHeader(200)
	w.Write([]byte(id))
}

func authenticateHandler(w http.ResponseWriter, r *http.Request) {
	// Authenticate client
	// message will be password, id
	msg, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %s\n", err)
		w.WriteHeader(500)
		return
	}

	pw := strings.Split(string(msg), ",")[0]
	id := strings.Split(string(msg), ",")[1]

	// Cringe hardcoding of password but who cares
	if pw == "temp pw" {
		clients_lock.Lock()
		if client, exists := clients[id]; exists {
			client.isAuthed = true
		}
		clients_lock.Unlock()
		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}
}
