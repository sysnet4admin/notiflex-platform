package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

var version = "v0.1.0"
var counter int64

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/id", idHandler)

	log.Printf("Notiflex API %s starting on :8080", version)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"version": version,
	})
}

func idHandler(w http.ResponseWriter, r *http.Request) {
	id := atomic.AddInt64(&counter, 1)
	pod := os.Getenv("HOSTNAME")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"pod":     pod,
		"version": version,
	})
	fmt.Printf("Generated ID: %d on pod %s\n", id, pod)
}
