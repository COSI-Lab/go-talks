package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

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
		log.Println("[WARN]", err)
	}
}

// /{week:[0-9]{8}}
func weekHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	week := vars["week"]

	// Validate week and get human readable version
	human, err := weekForHumans(week)

	if err != nil {
		log.Println("[INFO]", err)
		w.WriteHeader(400)
		return
	}

	// Prepare response
	talks := VisibleTalks(week)
	res := TemplateResponse{Talks: talks, Week: week, HumanWeek: human, NextWeek: addWeek(week), PrevWeek: subtractWeek(week)}

	// Render the template
	if isPast(nextWednesday(), week) {
		err = tmpls.ExecuteTemplate(w, "past.gohtml", res)
	} else {
		err = tmpls.ExecuteTemplate(w, "future.gohtml", res)
	}

	if err != nil {
		log.Println("[WARN]", err)
	}
}

func indexTalksHandler(w http.ResponseWriter, r *http.Request) {
	week := nextWednesday()

	// Validate week and get human readable version
	_, err := weekForHumans(week)

	if err != nil {
		log.Println("[INFO]", err)
		w.WriteHeader(400)
		return
	}

	talks := AllTalks(week)

	// Parse talks as JSON
	err = json.NewEncoder(w).Encode(talks)
	if err != nil {
		log.Println("[WARN]", err)
	}
}

// /{week:[0-9]{8}}/talks returns json of talks for a given week
func talksHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	week := vars["week"]

	// Validate week and get human readable version
	_, err := weekForHumans(week)

	if err != nil {
		log.Println("[INFO]", err)
		w.WriteHeader(400)
		return
	}

	talks := AllTalks(week)

	// Parse talks as JSON
	err = json.NewEncoder(w).Encode(talks)
	if err != nil {
		log.Println("[WARN]", err)
	}
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
		log.Println("[WARN]", err)
		return
	}

	client := &Client{conn: conn, send: make(chan []byte)}
	hub.register <- client

	// Run send and recieve in goroutines
	go client.write()
	go client.read()
}
