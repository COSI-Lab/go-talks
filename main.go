package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"golang.org/x/net/html"
)

var hub Hub
var tmpls *template.Template
var config Config
var trustedNetworks Networks

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

	configFile, err := os.Open("config.toml")
	if err != nil {
		log.Println("[WARN] Could not open config file, using defaults")
		config = DefaultConfig()
	} else {
		config = ParseConfig(configFile)
	}
	config.Validate()
	trustedNetworks = config.Network()

	// Connect to the database
	err = ConnectDB(&config)
	if err != nil {
		log.Fatalln("[ERROR] Failed to connect to the database", err)
	}

	// Set up all tables
	MakeDB()

	// Load templates and add markdown function
	tmpls = template.Must(template.New("").Funcs(template.FuncMap{"safe_markdown": markDownerSafe, "unsafe_markdown": markDownerUnsafe}).ParseGlob("templates/*.gohtml"))

	if err != nil {
		log.Fatalln("[ERROR] Failed to compile some template(s)", err)
	} else {
		log.Println("[INFO]", tmpls.DefinedTemplates())
	}

	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	// "api" endpoints
	r.HandleFunc("/talks", indexTalksHandler)
	r.HandleFunc("/{week:[0-9]{8}}/talks", talksHandler)
	r.HandleFunc("/ws", socketHandler)
	r.HandleFunc("/health", healthHandler)
	r.HandleFunc("/img/{id}", imageHandler)

	// static files
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	// templated pages
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/{week:[0-9]{8}}", weekHandler)

	// pages that are generated from markdown source in the posts directory
	posts := []string{"usage"}
	for _, post := range posts {
		r.HandleFunc("/"+post, markdownFactory(post))
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
	s.Every(1).Day().At("00:00").Do(invalidateCache)

	// Start servers
	var wg sync.WaitGroup
	for listener := range config.Listen {
		// Create http server
		srv := &http.Server{
			Handler:      r,
			Addr:         config.Listen[listener],
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}

		// Start http server
		log.Println("[INFO] Web server is now listening for connections on http://" + config.Listen[listener])
		wg.Add(1)
		go func() {
			log.Fatal(srv.ListenAndServe())
			wg.Done()
		}()
	}

	wg.Wait()
}

const extensions = blackfriday.NoIntraEmphasis | blackfriday.FencedCode | blackfriday.Autolink |
	blackfriday.Strikethrough | blackfriday.SpaceHeadings

// markDownerSafe accepts markdown and returns sanitized html
func markDownerSafe(args ...interface{}) template.HTML {

	s := blackfriday.Run([]byte(fmt.Sprintf("%s", args...)), blackfriday.WithExtensions(extensions))
	sanitized := bluemonday.UGCPolicy().SanitizeBytes(s)

	// Proxy any images and re-render the html
	doc, _ := html.Parse(bytes.NewReader(sanitized))
	findImagesAndCacheThem(doc)
	var buf bytes.Buffer
	html.Render(&buf, doc)
	return template.HTML(buf.String())
}

// markDowner accepts markdown and returns html
// NOT SAFE FOR USER INPUT
func markDownerUnsafe(args ...interface{}) template.HTML {
	s := blackfriday.Run([]byte(fmt.Sprintf("%s", args...)), blackfriday.WithExtensions(extensions))
	return template.HTML(s)
}
