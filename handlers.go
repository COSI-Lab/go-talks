package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/thanhpk/randstr"
)

var clients map[string]chan []byte  // map of clients and channels for each client
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
	id := randstr.Hex(16)
	// Create UUID but badly
	// Should work as we arent serving enough clients were psuedo random will mess us up

	clients_lock.Lock()
	clients[id] = make(chan []byte, 10)
	clients_lock.Unlock()
	log.Printf("new connection registered: %s\n", id)

	// Send id to client
	w.WriteHeader(200)
	w.Write([]byte(id))
}

func authenticateHandler(w http.ResponseWriter, r *http.Request) {

}
