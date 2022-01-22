package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var talks_password = ""
var hub Hub
var tmpls *template.Template

func main() {
	talks_password = os.Getenv("TALKS_PASSWORD")

	// Connect to the database
	err := ConnectDB("sqlite")
	if err != nil {
		log.Fatalln("Failed to connect to the database")
	}

	// Set up all tables
	MakeDB()

	// Load templates
	tmpls, err = template.ParseGlob("templates/*")

	if err != nil {
		log.Fatalln("Failed to compile some template(s)", err)
	} else {
		log.Println(tmpls.DefinedTemplates())
	}

	r := mux.NewRouter()

	// templated pages
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/all", allHandler)

	// "api" endpoints
	r.HandleFunc("/talks", talksHandler)
	r.HandleFunc("/socket", socketHandler)
	r.HandleFunc("/health", healthHandler)

	// static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	// Set up server listen address
	listenAddr, exists := os.LookupEnv("LISTEN")
	if !exists {
		listenAddr = ""
	}

	// Set up server port
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "5000"
	}

	// Create http server
	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%s", listenAddr, port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Start the hub
	hub = Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go hub.run()

	// Start server
	log.Println("Web server is now listening for connections on", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
