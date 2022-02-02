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

// Logs request Method and request URI
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("[INFO] " + r.Method + " " + r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func main() {
	talks_password = os.Getenv("TALKS_PASSWORD")

	// Connect to the database
	err := ConnectDB("sqlite")
	if err != nil {
		log.Fatalln("[ERROR] Failed to connect to the database", err)
	}

	// Set up all tables
	MakeDB()

	// Load templates
	tmpls, err = template.ParseGlob("templates/*")

	if err != nil {
		log.Fatalln("[ERROR] Failed to compile some template(s)", err)
	} else {
		log.Println("[INFO]", tmpls.DefinedTemplates())
	}

	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	// "api" endpoints
	r.HandleFunc("/talks", talksHandler)
	r.HandleFunc("/ws", socketHandler)
	r.HandleFunc("/health", healthHandler)

	// static files
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	// templated pages
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/{week:[0-9]{8}}", weekHandler)

	// Set up server listen address
	listenAddr, exists := os.LookupEnv("LISTEN")
	if !exists {
		listenAddr = ""
	}

	// Set up server port
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "5001"
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
	log.Println("[INFO] Web server is now listening for connections on http://localhost:" + port)
	log.Fatal(srv.ListenAndServe())
}
