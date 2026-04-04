package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
)

var version = "v0.1.0"
var counter int64

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/id", idHandler)

	port := "8080"
	fmt.Printf("Notiflex API %s starting on :%s\n", version, port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Server failed: %v\n", err)
		os.Exit(1)
	}
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
}
