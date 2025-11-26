package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"math/big"
	"net/http"
	"sync"
)

// ========================== CONFIG & GLOBALS ==========================
const (
	port            = ":8080"
	shortCodeLength = 7
	characters      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	// In-memory storage: shortCode → original long URL
	urlStore = make(map[string]string)
	// Thread-safe access to the map (multiple users at once)
	mu sync.RWMutex

	// Pre-parsed HTML template (faster than parsing every request)
	tmpl = template.Must(template.ParseFiles("templates/index.html"))
)

// ========================== SHORT CODE GENERATOR ==========================
func generateShortCode() string {
	for {
		code := make([]byte, shortCodeLength)
		for i := range code {
			// crypto/rand = truly random (not predictable like math/rand)
			idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(characters))))
			code[i] = characters[idx.Int64()]
		}
		candidate := string(code)

		// Check for extremely rare collision
		mu.RLock()
		if urlStore[candidate] != "" {
			mu.RUnlock()
			continue // try again
		}
		mu.RUnlock()

		// Claim this code
		mu.Lock()
		urlStore[candidate] = "" // reserve it temporarily
		mu.Unlock()
		return candidate
	}
}

// ========================== HANDLERS ==========================

// GET / → Show the homepage with the shorten form
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	tmpl.Execute(w, nil)
}

// POST /shorten → Accept long URL → return short link
func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	longURL := r.FormValue("url")
	if longURL == "" {
		http.Error(w, "Missing URL", http.StatusBadRequest)
		return
	}

	// Ensure URL starts with http:// or https://
	if !regex.MatchString(longURL) {
		longURL = "https://" + longURL
	}

	shortCode := generateShortCode()

	// Save the mapping
	mu.Lock()
	urlStore[shortCode] = longURL
	mu.Unlock()

	shortLink := fmt.Sprintf("http://localhost%s/%s", port, shortCode)

	// Support both HTML and JSON responses
	if r.Header.Get("Accept") == "application/json" {
		json.NewEncoder(w).Encode(map[string]string{
			"short": shortLink,
			"long":  longURL,
		})
	} else {
		// Show result on the same page
		data := struct {
			ShortLink string
			LongURL   string
		}{shortLink, longURL}
		tmpl.Execute(w, data)
	}
}

// GET /{shortCode} → Redirect to original URL
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := r.URL.Path[1:] // remove leading slash

	mu.RLock()
	longURL, exists := urlStore[shortCode]
	mu.RUnlock()

	if !exists || longURL == "" {
		http.Error(w, "Short link not found", http.StatusNotFound)
		return
	}

	// Permanent redirect (real shorteners use 301)
	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}

// ========================== MAIN ==========================
func main() {
	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/", redirectHandler) // catch all other paths for redirects

	fmt.Println("URL Shortener is LIVE!")
	fmt.Printf("Open your browser → http://localhost%s\n", port)
	fmt.Println("Ctrl+C to stop")

	// Start server
	http.ListenAndServe(port, nil)
}