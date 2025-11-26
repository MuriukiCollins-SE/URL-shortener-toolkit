package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/big"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

// ========================== CONFIG ==========================
const (
	port            = ":8080"
	shortCodeLength = 7
	characters      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	// In-memory store
	urlStore = make(map[string]string)
	mu       sync.RWMutex

	// Point to YOUR actual template location
	tmpl = template.Must(template.ParseFiles("Templates/index.html"))

	// URL validation
	validURLRegex = regexp.MustCompile(`^(https?://)?[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\.[a-zA-Z]{2,}$|^(https?://)?[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/.*)?$`)
	schemeRegex   = regexp.MustCompile(`^https?://`)
)

// ========================== SHORT CODE GENERATOR ==========================
func generateShortCode() string {
	for {
		code := make([]byte, shortCodeLength)
		for i := range code {
			idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(characters))))
			code[i] = characters[idx.Int64()]
		}
		candidate := string(code)

		mu.RLock()
		if urlStore[candidate] != "" {
			mu.RUnlock()
			continue
		}
		mu.RUnlock()

		mu.Lock()
		urlStore[candidate] = ""
		mu.Unlock()
		return candidate
	}
}

// ========================== HANDLERS ==========================
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		redirectHandler(w, r)
		return
	}
	tmpl.Execute(w, nil)
}

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

	if !validURLRegex.MatchString(longURL) {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	if !schemeRegex.MatchString(longURL) {
		longURL = "https://" + longURL
	}

	shortCode := generateShortCode()
	mu.Lock()
	urlStore[shortCode] = longURL
	mu.Unlock()

	// Build short link using the request host + scheme (works on Render / deployed hosts)
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}
	host := r.Host
	shortLink := fmt.Sprintf("%s://%s/%s", scheme, host, shortCode)

	if r.Header.Get("Accept") == "application/json" {
		json.NewEncoder(w).Encode(map[string]string{
			"short": shortLink,
			"long":  longURL,
		})
		return
	}

	data := struct {
		ShortLink string
		LongURL   string
	}{shortLink, longURL}
	tmpl.Execute(w, data)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	// Normalize path and trim any leading/trailing slashes
	shortCode := strings.Trim(r.URL.Path, "/")
	if shortCode == "" {
		// no code => treat as homepage
		homeHandler(w, r)
		return
	}

	mu.RLock()
	longURL, ok := urlStore[shortCode]
	mu.RUnlock()

	if !ok || longURL == "" {
		http.Error(w, "Short link not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}

// ========================== MAIN ==========================
func main() {
	// Serve CSS from Templates/css/
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("Templates/css"))))

	// Routes
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			homeHandler(w, r)
		} else {
			redirectHandler(w, r)
		}
	})

	log.Println("")
	log.Println("GoShort URL Shortener is LIVE!")
	log.Printf("Open → http://localhost%s", port)
	log.Println("Stunning design loaded from Templates/ + Templates/css/style.css")
	log.Println("Press Ctrl+C to stop")
	log.Println("")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
