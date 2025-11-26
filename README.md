# GoShort — Lightning-Fast URL Shortener Built with Pure Go & Generative AI


**Collins Muriuki** 

![GoShort Preview](App%20screenshort.png)  
*(Live screenshot of the url-shortener)*

---

##  LIVE DEMO

** Free HTTPS**  
[https://url-shortener-toolkit-1.onrender.com](https://url-shortener-toolkit-1.onrender.com)

---

## GitHub Repository

[https://github.com/MuriukiCollins-SE/URL-shortener-toolkit](https://github.com/MuriukiCollins-SE/URL-shortener-toolkit)

---

## 📝 Project Overview

A fully functional URL shortener built in Go, featuring a cyberpunk UI and deployed for free on Render.

**Tech Stack:**  
- **Language:** Go (Golang) — standard library only  
- **Design:** Hand-crafted cyberpunk UI (matrix rain, neon glows, blinking cursor)  
- **Deployment:** Render (free tier, auto-HTTPS, global CDN)  
- **Storage:** In-memory map with thread-safety (`sync.RWMutex`)  
- **Random codes:** `crypto/rand` — cryptographically secure

---

## ✨ Features

- Shorten any URL in 1 click (e.g., google.com → https://...onrender.com/AbC123)
- Smart protocol detection: `http://` on localhost, `https://` on live site
- Thread-safe for multiple users
- Beautiful hacker-movie interface
- 100% working live deployment
- Zero external dependencies — just `go run main.go`

---

##  How to Run Locally

### 1. Clone the repo

```bash
git clone https://github.com/MuriukiCollins-SE/URL-shortener-toolkit.git
cd URL-shortener-toolkit
```

### 2. Run it

```bash
go run main.go
```

---

### Step-by-Step Installation by OS

#### Windows

1. [Install Go](https://go.dev/dl/) (choose Windows installer)
2. Open PowerShell or CMD:
    ```powershell
    git clone https://github.com/MuriukiCollins-SE/URL-shortener-toolkit.git
    cd URL-shortener-toolkit
    go run main.go
    ```
3. Open [http://localhost:8080](http://localhost:8080)

---

#### macOS

1. Install Go (via Homebrew or official .pkg):
    ```bash
    brew install go
    # or download from https://go.dev/dl/
    ```
2. Clone and run:
    ```bash
    git clone https://github.com/MuriukiCollins-SE/URL-shortener-toolkit.git
    cd URL-shortener-toolkit
    go run main.go
    ```
3. Open [http://localhost:8080](http://localhost:8080)

---

#### Linux (Ubuntu/Debian/Fedora/Pop!_OS etc.)

1. Install Go:
    ```bash
    wget https://go.dev/dl/go1.23.4.linux-amd64.tar.gz
    sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc && source ~/.bashrc
    go version   # should show go1.23+
    ```
2. Clone and run:
    ```bash
    git clone https://github.com/MuriukiCollins-SE/URL-shortener-toolkit.git
    cd URL-shortener-toolkit
    go run main.go
    ```
3. Open [http://localhost:8080](http://localhost:8080)

---