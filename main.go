package main

import (
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

	// static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))
}
