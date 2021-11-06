package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// Connect to the database
	err := ConnectDB("sqlite")
	if err != nil {
		log.Fatalln("Failed to connect to the database")
	}

	// Set up all tables
	MakeDB()

	r := mux.NewRouter()

	// templated pages
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/all", allHandler)

	// "api" endpoints
	r.HandleFunc("/talks", talksHandler)
	r.HandleFunc("/register", registerHandler)
	r.HandleFunc("/authenticate", authenticateHandler)
	r.HandleFunc("/socket/{id}", socketHandler)
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

	// Set up ssl
	sslString, exists := os.LookupEnv("USE_SSL")
	var ssl bool
	if sslString == "true" || sslString == "yes" {
		ssl = true
	} else if !exists || sslString == "false" || sslString == "no" {
		ssl = false
	}

	// Create http server
	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%s", listenAddr, port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Start server
	log.Println("Web server is now listening for connections")
	if ssl {
		log.Fatal(srv.ListenAndServeTLS("certs/cert.pem", "certs/privkey.pem"))
	} else {
		log.Fatal(srv.ListenAndServe())
	}
}
