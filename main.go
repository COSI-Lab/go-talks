package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
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

	l := &http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	log.Printf("Serving on http://localhost:%d", 8000)
	log.Fatalf("%s", l.ListenAndServe())
}
