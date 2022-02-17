package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

var talks_password = ""
var hub Hub
var tmpls *template.Template

var TZ *time.Location

// Logs request Method and request URI
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("[INFO] " + r.Method + " " + r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func main() {
	var err error
	TZ, err = time.LoadLocation("America/New_York")

	if err != nil {
		log.Fatalln("[ERROR] Failed to load timezone:", err)
	}

	talks_password = os.Getenv("TALKS_PASSWORD")

	// Connect to the database
	err = ConnectDB()
	if err != nil {
		log.Fatalln("[ERROR] Failed to connect to the database", err)
	}

	// Set up all tables
	MakeDB()

	// Load templates with markdown func
	tmpls = template.Must(template.New("").Funcs(template.FuncMap{"markdown": markDowner}).ParseGlob("templates/*.gohtml"))

	if err != nil {
		log.Fatalln("[ERROR] Failed to compile some template(s)", err)
	} else {
		log.Println("[INFO]", tmpls.DefinedTemplates())
	}

	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	// "api" json encoded endpoints
	r.HandleFunc("/talks", indexTalksHandler)
	r.HandleFunc("/{week:[0-9]{8}}/talks", talksHandler)
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

	// Schedule backup tasks
	s := gocron.NewScheduler(TZ)
	s.Wednesday().At("23:59").Do(backup)

	// Start server
	log.Println("[INFO] Web server is now listening for connections on http://localhost:" + port)
	log.Fatal(srv.ListenAndServe())
}

// Processes the markdown and returns the html
func markDowner(args ...interface{}) template.HTML {
	s := blackfriday.MarkdownCommon([]byte(fmt.Sprintf("%s", args...)))
	santize := bluemonday.UGCPolicy().SanitizeBytes(s)
	return template.HTML(santize)
}
