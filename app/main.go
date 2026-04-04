package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

var version = "v0.2.0"

func main() {
	hostname, _ := os.Hostname()
	var counter int64

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"version": version,
			"pod":     hostname,
		})
	})

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		id := atomic.AddInt64(&counter, 1)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":  id,
			"pod": hostname,
		})
	})

	http.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "accepted",
			"message": "notification queued",
		})
	})

	fmt.Printf("Notiflex API %s starting on :8080\n", version)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
