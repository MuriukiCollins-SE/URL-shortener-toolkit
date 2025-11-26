package main

import (
	"crypto/rand"
	"encoding/json"
	"html/template"
	"log"
	"math/big"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

// ========================== CONFIG ==========================
const (
	shortCodeLength = 7
	characters      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	urlStore = make(map[string]string)
	mu       sync.RWMutex

	// Auto-detect port for Render / localhost
	listenAddr = func() string {
		if port := os.Getenv("PORT"); port != "" {
			return "0.0.0.0:" + port
		}
		return ":8080"
	}()

	tmpl = template.Must(template.ParseFiles("Templates/index.html"))

	validURLRegex = regexp.MustCompile(`^(https?://)?[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/.*)?$`)
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
		if _, exists := urlStore[candidate]; exists {
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

	longURL := strings.TrimSpace(r.FormValue("url"))
	if longURL == "" {
		http.Error(w, "Missing URL", http.StatusBadRequest)
		return
	}

	if !validURLRegex.MatchString(longURL) {
		http.Error(w, "Invalid URL – try google.com or https://youtube.com", http.StatusBadRequest)
		return
	}
	if !schemeRegex.MatchString(longURL) {
		longURL = "https://" + longURL
	}

	shortCode := generateShortCode()
	mu.Lock()
	urlStore[shortCode] = longURL
	mu.Unlock()

	// FINAL BULLETPROOF PROTOCOL DETECTION
	scheme := "http" // default for localhost

	// 1. Direct TLS connection (very rare in dev)
	if r.TLS != nil {
		scheme = "https"
	}
	// 2. Proxy header (most common on Render, Railway, etc.)
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}
	// 3. Hard fallback for Render (100% reliable)
	if strings.Contains(r.Host, "onrender.com") || strings.Contains(r.Host, "render.com") {
		scheme = "https"
	}

	shortLink := scheme + "://" + r.Host + "/" + shortCode

	if r.Header.Get("Accept") == "application/json" {
		json.NewEncoder(w).Encode(map[string]string{
			"short": shortLink,
			"long":  longURL,
		})
		return
	}

	data := struct{ ShortLink, LongURL string }{shortLink, longURL}
	tmpl.Execute(w, data)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		return
	}
	shortCode := strings.TrimPrefix(r.URL.Path, "/")

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
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("Templates/css"))))

	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			homeHandler(w, r)
		} else {
			redirectHandler(w, r)
		}
	})

	log.Println("GoShort by Collins Muriuki — LIVE!")
	log.Println("Local → http://localhost:8080")
	log.Println("Deployed → https://url-shortener-toolkit.onrender.com")
	log.Printf("Listening on %s", listenAddr)

	log.Fatal(http.ListenAndServe(listenAddr, nil))
}