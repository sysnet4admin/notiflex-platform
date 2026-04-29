package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

const appVersion = "v0.1.1"

type server struct {
	mu      sync.Mutex
	counter int64
	podName string
}

func newServer() *server {
	podName := os.Getenv("HOSTNAME")
	if podName == "" {
		podName = "unknown-pod"
	}

	return &server{podName: podName}
}

func (s *server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *server) idHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.mu.Lock()
	s.counter++
	id := s.counter
	s.mu.Unlock()

	writeJSON(w, http.StatusOK, map[string]string{
		"id":           strconv.FormatInt(id, 10),
		"generated_by": s.podName,
	})
}

func (s *server) versionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"version":      appVersion,
		"generated_by": s.podName,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to encode json response: %v", err)
	}
}

func main() {
	srv := newServer()
	mux := http.NewServeMux()
	mux.HandleFunc("/health", srv.healthHandler)
	mux.HandleFunc("/id", srv.idHandler)
	mux.HandleFunc("/version", srv.versionHandler)

	addr := ":8080"
	log.Printf("notiflex-api listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(fmt.Errorf("server stopped: %w", err))
	}
}
