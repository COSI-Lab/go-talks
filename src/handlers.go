package main

import (
	"fmt"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {

}

func allHandler(w http.ResponseWriter, r *http.Request) {

}

func talksHandler(w http.ResponseWriter, r *http.Request) {

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Return list of active clients for diagnostic purposes
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprint(hub.countConnections())))
}
