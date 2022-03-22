package main

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/html"
)

// Malious users could send ip-grabbers as the src for images
// so we need to download those images so others never have to leave the site

// Cache maps image ids to their downloaded content
var cache = make(map[string]Image)
var cacheLock = sync.RWMutex{}

type Image struct {
	ContentType string
	Data        []byte
}

// cacheImage downloads the image at the given url, caches it,
// and returns the sha256 of the url as the id
func cacheImage(url string) (string, error) {
	// Url safe hash the url
	hash := sha256.Sum256([]byte(url))
	id := base64.URLEncoding.EncodeToString(hash[:])

	// Check if the image is already cached
	cacheLock.RLock()
	if _, ok := cache[id]; ok {
		cacheLock.RUnlock()
		return id, nil
	}
	cacheLock.RUnlock()

	// Download the image
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Cache the image
	cacheLock.Lock()
	cache[id] = Image{
		ContentType: resp.Header.Get("Content-Type"),
		Data:        content,
	}
	cacheLock.Unlock()

	return id, nil
}

// findImagesAndCacheThem locates all images beneath the given node and
// downloads them to the cache and rewrites the src attribute to use the
// cached image
func findImagesAndCacheThem(doc *html.Node) {
	if doc.Type == html.ElementNode && doc.Data == "img" {
		for i, attr := range doc.Attr {
			if attr.Key == "src" {
				// Download the image and cache it
				id, err := cacheImage(attr.Val)
				if err != nil {
					log.Println("[ERROR] Failed to download image", err)
				}

				// Replace the src attribute with the cached image
				doc.Attr[i].Val = "/img/" + id
			}
		}
	}

	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		findImagesAndCacheThem(c)
	}
}
