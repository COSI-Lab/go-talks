package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
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

	range1 := net.ParseIP("128.0.0.0")
	range2 := net.ParseIP("153.0.0.0")

	id := randstr.Hex(16)
	authed := false

	strIp := r.Header.Get("X-Forwarded-For")
	ip := net.ParseIP(strIp)

	if bytes.Compare(ip, range1) >= 0 && bytes.Compare(ip, range2) <= 0 {
		authed = true
	}

	//TODO: Add Ipv6 support

	// Create UUID but badly
	// Should work as we arent serving enough clients were psuedo random will mess us up
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

}
