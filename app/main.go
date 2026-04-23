package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
)

var counter uint64

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func idHandler(w http.ResponseWriter, _ *http.Request) {
	nextID := atomic.AddUint64(&counter, 1)
	podName := os.Getenv("HOSTNAME")
	if podName == "" {
		podName = "unknown"
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"id":           strconv.FormatUint(nextID, 10),
		"generated_by": podName,
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/id", idHandler)

	addr := ":8080"
	log.Printf("notiflex api listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
