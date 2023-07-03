// Malicious users could send ip-grabbers (or other baddies) as the `src` for images
// Instead of exposing our users to this, we will proxy all images through our server
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

// Cache maps image ids to their downloaded content
var cache = make(map[[32]byte]Image)
var cacheLock = sync.RWMutex{}

// Image is what we store in the cache
type Image struct {
	ContentType string
	Data        []byte
}

// cacheImage downloads the image at the given url, caches it,
// and returns the sha256 of the url as the id
func cacheImage(url string) ([32]byte, error) {
	// Hash the image url for storage in the cache
	hash := sha256.Sum256([]byte(url))

	// Check if the image is already cached
	cacheLock.RLock()
	if _, ok := cache[hash]; ok {
		cacheLock.RUnlock()
		return hash, nil
	}
	cacheLock.RUnlock()

	// Download the image
	resp, err := http.Get(url)
	if err != nil {
		log.Println("[WARN] failed to make request to", url, err)
		return [32]byte{}, err
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[WARN] failed to read response body from", url, err)
		return [32]byte{}, err
	}

	log.Printf("[INFO] downloaded %s (%d bytes)\n", url, len(content))

	// Cache the image
	cacheLock.Lock()
	cache[hash] = Image{
		ContentType: resp.Header.Get("Content-Type"),
		Data:        content,
	}
	cacheLock.Unlock()

	return hash, nil
}

// findImagesAndCacheThem locates all images beneath the given html node and downloads
// them to the cache and rewrites the src attribute to use the cached image
func findImagesAndCacheThem(doc *html.Node) {
	if doc.Type == html.ElementNode && doc.Data == "img" {
		for i, attr := range doc.Attr {
			if attr.Key == "src" {
				// Download the image and cache it
				hash, err := cacheImage(attr.Val)
				if err != nil {
					log.Println("[WARN] Failed to download image", err)
				}
				id := base64.URLEncoding.EncodeToString(hash[:])

				// Rewrite the src attribute to use the cached image
				doc.Attr[i].Val = "/img/" + id
			}
		}
	}

	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		findImagesAndCacheThem(c)
	}
}

// Cache invalidation: remove all images from the cache
func invalidateCache() {
	// Log how much data we are removing from the cache
	var totalSize int
	cacheLock.RLock()
	for _, image := range cache {
		totalSize += len(image.Data)
	}
	cacheLock.RUnlock()

	log.Println("[INFO] clearing the cache, saving", totalSize, "bytes")

	cacheLock.Lock()
	cache = make(map[[32]byte]Image)
	cacheLock.Unlock()
}
