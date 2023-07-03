package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// TemplateResponse contains the information needed to render the past and future templates
type TemplateResponse struct {
	Talks     []Talk
	HumanWeek string
	Week      string
	NextWeek  string
	PrevWeek  string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	week := nextWednesday()

	// Validate week and get human readable version
	human, _ := weekForHumans(week)

	// Prepare response
	talks := VisibleTalks(week)
	res := TemplateResponse{Talks: talks, Week: week, HumanWeek: human, NextWeek: addWeek(week), PrevWeek: subtractWeek(week)}

	// Render the template
	err := tmpls.ExecuteTemplate(w, "future.gohtml", res)
	if err != nil {
		log.Println("[WARN] Failed to render template:", err)
	}
}

// /{week:[0-9]{8}}
func weekHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	week := vars["week"]

	// Validate week and get human readable version
	human, err := weekForHumans(week)
	if err != nil {
		log.Println("[WARN] Requested invalid week:", week)
		w.WriteHeader(400)
		return
	}

	// Prepare response
	res := TemplateResponse{Week: week, HumanWeek: human, NextWeek: addWeek(week), PrevWeek: subtractWeek(week)}

	// Render the template
	if isPast(nextWednesday(), week) {
		res.Talks = AllTalks(week)
		err = tmpls.ExecuteTemplate(w, "past.gohtml", res)
	} else {
		res.Talks = VisibleTalks(week)
		err = tmpls.ExecuteTemplate(w, "future.gohtml", res)
	}

	if err != nil {
		log.Println("[WARN] Failed to render template:", err)
	}
}

func indexTalksHandler(w http.ResponseWriter, r *http.Request) {
	week := nextWednesday()

	// Validate week and get human readable version
	_, err := weekForHumans(week)
	if err != nil {
		log.Println("[WARN] Invalid week:", week)
		w.WriteHeader(400)
		return
	}

	talks := AllTalks(week)

	// Parse talks as JSON
	err = json.NewEncoder(w).Encode(talks)
	if err != nil {
		log.Println("[WARN] Failed to encode talks:", err)
	}
}

// /{week:[0-9]{8}}/talks returns json of talks for a given week
func talksHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	week := vars["week"]

	// Validate week and get human readable version
	_, err := weekForHumans(week)
	if err != nil {
		log.Println("[WARN] Invalid week:", week)
		w.WriteHeader(400)
		return
	}

	talks := AllTalks(week)

	// Parse talks as JSON
	err = json.NewEncoder(w).Encode(talks)
	if err != nil {
		log.Println("[WARN] Failed to encode talks:", err)
	}
}

// /img/{id}
func imageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Lock the cache
	cacheLock.RLock()
	defer cacheLock.RUnlock()

	// Decode the id ([32]byte)
	hash, err := base64.URLEncoding.DecodeString(id)
	if err != nil {
		log.Println("[WARN] failed to decode image id:", err)
		w.WriteHeader(400)
		return
	}
	if len(hash) != 32 {
		log.Println("[WARN]", "Invalid hash length")
		w.WriteHeader(400)
		return
	}

	// Copy the hash into a new [32]byte
	var hash32 [32]byte
	copy(hash32[:], hash)

	// Get the image from the cache
	image, ok := cache[hash32]
	if !ok {
		w.WriteHeader(404)
		return
	}

	// Write the image to the response
	w.Header().Set("Content-Type", image.ContentType)
	w.Write(image.Data)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Return list of active clients for diagnostic purposes
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprint(hub.countConnections())))
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the connection to a websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[WARN] failed to upgrade connection:", err)
		return
	}

	addr, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
	ip := net.ParseIP(addr)

	authenticated := false
	if trustedNetworks.Contains(ip) {
		// Could be a load balancer, check the X-REAL-IP header
		if realIP := r.Header.Get("X-REAL-IP"); realIP != "" {
			// Check if the real IP is in the trusted network
			ip = net.ParseIP(realIP)
			if trustedNetworks.Contains(ip) {
				authenticated = true
			}
		} else {
			authenticated = true
		}
	}
	log.Printf("[INFO] New connection from %s (authenticated: %t)", ip, authenticated)

	client := &Client{conn: conn, send: make(chan []byte), auth: authenticated}
	hub.register <- client

	// Run send and receive in goroutines
	go client.write()
	go client.read()

	// Send an authentication response
	client.send <- authenticatedMessage(authenticated)
}

// Post contains the information needed to render a markdown post
type Post struct {
	Title   string
	Content string
}

// Creates a handler that serves static html after rendering markdown
func markdownFactory(post string) func(http.ResponseWriter, *http.Request) {
	path := "posts/" + post + ".md"

	// Read content from file
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("[FATAL] Failed to read markdown file:", err)
	}

	// Create a dummy writer to capture the output of the template
	buff := bytes.NewBuffer(nil)

	// Render the markdown
	err = tmpls.ExecuteTemplate(buff, "markdown.gohtml", Post{
		Title:   post,
		Content: string(content),
	})

	if err != nil {
		log.Fatal("[FATAL] Failed to render markdown template:", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(buff.Bytes())
	}
}
